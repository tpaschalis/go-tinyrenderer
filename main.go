package main

import "fmt"
import "os"

import "image"
import "image/color"
import "image/png"

func line(x0, y0, x1, y1 int, canvas *image.RGBA, c color.RGBA) {
	a := 0.01 // run it for a = 0.1 as well, to see the line thinning out
	for t := 0.0; t < 1.0; t += a {
		x := float64(x0)*(1.0-t) + float64(x1)*t
		y := float64(y0)*(1.0-t) + float64(y1)*t
		canvas.Set(int(x), int(y), c)
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
	img = flipVertically(img)

	png.Encode(f, img)

	fmt.Println("Success!")
}
