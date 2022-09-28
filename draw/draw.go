package draw

import (
	"image"
	"image/color"
	"math"

	"github.com/deeean/go-vector/vector2"
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

func barycentric(pts [3]vector3.Vector3, P *vector3.Vector3) vector3.Vector3 {
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

func Triangle(i *image.RGBA, c color.RGBA, v [3]vector3.Vector3, zb *[]float64, w, h int) {
	bboxmin := vector2.Vector2{X: math.Inf(1), Y: math.Inf(1)}
	bboxmax := vector2.Vector2{X: math.Inf(-1), Y: math.Inf(-1)}
	clamp := vector2.Vector2{X: float64(w - 1), Y: float64(h - 1)}

	for i := 0; i < len(v); i++ {
		bboxmin.X = math.Max(0.0, math.Min(bboxmin.X, float64(v[i].X)))
		bboxmin.Y = math.Max(0.0, math.Min(bboxmin.Y, float64(v[i].Y)))
		bboxmax.X = math.Min(clamp.X, math.Max(bboxmax.X, float64(v[i].X)))
		bboxmax.Y = math.Min(clamp.Y, math.Max(bboxmax.Y, float64(v[i].Y)))
	}

	p := &vector3.Vector3{}
	for p.X = bboxmin.X; p.X < bboxmax.X; p.X++ {
		for p.Y = bboxmin.Y; p.Y < bboxmax.Y; p.Y++ {
			bc := barycentric(v, p)

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
