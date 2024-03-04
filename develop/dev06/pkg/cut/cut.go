package cut

import (
	"errors"
	"os"
	"runtime"
	"strconv"
	"strings"

	flag "github.com/spf13/pflag"
)

var (
	ErrNoFileName      = errors.New("can't get filename")
	ErrIncorrectFields = errors.New("incorrect fields")
)

func Cut() (string, error) {
	// получаем параметры
	var fields string
	flag.StringVarP(&fields, "fields", "f", "", "choose fields")

	var dilimiter string
	flag.StringVarP(&dilimiter, "dilimiter", "d", " ", "choose dilimiter")

	var separated bool
	flag.BoolVarP(&separated, "separated", "s", false, "only separated")

	flag.Parse()
	filename := flag.Arg(0)
	if !strings.Contains(filename, ".") {
		return "", ErrNoFileName
	}

	// обрезаем
	res, err := cut(filename, fields, dilimiter, separated)
	if err != nil {
		return "", err
	}

	return res, nil
}

func cut(filename, fields, dilimiter string, separated bool) (string, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return "", err
	}
	if len(data) == 0 {
		return "", nil
	}

	// получаем строки
	var lines []string
	if runtime.GOOS == "windows" {
		lines = strings.Split(string(data), "\r\n")
	} else {
		lines = strings.Split(string(data), "\n")
	}

	if fields == "" && !separated {
		return string(data), nil
	}

	var fieldsInt []int
	if fields != "" {
		fieldsInt, err = getFields(fields)
		if err != nil {
			return "", err
		}
	}

	res := cutLines(lines, dilimiter, separated, fieldsInt)
	return res, nil
}

func getFields(fields string) ([]int, error) {
	res := make([]int, 0, len(fields))

	l := strings.Split(fields, " ")
	if len(l) == 1 {
		l = strings.Split(fields, "\t")
	}

	for _, str := range l {
		v, err := strconv.Atoi(str)
		if err != nil || v < 1 {
			return nil, ErrIncorrectFields
		}
		res = append(res, v)
	}

	return res, nil
}

func cutLines(lines []string, dilimiter string, separated bool, fields []int) string {
	var b strings.Builder
	for j, str := range lines {
		line := strings.Split(str, dilimiter)
		if len(line) == 1 {
			if !separated {
				b.WriteString(str + "\n")
			}
			continue
		}
		if len(fields) == 0 {
			b.WriteString(str + "\n")
			continue
		}
		// получаем последний индекс, чтобы не добавлять к элементу dilimiter
		maxI := findLastIndexOfField(line, fields)
		for i, field := range fields {
			if len(line) < field {
				continue
			}
			if i == maxI {
				b.WriteString(line[field-1])
				break
			}
			b.WriteString(line[field-1] + dilimiter)
		}

		if j != len(lines)-1 {
			b.WriteRune('\n')
		}
	}
	res := b.String()
	if res[len(res)-1] == '\n' {
		res = res[:len(res)-1]
	}
	return res
}

// if result == -1, then all fields are greater than len of ljne
func findLastIndexOfField(line []string, fields []int) int {
	maxField := len(line) - 1

	for j := len(fields) - 1; j >= 0; j-- {
		if fields[j]-1 <= maxField {
			return j
		}
	}
	return -1
}
