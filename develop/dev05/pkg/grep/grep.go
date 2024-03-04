package grep

import (
	"errors"
	"fmt"
	"os"
	"regexp"
	"runtime"
	"strconv"
	"strings"

	flag "github.com/spf13/pflag"
)

var (
	ErrNoFileName        = errors.New("can't get filename")
	ErrNoPattern         = errors.New("can't get pattern")
	ErrIncompatibleFlags = errors.New("incompatible flags")
)

func Grep() (string, error) {
	// получение параметров
	var after int
	flag.IntVarP(&after, "after", "A", -1, "print N lines after match")

	var before int
	flag.IntVarP(&before, "before", "B", -1, "print N lines before match")

	var context int
	flag.IntVarP(&context, "context", "C", -1, "print N lines near match")

	var count bool
	flag.BoolVarP(&count, "count", "c", false, "print count of matches")

	var ignoreCase bool
	flag.BoolVarP(&ignoreCase, "ignore-case", "i", false, "ignore case")

	var invert bool
	flag.BoolVarP(&invert, "invert", "v", false, "invert")

	var fixed bool
	flag.BoolVarP(&fixed, "fixed", "F", false, "exact match")

	var lineNum bool
	flag.BoolVarP(&lineNum, "line-num", "n", false, "print lines nums")

	flag.Parse()

	if context >= 0 && (after >= 0 || before >= 0) {
		return "", ErrIncompatibleFlags
	}

	filename := flag.Arg(1)
	if !strings.Contains(filename, ".") {
		return "", ErrNoFileName
	}

	pattern := flag.Arg(0)
	if pattern == "" {
		return "", ErrNoPattern
	}

	res, err := grep(filename, pattern, before, after, context, ignoreCase, invert, fixed, lineNum, count)
	if err != nil {
		return "", err
	}

	return res, nil
}

func grep(filename, pattern string, before, after, context int, ignoreCase, invert, fixed, lineNum, count bool) (string, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return "", err
	}

	if _, err = regexp.Compile(pattern); err != nil {
		return "", err
	}

	if len(data) == 0 {
		return "", nil
	}

	// в зависимости от платформы, разбиваем на строки
	var lines []string
	if runtime.GOOS == "windows" {
		lines = strings.Split(string(data), "\r\n")
	} else {
		lines = strings.Split(string(data), "\n")
	}

	// обработка
	res := grepLines(lines, pattern, before, after, context, ignoreCase, invert, fixed, lineNum)
	if count {
		return strconv.Itoa(len(res)), nil
	}

	var b strings.Builder
	for _, str := range res {
		b.WriteString(str + "\n")
	}

	return b.String(), nil
}

func grepLines(lines []string, pattern string, before, after, context int, ignoreCase, invert, fixed, lineNum bool) []string {
	// если нужно игнорировать регистр, приводим к нижнему
	if ignoreCase {
		pattern = strings.ToLower(pattern)
	}
	if fixed {
		pattern = fmt.Sprintf("^%s$", pattern)
	}
	r := regexp.MustCompile(pattern)
	res := make([]string, 0, len(lines))
	for i, str := range lines {
		// если нужно игнорировать регистр, приводим к нижнему
		if ignoreCase {
			str = strings.ToLower(str)
		}
		if !invert {
			if !r.MatchString(str) {
				continue
			}
			res = append(res, addMatch(lines, i, context, before, after, lineNum)...)
		} else {
			if r.MatchString(str) {
				continue
			}
			res = append(res, addMatch(lines, i, context, before, after, lineNum)...)
		}
	}
	return res
}

func addMatch(lines []string, i, context, before, after int, lineNum bool) []string {
	maxLen := 1
	if context > 0 {
		maxLen = context
	}
	if before > 0 {
		maxLen += before
	}
	if after > 0 {
		maxLen += after
	}

	res := make([]string, 0, maxLen)
	if context > 1 {
		// находим границы для контекста в зависимости от i
		before, after = getBeforeAfterByContext(lines, i, context)
		l := appendLines(lines, i, before, after, true, lineNum)
		res = append(res, l...)

	} else if after > 0 || before > 0 {
		l := appendLines(lines, i, before, after, false, lineNum)
		res = append(res, l...)
	} else {
		n := ""
		if lineNum {
			n = strconv.Itoa(i+1) + ": "
		}
		res = append(res, fmt.Sprintf("%s%s", n, lines[i]))
	}

	return res
}

func getBeforeAfterByContext(lines []string, i, context int) (int, int) {
	contextD2 := context / 2

	before := 0
	if i-contextD2 >= 0 {
		before = contextD2
	} else {
		before = i
	}

	after := 0
	if i+contextD2 < len(lines) {
		after = contextD2
	} else {
		after = len(lines) - 1 - i
	}

	if context%2 == 0 {
		if after+before == context {
			before--
			return after, before
		}
	} else {
		if after+before+1 == context {
			return after, before
		}
	}

	if before < contextD2 {
		after = context - before - 1
		return before, after
	}

	before = context - after - 1
	return before, after
}

func appendLines(lines []string, i, before, after int, include, lineNum bool) []string {
	jMin := i - before
	if jMin < 0 {
		jMin = 0
	}

	jMax := i + after
	if jMax > len(lines)-1 {
		jMax = len(lines) - 1
	}

	res := make([]string, 0, before+after+3)
	if before > 0 {
		for j := jMin; j < i; j++ {
			n := ""
			if lineNum {
				n = strconv.Itoa(j+1) + ": "
			}
			res = append(res, fmt.Sprintf("%s%s", n, lines[j]))
		}
	}
	if include {
		n := ""
		if lineNum {
			n = strconv.Itoa(i+1) + ": "
		}
		res = append(res, fmt.Sprintf("%s%s", n, lines[i]))
	}
	if after > 0 {
		for j := i + 1; j < jMax+1; j++ {
			n := ""
			if lineNum {
				n = strconv.Itoa(j+1) + ": "
			}
			res = append(res, fmt.Sprintf("%s%s", n, lines[j]))
		}
	}
	return res
}
