package render

import (
	"image"
	"image/color"
	"math"

	"github.com/deeean/go-vector/vector3"
	"github.com/hekatx/go-renderer/draw"
	"github.com/hekatx/go-renderer/obj"
)

func Model(model obj.Model, width, height int, image *image.RGBA, light_dir vector3.Vector3, texture *image.Image) {
	screen_coords := make([]vector3.Vector3, 3)
	var world_coords [3]vector3.Vector3
	var texture_vertices = make([]vector3.Vector3, 3)

	var zbuffer = make([]float64, width*height)
	for i := 0; i < len(zbuffer); i++ {
		zbuffer[i] = math.Inf(-1)
	}

	for i := 0; i < len(model.Faces); i++ {
		f := model.Faces[i]
		ft := model.FaceTexture[i]

		for j := 0; j < 3; j++ {
			v := model.Vertices[f[j]]
			screen_coords[j] = vector3.Vector3{X: float64(int((v.X+1.)*float64(width)/2. + .5)), Y: float64(int((v.Y+1.)*float64(height)/2. + .5)), Z: v.Z}
			world_coords[j] = v

			if len(ft) > 0 {
				tv := ft[j]
				texture_vertices[j] = model.Texture[tv]
			}
		}

		n := getNormal(world_coords)
		intensity := n.Dot(&light_dir)
		c := color.RGBA{uint8(intensity * 255), uint8(intensity * 255), uint8(intensity * 255), 255}

		if intensity > 0 {
			drawTriangle(image, c, screen_coords, &zbuffer, width, height, &texture_vertices, texture)
		}
	}
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
}

func getNormal(pts [3]vector3.Vector3) vector3.Vector3 {
	qr := pts[2].Sub(&pts[0])
	qs := pts[1].Sub(&pts[0])
	normal := qr.Cross(qs)
	return *normal.Normalize()
}

func FlipVertically(canvas *image.RGBA) *image.RGBA {
	bounds := canvas.Bounds()
	flipped := image.NewRGBA(image.Rect(0, 0, bounds.Max.X, bounds.Max.Y))
	for i := 0; i <= bounds.Max.X; i++ {
		for j := 0; j <= bounds.Max.Y; j++ {
			flipped.Set(i, bounds.Max.Y-j-1, canvas.At(i, j))
		}
	}
	return flipped
}

func getPixelValue(img *image.Image, x int, y int) (uint8, uint8, uint8) {
	r, g, b, _ := (*img).At(x, y).RGBA()
	return uint8(r / 257), uint8(g / 257), uint8(b / 257)
}

func getColorFromTexture(img *image.Image, vertexTextures *[]vector3.Vector3, barycentric *vector3.Vector3) color.RGBA {
	bo := (*img).Bounds()
	w := bo.Max.X
	h := bo.Max.Y

	x := (*vertexTextures)[0].X*barycentric.X + (*vertexTextures)[1].X*barycentric.Y + (*vertexTextures)[2].X*barycentric.Z
	y := (*vertexTextures)[0].Y*barycentric.X + (*vertexTextures)[1].Y*barycentric.Y + (*vertexTextures)[2].Y*barycentric.Z

	r, g, b := getPixelValue(img, int(x*float64(w)), int(y*float64(h)))
	return color.RGBA{r, g, b, 255}
}

func getBarycentricCoords(pts []vector3.Vector3, P *vector3.Vector3) vector3.Vector3 {
	v0 := vector3.Vector3{X: float64(pts[1].X - pts[0].X), Y: float64(pts[2].X - pts[0].X), Z: float64(pts[0].X - P.X)}
	v1 := vector3.Vector3{X: float64(pts[1].Y - pts[0].Y), Y: float64(pts[2].Y - pts[0].Y), Z: float64(pts[0].Y - P.Y)}
	u := v0.Cross(&v1)
	if math.Abs(u.Z) > 1e-2 {
		return vector3.Vector3{X: 1.0 - (u.X+u.Y)/u.Z, Y: u.X / u.Z, Z: u.Y / u.Z}
	}
	return vector3.Vector3{X: -1., Y: 1., Z: 1.}
}

func isPointOutsideTriangle(W1, W2, w3 float64) bool {
	return W1 < 0 || W2 < 0 || w3 < 0
}

func drawTriangle(i *image.RGBA, c color.RGBA, v []vector3.Vector3, zb *[]float64, w, h int, tvs *[]vector3.Vector3, t *image.Image) {
	p := vector3.Vector3{}
	for y := 0; y < i.Bounds().Max.Y; y++ {
		for x := 0; x < i.Bounds().Max.X; x++ {
			p.X = float64(x)
			p.Y = float64(y)
			bc := getBarycentricCoords(v, &p)

			if isPointOutsideTriangle(bc.X, bc.Y, bc.Z) {
				continue
			}

			p.Z = (bc.X * v[0].Z) + (bc.Y * v[1].Z) + (bc.Z * v[2].Z)

			if (*zb)[int(p.X+p.Y*float64(w))] < p.Z {
				(*zb)[int(p.X+p.Y*float64(w))] = p.Z
				cft := getColorFromTexture(t, tvs, &bc)
				i.Set(int(p.X), int(p.Y), cft)
			}
		}
	}
}
