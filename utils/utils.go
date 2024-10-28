package utils

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"os"
	"strings"
)

func IsExtensionAllowed(ext string) bool {
	disallowedExts := ".png.jpg.jpeg.webp.pdf.zip"
	return !strings.Contains(disallowedExts, strings.ToLower(ext))
}

func FileSystem(remote string) (fs.FS, error) {
	if remote == "" {
		return os.DirFS("."), nil
	}

	resp, err := http.Get("https://" + remote + "/archive/refs/heads/main.zip")
	if err != nil {
		return nil, fmt.Errorf("download failed: %w", err)
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read failed: %w", err)
	}

	return zip.NewReader(bytes.NewReader(data), int64(len(data)))
}
