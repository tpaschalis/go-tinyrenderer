package main

import "fmt"
import "os"

import "image"
import "image/color"
import "image/png"

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

	img.Set(52, 41, red)
	_ = white

	png.Encode(f, img)

	fmt.Println("Success!")
}
