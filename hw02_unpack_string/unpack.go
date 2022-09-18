package hw02unpackstring

import (
	"errors"
	"strings"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(input string) (string, error) {
	if input == "" {
		return input, nil
	}

	b := strings.Builder{}
	prev := rune(0)
	for _, current := range input {
		err := unpackRune(current, prev, &b)
		if err != nil {
			return "", err
		}
		prev = current
	}
	err := unpackRune(rune(0), prev, &b)
	if err != nil {
		return "", err
	}

	return b.String(), nil
}

func unpackRune(current rune, prev rune, b *strings.Builder) error {
	if unicode.IsDigit(current) {
		if unicode.IsDigit(prev) || prev == 0 {
			return ErrInvalidString
		}
		count := int(current - '0')
		if count > 0 {
			b.WriteString(strings.Repeat(string(prev), count))
		}
	} else if !unicode.IsDigit(prev) && prev != 0 {
		b.WriteRune(prev)
	}
	return nil
}
