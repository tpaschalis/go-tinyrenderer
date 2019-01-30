package main

import "fmt"
import "os"

import "image"
import "image/color"
import "image/png"

func line(x0, y0, x1, y1 int, canvas *image.RGBA, c color.RGBA) {
	steep := false
	if abs(x0-x1) < abs(y0-y1) {
		x0, y0 = y0, x0
		x1, y1 = y1, x1
		steep = true
	}

	if x0>x1 {
		x0, x1 = x1, x0
		y0, y1 = y1, y0
	}

	dx := x1-x0
	dy := y1-y0

	derror := absf(float64(dy)/float64(dx))
	errorCur := 0.0
	y := y0

	for x:=x0; x<=x1; x++ {
		if steep {
			canvas.Set(y, x, c)
		} else {
			canvas.Set(x, y, c)
		}
		errorCur += derror
		if errorCur > 0.5 {
			if y1>y0 {
				y+=1
				errorCur -= 1.
			} else {
				y+=-1
				errorCur -= 1.
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

func main() {
	w, h := 100, 100

	img := image.NewRGBA(image.Rect(0, 0, w, h))

	white := color.RGBA{255, 255, 255, 255}
	red := color.RGBA{255, 0, 0, 255}

	f, err := os.OpenFile("out.png", os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	for i := 0; i < w; i++ {
		for j := 0; j < h; j++ {
			img.Set(i, j, color.RGBA{0, 0, 0, 255})
		}
	}

	img.Set(52, 41, red)
	line(13, 20, 80, 40, img, white)
	line(20, 13, 40, 80, img, red)
	line(80, 40, 13, 20, img, red)
	img = flipVertically(img)

	png.Encode(f, img)

	fmt.Println("Success!")
}
