package mainsort

import (
	"errors"
	"fmt"
	"os"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"time"

	flag "github.com/spf13/pflag"
)

var (
	ErrIncompatibleFlags = errors.New("incompatible flags")
	ErrNoFileName        = errors.New("can't get filename")
)

func Sort() (bool, error) {
	// получение параметров
	var byNumber bool
	flag.BoolVarP(&byNumber, "n", "n", false, "sort by number")

	var colNum int
	flag.IntVarP(&colNum, "k", "k", -1, "sort by column")

	var desc bool
	flag.BoolVarP(&desc, "r", "r", false, "sort desc")

	var withoutRep bool
	flag.BoolVarP(&withoutRep, "u", "u", false, "without repeats")

	var tailSpace bool
	flag.BoolVarP(&tailSpace, "b", "b", false, "ignore tail spaces")

	var float bool
	flag.BoolVarP(&float, "h", "h", false, "check num suffix")

	var month bool
	flag.BoolVarP(&month, "M", "M", false, "sort by month")

	var sorted bool
	flag.BoolVarP(&sorted, "c", "c", false, "check sorted")

	flag.Parse()
	filename := flag.Arg(0)
	if !strings.Contains(filename, ".") {
		return false, ErrNoFileName
	}

	if month && (byNumber || float) {
		return false, ErrIncompatibleFlags
	}

	// сортируем
	res, err := sortFile(filename, byNumber, desc, withoutRep, tailSpace, float, month, sorted, colNum)
	if err != nil {
		return false, err
	}

	if res == "" {
		return true, nil
	}

	// пишем результат
	return false, writeRes(res, filename)
}

func writeRes(res, filename string) error {
	// получаение пути
	path := ""
	if strings.Contains(filename, "/") {
		paths := strings.Split(filename, "/")
		filename = paths[len(paths)-1]
		var b strings.Builder
		for i := 0; i < len(paths)-1; i++ {
			b.WriteString(paths[i] + "/")
		}
		path = b.String()
	}
	filesPart := strings.Split(filename, ".")
	if len(filesPart) > 2 {
		var b strings.Builder
		f1 := ""
		for i := 0; i < len(filesPart)-1; i++ {
			b.WriteString(filesPart[i] + ".")
		}
		f1 = b.String()
		filesPart = []string{f1, filesPart[len(filesPart)-1]}
	}
	//создание файла
	file, err := os.Create(fmt.Sprintf("%s_sorted.%s", path+filesPart[0], filesPart[1]))
	if err != nil {
		return err
	}

	defer file.Close()

	// запись
	_, err = file.Write([]byte(res))
	return err
}

func sortFile(filename string, byNumber, desc, withoutRep, tailSpace, float, month, sorted bool, colNum int) (string, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return "", err
	}

	// разбиваем на строки
	var lines []string
	if runtime.GOOS == "windows" {
		lines = strings.Split(string(data), "\r\n")
	} else {
		lines = strings.Split(string(data), "\n")
	}

	if len(lines) < 2 {
		return "", nil
	}

	// проверка на сортировку
	if sorted {
		if checkSorted(lines, desc, tailSpace, byNumber, float, month, withoutRep, colNum) {
			return "", nil
		}
	}

	// сортировка
	if colNum >= 0 {
		sortByCol(lines, colNum, desc, byNumber, float, month)
	} else {
		sortByLine(lines, desc, tailSpace, byNumber, float, month)
	}
	if !withoutRep {
		var b strings.Builder
		b.Grow(len(data))
		for _, c := range lines {
			b.WriteString(c + "\n")
		}
		return b.String(), nil
	}

	return removeRepeats(lines, len(data)), nil
}

func removeRepeats(lines []string, size int) string {
	m := make(map[string]bool, 0)
	var b strings.Builder
	b.Grow(size)

	for _, c := range lines {
		_, ok := m[c]
		if ok {
			continue
		}
		m[c] = true
		b.WriteString(c + "\n")
	}
	return b.String()
}

func sortByCol(lines []string, colNum int, desc, byNumber, float, month bool) {
	sort.Slice(lines, func(i, j int) bool {
		line1 := strings.Split(lines[i], " ")
		line2 := strings.Split(lines[j], " ")

		if len(line1)-1 < colNum {
			return !desc
		}
		if len(line2)-1 < colNum {
			return desc
		}

		x1 := line1[colNum]
		x2 := line2[colNum]

		if month {
			return sortByMonth(x1, x2, desc)
		}
		if byNumber {
			return sortByNumber(x1, x2, desc, float)
		}
		if desc {
			return x1 > x2
		}
		return x1 < x2

	})
}

func sortByLine(lines []string, desc, tailSpace, byNumber, float, month bool) {
	sort.Slice(lines, func(i, j int) bool {
		line1 := lines[i]
		line2 := lines[j]

		if tailSpace && []rune(line1)[len(line1)-1] == ' ' {
			line1 = line1[:len(line1)-1]
		}
		if tailSpace && []rune(line2)[len(line2)-1] == ' ' {
			line2 = line2[:len(line2)-1]
		}

		if month {
			return sortByMonth(line1, line2, desc)
		}
		if byNumber {
			return sortByNumber(line1, line2, desc, float)
		}
		if desc {
			return line1 > line2
		}
		return line1 < line2

	})
}

func checkSorted(lines []string, desc, tailSpace, byNumber, float, month, withoutRep bool, colNum int) bool {
	linesCopy := make([]string, len(lines))
	copy(linesCopy, lines)
	if colNum >= 0 {
		sortByCol(linesCopy, colNum, desc, byNumber, float, month)
	} else {
		sortByLine(linesCopy, desc, tailSpace, byNumber, float, month)
	}

	if !withoutRep {
		for i := 0; i < len(lines); i++ {
			if lines[i] != linesCopy[i] {
				return false
			}
		}
		return true
	}

	m := make(map[string]bool, 0)
	linesCopyRes := make([]string, 0, len(lines))
	for _, c := range lines {
		_, ok := m[c]
		if ok {
			continue
		}
		m[c] = true
		linesCopyRes = append(linesCopyRes, c)
	}

	for i := 0; i < len(lines); i++ {
		if lines[i] != linesCopyRes[i] {
			return false
		}
	}
	return true
}

func sortByNumber(x1, x2 string, desc, float bool) bool {
	var x1Num, x2Num float64
	var err error

	if float {
		x1Num, err = strconv.ParseFloat(x1, 64)
	} else {
		var x1NumInt int
		x1NumInt, err = strconv.Atoi(x1)
		x1Num = float64(x1NumInt)
	}
	if err != nil {
		return !desc
	}

	if float {
		x2Num, err = strconv.ParseFloat(x2, 64)
	} else {
		var x2NumInt int
		x2NumInt, err = strconv.Atoi(x2)
		x2Num = float64(x2NumInt)
	}
	if err != nil {
		return desc
	}

	if desc {
		return x1Num > x2Num
	}
	return x1Num < x2Num
}

func sortByMonth(x1, x2 string, desc bool) bool {
	time1, err := time.Parse("Jan", x1)
	if err != nil {
		time1, err = time.Parse("January", x1)
		if err != nil {
			return !desc
		}
	}

	time2, err := time.Parse("Jan", x2)
	if err != nil {
		time2, err = time.Parse("January", x2)
		if err != nil {
			return desc
		}
	}

	if desc {
		return time1.Month() > time2.Month()
	}
	return time1.Month() < time2.Month()
}
