package render

import (
	"image"
	"image/color"

	"github.com/deeean/go-vector/vector3"
	"github.com/hekatx/go-renderer/draw"
	"github.com/hekatx/go-renderer/obj"
)

func Model(model obj.Model, width, height int, image *image.RGBA, light_dir vector3.Vector3) {
	var screen_coords [3][2]float64
	var world_coords [3]vector3.Vector3

	for i := 0; i < len(model.Faces); i++ {
		face := model.Faces[i]
		for j := 0; j < 3; j++ {
			v := model.Vertices[face[j]]
			screen_coords[j] = [2]float64{(v.X + 1.) * float64(width) / 2., (v.Y + 1.) * float64(height) / 2.}
			world_coords[j] = v
		}
		n := calculateNormal(world_coords[0], world_coords[1], world_coords[2])

		intensity := n.Dot(&light_dir)
		color := color.RGBA{uint8(intensity * 255), uint8(intensity * 255), uint8(intensity * 255), 255}
		triangle_vertices := []draw.Point{
			{
				X: float64(screen_coords[0][0]),
				Y: float64(screen_coords[0][1]),
			},
			{
				X: float64(screen_coords[1][0]),
				Y: float64(screen_coords[1][1]),
			},
			{
				X: float64(screen_coords[2][0]),
				Y: float64(screen_coords[2][1]),
			},
		}

		if intensity > 0 {
			draw.Triangle(image, color, triangle_vertices)
		}
	}
	flipVertically(image)
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

			draw.Line(x0, y0, x1, y1, image, white)
		}
	}
	flipVertically(image)
}

func calculateNormal(a, b, c vector3.Vector3) vector3.Vector3 {
	qr := c.Sub(&a)
	qs := b.Sub(&a)
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
