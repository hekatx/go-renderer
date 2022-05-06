package main

import (
	"fmt"
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

func main() {
	width := 600
	height := 600
	base := img.New(width, height)

	white := color.RGBA{255, 255, 255, 0xff}

	d := obj.Decode("./african_head.obj")

	fmt.Printf("%+v\n", len(d.Faces))

	for i := 0; i < len(d.Faces); i++ {
		face := d.Faces[i]
		for j := 0; j < len(face); j++ {
			v0 := d.Vertices[face[j]]
			v1 := d.Vertices[face[(j+1)%len(face)]]
			x0 := float64((v0[0] + 1.) * float32(width) / 2.)
			y0 := float64((v0[1] + 1.) * float32(height) / 2.)
			x1 := float64((v1[0] + 1.) * float32(width) / 2.)
			y1 := float64((v1[1] + 1.) * float32(height) / 2.)
			img.Line(x0, y0, x1, y1, base, white)
		}
	}

	f, _ := os.Create("image.png")
	err := png.Encode(f, base)

	if err != nil {
		fmt.Fprintf(os.Stdout, "fucked up %v", err)
	}

}
