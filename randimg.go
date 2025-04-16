package main

import (
	"bufio"
	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/gomono"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/math/fixed"
	"image"
	"image/color"
	"image/jpeg"
	"log"
	"math/rand/v2"
	"os"
	"strings"
	"time"
)

func main() {
	width := 1000
	height := 500
	img := image.NewRGBA(image.Rect(0, 0, width, height))

	cols := make([]color.RGBA, 5)
	white := color.RGBA{
		R: 0xff,
		G: 0xff,
		B: 0xff,
		A: 0xff,
	}

	for i := range cols {
		col := color.RGBA{
			R: rand.N[uint8](0xff),
			G: rand.N[uint8](0xff),
			B: rand.N[uint8](0xff),
			A: 0xff,
		}
		cols[i] = col
	}

	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			if y < 420 {
				img.Set(x, y, cols[x/200])
			} else {
				img.Set(x, y, white)
			}
		}
	}

	f, err := opentype.Parse(gomono.TTF)
	if err != nil {
		log.Fatalf("Failed to open packaged TTF: %v", err)
	}
	face, err := opentype.NewFace(f, &opentype.FaceOptions{
		Size:    32,
		DPI:     72,
		Hinting: font.HintingNone,
	})

	d := &font.Drawer{
		Dst:  img,
		Src:  image.NewUniform(color.RGBA{0, 0, 0, 0xff}),
		Face: face,
		Dot:  fixed.Point26_6{fixed.Int26_6(1000), fixed.Int26_6(29000)},
	}
	t := time.Now()
	d.DrawString("Generated: " + t.Format(time.DateTime))
	d.Dot = fixed.Point26_6{X: fixed.Int26_6(1000), Y: fixed.Int26_6(31000)}
	if len(os.Args) > 1 {
		d.DrawString(strings.Join(os.Args[1:], " "))
	} else {
		d.DrawString("[no title]")
	}

	newFile, err := os.Create("out.jpg")
	if err != nil {
		log.Fatalf("Failed to create file %v", err)
	}
	defer newFile.Close()

	if err := jpeg.Encode(bufio.NewWriter(newFile), img, &jpeg.Options{Quality: 100}); err != nil {
		log.Fatalf("Failed to encode image %v", err)
	}
	log.Println("Done")
}
