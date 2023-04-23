package hw03frequencyanalysis

import (
	"regexp"
	"sort"
	"strings"
	"unicode"
)

type Freq struct {
	Word  string
	Count int
}

func Top10(input string) []string {
	prepare := []Freq{}
	res := []string{}
	validate := make(map[string]struct{})

	// words := strings.Split(input, " ")
	// words := strings.Fields(input)
	f := func(c rune) bool {
		return !unicode.IsLetter(c)
	}
	words := strings.FieldsFunc(input, f)

	if len(words) == 0 {
		return res
	}

	for _, word := range words {
		word = strings.ToLower(word)
		if _, ok := validate[word]; ok {
			continue
		}

		re, err := regexp.Compile(`(?mi)(\s|\.|\,)` + word + `(\s|\.|\,)`)
		if err != nil {
			return res
		}
		count := re.FindAllString(input, -1)

		prepare = append(prepare, Freq{Word: word, Count: len(count)})
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

	for i := 0; i < 10; i++ {
		res = append(res, prepare[i].Word)
	}

	return res
}
