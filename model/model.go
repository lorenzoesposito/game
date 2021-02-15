package model

import (
	"fmt"
	"io"
	"os"

	. "game.com/lorenzo/game/utils"
)

type Model struct {
	Vecs       []Vec3f
	VecIndices []float32
	Materials  []string
}

func NewModel(file string) Model {
	// Open the file for reading and check for errors.
	objFile, err := os.Open(file)
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
			// Create a vec to assign digits to.
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
					model.Materials = append(model.Materials, material)
				}
			}
		}
	}

	return model
}
