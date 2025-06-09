package main

import (
	"image/color"
	"image/jpeg"
	"image/png"
	"os"
	"testing"
)

func TestJPEGGeneration(t *testing.T) {
	testFile := "test_output.jpg"
	defer os.Remove(testFile)

	img, err := generateImage([]string{"Test JPEG"})
	if err != nil {
		t.Fatalf("Failed to generate image: %v", err)
	}

	err = saveImage(img, testFile, "jpg")
	if err != nil {
		t.Fatalf("Failed to save JPEG: %v", err)
	}

	// Validate the saved file
	file, err := os.Open(testFile)
	if err != nil {
		t.Fatalf("Failed to open test file: %v", err)
	}
	defer file.Close()

	_, err = jpeg.Decode(file)
	if err != nil {
		t.Fatalf("Failed to decode JPEG: %v", err)
	}
}

func TestPNGGeneration(t *testing.T) {
	testFile := "test_output.png"
	defer os.Remove(testFile)

	img, err := generateImage([]string{"Test PNG"})
	if err != nil {
		t.Fatalf("Failed to generate image: %v", err)
	}

	err = saveImage(img, testFile, "png")
	if err != nil {
		t.Fatalf("Failed to save PNG: %v", err)
	}

	file, err := os.Open(testFile)
	if err != nil {
		t.Fatalf("Failed to open test file: %v", err)
	}
	defer file.Close()

	_, err = png.Decode(file)
	if err != nil {
		t.Fatalf("Failed to decode PNG: %v", err)
	}
}

func TestImageDimensions(t *testing.T) {
	img, err := generateImage([]string{"Test"})
	if err != nil {
		t.Fatalf("Failed to generate image: %v", err)
	}

	bounds := img.Bounds()
	if bounds.Dx() != 1000 {
		t.Errorf("Expected width 1000, got %d", bounds.Dx())
	}
	if bounds.Dy() != 500 {
		t.Errorf("Expected height 500, got %d", bounds.Dy())
	}
}

func TestColorStripes(t *testing.T) {
	img, err := generateImage([]string{"Test"})
	if err != nil {
		t.Fatalf("Failed to generate image: %v", err)
	}

	// Check that different stripes have different colors
	colors := make(map[color.RGBA]bool)

	for i := 0; i < 5; i++ {
		x := i*200 + 100
		y := 200
		c := img.At(x, y)
		if rgba, ok := c.(color.RGBA); ok {
			colors[rgba] = true
		}
	}

	if len(colors) < 2 {
		t.Error("Expected color variation in stripes")
	}
}

func TestFooterIsWhite(t *testing.T) {
	img, err := generateImage([]string{"     "})
	if err != nil {
		t.Fatalf("Failed to generate image: %v", err)
	}

	white := color.RGBA{R: 0xff, G: 0xff, B: 0xff, A: 0xff}

	for x := 0; x < 1000; x += 100 {
		y := 480 // Footer row
		c := img.At(x, y)
		if rgba, ok := c.(color.RGBA); ok {
			if rgba != white {
				t.Errorf("Expected white in footer at (%d,%d), got %+v", x, y, rgba)
				return
			}
		}
	}
}

func TestOutputFilename(t *testing.T) {
	tests := []struct {
		fileType, input, expected string
	}{
		{"jpg", "test", "test.jpg"},
		{"png", "test.jpg", "test.png"},
		{"jpg", "test.jpeg", "test.jpeg"},
		{"png", "image", "image.png"},
	}

	for _, test := range tests {
		result := determineOutputFilename(test.fileType, test.input)
		if result != test.expected {
			t.Errorf("determineOutputFilename(%q, %q) = %q, want %q",
				test.fileType, test.input, result, test.expected)
		}
	}
}
