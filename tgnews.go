package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"tgnews/tools"
)

var FileCounter int = 0

func main() {
	sysPath := os.Args[1] // Get path to file from command line arguments

	if err := filepath.Walk(sysPath, walkFn); err != nil {
		return
	}

	fmt.Printf("\n%d files have been cleared\n\n", FileCounter)
}

// walkFn function passed in filepath.Walk for recoursive find
// files in :path: directory for work with each files
func walkFn(path string, info os.FileInfo, err error) error {
	fi, err := os.Stat(path)
	if fi.Mode().IsRegular() {
		htmlData, _ := ioutil.ReadFile(path)
		words := tools.Dvornik(string(htmlData))
		_ = tools.BagOfWords(words, 30)
		// fmt.Printf("%+v\n", freq)
		tools.DetectLanguage(words, -1)

		FileCounter += 1
	}
	return err
}
