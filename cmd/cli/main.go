package main

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/axzilla/stackpack/utils"
)

var path string

func main() {
	output := flag.String("output", "stackpack.txt", "Set a output e.g. myfile.txt, default is stackpack.txt")
	includePattern := flag.String("include", "*", "Glob patterns to include (comma-separated)")
	excludePattern := flag.String("exclude", "", "Glob patterns to ignore (comma-separated)")
	repoUrl := flag.String("remote", "", "Get a remote repository, e.g. https://github.com/axzilla/stackpack")
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

	if *repoUrl != "" && !utils.IsValidGithubURL(*repoUrl) {
		fmt.Println("Not a valid GitHub repo URL")
		return
	}
	fileSystem, err := utils.FileSystem(*repoUrl)
	if err != nil {
		fmt.Println(err)
	}

	err = utils.WriteToBuffer(fileSystem, repoUrl, includePatterns, excludePatterns, &filetreeBuffer, &filecontentsBuffer)
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
