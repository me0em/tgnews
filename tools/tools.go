package tools

import (
	"encoding/json"
	"io/ioutil"
	"sort"
	// "fmt"
)

/*
	TYPES
*/

// Type represents language profilies data
type profiles struct {
	Data map[string][]string
}


/*
	FUNCTIONS
*/

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

// LoadProfile loads json file to type :profiles:
func LoadProfile(filename string) profiles {
	data, _ := ioutil.ReadFile(filename)
	var prf profiles
	_ = json.Unmarshal(data, &prf.Data)
	return prf
}

func Min(x, y float64) float64 {
	if x <= y {
		return x
	}
	return y
}

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

// Punctuation symbols are noise for our algorithms
func IsPunctuation(char string) bool {
	if char == "\n" || char == "+" || char == "-" || char == "–" || char == "“"{
		return true
	}

	if char == "/" || char == "(" || char == ")" || char == "," || char == "’" {
		return true
	}

	if char == "”" || char == "*" || char == "@" || char == "." || char == "'" {
		return true
	}

	if char == "!" || char == "?" || char == "[" || char == "]" || char == "{" {
		return true
	}

	if char == "}" || char == "’" || char == "`" || char == "%" || char == "#" {
		return true
	}

	if char == ":" || char == ";" || char == "&" || char == "1" || char == "2" {
		return true
	}

	if char == "3" || char == "4" || char == "5" || char == "6" || char == "7" {
		return true
	}

	if char == "8" || char == "9" || char == "0" ||  char == "\""{
		return true
	}
	return false
}

/*
	SETS
*/

// Sets and operations with them
// Sets are used for find a most similarity text
// In particlular, union and intersection operations
// used for SimilarityMeasure()
type void struct{}

type Set struct {
	Data map[string]*void // empty set
}

func (s Set) Add(value string) {
	var illusion void
	s.Data[value] = &illusion
}

func (s Set) Delete(value string) {
	delete(s.Data, value)
}

func (s Set) Size() int {
	return len(s.Data)
}

func (s Set) IsExists(value string) bool {
	_, result := s.Data[value]
	return result
}

func (a Set) Union(b Set) Set {
	c := a
	for value := range(b.Data) {
		if c.IsExists(value) == false {
			c.Add(value)
		}
	}
	return c
}

func (a Set) Intersection(b Set) Set {
	var c Set
	c.Data = make(map[string]*void)
	for value := range(b.Data) {
		if a.IsExists(value) {
			c.Add(value)
		}
	}
	return c
}
