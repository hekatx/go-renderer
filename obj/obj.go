package obj

import (
	"bufio"
	"log"
	"os"
	"strconv"
	"strings"
)

type ObjData struct {
	Vertices [][]float64
	Faces    [][]int
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func Decode(path string) ObjData {
	// Max capacity of the buffer. Pump up if .obj file is bigger
	const max = 512 * 1024

	d := ObjData{}

	f, e := os.Open(path)
	check(e)
	defer f.Close()

	s := bufio.NewScanner(f)

	b := make([]byte, max)
	s.Buffer(b, max)

	for s.Scan() {
		if len(s.Text()) != 0 {
			line := strings.Fields(s.Text())
			if line[0] == "v" {
				vstring := line[1:]
				vcoords, e := stringToFloat(vstring)

				check(e)

				d.Vertices = append(d.Vertices, vcoords)
			}

			if line[0] == "f" {
				fstring := line[1:]

				parsedfstring := make([]string, 0, len(fstring))

				for _, fs := range fstring {
					parsedfstring = append(parsedfstring, strings.Split(string(fs), "/")[0])
				}
				f, e := stringToInt(parsedfstring)

				check(e)

				d.Faces = append(d.Faces, f)
			}
		}
	}

	if err := s.Err(); err != nil {
		log.Fatal(err)
	}

	return d
}

func stringToFloat(arr []string) ([]float64, error) {
	fa := make([]float64, 0, len(arr))
	for _, a := range arr {
		f, e := strconv.ParseFloat(a, 32)
		if e != nil {
			return fa, e
		}
		fa = append(fa, float64(f))
	}

	return fa, nil
}

func stringToInt(arr []string) ([]int, error) {
	fa := make([]int, 0, len(arr))
	for _, a := range arr {
		f, e := strconv.Atoi(a)
		if e != nil {
			return fa, e
		}
		fa = append(fa, int(f)-1)
	}

	return fa, nil
}
