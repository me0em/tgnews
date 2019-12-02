package tools

import (
	"io/ioutil"
	"fmt"
	"math"
	"path/filepath"
	// "regexp"
	"strings"
)

// Delete all html tags
// return array with words
func Dvornik(article string) []string {
	length := len(article)
	memory_carrage := -1
	carrage := 0

	for true {
		char := string(article[carrage])

		if IsPunctuation(char) {
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
// to frequency
func BiGrams(words []string) []string {
	freqMap := make(map[string]int)
	var length int

	for _, word := range words {
		runeWord := []rune(word)
		length = len(runeWord)

		for i := 0; i < length-1; i++ {
			a := runeWord[i]
			b := runeWord[i+1]
			currStr := fmt.Sprintf("%s%s", string(a), string(b))

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

// DetectLanguage predicts the language of the text.
// If you don't want to determine :amount:, then pass amount = -1
func DetectLanguage(words []string, amount int) string {
	lgProfiles := LoadProfile("language_profiles.json") // TODO
	biGramsData := BiGrams(words)

	length := len(biGramsData)

	// Calculate amount of bi_grams that will be passed to OutOfPlaceMeasure

	if length > 445 {
		biGramsData = biGramsData[:445]
	}

	if length < 445 && amount > length {
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

	for lang, profile := range lgProfiles.Data {
		distance := OutOfPlaceMeasure(profile[:amount], biGramsData[:amount])
		if distance < measure {
			measure = distance
			predictedLang = lang
		}
	}
	return predictedLang
}

// TODO: мб удалять HashTable если не нужен
// Type represents frequency analysis
type Frequency struct {
	Filename        string
	HashTable       map[string]int // all bag-ow-words data
	Top             []string       // Only top of words
	CuttedHashTable map[string]int // bag-ow-words data cutted with respect to Top
}

func (f Frequency) ReduceMap() {
	f.CuttedHashTable = make(map[string]int)

	for _, word := range f.Top {
		if amount, isKeyExists := f.HashTable[word]; isKeyExists {
			f.CuttedHashTable[word] = amount
		}
	}
}

func BagOfWords(words []string, top int) Frequency {
	var freq Frequency
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

// en.wikipedia.org/wiki/Tf%E2%80%93idf
func TFIDF(textFreq Frequency, corpusBagOfWords map[string]int, corpusLength int) map[string]float64 {
	result := make(map[string]float64)
	topAmount := 0
	for _, word := range textFreq.Top {
		topAmount += textFreq.HashTable[word]
	}

	for _, word := range textFreq.Top {
		result[word] = float64(textFreq.HashTable[word]) / float64(topAmount) * math.Log(float64(corpusLength)/float64(corpusBagOfWords[word]))
	}
	return result
}

func SimilarityMeasure(hashTableA, hashTableB map[string]float64) float64 {
	var (
		a Set
		b Set
	)

	a.Data = make(map[string]*void)
	b.Data = make(map[string]*void)
	// c.Data = make(map[string]*void)

	for k, _ := range hashTableA {
		a.Add(k)
	}
	for k, _ := range hashTableB {
		b.Add(k)
	}

	c := a.Intersection(b)

	intersectionSum := 0.0
	for v, _ := range c.Data {
		intersectionSum += Min(hashTableA[v], hashTableB[v])
	}

	c = a.Union(b)
	unionSum := 0.0
	for v, _ := range c.Data {
		unionSum += hashTableA[v] + hashTableB[v]
	}

	return float64(intersectionSum / (unionSum - intersectionSum))
}

// Property of pages groups with respect to similarity measure
type GroupProperties struct {
	Class    int
	Distance float64
}

func MakeThreads(textObjects []Frequency, bowOverFiles map[string]int, paths []string, lang string) map[int]*output {
	// var groupedTexts = make(map[string]*GroupProperties) // map with groups
	class := 1 // similar news will have the same class
	length := len(paths)
	var classArr = make([]int, length)
	var distances = make([]float64, length)
	// var titles []float64

	for i, _ := range classArr {
		classArr[i] = -1
	}

	var similarityCoeff float64
	if lang == "ru" {
		similarityCoeff = 0.16
	}
	if lang == "en" {
		similarityCoeff = 0.18
	}

	for counterUp, _ := range textObjects {

		for counter, _ := range textObjects {

			if counter > counterUp {

				tfidf1 := TFIDF(textObjects[counterUp], bowOverFiles, len(paths))
				tfidf2 := TFIDF(textObjects[counter], bowOverFiles, len(paths))
				measure := SimilarityMeasure(tfidf1, tfidf2)

				if measure > similarityCoeff {

					if classArr[counterUp] == -1 && classArr[counter] == -1 {
						classArr[counterUp] = class
						classArr[counter] = class
						distances[counterUp] = measure
						distances[counter] = measure
						class += 1
					}

					if classArr[counterUp] != -1 && classArr[counter] == -1 {
						classArr[counter] = classArr[counterUp]
						distances[counter] = measure
					}

					if classArr[counterUp] == -1 && classArr[counter] != -1 {
						classArr[counterUp] = classArr[counter]
						distances[counterUp] = measure
					}

					if classArr[counterUp] != -1 && classArr[counter] != -1 {
						if classArr[counterUp] != classArr[counter] {

							if measure > distances[counter] || measure > distances[counterUp] {
								if distances[counterUp] > distances[counter] {
									classArr[counter] = classArr[counterUp]
									distances[counter] = (distances[counter] + measure) / 2
									distances[counterUp] = (distances[counterUp] + measure) / 2
								} else {
									classArr[counterUp] = classArr[counter]
									distances[counterUp] = (distances[counterUp] + measure) / 2
									distances[counter] = (distances[counter] + measure) / 2
								}
							}
						}
					}
				}
			}
		}
	}

	var result = make(map[int]*output)
	var maxDistances = make(map[int]int)

	for ind, v := range classArr {
		if v != -1 {
			if result[v] == nil {
				result[v] = &output{"Test", nil}
			}
			currInd, _ := maxDistances[v]
			if distances[ind] > distances[currInd] {
				maxDistances[v] = ind
			}
			result[v].FilePaths = append(result[v].FilePaths, filepath.Base(paths[ind]))
		}
	}

	for k, v := range result {
		// re := regexp.MustCompile(`title" content="[a-zA-Zа-яА-Я0-9- \"\&\*’ ' :,. !?]*/>`)
		fmt.Println(paths[maxDistances[k]])
		data, _ := ioutil.ReadFile(paths[maxDistances[k]])
		// title := re.Find(data)
			// if title == "" {
					// fmt.Println("\n\ntitle", string(title))
					// fmt.Println(":", string(title))
			// }
		title := GetInnerSubstring(string(data), `title" content="`, "\"/>")

		if string(title) != "" {
			v.Title = string(title)
		}
	}

	return result
}

type output struct {
	Title     string
	FilePaths []string
}


func GetInnerSubstring(str string, prefix string, suffix string) string {
	var beginIndex, endIndex int
	beginIndex = strings.Index(str, prefix)
	if beginIndex == -1 {
		beginIndex = 0
		endIndex = 0
	} else if len(prefix) == 0 {
		beginIndex = 0
		endIndex = strings.Index(str, suffix)
		if endIndex == -1 || len(suffix) == 0 {
			endIndex = len(str)
		}
	} else {
		beginIndex += len(prefix)
		endIndex = strings.Index(str[beginIndex:], suffix)
		if endIndex == -1 {
			if strings.Index(str, suffix) < beginIndex {
				endIndex = beginIndex
			} else {
				endIndex = len(str)
			}
		} else {
			if len(suffix) == 0 {
				endIndex = len(str)
			} else {
				endIndex += beginIndex
			}
		}
	}

	return str[beginIndex:endIndex]
}
