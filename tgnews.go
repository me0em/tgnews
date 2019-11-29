package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

var FileCounter int = 0
var paths []string // FIXME:  для распараллеливания, иначе на лету делать

func cleanPaths(paths []string) {
	for _, path := range paths {
		htmlData, _ := ioutil.ReadFile(path)
		_ = dvornik(string(htmlData))
	}
}

func main() {
	// Get path to file from command line arguments
	sysPath := os.Args[1]

	err := filepath.Walk(sysPath, walkFn)
	if err != nil {
		return
	}

	cleanPaths(paths)

	fmt.Printf("\n%d files have been cleared\n\n", FileCounter)
}

// walkFn function passed in filepath.Walk for recoursive find
// files in :path: directory for work with each files
func walkFn(path string, info os.FileInfo, err error) error {
	fi, err := os.Stat(path)
	if fi.Mode().IsRegular() {
		paths = append(paths, path)
		// htmlData, _ := ioutil.ReadFile(path)
		// _ = dvornik(string(htmlData))
		FileCounter += 1
	}
	return err
}

// Delete all html tags
// return array with words
func dvornik(article string) []string {
	length := len(article)
	memory_carrage := -1
	carrage := 0

	for true {
		char := string(article[carrage])

		if char == "\n" || char == "+" || char == "-" || char == "–" || char == "," || char == "’" || char == "“" {
			article = article[:carrage] + article[carrage+1:]
			length -= 1
			carrage -= 1
		}

		if char == "”" || char == "*" || char == "@" || char == "." || char == "/" || char == "(" || char == ")" {
			article = article[:carrage] + article[carrage+1:]
			length -= 1
			carrage -= 1
		}

		if char == "!" || char == "?" || char == "[" || char == "]" || char == "{" || char == "}" || char == "'" {
			article = article[:carrage] + article[carrage+1:]
			length -= 1
			carrage -= 1
		}

		if char == "’" || char == "`" || char == "%" || char == "#" || char == ":" || char == ";" || char == "&" || char == "1" || char == "2" || char == "3" || char == "4" || char == "5" || char == "6" || char == "7" || char == "8" || char == "9" || char == "0" {
			article = article[:carrage] + article[carrage+1:]
			length -= 1
			carrage -= 1
		}

		if char == "<" {
			memory_carrage = carrage
		}

		if char == ">" && memory_carrage != -1 {
			article = article[:memory_carrage] + article[carrage+1:]
			length -= carrage - memory_carrage + 1
			carrage -= carrage - memory_carrage + 1
			memory_carrage = -1
		}

		carrage += 1
		if carrage == length {
			break
		}
	}

	article = strings.ToLower(article)
	return strings.Fields(article)
}

// Construct array of bi-grams, sorted with respect
// on frequency
func biGrams(words []string) []string {
	var freqArr []string
	freqMap := make(map[string]int)
	var length int

	for _, word := range words {
		length = len(word)

		for i := 0; i < length-1; i++ {
			currStr := word[i : i+2]

			if _, isKeyExists := freqMap[currStr]; isKeyExists {
				freqMap[currStr] += 1
			} else {
				freqMap[currStr] = 1
			}
		}
	}

	// Sort map by value for frequency analysis
	type kv struct {
		Key   string
		Value int
	}
	var ss []kv
	for k, v := range freqMap {
		ss = append(ss, kv{k, v})
	}

	sort.Slice(ss, func(i, j int) bool {
		return ss[i].Value > ss[j].Value
	})

	for _, kv := range ss {
		freqArr = append(freqArr, kv.Key)
	}
	return freqArr
}

func bagOfWords(words []string) map[string]int {
	freqMap := make(map[string]int)
	var length int

	for _, word := range words {
		if _, isKeyExists := freqMap[word]; isKeyExists {
			freqMap[word] += 1
		} else {
			freqMap[word] = 1
		}
	}
	return freqMap
}

func bagOfWordsOverFiles(filePaths []string) map[string]int {
	freqMap := make(map[string]int)

	for _,filePath := range(filePaths) {
		htmlData, _ := ioutil.ReadFile(path)
		words = dvornik(string(htmlData))

		for _, word := range words {
			if _, isKeyExists := freqMap[word]; isKeyExists {
				freqMap[word] += 1
			} else {
				freqMap[word] = 1
			}
		}
	}
	return freqMap
}
