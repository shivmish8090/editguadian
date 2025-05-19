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
	reader, err := os.Open(webpPath)
	if err != nil {
		return err
	}
	defer reader.Close()

	img, err := webp.Decode(reader)
	if err != nil {
		return fmt.Errorf("Failed to convert webp2png\nWebpPath: %s\nPngPath: %s\nError: %w", webpPath, "", err)
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
		return fmt.Errorf("Failed to convert webp2png\nWebpPath: %s\nPngPath: %s\nError: %w", webpPath, pngPath, err)
	}

	return nil
}

func ExtractFrames(inputPath string, count int) ([]string, error) {
  
	base := filepath.Base(inputPath)
	name := base[:len(base)-len(filepath.Ext(base))]

	rand.Seed(time.Now().UnixNano())

	var extracted []string
	for i := 0; i < count; i++ {
		randomTime := rand.Float64() * 3 // random between 0s and 3s
		imageName := fmt.Sprintf("%s_%d.png", name, i)

		ffmpeg.Input(inputPath, ffmpeg.KwArgs{"ss": fmt.Sprintf("%.2f", randomTime)}).
			Output(imageName, ffmpeg.KwArgs{"frames:v": 1}).
			OverWriteOutput().
			Run()

		// Add only if file was created
		if _, statErr := os.Stat(imageName); statErr == nil {
			extracted = append(extracted, imageName)
		}
	}

	if len(extracted) == 0 {
		return nil, fmt.Errorf("no images were extracted")
	}

	return extracted, nil
}
