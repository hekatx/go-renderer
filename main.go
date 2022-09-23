package main

import (
	"fmt"
	"image/png"
	"os"

	"github.com/deeean/go-vector/vector3"
	"github.com/hekatx/go-renderer/draw"
	"github.com/hekatx/go-renderer/obj"
	"github.com/hekatx/go-renderer/render"
)

func main() {
	width := 600
	height := 600
	base := draw.NewImage(width, height)

	head_model := obj.Decode("./african_head.obj")

	light_dir := vector3.New(0., 0., -1.)
	render.Model(head_model, width, height, base, *light_dir)

	f, _ := os.Create("image.png")
	err := png.Encode(f, base)

	if err != nil {
		fmt.Fprintf(os.Stdout, "fucked up %v", err)
	}
}
