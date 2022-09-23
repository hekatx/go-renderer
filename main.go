package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"

	"github.com/deeean/go-vector/vector3"
	"github.com/hekatx/go-renderer/draw"
	"github.com/hekatx/go-renderer/obj"
)

// TODO: Implement color based on light's direction (end of lesson 2 of tinyrenderer)
// Make use of vector3's for vertices coords as it will make it easier to calculate the
// normals
func renderModel(model obj.Model, width, height int, image *image.RGBA) {
	var screen_coords [3][2]float64
	var world_coords [3]vector3.Vector3
	light_dir := vector3.New(0., 0., -1.)
	for i := 0; i < len(model.Faces); i++ {
		face := model.Faces[i]
		for j := 0; j < 3; j++ {
			v := model.Vertices[face[j]]
			screen_coords[j] = [2]float64{(v.X + 1.) * float64(width) / 2., (v.Y + 1.) * float64(height) / 2.}
			world_coords[j] = v
		}
		a := world_coords[0]
		b := world_coords[1]
		c := world_coords[2]
		qr := c.Sub(&a)
		qs := b.Sub(&a)
		normal := qr.Cross(qs)
		normal = normal.Normalize()

		intensity := normal.Dot(light_dir)
		random_color := color.RGBA{uint8(intensity * 255), uint8(intensity * 255), uint8(intensity * 255), 255}

		if intensity > 0 {
			draw.Triangle(
				image,
				random_color,
				[]draw.Point{
					{
						X: float32(screen_coords[0][0]),
						Y: float32(screen_coords[0][1]),
					},
					{
						X: float32(screen_coords[1][0]),
						Y: float32(screen_coords[1][1]),
					},
					{
						X: float32(screen_coords[2][0]),
						Y: float32(screen_coords[2][1]),
					},
				},
			)
		}
	}
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

func main() {
	width := 600
	height := 600
	base := draw.NewImage(width, height)

	head_model := obj.Decode("./african_head.obj")

	renderModel(head_model, width, height, base)
	base = flipVertically(base)

	f, _ := os.Create("image.png")
	err := png.Encode(f, base)

	if err != nil {
		fmt.Fprintf(os.Stdout, "fucked up %v", err)
	}
}
