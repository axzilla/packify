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

		layers := strings.Count(path, "/")
		for range layers {
			_, err := writer.WriteString(" ")
			if err != nil {
				return err
			}
		}

		_, err = writer.WriteString(path + "\n")
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		fmt.Println(err)
	}
}
