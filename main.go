package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"

	"github.com/hekatx/go-renderer/img"
	"github.com/hekatx/go-renderer/obj"
)

//import (
//	"fmt"
//	"image/color"
//	"image/png"
//	"os"
//
//	"github.com/hekatx/go-renderer/img"
//)

func renderWireframe(model obj.ObjData, width, height int, image *image.RGBA) {
	white := color.RGBA{255, 255, 255, 0xff}
	for i := 0; i < len(model.Faces); i++ {
		face := model.Faces[i]
		for j := 0; j < len(face); j++ {
			v0 := model.Vertices[face[j]]
			v1 := model.Vertices[face[(j+1)%len(face)]]

			x0 := (v0[0] + 1.) * float64(width) / 2.
			y0 := (v0[1] + 1.) * float64(height) / 2.
			x1 := (v1[0] + 1.) * float64(width) / 2.
			y1 := (v1[1] + 1.) * float64(height) / 2.

			img.Line(x0, y0, x1, y1, image, white)

		}
	}
}

func main() {
	width := 600
	height := 600
	base := img.New(width, height)

	d := obj.Decode("./african_head.obj")

	fmt.Printf("%+v\n", len(d.Faces))

	renderWireframe(d, width, height, base)

	f, _ := os.Create("image.png")
	err := png.Encode(f, base)

	if err != nil {
		fmt.Fprintf(os.Stdout, "fucked up %v", err)
	}

}
