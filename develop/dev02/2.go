package two

import (
	"errors"
	"regexp"
	"strconv"
	"strings"
	"unicode/utf8"
)

var (
	ErrInvalidLine  = errors.New("invalid line")
	ErrInvalidValue = errors.New("invalid value")
)

func StringUnpacking(str string) (string, error) {
	// проверка, что escape неккоретнен
	err := isIncorrectEsc(str)
	if err != nil {
		return "", err
	}
	// если специальных символов нет, то возвращаем строку
	r, _ := regexp.Compile(`[0-9]||\\`)
	if !r.Match([]byte(str)) {
		return str, nil
	}
	var b strings.Builder
	b.Grow(len(str) * 5)
	size := utf8.RuneCountInString(str)
	rep := ""
	X := -1
	isEsc := false

	for i, c := range str {
		// если предыдущий символ \, то следующий воспринимаем как обычный символ
		if string(c) == "\\" && !isEsc {
			if size-1 == i {
				return "", ErrInvalidLine
			}
			isEsc = true
			continue
		}

		x, err := strconv.Atoi(string(c))
		// добавление символов
		if err != nil || isEsc {
			if string(c) != "" && X == -1 {
				b.WriteString(rep)
			} else {
				b.WriteString(strings.Repeat(rep, X))
				X = -1
			}
			if size-1 == i {
				b.WriteRune(c)
			}
			rep = string(c)
			isEsc = false
			continue

		}

		if rep == "" {
			return "", ErrInvalidLine
		}
		// на случай, если рядом несколько чисел
		X, err = addNextValue(x, X)
		if err != nil {
			return "", err
		}
		if size-1 == i {
			b.WriteString(strings.Repeat(rep, X))
		}
	}
	return b.String(), nil
}

// if old < 0, then the other value will come first
func addNextValue(x, old int) (int, error) {
	if old < 0 {
		return x, nil
	}
	if old == 0 {
		return 0, ErrInvalidValue
	}
	return old*10 + x, nil
}

func isIncorrectEsc(str string) error {
	if str == "" {
		return nil
	}
	strR := []rune(str)
	cEsc := 0
	for i := len(strR) - 1; i > 0; i-- {
		if string(strR[i]) == "\\" {
			cEsc += 1
			continue
		}
		break
	}
	if cEsc%2 != 0 {
		return ErrInvalidLine
	}
	return nil
}
