package cut

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCut(t *testing.T) {
	type res struct {
		data      string
		fields    []int
		dilimiter string
		separated bool
		res       string
	}
	l := []res{
		{"1:2 3 4 5 6 7 8\n1:2 3 4 5 6 7 8",
			[]int{2}, ":", true,
			"2 3 4 5 6 7 8\n2 3 4 5 6 7 8"},
		{"1:2 3 4 5 6 7 8\n1:2 3 4 5 6 7 8\n1 2 3 4 5 6 7 8",
			[]int{1}, ":", true,
			"1\n1"},
		{"1 2 3 4 5 6 7 8\n1 2 3 4 5 6 7 8\n1 2 3 4 5 6 7 8",
			[]int{4, 1}, " ", true,
			"4 1\n4 1\n4 1"},
	}
	for _, v := range l {
		var lines = strings.Split(v.data, "\n")

		res := cutLines(lines, v.dilimiter, v.separated, v.fields)
		assert.Equal(t, res, v.res)
	}
}
