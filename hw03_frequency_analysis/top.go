package hw03frequencyanalysis

import (
	"regexp"
	"sort"
)

const maxResultSize = 10

func Top10(input string) (result []string) {
	result = make([]string, 0, maxResultSize)
	if input == "" {
		return
	}

	splitter := regexp.MustCompile(`[\n\s\t]+`)
	freqMap := make(map[string]int)
	for _, s := range splitter.Split(input, -1) {
		if s == "" {
			continue
		}
		count := freqMap[s]
		count++
		freqMap[s] = count
	}

	type frequency struct {
		str   string
		count int
	}
	freqSlice := make([]frequency, 0, len(freqMap))
	for k, v := range freqMap {
		freqSlice = append(freqSlice, frequency{
			str:   k,
			count: v,
		})
	}
	sort.Slice(freqSlice, func(i, j int) bool {
		if freqSlice[i].count != freqSlice[j].count {
			return freqSlice[i].count > freqSlice[j].count
		}
		return freqSlice[i].str < freqSlice[j].str
	})
	sliceLength := len(freqSlice)
	if sliceLength > maxResultSize {
		sliceLength = maxResultSize
	}
	for _, freq := range freqSlice[:sliceLength] {
		result = append(result, freq.str)
	}

	return result
}
