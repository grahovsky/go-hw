package hw03frequencyanalysis

import (
	"regexp"
	"sort"
	"strings"
)

type Freq struct {
	Word  string
	Count int
}

// find words.
// var re = regexp.MustCompile(`(?m)([\p{L}\d_\-\.])+`)

// find delimeter.
var re = regexp.MustCompile(`(?m)(!|\.|,|\s-|\(|\))*\s+`)

func Top10(input string) []string {
	res := make([]string, 0, 10)

	prepare := TopStruct(input)

	for i := 0; i < 10 && i < len(prepare); i++ {
		res = append(res, prepare[i].Word)
	}

	return res
}

func TopStruct(input string) []Freq {
	prepare := []Freq{}
	validate := make(map[string]struct{})

	words := re.Split(input, -1)

	if len(words) == 0 {
		return prepare
	}

	for _, word := range words {
		word = strings.ToLower(word)
		if _, ok := validate[word]; ok || word == "" {
			continue
		}

		prepare = append(prepare, Freq{Word: word, Count: CountWords(words, word)})
		validate[word] = struct{}{}
	}

	sort.Slice(prepare, func(i, j int) bool {
		iv, jv := prepare[i], prepare[j]
		switch {
		case iv.Count != jv.Count:
			return iv.Count > jv.Count
		default:
			return iv.Word < jv.Word
		}
	})

	return prepare
}

func CountWords(words []string, valid string) (count int) {
	count = 0

	for _, word := range words {
		if strings.ToLower(word) == valid {
			count++
		}
	}
	return
}
