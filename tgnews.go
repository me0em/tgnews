package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"tgnews/tools"
)

var (
	FileCounter int = 0

 	COMMAND = os.Args[1] // Get mode wich in program has been launched
	SYSPATH = os.Args[2] // Get path to file from command line arguments
)

func main() {
	if err := filepath.Walk(SYSPATH, walkFn); err != nil {
		return
	}

	fmt.Printf("\nProgram completed successfully. %d files have been processed\n\n", FileCounter)
}

// walkFn function passed in filepath.Walk for recoursive find
// files in :path: directory for work with each files
func walkFn(path string, info os.FileInfo, err error) error {
	fi, err := os.Stat(path)

	if fi.Mode().IsRegular() {
		FileCounter += 1

		if COMMAND == "languages" {
			htmlData, _ := ioutil.ReadFile(path)
			words := tools.Dvornik(string(htmlData))

			tools.DetectLanguage(words, 400)
		}

		if COMMAND == "news" {}
		if COMMAND == "categories" {
			// _ = tools.BagOfWords(words, 30)
		}
		if COMMAND == "threads" {}
		if COMMAND == "top" {}
	}

	return err
}
