package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(s string) (string, error) {
	str := []rune(s)
	var result strings.Builder
	countStr := len(str) // Количество символов
	for i, item := range str {
		if unicode.IsDigit(item) && i == 0 {
			return "", ErrInvalidString
		}
		if unicode.IsDigit(item) && unicode.IsDigit(str[i-1]) {
			return "", ErrInvalidString
		}
		// Проверяем следующий символ
		if countStr > i+1 && unicode.IsDigit(str[i+1]) {
			// Преобразуем символ в число
			strNextCount, err := strconv.Atoi(string(str[i+1]))
			if err != nil {
				return "", err
			}
			if strNextCount == 0 {
				continue
			}
		}
		// Если текущий символ - это цифра
		if unicode.IsDigit(item) {
			// Преобразуем символ в число
			count, err := strconv.Atoi(string(item))
			if err != nil {
				return "", err
			}
			if count == 0 {
				continue
			}
			// Добавляем символ в builder count раз
			result.WriteString(strings.Repeat(string(str[i-1]), count-1))
			continue
		}
		// Добавляем текущий символ в builder
		result.WriteRune(item)
	}
	return result.String(), nil
}
