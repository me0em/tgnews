package tools

import (
	"io/ioutil"
	"math"
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
type frequency struct {
	Filename        string
	HashTable       map[string]int // all bag-ow-words data
	Top             []string       // Only top of words
	CuttedHashTable map[string]int // bag-ow-words data cutted with respect to Top
}

func (f frequency) ReduceMap() {
	for _, word := range f.Top {
		if amount, isKeyExists := f.HashTable[word]; isKeyExists {
			f.CuttedHashTable[word] = amount
		}
	}
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
		tmpMap := make(map[string]int)
		htmlData, _ := ioutil.ReadFile(filePath)
		words := Dvornik(string(htmlData))

		for _, word := range words {
			if _, isKeyExists := tmpMap[word]; isKeyExists {
				tmpMap[word] += 1
			} else {
				tmpMap[word] = 1
			}
		}

		for key, _ := range tmpMap {
			if _, isKeyExists := freq.HashTable[key]; isKeyExists {
				freq.HashTable[key] += 1
			} else {
				freq.HashTable[key] = 1
			}
		}
	}
	return freq
}

// en.wikipedia.org/wiki/Tf%E2%80%93idf
func TFIDF(textFreq frequency, corpusBagOfWords map[string]int, corpusLength int) map[string]float64 {
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

func SimilarityMeasure(frequencyA, frequencyB frequency) float64 {
	var (
		a Set
		b Set
	)

	for _, v := range frequencyA.Top {
		a.Add(v)
	}
	for _, v := range frequencyB.Top {
		b.Add(v)
	}
	c := a.Intersection(b)

	intersectionSum := 0
	for v, _ := range c.Data {
		intersectionSum += Min(frequencyA.HashTable[v], frequencyB.HashTable[v])
	}

	c = a.Union(b)
	unionSum := 0
	for v, _ := range c.Data {
		unionSum += frequencyA.HashTable[v] + frequencyB.HashTable[v]
	}

	return float64(intersectionSum / (unionSum - intersectionSum))
}

// Property of pages groups with respect to similarity measure
type GroupProperties struct {
	Class    int
	Distance float64
}

func makeThreads(textObjects []frequency, lang string) {
	var groupedTexts map[string]GroupProperties // map with groups
	class := 1                                  // similar news will have the same class

	if lang == "ru" {
		similarityCoeff := 0.16
	}
	if lang == "en" {
		similarityCoeff := 0.22
	}

	bowOverFiles := BagOfWordsOverFiles(paths, 30)

	for counterUp, textObjectUp := range textObjects {

		for counter, textObject := range textObjects {

			if counter > counterUp {
				tfidf1 := TFIDF(textObjectUp, bowOverFiles)
				tfidf2 := TFIDF(textObject, bowOverFiles)
				measure := SimilarityMeasure(tfidf1, tfidf2)

				if measure > similarityCoeff {

					if groupedTexts[textObjectUp.Filename] == (GroupProperties{}) && groupedTexts[textObject.Filename] == (GroupProperties{}) {
						groupedTexts[textObjectUp.Filename] == GroupProperties{class, measure}
						groupedTexts[textObject.Filename] == GroupProperties{class, measure}
						class += 1
					}

					if groupedTexts[textObjectUp.Filename] != (GroupProperties{}) && groupedTexts[textObject.Filename] != (GroupProperties{}) {

						if groupedTexts[textObjectUp.Filename].Class != groupedTexts[textObject.Filename].Class {

							if measure > groupedTexts[textObjectUp.Filename].Distance || measure > groupedTexts[textObjectUp.Filename].Distance {
								if groupedTexts[textObjectUp.Filename].Distance > groupedTexts[textObjectUp.Filename].Distance {
									groupedTexts[textObjectUp.Filename].Distance = (groupedTexts[textObjectUp.Filename].Distance + measure) / 2
									groupedTexts[textObject.Filename] = GroupProperties{
										groupedTexts[textObjectUp.Filename].Class,
										(groupedTexts[textObject.Filename].Distance + measure) / 2,
									}
								} else {
									groupedTexts[textObjectUp.Filename] = GroupProperties{
										groupedTexts[textObject.Filename].Class,
										(groupedTexts[textObjectUp.Filename].Distance + measure) / 2,
									}
									groupedTexts[textObject.Filename].Distance = (groupedTexts[textObject.Filename].Distance + measure) / 2
								}
							}
						}

						if groupedTexts[textObjectUp.Filename] != (GroupProperties{}) && groupedTexts[textObject.Filename] == (GroupProperties{}) {
							groupedTexts[textObject.Filename] = GroupProperties{
								groupedTexts[textObjectUp.Filename].Class,
								measure,
							}
						}

						if groupedTexts[textObjectUp.Filename] == (GroupProperties{}) && groupedTexts[textObject.Filename] != (GroupProperties{}) {
							groupedTexts[textObjectUp.Filename] = GroupProperties{
								groupedTexts[textObject.Filename].Class,
								measure,
							}
						}
					}
				}
			}
		}
	}
	return groupedTexts
}
