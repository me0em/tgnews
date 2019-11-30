package tools

import (
	"io/ioutil"
	"sort"
	"strings"
	"encoding/json"
)

// Abs for Out-of-Place measure
func Abs(x int) int {
	if x < 0 {
		return -x
	} else {
		return x
	}
}

// IndexByValue func for Out-of-Place measure
// If value does not exists, return len(arr)
func IndexByValue(arr []string, value string) int {
	for ind, val := range arr {
		if value == val {
			return ind
		}
	}
	return len(arr)
}

// Delete all html tags
// return array with words
func Dvornik(article string) []string {
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
	// TODO:
	// if len(article) > 450 {
		// return strings.Fields(article[:450])
	// }
	return strings.Fields(article)
}

// Construct array of bi-grams, sorted with respect
// on frequency
func biGrams(words []string) []string {
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
	// Take only 450 top of bi-grams becouse of language profiles json
	// and linguistic a posteriori laws
	return MapToSortedCuttedArray(freqMap, 450)
}

// The measure between two sorted array of bi-grams.
// The size of two arrays must be equal. The result will be
// equal len(phonemes)^2 in the worst-case scenario.
func OutOfPlaceMeasure(x, y []string) int {
	distance := 0
	length := len(x)

	for ind, gram := range x {
		delta := IndexByValue(y, gram)
		if delta != length {
			distance += Abs(ind - delta)
		} else {
			distance += delta
		}
	}
	return distance
}

// Type represents language profilies data
type profiles struct {
	Data map[string][]string
}

// LoadProfile loads json file to type :profiles:
func LoadProfile(filename string) profiles {
	data, _ := ioutil.ReadFile(filename)
	var prf profiles
	_ = json.Unmarshal(data, &prf.Data)
	return prf
}

// DetectLanguage predicts the language of the text.
// If you don't want to determine :amount:, then pass amount = -1
func DetectLanguage(words []string, amount int) string {
	lgProfiles := LoadProfile("language_profiles.json")
	biGramsData := biGrams(words)
	length := len(biGramsData)

	// Calculate amount of bi_grams that will be passed to OutOfPlaceMeasure

	if length > 445 {
		biGramsData = biGramsData[:445]
	}

	if length < 445 && amount > length{
		amount = length
	}

	if amount > length {
		amount = length
	}

	if amount == -1 {
		amount = length
	}

	measure := amount * amount
    predictedLang := "other"

	for lang, profile := range(lgProfiles.Data) {
		distance := OutOfPlaceMeasure(profile[:amount], biGramsData[:amount])
		if distance < measure {
            measure = distance
            predictedLang = lang
		}
	}
	return predictedLang
}

// Type represents frequency analysis
type frequency struct {
	HashTable map[string]int // all bag-ow-words data
	Top       []string       // Only top of words
}

func BagOfWords(words []string, top int) frequency {
	var freq frequency
	freq.HashTable = make(map[string]int)

	for _, word := range words {
		if _, isKeyExists := freq.HashTable[word]; isKeyExists {
			freq.HashTable[word] += 1
		} else {
			freq.HashTable[word] = 1
		}
	}

	freq.Top = MapToSortedCuttedArray(freq.HashTable, top)
	return freq
}

func BagOfWordsOverFiles(filePaths []string, top int) frequency {
	var freq frequency
	freq.HashTable = make(map[string]int)

	for _, filePath := range filePaths {
		htmlData, _ := ioutil.ReadFile(filePath)
		words := Dvornik(string(htmlData))

		for _, word := range words {
			if _, isKeyExists := freq.HashTable[word]; isKeyExists {
				freq.HashTable[word] += 1
			} else {
				freq.HashTable[word] = 1
			}
		}
	}

	freq.Top = MapToSortedCuttedArray(freq.HashTable, top)
	return freq
}

// Map -> Sorted Array with only N=top top elements
func MapToSortedCuttedArray(someMap map[string]int, top int) []string {
	var topArray []string
	type kv struct {
		Key   string
		Value int
	}
	var ss []kv

	lengthOfMap := len(someMap)
	if top > lengthOfMap {
		top = lengthOfMap
	}

	for k, v := range someMap {
		ss = append(ss, kv{k, v})
	}

	sort.Slice(ss, func(i, j int) bool {
		return ss[i].Value > ss[j].Value
	})

	for _, kv := range ss {
		topArray = append(topArray, kv.Key)
		if len(topArray) == top {
			break
		}
	}
	return topArray
}
