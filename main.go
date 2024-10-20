package main

import (
	"bufio"
	"fmt"
	"io/fs"
	"os"
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

	writer := bufio.NewWriter(file)
	defer writer.Flush()

	err = fs.WalkDir(os.DirFS("."), ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			fmt.Println(err)
		}

		if strings.Contains(path, ".git") || strings.Contains(path, ".DS_Store") {
			return nil
		}

		layers := strings.Count(path, "/")
		for range layers {
			_, err := writer.WriteString(" ")
			if err != nil {
				return err
			}
		}

		_, err = writer.WriteString(d.Name() + "\n")
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		fmt.Println(err)
	}
}
