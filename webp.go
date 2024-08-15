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

	buf := new(bytes.Buffer)
	err = png.Encode(buf, combined)
	if err != nil {
		panic(err)
	}

	f, err := os.Create("ok.png")
	if err != nil {
		panic(err)
	}
	b := buf.Bytes()
	for i := 0; i < len(b); i += 1024 {
		_, err = f.Write(b[i : i+1024])
		if err != nil {
			panic(err)
		}

		slog.Info("write file", "i", i)
	}

	// Save the combined image to a file
	if err = os.WriteFile("ok.png", buf.Bytes(), 0644); err != nil {
		panic(err)
	}
}
