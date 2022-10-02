package draw

import (
	"image"
	"image/color"
	"math"
)

type Point struct {
	X float64
	Y float64
	Z float64
}

type BarycentricCoordinates struct {
	W1 float64
	W2 float64
}

func Line(x0, y0, x1, y1 float64, i *image.RGBA, c color.RGBA) {
	steep := false

	if math.Abs(x0-x1) < math.Abs(y0-y1) {
		x0, y0 = y0, x0
		x1, y1 = y1, x1
		steep = true
	}

	if x0 > x1 {
		x0, x1 = x1, x0
		y0, y1 = y1, y0
	}

	dx := x1 - x0
	dy := y1 - y0

	derr := math.Abs(dy) * 2
	error := 0.0
	y := y0

	for x := x0; x <= x1; x++ {
		if steep {
			i.Set(int(y), int(x), c)
		} else {
			i.Set(int(x), int(y), c)
		}
		error += derr
		if error > dx {
			if y1 > y0 {
				y += 1
			} else {
				y -= 1
			}
			error -= dx * 2
		}

	}
}

func NewImage(w, h int) *image.RGBA {
	topLeft := image.Point{0, 0}
	bottomRight := image.Point{w, h}

	img := image.NewRGBA(image.Rectangle{topLeft, bottomRight})

	black := color.RGBA{0, 0, 0, 0xff}

	for y := 0; y < w; y++ {
		for x := 0; x < h; x++ {
			img.Set(x, y, image.NewUniform(black))
		}
	}

	return img
}
