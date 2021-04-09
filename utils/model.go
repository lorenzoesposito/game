package utils

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
)

type Model struct {
	Vecs                                 []Vec3f
	VecIndices, NormalIndices, UvIndices []float32
	Materials                            []Vec3f
}

func Mesh(file string) []Tri3f {
	var m = NewModel(file)
	var mesh []Tri3f
	for i := 0; i < len(m.VecIndices); i += 3 {
		var face Tri3f
		face = Tri3f{
			m.Vecs[int(m.VecIndices[i])-1],
			m.Vecs[int(m.VecIndices[i+1])-1],
			m.Vecs[int(m.VecIndices[i+2])-1]}
		mesh = append(mesh, face)
	}
	return mesh
}

func Color(file string) []Vec3f {
	var m = NewModel(file)
	return m.Materials
}

func readMTL(file, material string) Vec3f {
	fila, err := ioutil.ReadFile(fmt.Sprintf("%s.mtl", file))
	//mtlFile, err := os.Open(fmt.Sprintf("%s.mtl", file))
	if err != nil {
		panic(err)
	}
	lines := strings.Split(string(fila), "\n")[3:]
	txt := strings.Split(strings.Join(lines, ","), "newmtl")
	for _, i := range txt {
		mtl := strings.Split(i, ",")
		if len(mtl) > 3 && strings.Trim(mtl[0], " ") == material {
			kd := strings.Split(mtl[3], " ")[1:]
			return Vec3f{StringToFloat(kd[0]), StringToFloat(kd[1]), StringToFloat(kd[2])}
		}
	}
	return Vec3f{1, 1, 1}
}

func NewModel(file string) Model {
	// Open the file for reading and check for errors.
	objFile, err := os.Open(fmt.Sprintf("%s.obj", file))
	if err != nil {
		panic(err)
	}

	defer objFile.Close()

	model := Model{}

	for {
		var lineType string

		// Scan the type field.
		_, err := fmt.Fscanf(objFile, "%s", &lineType)
		if err != nil {
			if err == io.EOF {
				break
			}
		}

		// Check type.
		switch lineType {
		// VERTICES.
		case "v":
			vec := Vec3f{}

			fmt.Fscanf(objFile, "%f %f %f\n", &vec.X, &vec.Y, &vec.Z)

			model.Vecs = append(model.Vecs, vec)

		case "f":
			/*
				norm := make([]float32, 3)
				vec := make([]float32, 3)
				uv := make([]float32, 3)

				// Get the digits from the file.
				matches, _ := fmt.Fscanf(objFile, "%f/%f/%f %f/%f/%f %f/%f/%f\n", &vec[0], &uv[0], &norm[0], &vec[1], &uv[1], &norm[1], &vec[2], &uv[2], &norm[2])

				if matches != 9 {
					panic("Cannot read your file")
				}

				model.VecIndices = append(model.VecIndices, vec[0])
				model.VecIndices = append(model.VecIndices, vec[1])
				model.VecIndices = append(model.VecIndices, vec[2])*/

		}
		if string(lineType[0]) == "M" {
			for {
				var material = lineType
				var line string
				// Scan the type field.
				_, err := fmt.Fscanf(objFile, "%s", &line)

				if err != nil || line == "usemtl" {
					break
				}
				if line == "f" {
					norm := make([]float32, 3)
					vec := make([]float32, 3)
					uv := make([]float32, 3)

					// Get the digits from the file.
					matches, _ := fmt.Fscanf(objFile, "%f/%f/%f %f/%f/%f %f/%f/%f\n", &vec[0], &uv[0], &norm[0], &vec[1], &uv[1], &norm[1], &vec[2], &uv[2], &norm[2])

					if matches != 9 {
						panic("Cannot read your file")
					}

					model.VecIndices = append(model.VecIndices, vec[0])
					model.VecIndices = append(model.VecIndices, vec[1])
					model.VecIndices = append(model.VecIndices, vec[2])
					model.Materials = append(model.Materials, readMTL(file, material))
				}
			}
		}
	}

	return model
}
