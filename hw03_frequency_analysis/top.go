package hw03frequencyanalysis

import (
	"regexp"
	"sort"
	"strings"
)

type FrequencyElement struct {
	word  string
	count int
}

var rePunctuation = regexp.MustCompile(`(\p{L})[,!.;:]$`)

var reQuotes = regexp.MustCompile(`'(\p{L})'`)

func Top10(s string) []string {
	words := strings.Fields(s)

	freqMap := buildFreqMap(words)

	resultSlice := transformMapToSlice(freqMap)

	sortFreqSlice(resultSlice)

	resultSize := 10
	if len(resultSlice) < 10 {
		resultSize = len(resultSlice)
	}
	result := make([]string, resultSize)
	i := 0
	for i < resultSize {
		result[i] = resultSlice[i].word
		i++
	}

	return result
}

func sortFreqSlice(resultSlice []FrequencyElement) {
	sort.Slice(resultSlice, func(i, j int) bool {
		if resultSlice[i].count == resultSlice[j].count {
			return resultSlice[i].word < resultSlice[j].word
		}
		return resultSlice[i].count > resultSlice[j].count
	})
}

func buildFreqMap(words []string) map[string]int {
	freqMap := make(map[string]int)

	for _, word := range words {
		wordLower := strings.ToLower(word)

		r1 := rePunctuation.ReplaceAllString(wordLower, "$1")
		r2 := reQuotes.ReplaceAllString(r1, "$2")

		if r2 == "" || r2 == "-" {
			continue
		}

		freqMap[r2]++
	}

	return freqMap
}

func transformMapToSlice(freqMap map[string]int) []FrequencyElement {
	res := make([]FrequencyElement, len(freqMap))

	i := 0
	for k, v := range freqMap {
		res[i] = FrequencyElement{
			word:  k,
			count: v,
		}
		i++
	}

	return res
}
