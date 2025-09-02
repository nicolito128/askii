package main

import (
	"flag"
	"fmt"
	"image"
	"image/color"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"os"
)

const resetColor = "\x1b[0m"

func truecolor(c color.Color) string {
	if os.Getenv("NO_COLOR") == "1" {
		return ""
	}
	r, g, b, _ := c.RGBA()
	return fmt.Sprintf("\x1b[38;2;%d;%d;%dm", r>>8, g>>8, b>>8)
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: askii <image>")
		return
	}

	f, err := os.Open(os.Args[1])
	if err != nil {
		panic(err)
	}
	defer f.Close()

	img, _, err := image.Decode(f)
	if err != nil {
		panic(err)
	}

	palette := []rune(" .:-=+*#%@")
	pw := len(palette) - 1

	b := img.Bounds()

	xStep := *flag.Int("x", 4, "x axis step")
	yStep := *flag.Int("y", 8, "y axis step")

	for y := b.Min.Y; y < b.Max.Y; y += yStep {
		for x := b.Min.X; x < b.Max.X; x += xStep {
			r, g, b, _ := img.At(x, y).RGBA()
			gray := (r*30 + g*59 + b*11) / 100
			gray >>= 8
			idx := int(gray) * pw / 255
			color := truecolor(img.At(x, y))
			fmt.Printf("%s%c%s", color, palette[idx], resetColor)
		}
		fmt.Println()
	}
}
