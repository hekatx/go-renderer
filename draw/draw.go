package draw

import (
	"image"
	"image/color"
	"math"
	"sort"

	"github.com/deeean/go-vector/vector3"
)

type Point struct {
	X float64
	Y float64
	Z float64
}

type BarycentricCoordinates struct {
	w1 float64
	w2 float64
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

func barycentric(pts []vector3.Vector3, P *vector3.Vector3) vector3.Vector3 {
	v0 := vector3.Vector3{X: float64(pts[1].X - pts[0].X), Y: float64(pts[2].X - pts[0].X), Z: float64(pts[0].X - P.X)}
	v1 := vector3.Vector3{X: float64(pts[1].Y - pts[0].Y), Y: float64(pts[2].Y - pts[0].Y), Z: float64(pts[0].Y - P.Y)}
	u := v0.Cross(&v1)
	if math.Abs(u.Z) > 1e-2 {
		return vector3.Vector3{X: 1.0 - (u.X+u.Y)/u.Z, Y: u.X / u.Z, Z: u.Y / u.Z}
	}
	return vector3.Vector3{X: -1., Y: 1., Z: 1.}
}

func getBarycentricCoords(p, a, b, c Point) BarycentricCoordinates {
	var bc BarycentricCoordinates
	h := c.Y - a.Y
	bc.w1 = ((a.X * h) + ((p.Y - a.Y) * (c.X - a.X)) - (p.X * h)) / (((b.Y - a.Y) * (c.X - a.X)) - ((b.X - a.X) * h))
	bc.w2 = (p.Y - a.Y - (bc.w1 * (b.Y - a.Y))) / h
	return bc
}

func isPointOutsideTriangle(w1, w2, w3 float64) bool {
	return w1 < 0 || w2 < 0 || w3 < 0
}

func Triangle(i *image.RGBA, c color.RGBA, v []vector3.Vector3, zb *[]float64, w, h int) {
	sort.Slice(v, func(i, j int) bool {
		return v[i].Y < v[j].Y
	})

	p := vector3.Vector3{}
	for y := 0; y < i.Bounds().Max.Y; y++ {
		for x := 0; x < i.Bounds().Max.X; x++ {
			p.X = float64(x)
			p.Y = float64(y)
			bc := barycentric(v, &p)

			if isPointOutsideTriangle(bc.X, bc.Y, bc.Z) {
				continue
			}

			p.Z = (bc.X * v[0].Z) + (bc.Y * v[1].Z) + (bc.Z * v[2].Z)

			if (*zb)[int(p.X+p.Y*float64(w))] < p.Z {
				(*zb)[int(p.X+p.Y*float64(w))] = p.Z
				i.Set(int(p.X), int(p.Y), image.NewUniform(c))
			}
		}
	}
}
