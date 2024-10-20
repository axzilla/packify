package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

type TreePoint struct {
	name        string
	isDirectory bool
	children    []string
}

func main() {
	file, err := os.Create("packify.txt")
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()

	packify := bufio.NewWriter(file)

	var filesBuffer bytes.Buffer

	err = fs.WalkDir(os.DirFS("."), ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			fmt.Println(err)
		}

		if strings.Contains(path, ".git") || strings.Contains(path, ".DS_Store") {
			return nil
		}

		depth := strings.Count(path, "/")
		for range depth {
			_, err := packify.WriteString(" ")
			if err != nil {
				return err
			}
		}

		_, err = packify.WriteString(d.Name() + "\n")
		if err != nil {
			return err
		}

		ext := filepath.Ext(d.Name())
		if !d.IsDir() && isExtensionAllowed(ext) {
			openedFile, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			_, err = filesBuffer.Write(openedFile)
			if err != nil {
				return nil
			}
		}

		return nil
	})
	if err != nil {
		fmt.Println(err)
	}

	err = packify.Flush()
	if err != nil {
		fmt.Println(err)
	}

	byte, err := filesBuffer.WriteTo(file)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(byte, "Bytes written!")
}

func isExtensionAllowed(ext string) bool {
	disallowedExts := ".png.jpg.jpeg.webp.pdf.zip"
	return !strings.Contains(disallowedExts, strings.ToLower(ext))
}
