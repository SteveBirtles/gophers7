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
	WIDTH    = 1024
	HEIGHT   = 768
	THREADS  = 4
	INFINITY = 1e+75
)

type coordinate struct {
	x, y  int
	value float64
}

var (
	iterations = 500.0
	x, y       = -0.5, 0.0
	scale      = 0.002
	colors     = [][][3]float64{
		{{0, 0, 255}, {255, 0, 255}, {255, 0, 0}, {255, 255, 0}, {0, 255, 0}, {0, 255, 255}, {0, 0, 0}},
		{{0, 0, 0}, {255, 255, 255}, {0, 0, 0}},
		{{255, 255, 255}, {0, 0, 0}},
		{{0, 0, 0}, {255, 0, 0}, {255, 255, 0}, {255, 0, 0}, {255, 255, 0}, {0, 0, 0}},
		{{255, 255, 255}, {255, 0, 0}, {0, 0, 0}, {255, 255, 0}, {255, 255, 255}, {0, 0, 255}, {0, 0, 64}},
	}
	palette = 0
	power   = 1.0
)

func worker(inputs chan coordinate, outputs chan coordinate) {

	for point := range inputs {

		/* --- Calculate value for point --- */

		c := complex(float64(point.x-WIDTH/2)*scale+x, float64(point.y-HEIGHT/2)*scale+y)
		z := complex(0, 0)
		var i int
		for i = 0; i < int(iterations); i++ {
			z = z*z + c
			if imag(z) > INFINITY || real(z) > INFINITY {
				break
			}
		}
		point.value = float64(i) / float64(iterations)

		outputs <- point

	}

}

func resetInputs(inputs chan coordinate) {

	/* --- Clear the input channel --- */

clear:
	for {
		select {
		case <-inputs:
		default:
			break clear
		}
	}

	/* --- Fill the input channel --- */

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
	window, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	img := image.NewNRGBA(image.Rect(0, 0, WIDTH, HEIGHT))
	pic := pixel.PictureDataFromImage(img)
	sprite := pixel.NewSprite(pic, pic.Bounds())

	/*-----------------------------------------------*/

	/* --- Prepare the input and output channels --- */

	inputs := make(chan coordinate, WIDTH*HEIGHT)
	outputs := make(chan coordinate, WIDTH*HEIGHT)
	resetInputs(inputs)

	/* --- Launch the workers! --- */

	for w := 0; w < THREADS; w++ {
		go worker(inputs, outputs)
	}

	/*-----------------------------------------------*/

	second := time.Tick(time.Second)

	for !window.Closed() {

		if window.JustPressed(pixelgl.KeyEscape) {
			window.SetClosed(true)
		}

		if window.JustPressed(pixelgl.MouseButtonLeft) {
			mouse := window.MousePosition()
			x = (mouse.X-WIDTH/2)*scale + x
			y = (HEIGHT/2-mouse.Y)*scale + y
			scale /= 4
			resetInputs(inputs)
		}
		if window.JustPressed(pixelgl.MouseButtonRight) {
			scale *= 4
			resetInputs(inputs)
		}

		if window.JustPressed(pixelgl.KeyPageDown) {
			scale *= 1.25
			resetInputs(inputs)
		}
		if window.JustPressed(pixelgl.KeyPageUp) {
			scale /= 1.25
			resetInputs(inputs)
		}

		if window.JustPressed(pixelgl.KeyW) {
			iterations *= 1.25
			resetInputs(inputs)
		}
		if window.JustPressed(pixelgl.KeyS) {
			iterations /= 1.25
			resetInputs(inputs)
		}

		if window.JustPressed(pixelgl.KeyD) {
			power *= 1.25
			resetInputs(inputs)
		}
		if window.JustPressed(pixelgl.KeyA) {
			power /= 1.25
			resetInputs(inputs)
		}

		if window.JustPressed(pixelgl.KeyUp) {
			y -= (WIDTH / 10) * scale
			resetInputs(inputs)
		}
		if window.JustPressed(pixelgl.KeyDown) {
			y += (WIDTH / 10) * scale
			resetInputs(inputs)
		}
		if window.JustPressed(pixelgl.KeyLeft) {
			x -= (WIDTH / 10) * scale
			resetInputs(inputs)
		}
		if window.JustPressed(pixelgl.KeyRight) {
			x += (WIDTH / 10) * scale
			resetInputs(inputs)
		}

		if window.JustPressed(pixelgl.Key1) {
			palette = 0
			resetInputs(inputs)
		}
		if window.JustPressed(pixelgl.Key2) {
			palette = 1
			resetInputs(inputs)
		}
		if window.JustPressed(pixelgl.Key3) {
			palette = 2
			resetInputs(inputs)
		}
		if window.JustPressed(pixelgl.Key4) {
			palette = 3
			resetInputs(inputs)
		}
		if window.JustPressed(pixelgl.Key5) {
			palette = 4
			resetInputs(inputs)
		}

		/*-----------------------------------------------*/

		select {
		case <-second:

			/* --- Every tick (once per second) drain the output channel and update image --- */

		receiver:
			for {
				select {
				case p := <-outputs:
					img.Set(p.x, p.y, calculateColor(p.value, palette))
				default:
					break receiver
				}
			}

			pic := pixel.PictureDataFromImage(img)
			sprite = pixel.NewSprite(pic, pic.Bounds())

			window.SetTitle(fmt.Sprintf("X: %16.16f Y: %16.16f Scale: %16.16f Iterations: %d Power: %f", x, y, scale, int(iterations), power))

		default:
		}

		window.Clear(color.RGBA{0, 0, 0, 0xff})

		sprite.Draw(window, pixel.IM.Moved(window.Bounds().Center()))

		window.Update()
	}
}

func main() {
	pixelgl.Run(run)
}
