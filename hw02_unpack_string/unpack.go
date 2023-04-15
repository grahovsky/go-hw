package hw02unpackstring

import (
	"errors"
	"strings"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(input string) (string, error) {
	var bf strings.Builder
	runes := []rune(reverseString(input))

	i := 0

	for {
		if i > len(runes)-1 {
			break
		}

		count := getCount(&i, runes)

		// переделать
		if i > len(runes)-1 || unicode.IsDigit(runes[i]) {
			return "", ErrInvalidString
		}

		char := getChar(&i, runes)

		bf.WriteString(strings.Repeat(char, count))
	}

	return reverseString(bf.String()), nil
}

func getCount(i *int, runes []rune) (count int) {
	r := runes[*i]
	count = 1

	if unicode.IsDigit(r) {
		count = int(r - '0')
		*i++
	}

	return count
}

func getChar(i *int, runes []rune) (char string) {
	char = string(runes[*i])

	if unicode.IsDigit(runes[*i]) {
		char = ""
	}

	*i++

	return char
}

func reverseString(str string) (result string) {
	for _, v := range str {
		result = string(v) + result
	}
	return
}
