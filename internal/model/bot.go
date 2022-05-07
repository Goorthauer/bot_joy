package model

import (
	"image/jpeg"
	"os"

	"github.com/disintegration/imaging"
)

func cropImage(filename string) error {
	src, err := imaging.Open(filename)
	if err != nil {
		return err
	}
	var imageWith = src.Bounds().Dx()
	var imageHeight = src.Bounds().Dy()
	src = imaging.CropAnchor(src, imageWith, imageHeight-15, imaging.Top)
	image, err := os.Create(filename)

	if err != nil {
		panic(err)
	}
	defer image.Close()
	jpeg.Encode(image, src, nil)
	return nil
}
