package main

import (
	"fmt"
	"image/color"
	"image/png"
	"os"

	"github.com/hekatx/go-renderer/img"
)

func main() {
	base := img.New(200, 200)

	white := color.RGBA{255, 255, 255, 0xff}
	red := color.RGBA{255, 0, 0, 0xff}

	img.Line(13, 20, 80, 40, base, white)
	img.Line(20, 13, 40, 80, base, red)
	img.Line(80, 40, 13, 20, base, red)

	f, _ := os.Create("image.png")
	err := png.Encode(f, base)

	if err != nil {
		fmt.Fprintf(os.Stdout, "fucked up %v", err)
	}
}
