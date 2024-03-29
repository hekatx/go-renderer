package obj

import (
	"bufio"
	"errors"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/deeean/go-vector/vector3"
)

type Model struct {
	Vertices    []vector3.Vector3
	Faces       [][]int
	Texture     []vector3.Vector3
	FaceTexture [][]int
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func Decode(path string) Model {
	var geometryTypes = struct {
		Vertex  string
		Face    string
		Texture string
	}{
		Vertex:  "v",
		Face:    "f",
		Texture: "vt",
	}

	// Max capacity of the buffer. Pump up if .obj file is bigger
	const max = 512 * 1024

	model := Model{}

	file, err := os.Open(path)
	check(err)
	defer file.Close()

	scanner := bufio.NewScanner(file)

	buffer := make([]byte, max)
	scanner.Buffer(buffer, max)

	for scanner.Scan() {
		// Check for empty lines
		if len(scanner.Text()) != 0 {
			// Split geometric data by its type and values
			line := strings.Fields(scanner.Text())
			dataType := line[0]
			if dataType == geometryTypes.Vertex {
				vertices, e := stringToFloatVector3(line[1:])

				check(e)

				model.Vertices = append(model.Vertices, vertices)
			}

			if dataType == geometryTypes.Face {
				vi, vt, e := parseFaces(line[1:])

				if e != nil {
					check(e)
				}

				model.Faces = append(model.Faces, vi)
				model.FaceTexture = append(model.FaceTexture, vt)
			}

			if dataType == geometryTypes.Texture {
				v, e := stringToFloatVector3(line[1:])

				check(e)

				model.Texture = append(model.Texture, v)
			}
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return model
}

func stringToFloatVector3(arr []string) (vector3.Vector3, error) {
	fa := make([]float64, 0, len(arr))
	for _, a := range arr {
		f, e := strconv.ParseFloat(a, 64)
		if e != nil {
			return *vector3.New(0, 0, 0), e
		}
		fa = append(fa, f)
	}

	return *vector3.New(fa[0], fa[1], fa[2]), nil
}

func parseFaces(vNormalIndices []string) ([]int, []int, error) {
	vIndices := make([]int, 0, len(vNormalIndices))
	var vtextures []int

	for _, vn := range vNormalIndices {
		indices := strings.Split(string(vn), "/")
		vIndex := indices[0]
		if string(vIndex[0]) == "/" {
			return nil, nil, errors.New("vNormalIndices have wrong formatting")
		}
		vi, _ := strconv.Atoi(vIndex)
		if len(indices) > 1 {
			vt, _ := strconv.Atoi(indices[1])
			vtextures = append(vtextures, vt-1)
		}
		// Vertex indices are not 1-based so we must subtract 1 for every index
		vIndices = append(vIndices, vi-1)
	}

	return vIndices, vtextures, nil
}
