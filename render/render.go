package render

import (
	"image"
	"image/color"
	"math"

	"github.com/deeean/go-vector/vector3"
	"github.com/hekatx/go-renderer/draw"
	"github.com/hekatx/go-renderer/obj"
)

func Model(model obj.Model, width, height int, image *image.RGBA, light_dir vector3.Vector3) {
	var screen_coords [3]vector3.Vector3
	var world_coords [3]vector3.Vector3

	var zbuffer = make([]float64, width*height)
	for i := 0; i < len(zbuffer); i++ {
		zbuffer[i] = math.Inf(-1)
	}

	for i := 0; i < len(model.Faces); i++ {
		f := model.Faces[i]

		for j := 0; j < 3; j++ {
			v := model.Vertices[f[j]]
			screen_coords[j] = vector3.Vector3{X: float64(int((v.X+1.)*float64(width)/2. + .5)), Y: float64(int((v.Y+1.)*float64(height)/2. + .5)), Z: v.Z}
			world_coords[j] = v
		}

		n := calculateNormal(world_coords)
		intensity := n.Dot(&light_dir)
		c := color.RGBA{uint8(intensity * 255), uint8(intensity * 255), uint8(intensity * 255), 255}

		if intensity > 0 {
			draw.Triangle(image, c, screen_coords, &zbuffer, width, height)
		}
	}
	flipVertically(image)
}

func Wireframe(model obj.Model, width, height int, image *image.RGBA) {
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

			draw.Line(x0, y0, x1, y1, image, white)
		}
	}
	flipVertically(image)
}

func calculateNormal(pts [3]vector3.Vector3) vector3.Vector3 {
	qr := pts[2].Sub(&pts[0])
	qs := pts[1].Sub(&pts[0])
	normal := qr.Cross(qs)
	return *normal.Normalize()
}

func flipVertically(canvas *image.RGBA) *image.RGBA {
	bounds := canvas.Bounds()
	flipped := image.NewRGBA(image.Rect(0, 0, bounds.Max.X, bounds.Max.Y))
	for i := 0; i <= bounds.Max.X; i++ {
		for j := 0; j <= bounds.Max.Y; j++ {
			flipped.Set(i, bounds.Max.Y-j-1, canvas.At(i, j))
		}
	}
	return flipped
}
