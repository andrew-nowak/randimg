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

	// Determine output filename
	outputFile := determineOutputFilename(fileType, output)

	// Generate the image
	img, err := generateImage(flag.Args())
	if err != nil {
		log.Fatalf("Failed to generate image: %v", err)
	}

	// Save the image
	err = saveImage(img, outputFile, fileType)
	if err != nil {
		log.Fatalf("Failed to save image: %v", err)
	}

	log.Println("Done")
}

// determineOutputFilename adjusts the output filename based on file type
func determineOutputFilename(fileType, output string) string {
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
	return output
}

// generateImage creates a random striped image with text overlay
func generateImage(args []string) (*image.RGBA, error) {
	width := 1000
	height := 500
	img := image.NewRGBA(image.Rect(0, 0, width, height))

	// Generate random colors for stripes
	cols := make([]color.RGBA, 5)
	white := color.RGBA{R: 0xff, G: 0xff, B: 0xff, A: 0xff}

	for i := range cols {
		col := color.RGBA{
			R: rand.N[uint8](0xff),
			G: rand.N[uint8](0xff),
			B: rand.N[uint8](0xff),
			A: 0xff,
		}
		cols[i] = col
	}

	// Draw colored stripes and white footer
	for x := range width {
		for y := range height {
			if y < 420 {
				img.Set(x, y, cols[x/200])
			} else {
				img.Set(x, y, white)
			}
		}
	}

	// Add text overlay
	err := addTextOverlay(img, args)
	if err != nil {
		return nil, err
	}

	return img, nil
}

// addTextOverlay adds timestamp and title text to the image
func addTextOverlay(img *image.RGBA, args []string) error {
	f, err := opentype.Parse(gomono.TTF)
	if err != nil {
		return err
	}

	face, err := opentype.NewFace(f, &opentype.FaceOptions{
		Size:    32,
		DPI:     72,
		Hinting: font.HintingNone,
	})
	if err != nil {
		return err
	}

	d := &font.Drawer{
		Dst:  img,
		Src:  image.NewUniform(color.RGBA{0, 0, 0, 0xff}),
		Face: face,
		Dot:  fixed.Point26_6{X: fixed.Int26_6(1000), Y: fixed.Int26_6(29000)},
	}

	t := time.Now()
	d.DrawString("Generated: " + t.Format(time.DateTime))
	d.Dot = fixed.Point26_6{X: fixed.Int26_6(1000), Y: fixed.Int26_6(31000)}

	if len(args) >= 1 {
		d.DrawString(strings.Join(args, " "))
	} else {
		d.DrawString("[no title]")
	}

	return nil
}

// saveImage writes the image to a file in the specified format
func saveImage(img image.Image, filename, format string) error {
	newFile, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer newFile.Close()

	writer := bufio.NewWriter(newFile)
	defer writer.Flush()

	if strings.ToLower(format) == "png" {
		return png.Encode(writer, img)
	}
	return jpeg.Encode(writer, img, &jpeg.Options{Quality: 100})
}
