package hw02unpackstring

import (
	"errors"
	"strings"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(input string) (string, error) {
	// Place your code here.
	var bf strings.Builder

	for i := 0; i < len([]rune(input))-1; i++ {
		r := []rune(input)[i]
		nr := []rune(input)[i+1]

		if unicode.IsDigit(nr) {
			if unicode.IsDigit(r) {
				return "", ErrInvalidString
			}

			bf.WriteString(strings.Repeat(string(r), int(nr-'0')))
		} else {
			if unicode.IsDigit(r) {
				if i == 0 {
					return "", ErrInvalidString
				}
				if i == len([]rune(input))-2 {
					bf.WriteString(string(nr))
				}
				continue
			}

			bf.WriteString(string(r))

			if i == len([]rune(input))-2 {
				bf.WriteString(string(nr))
			}
		}
	}
	return bf.String(), nil
}
