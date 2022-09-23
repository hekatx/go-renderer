package draw

import (
	"image"
	"image/color"
	"math"
	"sort"

	"github.com/hekatx/go-renderer/obj"
)

type Point struct {
	X float32
	Y float32
}

type BarycentricCoordinates struct {
	w1 float32
	w2 float32
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

func RenderWireframe(model obj.Model, width, height int, image *image.RGBA) {
	white := color.RGBA{255, 255, 255, 0xff}
	for i := 0; i < len(model.Faces); i++ {
		face := model.Faces[i]
		for j := 0; j < len(face); j++ {
			v0 := model.Vertices[face[j]]
			v1 := model.Vertices[face[(j+1)%len(face)]]

			x0 := (v0.X + 1.) * float64(width) / 2.
			y0 := (v0.Y + 1.) * float64(height) / 2.
			x1 := (v1.X + 1.) * float64(width) / 2.
			y1 := (v1.Y + 1.) * float64(height) / 2.

			Line(x0, y0, x1, y1, image, white)
		}
	}
}

func getBarycentricCoords(p, a, b, c Point) BarycentricCoordinates {
	var bc BarycentricCoordinates
	h := c.Y - a.Y
	bc.w1 = ((a.X * h) + ((p.Y - a.Y) * (c.X - a.X)) - (p.X * h)) / (((b.Y - a.Y) * (c.X - a.X)) - ((b.X - a.X) * h))
	bc.w2 = (p.Y - a.Y - (bc.w1 * (b.Y - a.Y))) / h
	return bc
}

func isPointInsideTriangle(w1, w2 float64) bool {
	return w1 >= 0 && w2 >= 0 && (w1+w2) <= 1.
}

func Triangle(i *image.RGBA, c color.RGBA, v []Point) {
	sort.Slice(v, func(i, j int) bool {
		return v[i].Y < v[j].Y
	})

	for y := 0; y < i.Bounds().Max.X; y++ {
		for x := 0; x < i.Bounds().Max.Y; x++ {
			ip := Point{float32(x), float32(y)}
			bc := getBarycentricCoords(ip, v[2], v[0], v[1])
			if isPointInsideTriangle(float64(bc.w1), float64(bc.w2)) {
				i.Set(x, y, image.NewUniform(c))
			}
		}
	}
}
