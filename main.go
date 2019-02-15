package main

import "fmt"
import "os"

import "image"
import "image/color"
import "image/png"

import "github.com/golang/geo/r3"

import "math/rand"

func flipVertically(canvas *image.RGBA) *image.RGBA {
	bounds := canvas.Bounds()
	flipped := image.NewRGBA(image.Rect(0, 0, bounds.Max.X, bounds.Max.Y))
	for i := 0; i <= bounds.Max.X; i++ {
		for j := 0; j <= bounds.Max.Y; j++ {
			flipped.Set(i, bounds.Max.Y-j-1, canvas.At(i, j))
		}
	}
	return flipped
}

func abs(x int) int {
	if x < 0 {
		return -1 * x
	}
	return x
}

func absf(x float64) float64 {
	if x < 0 {
		return -1 * x
	}
	return x
}

func swap(a, b int) (int, int) {
	return b, a
}

func swapf(a, b float64) (float64, float64) {
	return b, a
}

func triangle(t0, t1, t2 r3.Vector, canvas *image.RGBA, c color.RGBA) {
	if t0.Y == t1.Y && t0.Y == t2.Y {
		return
	}

	// bubble-sort vectors, according to their y-coordinate
	if t0.Y > t1.Y {
		t0, t1 = t1, t0
	}

	if t0.Y > t2.Y {
		t0, t2 = t2, t0
	}

	if t1.Y > t2.Y {
		t1, t2 = t2, t1
	}
	total_height := t2.Y - t0.Y

	var seg_height float64
	var alpha, beta float64
	var A, B r3.Vector

	for i := 0.; i <= total_height; i++ {
		second_half := i > t1.Y-t0.Y || t1.Y == t0.Y // a boolean value
		if second_half {
			seg_height = t2.Y - t1.Y
		} else {
			seg_height = t1.Y - t0.Y
		}
		alpha = i / total_height
		if second_half {
			beta = (i - (t1.Y - t0.Y)) / seg_height
		} else {
			beta = i / seg_height
		}
		A = r3.Vector.Add(t0, r3.Vector.Mul(r3.Vector.Sub(t2, t0), alpha))
		if second_half {
			B = r3.Vector.Add(t1, r3.Vector.Mul(r3.Vector.Sub(t2, t1), beta))
		} else {
			B = r3.Vector.Add(t0, r3.Vector.Mul(r3.Vector.Sub(t1, t0), beta))
		}

		if A.X > B.X {
			A, B = B, A
		}
		for j := int(A.X); j <= int(B.X); j++ {
			canvas.Set(j, int(t0.Y+i), c)
		}
	}
}

func main() {
	w, h := 800, 800
	fw, fh := 800., 800.
	//w, h := 150, 150
	//fw, fh := 150., 150.
	_, _, _, _ = w, h, fw, fh

	img := image.NewRGBA(image.Rect(0, 0, w, h))

	white := color.RGBA{255, 255, 255, 255}
	_ = white
	red := color.RGBA{255, 0, 0, 255}
	_ = red
	green := color.RGBA{0, 255, 0, 255}
	_ = green

	f, err := os.OpenFile("out.png", os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	// Set black background, otherwise it's transparert and appears like a checkerboard pattern.
	for i := 0; i < w; i++ {
		for j := 0; j < h; j++ {
			img.Set(i, j, color.RGBA{0, 0, 100, 255})
		}
	}

	Model := readObj("obj/human_head.obj")
	rand.Seed(1)

	lightDir := r3.Vector{0., 0., -1.}

	for i := 0; i < Model.Nfaces; i++ {
		face := Model.Faces[i]
		var screen_coords, world_coords []r3.Vector

		var tmp []int
		for i := range face.components {
			tmp = append(tmp, face.components[i][0])
		}

		// Took me a while to find out, but screen coordinates must be integers (well, duuh?!),
		// Otherwise weird artifacts of points not rendering between vertices are encountered -.-
		for j := 0; j < 3; j++ {
			v := Model.Verts[tmp[j]-1].coords
			screen_coords = append(screen_coords, r3.Vector{
				float64(int((v.X + 1.) * fw / 2.)),
				float64(int((v.Y + 1.) * fh / 2.)),
				0.})
			world_coords = append(world_coords, v)
		}
		a := world_coords[0]
		b := world_coords[1]
		c := world_coords[2]
		n := r3.Vector.Cross(r3.Vector.Sub(c, b), r3.Vector.Sub(b, a))
		n = r3.Vector.Normalize(n)
		intensity := r3.Vector.Dot(n, lightDir)
		if intensity >= 0 {
			triangle(screen_coords[0], screen_coords[1], screen_coords[2], img, color.RGBA{uint8(intensity * 255), uint8(intensity * 255), uint8(intensity * 255), 255})
		}
	}

	img = flipVertically(img)
	png.Encode(f, img)
	fmt.Println("Success!")
}
