package utils

import (
	"fmt"
	"image/png"
	"os"
	"path/filepath"
	"strings"

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

func ExtractFrames(inputPath string) ([]string, error) {
	// Prepare output paths in the same directory
	base := filepath.Base(inputPath)
	name := base[:len(base)-len(filepath.Ext(base))]

	midImage := fmt.Sprintf("%s_mid.png", name)
	endImage := fmt.Sprintf("%s_end.png", name)

	// Extract frame around 0.8s (middle zone)
	err := ffmpeg.Input(inputPath, ffmpeg.KwArgs{"ss": 0.8}).
		Output(midImage, ffmpeg.KwArgs{"frames:v": 1}).
		OverWriteOutput().
		Run()
	if err != nil {
		return nil, fmt.Errorf("failed to extract middle frame: %w", err)
	}

	// Extract frame near end (~0.1s before)
	err = ffmpeg.Input(inputPath, ffmpeg.KwArgs{"sseof": -0.1}).
		Output(endImage, ffmpeg.KwArgs{"frames:v": 1}).
		OverWriteOutput().
		Run()
	if err != nil {
		return nil, fmt.Errorf("failed to extract end frame: %w", err)
	}

	// Check both files exist
	if _, err := os.Stat(midImage); os.IsNotExist(err) {
		return nil, fmt.Errorf("middle image file not created")
	}
	if _, err := os.Stat(endImage); os.IsNotExist(err) {
		return nil, fmt.Errorf("end image file not created")
	}

	return []string{midImage, endImage}, nil
}
