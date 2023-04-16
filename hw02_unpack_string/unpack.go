package hw02unpackstring

import (
	"errors"
	"strings"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(input string) (string, error) {
	var bf strings.Builder
	runes := []rune(input)

	i := 0

	for {
		if i > len(runes)-1 {
			break
		}

		symb, err := getSymb(&i, runes)
		if err != nil {
			return "", ErrInvalidString
		}
		count := getCount(&i, runes)

		bf.WriteString(strings.Repeat(symb, count))
	}

	return bf.String(), nil
}

func getCount(i *int, runes []rune) (count int) {
	count = 1

	if *i > len(runes)-1 {
		return count
	}

	r := runes[*i]
	if unicode.IsDigit(r) {
		count = int(r - '0')
		*i++
	}

	return count
}

func getSymb(i *int, runes []rune) (symb string, _ error) {
	if unicode.IsDigit(runes[*i]) {
		return "", ErrInvalidString
	}

	if *i < len(runes)-1 && string(runes[*i]) == "\\" {
		symb = string(runes[*i+1])
		if !(symb == "\\" || unicode.IsDigit(runes[*i+1])) {
			return "", ErrInvalidString
		}
		*i += 2
	} else {
		symb = string(runes[*i])
		*i++
	}

	return symb, nil
}
