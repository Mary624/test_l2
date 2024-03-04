package four

import (
	"sort"
	"strings"
)

func AnagramSearch(input []string) map[string][]string {
	// результат
	res := make(map[string][]string, len(input))
	// ключ - отсортированное слово, а значение - первое слово (для поиска по мапе результатов)
	search := make(map[string]string, len(input))
	for _, str := range input {
		if str == "" {
			continue
		}
		str = strings.ToLower(str)
		sorted := string(sortAscStr(str))
		key, ok := search[sorted]
		// если это новое множество
		if !ok {
			res[str] = make([]string, 0, len(input))
			res[str] = append(res[str], str)
			search[sorted] = str
			continue
		}
		// если анаграма, то добавляем
		if isAnagram(str, res[key][0]) {
			res[key] = append(res[key], str)
		}
	}
	// убираем повторы
	for k, v := range res {
		if len(v) == 1 {
			delete(res, k)
			continue
		}
		// сортируем
		res[k] = sortStringSet(v)
	}
	return res
}

func sortAscStr(str string) []rune {
	x1R := []rune(str)
	sort.Slice(x1R, func(i, j int) bool {
		return x1R[i] < x1R[j]
	})
	return x1R
}

func isAnagram(x1, x2 string) bool {
	if len(x1) != len(x2) {
		return false
	}

	x1R, x2R := sortAscStr(x1), sortAscStr(x2)
	for i := 0; i < len(x1R); i++ {
		if x1R[i] != x2R[i] {
			return false
		}
	}
	return true

}

// func isAnagram1(x1, x2 string) bool {
// 	if len(x1) != len(x2) {
// 		return false
// 	}

// 	m := make(map[rune]bool, len(x1))

// }

func sortStringSet(input []string) []string {
	sort.Slice(input, func(i, j int) bool {
		return input[i] < input[j]
	})
	res := make([]string, 0, len(input))
	prev := ""
	for _, str := range input {
		if str != prev {
			res = append(res, str)
			prev = str
		}
	}
	return res
}
