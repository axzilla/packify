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

	"github.com/axzilla/packify/utils"
)

var path string

func main() {
	output := flag.String("output", "packify.txt", "Set a output e.g. myfile.txt, default is packify.txt")
	includePattern := flag.String("include", "*", "Glob patterns to include (comma-separated)")
	excludePattern := flag.String("exclude", "", "Glob patterns to ignore (comma-separated)")
	remote := flag.String("remote", "", "Get a remote repository, e.g. github.com/axzilla/packify")
	flag.Parse()

	includePatterns := strings.Split(*includePattern, ",")
	excludePatterns := strings.Split(*excludePattern, ",")

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

	fileSystem, err := utils.FileSystem(*remote)
	if err != nil {
		fmt.Printf("Failed to get filesystem: %v\n", err)
		return
	}

	// Iterate recursive
	err = fs.WalkDir(fileSystem, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			fmt.Println(err)
		}

		// Ignore some defaults
		if strings.Contains(path, ".git") || strings.Contains(path, ".DS_Store") {
			return nil
		}

		if !d.IsDir() {
			// Include files
			included := false
			for _, pattern := range includePatterns {
				m, err := filepath.Match(pattern, d.Name())
				if err != nil {
					return err
				}
				if m {
					included = true
					break
				}
			}
			if !included {
				return nil
			}

			// Exclude files
			excluded := false
			for _, pattern := range excludePatterns {
				m, err := filepath.Match(pattern, d.Name())
				if err != nil {
					return err
				}
				if m {
					excluded = true
					break
				}
			}
			if excluded {
				return nil
			}
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
		if !d.IsDir() && utils.IsExtensionAllowed(ext) {
			// use fs. instead os. because fs works with every filesystem not only with local one on HDD
			openedFile, err := fs.ReadFile(fileSystem, path)
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
