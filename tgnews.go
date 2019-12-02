package main

import (
	// "encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"tgnews/tools"
)

var (
	FileCounter         = 0
	COMMAND             = os.Args[1] // Get mode wich in program has been launched
	SYSPATH             = os.Args[2] // Get path to file from command line arguments
	languageClusters    = make(map[string][]string)
	frequencyArray      = make(map[string][]tools.Frequency)
	bagOfWordsOverFiles = make(map[string]map[string]int)
)

func main() {

	if COMMAND == "test" {
		// t := tools.LoadProfile("language_profiles.json")
		// fmt.Println(t.Data["ru"])

		htmlData, _ := ioutil.ReadFile("/Users/atg/Desktop/Telegram Contest/DataClusteringSample0817/20191116/10/3621426919099766.html")
		words := tools.Dvornik(string(htmlData))
		for _, j := range tools.BiGrams(words) {
			fmt.Println(j)
		}
		language := tools.DetectLanguage(words, 2400)
		fmt.Println(language)
	}

	if COMMAND == "languages" {
		if err := filepath.Walk(SYSPATH, languagesWalkFn); err != nil {
			return
		}

		fmt.Printf("len of languageClusters['en']: %d\n", len(languageClusters["en"]))
		fmt.Printf("len of languageClusters['ru']: %d\n", len(languageClusters["ru"]))
		fmt.Printf("FileCounter: %d\n", FileCounter)

		// output, _ := json.Marshal(languageClusters)
		// fmt.Println(string(output))
	}

	if COMMAND == "threads" {
		bagOfWordsOverFiles["en"] = make(map[string]int)
		bagOfWordsOverFiles["ru"] = make(map[string]int)

		if err := filepath.Walk(SYSPATH, threadsWalkFn); err != nil {
			return
		}

		fmt.Printf("\n\nLen of bagOfWordsOverFiles['en']: %d\n", len(bagOfWordsOverFiles["en"]))
		fmt.Printf("Len of bagOfWordsOverFiles['ru']: %d\n", len(bagOfWordsOverFiles["ru"]))
		fmt.Printf("Len of languageClusters['en']: %d\n", len(languageClusters["en"]))
		fmt.Printf("Len of languageClusters['ru']: %d\n", len(languageClusters["ru"]))

		en_res := tools.MakeThreads(frequencyArray["en"], bagOfWordsOverFiles["en"], languageClusters["en"], "en")
		for k, v := range(en_res) {
			fmt.Println("\n", k, *v, "\n")
		}

		ru_res := tools.MakeThreads(frequencyArray["ru"], bagOfWordsOverFiles["ru"], languageClusters["ru"], "ru")
		for k, v := range(ru_res) {
			fmt.Println("\n", k, *v, "\n")
		}

		fmt.Printf("\n%d files have been proccessed\n", FileCounter)
	}

}

// threads clustering
// walkFn function passed in filepath.Walk for recoursive find
// files in :path: directory for work with each files
func threadsWalkFn(path string, info os.FileInfo, err error) error {
	FileCounter += 1

	fi, err := os.Stat(path)
	if fi.Mode().IsRegular() {

		htmlData, _ := ioutil.ReadFile(path)
		words := tools.Dvornik(string(htmlData))
		language := tools.DetectLanguage(words, 400)

		if language == "ru" {
			languageClusters["ru"] = append(languageClusters["ru"], path)
			freq := tools.BagOfWords(words, 30)
			freq.ReduceMap()
			frequencyArray["ru"] = append(frequencyArray["ru"], freq)

			for _, word := range words {
				if _, isKeyExists := bagOfWordsOverFiles["ru"][word]; isKeyExists {
					bagOfWordsOverFiles["ru"][word] += 1
				} else {
					bagOfWordsOverFiles["ru"][word] = 1
				}
			}
		}

		if language == "en" {
			languageClusters["en"] = append(languageClusters["en"], path)
			freq := tools.BagOfWords(words, 30)
			freq.ReduceMap()
			frequencyArray["en"] = append(frequencyArray["en"], freq)

			for _, word := range words {
				if _, isKeyExists := bagOfWordsOverFiles["en"][word]; isKeyExists {
					bagOfWordsOverFiles["en"][word] += 1
				} else {
					bagOfWordsOverFiles["en"][word] = 1
				}
			}
		}
	}
	return err
}

// language clustering
// walkFn function passed in filepath.Walk for recoursive find
// files in :path: directory for work with each files
func languagesWalkFn(path string, info os.FileInfo, err error) error {
	FileCounter += 1
	fi, err := os.Stat(path)
	if fi.Mode().IsRegular() {

		htmlData, _ := ioutil.ReadFile(path)
		words := tools.Dvornik(string(htmlData))
		language := tools.DetectLanguage(words, 350)

		if language == "ru" {
			languageClusters["ru"] = append(languageClusters["ru"], path)
		}
		if language == "en" {
			languageClusters["en"] = append(languageClusters["en"], path)
		}
	}
	return err
}
