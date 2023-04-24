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

var re = regexp.MustCompile(`(?m)([\p{L}\d_\-\.])+`)

func Top10(input string) []string {
	res := []string{}

	prepare := TopStruct(input)

	for i := 0; i < 10 && i < len(prepare); i++ {
		res = append(res, prepare[i].Word)
	}

	return res
}

func TopStruct(input string) []Freq {
	prepare := []Freq{}

	validate := make(map[string]struct{})

	// - так же работает
	// f := func(c rune) bool {
	// 	return !unicode.IsLetter(c) && c != '-' && c != '.'
	// }
	// words := strings.FieldsFunc(input, f)

	words := re.FindAllString(input, -1)

	if len(words) == 0 {
		return prepare
	}

	for _, word := range words {
		word = strings.ToLower(strings.Trim(word, "."))
		if _, ok := validate[word]; ok {
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

	if valid == "-" || valid == "." {
		return
	}

	for _, word := range words {
		if strings.ToLower(strings.Trim(word, ".")) == valid {
			count++
		}
	}
	return
}
