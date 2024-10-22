package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

var path string

func main() {
	output := flag.String("output", "packify.txt", "Set a output e.g. myfile.txt, default is packify.txt")
	includePattern := flag.String("include", "*", "Glob patterns to include (comma-separated)")
	// excludePattern := flag.String("exclude", "", "Glob patterns to ignore (comma-separated)")
	flag.Parse()

	includePatterns := strings.Split(*includePattern, ",")
	// excludePatterns := strings.Split(*excludePattern, ",")

	file, err := os.Create(*output)
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()

	filetreeBuffer := bufio.NewWriter(file)

	var filecontentsBuffer bytes.Buffer

	// Filetree header
	_, err = filetreeBuffer.WriteString("=========================\n")
	_, err = filetreeBuffer.WriteString("Filetree\n")
	_, err = filetreeBuffer.WriteString("=========================\n")

	// Iterate recursive
	err = fs.WalkDir(os.DirFS("."), ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			fmt.Println(err)
		}

		// Ignore some defaults
		if strings.Contains(path, ".git") || strings.Contains(path, ".DS_Store") {
			return nil
		}

		// Include files
		matched := false
		for _, pattern := range includePatterns {
			m, err := filepath.Match(pattern, d.Name())
			if err != nil {
				return err
			}
			if m {
				matched = true
				break
			}
		}
		if !matched {
			return nil
		}

		// Write depth indent for filetree into buffer
		depth := strings.Count(path, "/")
		for range depth {
			_, err := filetreeBuffer.WriteString(" ")
			if err != nil {
				return err
			}
		}

		// Write file/folder names for filetree into buffer
		_, err = filetreeBuffer.WriteString(d.Name())
		if err != nil {
			return err
		}
		if d.IsDir() {
			_, err = filetreeBuffer.WriteString("/")
			if err != nil {
				return err
			}
		}
		_, err = filetreeBuffer.WriteString("\n")
		if err != nil {
			return err
		}

		// Write filecontents into buffer
		ext := filepath.Ext(d.Name())
		if !d.IsDir() && isExtensionAllowed(ext) {
			openedFile, err := os.ReadFile(path)
			if err != nil {
				return err
			}

			_, err = filecontentsBuffer.WriteString("=========================\n")
			_, err = filecontentsBuffer.WriteString("File: " + path + "\n")
			_, err = filecontentsBuffer.WriteString("=========================\n")
			_, err = filecontentsBuffer.Write(openedFile)
			_, err = filecontentsBuffer.WriteString("\n")
			if err != nil {
				return nil
			}
		}

		return nil
	})
	if err != nil {
		fmt.Println(err)
	}

	// Finally write filetree to disk
	err = filetreeBuffer.Flush()
	if err != nil {
		fmt.Println(err)
	}

	// Finally writes filecontents to disk
	byte, err := filecontentsBuffer.WriteTo(file)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(byte, "Bytes written!")
}

func isExtensionAllowed(ext string) bool {
	disallowedExts := ".png.jpg.jpeg.webp.pdf.zip"
	return !strings.Contains(disallowedExts, strings.ToLower(ext))
}
