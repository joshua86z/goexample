package main

import (
	"bytes"
	"image"
	"image/color"
	"image/png"
	"log/slog"
	"os"
	"path"
	"path/filepath"

	"golang.org/x/image/draw"
	"golang.org/x/image/webp"
)

func main() {

	var imgs []image.Image
	var width, height int

	err := filepath.Walk(".", func(fileName string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Check if the current path is a file
		if !info.IsDir() {
			if ".webp" != path.Ext(fileName) {
				return err
			}

			file1, err := os.Open(fileName)
			if err != nil {
				return err
			}
			defer file1.Close()

			img1, err := webp.Decode(file1)
			if err != nil {
				return err
			}

			width = max(width, img1.Bounds().Dx())
			height = height + img1.Bounds().Dy()

			imgs = append(imgs, img1)
			slog.Info("read image", "name", fileName)
		}

		return nil
	})

	if err != nil {
		panic(err)
	}

	rect := image.Rect(0, 0, width, height)
	combined := image.NewRGBA(rect)

	// Fill the background with white color
	draw.Draw(combined, combined.Bounds(), &image.Uniform{C: color.White}, image.Point{}, draw.Src)

	y := 0

	for k, img := range imgs {
		x := (width - img.Bounds().Dx()) / 2
		y1 := y + img.Bounds().Dy()

		rect1 := image.Rectangle{
			Min: image.Point{X: x, Y: y},
			Max: image.Point{X: x + img.Bounds().Dx(), Y: y1},
		}

		draw.Draw(combined, rect1, img, image.Point{}, draw.Src)

		y = y1
		slog.Info("write image", "k", k, "x", rect1.Max.X, "y", rect1.Max.Y)
	}
	// Draw the first image on the left side

	// Draw the second image on the right side

	buf := new(bytes.Buffer)
	err = png.Encode(buf, combined)
	if err != nil {
		panic(err)
	}

	// Save the combined image to a file
	if err = os.WriteFile("ok.png", buf.Bytes(), 0644); err != nil {
		panic(err)
	}
}

func joinImages(img1, img2 image.Image) image.Image {
	rect := image.Rect(0, 0, max(img1.Bounds().Dx(), img2.Bounds().Dx()), img1.Bounds().Dy()+img2.Bounds().Dy())
	combined := image.NewRGBA(rect)

	// Fill the background with white color
	draw.Draw(combined, combined.Bounds(), &image.Uniform{C: color.White}, image.Point{}, draw.Src)

	// Draw the first image on the left side
	draw.Draw(combined, img1.Bounds().Add(image.Pt(0, 0)), img1, image.Point{}, draw.Src)

	// Draw the second image on the right side
	draw.Draw(combined, img2.Bounds().Add(image.Pt(0, img1.Bounds().Dy())), img2, image.Point{}, draw.Src)

	buf := new(bytes.Buffer)
	_ = png.Encode(buf, combined)

	ret, _ := png.Decode(buf)
	return ret
}
