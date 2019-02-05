package main

import "os"
import "bufio"
import "strings"
import "github.com/golang/geo/r3"
import "strconv"

type Model struct {
	Nverts int
	Nfaces int
	Verts []Vertex
	Faces []Face
}

type Vertex struct {
	id int
	coords r3.Vector
}

type Face struct {
	id int
	components [][]int
}


func readObj(filepath string) Model {
	file, err := os.Open(filepath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	var lines []string
	var current string
	var verts []Vertex
	var faces []Face
	scanner := bufio.NewScanner(file)
	vertNo:=0
	faceNo:=0
	for scanner.Scan() {
		current = scanner.Text()
		lines = append(lines, current)

		// Handle Vertices
		if len(current) > 2 && current[0:2] == "v " {
			vertNo += 1
			tmp := strings.Fields(current)
			var err error
			// TODO handle errors here, as well as len(tmp) != 4
			vertX, err := strconv.ParseFloat(tmp[1], 32)
			vertY, err := strconv.ParseFloat(tmp[2], 32)
			vertZ, err := strconv.ParseFloat(tmp[3], 32)
			if err != nil {
				panic(err)
			}
			verts = append(verts, Vertex{vertNo, r3.Vector{vertX, vertY, vertZ}})
		}

		// Handle Faces
		if len(current) > 2 && current[0:2] == "f " {
			faceNo += 1
			tmp := strings.Fields(current)
			face := make([][]int, len(tmp)-1)
			for i:=0; i<len(face); i++ {
				tmpF := strings.Split(tmp[i+1], "/") // First entry is the "f " tag
				face[i] = make([]int, len(tmpF))
				for j, s := range tmpF {
					face[i][j], err = strconv.Atoi(s)
					if err != nil {
						panic(err)
					}
				}
			}
			faces = append(faces, Face{faceNo, face})
		}
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}

	return Model{len(verts), len(faces), verts, faces}
}


