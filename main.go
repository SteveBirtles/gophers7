package main

import (
	"fmt"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"image"
	"image/color"
	"math"
	"time"
)

const (
	WIDTH   = 1024
	HEIGHT  = 768
	THREADS = 4
)

type coordinate struct {
	x, y  int
	value float64
}

var (
	fractal    [WIDTH][HEIGHT]float64
	iterations = 500.0
	x, y       = -0.5, 0.0
	scale      = 0.002
	infinity   = 1e+75
	startTime  time.Time
	done       bool
	window     *pixelgl.Window
	colors     = [][][3]float64{
		{{0, 0, 255}, {255, 0, 255}, {255, 0, 0}, {255, 255, 0}, {0, 255, 0}, {0, 255, 255}, {0, 0, 0}},
		{{0, 0, 0}, {255, 255, 255}, {0, 0, 0}},
		{{255, 255, 255}, {0, 0, 0}},
		{{0, 0, 0}, {255, 0, 0}, {255, 255, 0}, {255, 0, 0}, {255, 255, 0}, {0, 0, 0}},
		{{255, 255, 255}, {255, 0, 0}, {0, 0, 0}, {255, 255, 0}, {255, 255, 255}, {0, 0, 255}, {0, 0, 64}},
	}
	palette = 0
	power   = 1.0
	inputs  = make(chan coordinate, WIDTH*HEIGHT)
	outputs = make(chan coordinate, WIDTH*HEIGHT)
)

func worker(inputs <-chan coordinate, outputs chan<- coordinate) {

	for point := range inputs {

		c := complex(float64(point.x-WIDTH/2)*scale+x, float64(point.y-HEIGHT/2)*scale+y)
		z := complex(0, 0)
		var i int
		for i = 0; i < int(iterations); i++ {
			z = z*z + c
			if imag(z) > infinity || real(z) > infinity {
				break
			}
		}
		point.value = float64(i) / float64(iterations)

		outputs <- point

	}

}

func clearInputs() {
	for {
		select {
		case <-inputs:
		default:
			return
		}
	}
}

func resetInputs() {
	done = false
	startTime = time.Now()
	window.SetTitle(fmt.Sprintf("X: %16.16f Y: %16.16f Scale: %16.16f Iterations: %d Power: %f", x, y, scale, int(iterations), power))

	clearInputs()

	for i := 0; i < WIDTH/2; i++ {
		for j := 0; j < HEIGHT; j++ {
			inputs <- coordinate{WIDTH/2 + i, j, 0}
			inputs <- coordinate{WIDTH/2 - i - 1, j, 0}
		}
	}
}

func calculateColor(value float64, p int) color.RGBA {

	if p >= len(colors) {
		panic("Palette number too high!")
	}

	value = math.Pow(value, power)

	base, frac := math.Modf(value * float64(len(colors[p])-1))

	index := int(base) + 1

	r := uint8(colors[p][index-1][0])
	g := uint8(colors[p][index-1][1])
	b := uint8(colors[p][index-1][2])

	if frac != 0 {
		r += uint8((colors[p][index][0] - colors[p][index-1][0]) * frac)
		g += uint8((colors[p][index][1] - colors[p][index-1][1]) * frac)
		b += uint8((colors[p][index][2] - colors[p][index-1][2]) * frac)
	}

	return color.RGBA{r, g, b, 0xff}

}

func run() {

	cfg := pixelgl.WindowConfig{
		Title:  "Pixel Rocks!",
		Bounds: pixel.R(0, 0, WIDTH, HEIGHT),
		VSync:  true,
	}
	var err error
	window, err = pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	image := image.NewNRGBA(image.Rect(0, 0, WIDTH, HEIGHT))

	picture := pixel.PictureDataFromImage(image)
	sprite := pixel.NewSprite(picture, picture.Bounds())

	second := time.Tick(time.Second)

	resetInputs()

	for w := 0; w < THREADS; w++ {
		go worker(inputs, outputs)
	}

	for !window.Closed() {

		window.Clear(color.RGBA{0, 0, 0, 0xff})

		if window.JustPressed(pixelgl.KeyEscape) {
			window.SetClosed(true)
		}

		if window.JustPressed(pixelgl.MouseButtonLeft) {
			mouse := window.MousePosition()
			x = (mouse.X-WIDTH/2)*scale + x
			y = (HEIGHT/2-mouse.Y)*scale + y
			scale /= 4
			resetInputs()
		}

		if window.JustPressed(pixelgl.MouseButtonRight) {
			scale *= 4
			resetInputs()
		}

		if window.JustPressed(pixelgl.KeyPageUp) {
			iterations *= 1.25
			resetInputs()
		}

		if window.JustPressed(pixelgl.KeyPageDown) {
			iterations /= 1.25
			resetInputs()
		}

		if window.JustPressed(pixelgl.KeyUp) {
			power *= 1.25
			resetInputs()
		}

		if window.JustPressed(pixelgl.KeyDown) {
			power /= 1.25
			resetInputs()
		}

		if window.JustPressed(pixelgl.Key1) {
			palette = 0
			resetInputs()
		}
		if window.JustPressed(pixelgl.Key2) {
			palette = 1
			resetInputs()
		}
		if window.JustPressed(pixelgl.Key3) {
			palette = 2
			resetInputs()
		}
		if window.JustPressed(pixelgl.Key4) {
			palette = 3
			resetInputs()
		}
		if window.JustPressed(pixelgl.Key5) {
			palette = 4
			resetInputs()
		}

	receiver:
		for !done {
			select {
			case c := <-outputs:
				fractal[c.x][c.y] = c.value
			default:
				if len(inputs) == 0 {
					fmt.Println("Processing complete in", time.Since(startTime).Seconds(), "seconds.")
					done = true
				}
				break receiver
			}
		}

		select {
		case <-second:

			for i := 0; i < WIDTH; i++ {
				for j := 0; j < HEIGHT; j++ {
					c := calculateColor(fractal[i][j], palette)
					image.Set(i, j, c)
				}
			}

			picture := pixel.PictureDataFromImage(image)
			sprite = pixel.NewSprite(picture, picture.Bounds())

		default:
		}

		sprite.Draw(window, pixel.IM.Moved(window.Bounds().Center()))

		window.Update()
	}
}

func main() {
	pixelgl.Run(run)
}
