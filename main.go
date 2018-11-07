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



}

func resetInputs(inputs chan coordinate) {



}

func calculateColor(value float64, p int) color.RGBA {

	if p >= len(colors) {
		panic("Palette number too high!")
	}

	value = math.Pow(value, power)
	base, frac := math.Modf(value * float64(len(colors[p])-1))

	index := int(base) + 1

	if index < 0 || index >= len(colors[p]) {
		return color.RGBA{0, 0, 0, 0}
	}

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

	/* --- Prepare the input and output channels --- */



	/*-----------------------------------------------*/


	/* --- Launch the workers! --- */



	/*-----------------------------------------------*/

	ticker := time.Tick(time.Millisecond * 100)

	for !window.Closed() {

		if window.JustPressed(pixelgl.KeyEscape) {
			window.SetClosed(true)
		}

		if window.JustPressed(pixelgl.MouseButtonLeft) {
			mouse := window.MousePosition()
			x = (mouse.X-WIDTH/2)*scale + x
			y = (HEIGHT/2-mouse.Y)*scale + y
			scale /= 4
			//resetInputs(inputs)
		}
		if window.JustPressed(pixelgl.MouseButtonRight) {
			scale *= 4
			//resetInputs(inputs)
		}

		if window.JustPressed(pixelgl.KeyPageDown) {
			scale *= 1.25
			//resetInputs(inputs)
		}
		if window.JustPressed(pixelgl.KeyPageUp) {
			scale /= 1.25
			//resetInputs(inputs)
		}

		if window.JustPressed(pixelgl.KeyW) {
			iterations *= 1.25
			//resetInputs(inputs)
		}
		if window.JustPressed(pixelgl.KeyS) {
			iterations /= 1.25
			//resetInputs(inputs)
		}

		if window.JustPressed(pixelgl.KeyD) {
			power *= 1.25
			//resetInputs(inputs)
		}
		if window.JustPressed(pixelgl.KeyA) {
			power /= 1.25
			//resetInputs(inputs)
		}

		if window.JustPressed(pixelgl.KeyUp) {
			y -= (WIDTH / 10) * scale
			//resetInputs(inputs)
		}
		if window.JustPressed(pixelgl.KeyDown) {
			y += (WIDTH / 10) * scale
			//resetInputs(inputs)
		}
		if window.JustPressed(pixelgl.KeyLeft) {
			x -= (WIDTH / 10) * scale
			//resetInputs(inputs)
		}
		if window.JustPressed(pixelgl.KeyRight) {
			x += (WIDTH / 10) * scale
			//resetInputs(inputs)
		}

		if window.JustPressed(pixelgl.Key1) {
			palette = 0
			//resetInputs(inputs)
		}
		if window.JustPressed(pixelgl.Key2) {
			palette = 1
			//resetInputs(inputs)
		}
		if window.JustPressed(pixelgl.Key3) {
			palette = 2
			//resetInputs(inputs)
		}
		if window.JustPressed(pixelgl.Key4) {
			palette = 3
			//resetInputs(inputs)
		}
		if window.JustPressed(pixelgl.Key5) {
			palette = 4
			//resetInputs(inputs)
		}

		select {
		case <-ticker:

			/* --- Every tick (100ms) drain the output channel and update image --- */



			/*-----------------------------------------------*/

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
