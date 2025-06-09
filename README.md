# randimg

A simple command-line tool that generates random colourful images with customizable text overlays.

## Features

- Generates images with 5 random coloured vertical stripes
- Adds timestamp and custom text overlay
- Supports both JPEG and PNG output formats
- Configurable output filename

## Usage

```bash
# Generate a basic image
randimg

# Generate with custom title
randimg "My Custom Title"

# Generate PNG instead of JPEG
randimg -type png "My Image"

# Custom output filename
randimg -output my-image.jpg "Custom Image"
```

## Command Line Options

- `-type` or `-t`: Output format (`jpg` or `png`, default: `jpg`)
- `-output` or `-o`: Output filename (default: `out.jpg`)

## Dependencies

- Go 1.24.1+
- `golang.org/x/image` for font rendering
