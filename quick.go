package main

// go fmt ./... && go vet ./... && go test && go run quick.go -cpuprofile cpu.prof && echo top | go tool pprof cpu.prof

import (
	"flag"
	"fmt"
	"log"
	"os"
	"runtime/pprof"
	"strings"

	"github.com/erikbryant/dictionaries"
)

var (
	cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")
)

// This program looks at a word list and determines whether there are
// any sets of 5 words that, taken together, would eliminate 25 unique
// letters.

// loadDict returns the words from the dictionary
func loadDict(wordLen int) []string {
	// words := dictionaries.LoadFile("../dictionaries/wordleGuessable.dict")
	words := dictionaries.LoadFile("../dictionaries/merged.dict")
	words = dictionaries.FilterByLen(words, wordLen)
	words = dictionaries.SortUnique(words)

	return words
}

// removeDupleLetters returns a list of words that have no duplicate letters
func removeDupeLetters(words []string) []string {
	pruned := []string{}

	for _, word := range words {
		used := map[rune]bool{}
		double := false
		for _, val := range word {
			if used[val] {
				double = true
				break
			}
			used[val] = true
		}
		if !double {
			pruned = append(pruned, word)
		}
	}

	return pruned
}

// filter removes any words that contain any of the letters in 'letters'
func filter(words []string, letters string) []string {
	pruned := []string{}

	for _, word := range words {
		if strings.ContainsAny(word, letters) {
			continue
		}
		pruned = append(pruned, word)
	}

	return pruned
}

// letterElimination returns the five words that eliminate 25 distinct letters,
// or []string{} if no set of five exists
func letterElimination(words []string) ([]string, int) {
	words = removeDupeLetters(words)

	// Words that only have one vowel
	aWords := filter(words, "eiou")
	eWords := filter(words, "aiou")
	iWords := filter(words, "aeou")
	oWords := filter(words, "aeiu")
	uWords := filter(words, "aeio")

	maxUsed := 0
	maxGuesses := []string{}

	for _, a := range aWords {
		for _, e := range eWords {
			if strings.ContainsAny(e, a) {
				continue
			}
			for _, i := range iWords {
				if strings.ContainsAny(i, a+e) {
					continue
				}
				for _, o := range oWords {
					if strings.ContainsAny(o, a+e+i) {
						continue
					}
					for _, u := range uWords {
						if strings.ContainsAny(u, a+e+i+o) {
							continue
						}
						// If we got here, all guess letters were distinct
						maxUsed = 25
						maxGuesses = []string{a, e, i, o, u}
						fmt.Println(maxGuesses, maxUsed)
					}
				}
			}
		}
	}

	return maxGuesses, maxUsed
}

func main() {
	fmt.Printf("Welcome to quick\n\n")

	flag.Parse()
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	wordLen := 5
	words := loadDict(wordLen)

	guesses, eliminated := letterElimination(words)
	fmt.Println()
	fmt.Println("Guessing:", guesses, "eliminates", eliminated, "distinct letters")
}
