package main

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"image/color"
	"image"
	"time"
	"fmt"
	"math"
)

const (
	WIDTH = 1024
	HEIGHT = 768
)

type fractalPoint struct {
	x, y int
	value float64
}

var (
	fractal    [WIDTH][HEIGHT]float64
	iterations  = 500
	x, y              = -0.5, 0.0
	scale            = 0.002
	inputs          = make(chan fractalPoint, WIDTH*HEIGHT)
	outputs        = make(chan fractalPoint, WIDTH*HEIGHT)
	infexp          = 50.0
	infinity      = math.Pow(10, infexp)
	startTime  time.Time
	done       bool
	window     *pixelgl.Window
)


func worker(inputs <-chan fractalPoint, outputs chan<- fractalPoint) {

	for point := range inputs {

		c := complex(float64(point.x-WIDTH/2)*scale+x, float64(point.y-HEIGHT/2)*scale+y)

		z := complex(0, 0)

		var i int

		for i = 0; i < iterations; i++ {
			z = z*z + c
			if imag(z) > infinity || real(z) > infinity {
				break
			}
		}

		point.value = float64(i) / float64(iterations)

		outputs <- point

	}

}

func resetInputs() {
	done = false
	startTime = time.Now()
	window.SetTitle(fmt.Sprintf("X: %16.16f Y: %16.16f Scale: %16.16f Iterations: %d Infinity: 10^%d", x, y, scale, iterations, int(infexp)))
	for len(inputs) > 0 { <- inputs }
	for i := 0; i < WIDTH/2; i++ {
		for j := 0; j < HEIGHT; j++ {
			inputs <- fractalPoint{WIDTH/2+ i, j, 0}
			inputs <- fractalPoint{WIDTH/2- i -1, j, 0}
		}
	}
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

	for w := 0; w < 4; w++ {
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
			mouse := window.MousePosition()
			x = (mouse.X-WIDTH/2)*scale + x
			y = (HEIGHT/2-mouse.Y)*scale + y
			scale *= 4
			resetInputs()
		}

		if window.JustPressed(pixelgl.KeyPageUp) {
			iterations *= 2
			resetInputs()
		}

		if window.JustPressed(pixelgl.KeyPageDown) {
			iterations /= 2
			resetInputs()
		}

		if window.JustPressed(pixelgl.KeyEqual) {
			infexp *= 1.5
			infinity = math.Pow(10, infexp)
			resetInputs()
		}

		if window.JustPressed(pixelgl.KeyMinus) {
			infexp /= 1.5
			infinity = math.Pow(10, infexp)
			resetInputs()
		}

	receiver:
		for !done {
			select {
			case point := <-outputs:
				fractal[point.x][point.y] = point.value
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

					value := fractal[i][j] * 6
					var c color.RGBA

					if value < 1 {
						c = color.RGBA{uint8(255 * value), 0, 255, 0xff}
					} else if value < 2 {
						value -= 1
						c = color.RGBA{255, uint8(255 * value), uint8(255 * (1.0 - value)), 0xff}
					} else if value < 3 {
						value -= 2
						c = color.RGBA{uint8(255 * (1.0 - value)), 255, 0, 0xff}
					} else if value < 4 {
						value -= 3
						c = color.RGBA{0, 255, uint8(255 * value), 0xff}
					} else if value < 5 {
						value -= 4
						c = color.RGBA{0, uint8(255 * (1.0 - value)), 255, 0xff}
					} else {
						value -= 5
						c = color.RGBA{0, 0, uint8(255 * (1.0 - value)), 0xff}
					}

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