package utilities

import (
	"errors"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"os"
	"path/filepath"

	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/goregular"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/math/fixed"
)

func GenerateImage(text string) error {
	img := image.NewRGBA(image.Rect(0, 0, 400, 150))
	draw.Draw(img, img.Bounds(), &image.Uniform{color.RGBA{30, 40, 50, 255}}, image.Point{}, draw.Src)

	parsedFont, err := opentype.Parse(goregular.TTF)
	if err != nil {
		return errors.New(fmt.Sprintf("Failed to parse font: %v", err))
	}

	face, err := opentype.NewFace(parsedFont, &opentype.FaceOptions{Size: 28, DPI: 72, Hinting: font.HintingNone})
	if err != nil {
		return errors.New(fmt.Sprintf("Failed to create font face: %v", err))
	}

	drawer := &font.Drawer{Dst: img, Src: image.NewUniform(color.RGBA{255, 255, 255, 255}), Face: face, Dot: fixed.Point26_6{X: fixed.I(50), Y: fixed.I(85)}}
	drawer.DrawString(text)

	directory, _ := os.Getwd()
	imagePath := filepath.Join(directory, "images", fmt.Sprintf("%s.png", text))
	file, err := os.Create(imagePath)
	if err != nil {
		return errors.New(fmt.Sprintf("There was a problem creating the image: %v", err))
	}

	defer file.Close()

	if err := png.Encode(file, img); err != nil {
		return errors.New(fmt.Sprintf("Failed to encode PNG: %v", err))
	}

	return nil
}
