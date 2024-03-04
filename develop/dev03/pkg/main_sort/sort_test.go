package mainsort

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRemoveRepeats(t *testing.T) {
	lines := []string{"a", "a", "a", "a", "a", "a", "a"}
	assert.Equal(t, "a\n", removeRepeats(lines, len(lines)))
	lines1 := []string{"a", "b", "a", "a", "a", "a", "a"}
	assert.Equal(t, "a\nb\n", removeRepeats(lines1, len(lines1)))
}

func TestSortByCol(t *testing.T) {
	type res struct {
		lines    []string
		linesRes []string
		colNum   int
		desc     bool
		byNumber bool
		float    bool
		month    bool
	}
	l := []res{
		{
			lines:    []string{"a b", "a a"},
			linesRes: []string{"a a", "a b"},
			colNum:   1,
			desc:     false,
			byNumber: false,
			float:    false,
			month:    false,
		},
		{
			lines:    []string{"a b", "a a"},
			linesRes: []string{"a b", "a a"},
			colNum:   1,
			desc:     true,
			byNumber: false,
			float:    false,
			month:    false,
		},
		{
			lines:    []string{"a 2", "a 1"},
			linesRes: []string{"a 1", "a 2"},
			colNum:   1,
			desc:     false,
			byNumber: true,
			float:    false,
			month:    false,
		},
		{
			lines:    []string{"a 2", "a 1", "a a", "a 0", "a 3"},
			linesRes: []string{"a a", "a 0", "a 1", "a 2", "a 3"},
			colNum:   1,
			desc:     false,
			byNumber: true,
			float:    false,
			month:    false,
		},
		{
			lines:    []string{"a 2", "a 1", "a a", "a 0"},
			linesRes: []string{"a a", "a 0", "a 1", "a 2"},
			colNum:   1,
			desc:     false,
			byNumber: true,
			float:    false,
			month:    false,
		},
		{
			lines:    []string{"a 2.2", "a 2", "a a", "a 0"},
			linesRes: []string{"a a", "a 0", "a 2", "a 2.2"},
			colNum:   1,
			desc:     false,
			byNumber: true,
			float:    true,
			month:    false,
		},
		{
			lines:    []string{"a Jan", "a 2", "a Feb", "a Jan"},
			linesRes: []string{"a 2", "a Jan", "a Jan", "a Feb"},
			colNum:   1,
			desc:     false,
			byNumber: false,
			float:    false,
			month:    true,
		},
		{
			lines:    []string{"a January", "a 2", "a February", "a January"},
			linesRes: []string{"a 2", "a January", "a January", "a February"},
			colNum:   1,
			desc:     false,
			byNumber: false,
			float:    false,
			month:    true,
		},
	}
	for _, r := range l {
		sortByCol(r.lines, r.colNum, r.desc, r.byNumber, r.float, r.month)
		assert.Equal(t, r.lines, r.linesRes)
		fmt.Println()
	}
}
