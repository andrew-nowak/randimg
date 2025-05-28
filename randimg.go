package main

import (
	"bufio"
	"flag"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"log"
	"math/rand/v2"
	"os"
	"strings"
	"time"

	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/gomono"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/math/fixed"
)

var fileType string
var output string

func init() {

	flag.StringVar(&fileType, "type", "jpg", "file type to output")
	flag.StringVar(&fileType, "t", "jpg", "shorthand for type")

	flag.StringVar(&output, "output", "out.jpg", "path to output file to")
	flag.StringVar(&output, "o", "out.jpg", "shorthand for output")
}

func main() {
	flag.Parse()

	if strings.ToLower(fileType) == "png" {
		output = strings.Replace(output, ".jpg", ".png", -1)
		if !strings.HasSuffix(output, ".png") {
			output = output + ".png"
		}
	} else {
		if !(strings.HasSuffix(output, ".jpg") || strings.HasSuffix(output, ".jpeg")) {
			output = output + ".jpg"
		}
	}

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

	for x := range width {
		for y := range height {
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
		Dot:  fixed.Point26_6{X: fixed.Int26_6(1000), Y: fixed.Int26_6(29000)},
	}
	t := time.Now()
	d.DrawString("Generated: " + t.Format(time.DateTime))
	d.Dot = fixed.Point26_6{X: fixed.Int26_6(1000), Y: fixed.Int26_6(31000)}
	if len(flag.Args()) >= 1 {
		d.DrawString(strings.Join(flag.Args(), " "))
	} else {
		d.DrawString("[no title]")
	}

	newFile, err := os.Create(output)
	if err != nil {
		log.Fatalf("Failed to create file %v", err)
	}
	defer newFile.Close()

	if strings.ToLower(fileType) == "png" {
		if err := png.Encode(bufio.NewWriter(newFile), img); err != nil {
			log.Fatalf("Failed to encode image %v", err)
		}
	} else if err := jpeg.Encode(bufio.NewWriter(newFile), img, &jpeg.Options{Quality: 100}); err != nil {
		log.Fatalf("Failed to encode image %v", err)
	}
	log.Println("Done")
}
