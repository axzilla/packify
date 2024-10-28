package main

import (
	"bytes"
	"flag"
	"fmt"
	"os"
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

	var filetreeBuffer bytes.Buffer
	var filecontentsBuffer bytes.Buffer

	err = utils.WriteToBuffer(remote, includePatterns, excludePatterns, &filetreeBuffer, &filecontentsBuffer)
	if err != nil {
		fmt.Println(err)
	}

	// Finally write filetree to disk
	treeBytes, err := filetreeBuffer.WriteTo(file)
	if err != nil {
		fmt.Println(err)
	}

	// Finally writes filecontents to disk
	contentBytes, err := filecontentsBuffer.WriteTo(file)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(treeBytes+contentBytes, "Bytes written!")
}
