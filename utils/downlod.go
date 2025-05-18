package utils

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

func DownloadFile(url string) (string, error) {
	fileName := filepath.Base(url)
	tempPath := filepath.Join(os.TempDir(), fileName)

	out, err := os.Create(tempPath)
	if err != nil {
		return "", err
	}
	defer out.Close()

	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("bad status: %s", resp.Status)
	}

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return "", err
	}

	return tempPath, nil
}
