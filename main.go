package main

import (
	"errors"
	"flag"
	"fmt"
	"image"
	"image/color"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"os"
	"strings"
)

const (
	defaultPalette = " .:-=+*#%@"
	defaultX       = 4
	defaultY       = 8
	resetColor     = "\x1b[0m"
)

var (
	optPalette = flag.String("p", defaultPalette, "Characters used as palette for the ASCII art, ordered from light to dark.")
	optX       = flag.Int("x", defaultX, "Step size in the x direction (width).")
	optY       = flag.Int("y", defaultY, "Step size in the y direction (height).")
)

func truecolor(c color.Color) string {
	if os.Getenv("NO_COLOR") == "1" {
		return ""
	}
	r, g, b, _ := c.RGBA()
	return fmt.Sprintf("\x1b[38;2;%d;%d;%dm", r>>8, g>>8, b>>8)
}

func main() {
	if len(os.Args[1:]) < 1 {
		printUsage()
		return
	}

	path := os.Args[1]
	if strings.TrimSpace(path) == "" {
		fatalf("Image path cannot be empty\n")
	}

	// Remove the first argument (image path) so that flag package can parse the rest
	os.Args = append([]string{os.Args[0]}, os.Args[2:]...)
	flag.Parse()

	if strings.TrimSpace(*optPalette) == "" {
		printUsage("Error: palette cannot be empty")
		return
	}
	palette := []rune(*optPalette)

	xStep := *optX
	if xStep <= 0 {
		printUsage("Error: x step must be greater than 0")
		return
	}

	yStep := *optY
	if yStep <= 0 {
		printUsage("Error: y step must be greater than 0")
		return
	}

	f, err := os.Open(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			fatalf("Error: file does not exist: %s\n", path)
		}
		fatalf("Error: failed to open file: %v\n", err)
	}
	defer f.Close()

	img, _, err := image.Decode(f)
	if err != nil {
		fatalf("Error: failed to decode image: %v\n", err)
	}

	paletteSize := len(palette)
	rect := img.Bounds()

	if xStep > rect.Dx() {
		printUsage("Error: x step is larger than image width")
		return
	}
	if yStep > rect.Dy() {
		printUsage("Error: y step is larger than image height")
		return
	}

	var builder strings.Builder
	for y := rect.Min.Y; y < rect.Max.Y; y += yStep {
		for x := rect.Min.X; x < rect.Max.X; x += xStep {
			r, g, b, _ := img.At(x, y).RGBA()

			gray := (r*30 + g*59 + b*11) / 100
			gray >>= 8

			idx := int(gray) * paletteSize / 255

			color := truecolor(img.At(x, y))

			// Better than fmt.Sprintf
			builder.WriteString(color)
			builder.WriteRune(palette[idx])
			builder.WriteString(resetColor)
		}
		builder.WriteByte('\n')
	}

	fmt.Print(builder.String())
}

func printUsage(messages ...string) {
	for _, m := range messages {
		if m != "" {
			fmt.Println(m)
		}
	}
	fmt.Printf("\t\naskii <image-path> [options]\t\n\n")
	flag.Usage()
}

func fatalf(format string, args ...any) {
	fmt.Printf(format, args...)
	os.Exit(1)
}
