package main

import "fmt"
import "os"

import "image"
import "image/color"
import "image/png"

import "github.com/golang/geo/r3"

func line(p0, p1 r3.Vector, canvas *image.RGBA, c color.RGBA) {
	steep := false
	if absf(p0.X-p1.X) < absf(p0.Y-p1.Y) {
		p0.X, p0.Y = p0.Y, p0.X
		p1.X, p1.Y = p1.Y, p1.X
		steep = true
	}

	if p0.X > p1.X {
		p0, p1 = p1, p0
	}

	dx := int(p1.X - p0.X)
	dy := int(p1.Y - p0.Y)
	derror2 := 2 * abs(dy)
	error2 := 0
	y := int(p0.Y)

	if steep {
		for x := int(p0.X); x <= int(p1.X); x++ {
			canvas.Set(y, x, c)
			error2 += derror2
			if error2 > dx {
				if p1.Y > p0.Y {
					y += 1
					error2 -= 2 * dx
				} else {
					y += -1
					error2 -= 2 * dx
				}
			}
		}
	} else {
		for x := int(p0.X); x <= int(p1.X); x++ {
			canvas.Set(x, y, c)
			error2 += derror2
			if error2 > dx {
				if p1.Y > p0.Y {
					y += 1
					error2 -= 2 * dx
				} else {
					y += -1
					error2 -= 2 * dx
				}
			}
		}
	}

}

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

func triangle(t0, t1, t2 r3.Vector, canvas *image.RGBA, c color.RGBA) {
	line(t0, t1, canvas, c)
	line(t1, t2, canvas, c)
	line(t2, t0, canvas, c)
}

func main() {
	//w, h := 800, 800
	//fw, fh := 800., 800.
	w, h := 200, 200
	fw, fh := 200., 200.
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
			img.Set(i, j, color.RGBA{0, 0, 0, 255})
		}
	}

	//Model := readObj("obj/human_head.obj")

	t0 := []r3.Vector{
		{10, 70, 0},
		{50, 160, 0},
		{70, 80, 0},
	}
	t1 := []r3.Vector{
		{180, 50, 0},
		{150, 1, 0},
		{70, 180, 0},
	}
	t2 := []r3.Vector{
		{180, 150, 0},
		{120, 160, 0},
		{130, 180, 0},
	}

	triangle(t0[0], t0[1], t0[2], img, red)
	triangle(t1[0], t1[1], t1[2], img, white)
	triangle(t2[0], t2[1], t2[2], img, green)

	img = flipVertically(img)
	png.Encode(f, img)
	fmt.Println("Success!")
}
