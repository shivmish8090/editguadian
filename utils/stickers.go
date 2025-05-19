package utils


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
