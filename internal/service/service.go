package service

import (
	"errors"
	"strings"

	"github.com/Yandex-Practicum/go1fl-sprint6-final/pkg/morse"
)
// Ошибки
var (
	ErrEmptyInput  = errors.New("input string is empty")
	ErrDecodeInput = errors.New("failed to decode morse code: invalid morse sequence")
	ErrEncodeInput = errors.New("failed to encode text to morse code: text contains unsupported characters")
)

func Automatic_Detection(input string) (string, error) {
	if input == "" {
		return "", ErrEmptyInput
	}

	if strings.ContainsAny(input, morse.Period+"-") {
		result := morse.ToText(input)
		if result == "" {
			return "", ErrDecodeInput
		}
		return result, nil
	}

	result := morse.ToMorse(input)
	if result == "" {
		return "", ErrEncodeInput
	}
	return result, nil
}
