package utils

import (
        "image/gif"
        "image/jpeg"
        "image/png"
        "os"
        "path/filepath"

        "github.com/M3chD09/tgsconverter/libtgsconverter"
        ffmpeg "github.com/u2takey/ffmpeg-go"
        "golang.org/x/image/webp"
)
func Webp2Png(webpPath string) error {
	if filepath.Ext(webpPath) != ".webp" {
		panic("input file must be .webp")
	}

	reader, err := os.Open(webpPath)
	if err != nil {
		return err
	}
	defer reader.Close()

	img, err := webp.Decode(reader)
	if err != nil {
		return NewConvertError("webp2png", webpPath, "", err)
	}

	// Change extension to .png
	pngPath := strings.TrimSuffix(webpPath, filepath.Ext(webpPath)) + ".png"
	writer, err := os.Create(pngPath)
	if err != nil {
		return err
	}
	defer writer.Close()

	err = png.Encode(writer, img)
	if err != nil {
		return NewConvertError("webp2png", webpPath, pngPath, err)
	}

	return nil
}
