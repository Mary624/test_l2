package grep

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGrep(t *testing.T) {
	type param struct {
		linesRes   []string
		lines      []string
		pattern    string
		before     int
		after      int
		context    int
		ignoreCase bool
		invert     bool
		fixed      bool
		lineNum    bool
	}
	l := []param{
		{
			lines:      []string{"gg", "3", "dsfs", "", "asffsa", "5", "afafg"},
			linesRes:   []string{"4: "},
			pattern:    ".+",
			before:     -1,
			after:      -1,
			context:    -1,
			ignoreCase: true,
			invert:     true,
			fixed:      true,
			lineNum:    true,
		},
		{
			lines:      []string{"gg", "3", "dsfs", "", "asffsa", "5", "afafg"},
			linesRes:   []string{"dsfs", "", "asffsa"},
			pattern:    ".+",
			before:     -1,
			after:      -1,
			context:    3,
			ignoreCase: true,
			invert:     true,
			fixed:      true,
			lineNum:    false,
		},
		{
			lines:      []string{"gg", "3", "dsfs", "", "asffsa", "5", "afafg"},
			linesRes:   []string{"gg", "3", "dsfs"},
			pattern:    "gg",
			before:     -1,
			after:      -1,
			context:    3,
			ignoreCase: true,
			invert:     false,
			fixed:      true,
			lineNum:    false,
		},
		{
			lines:      []string{"gg", "3", "dsfs", "", "asffsa", "5", "afafg"},
			linesRes:   []string{"", "asffsa", "5", "afafg"},
			pattern:    "5",
			before:     -1,
			after:      -1,
			context:    4,
			ignoreCase: true,
			invert:     false,
			fixed:      true,
			lineNum:    false,
		},
		{
			lines:      []string{"gg", "3", "dsfs", "", "asffsa", "5", "afafg"},
			linesRes:   []string{"dsfs", "asffsa"},
			pattern:    ".+",
			before:     1,
			after:      1,
			context:    -1,
			ignoreCase: true,
			invert:     true,
			fixed:      true,
			lineNum:    false,
		},
		{
			lines:      []string{"gg", "3", "dsfs", "", "asffsa", "5", "afafg"},
			linesRes:   []string{"gg", "3", "dsfs", "asffsa", "5", "afafg"},
			pattern:    ".+",
			before:     -1,
			after:      -1,
			context:    -1,
			ignoreCase: true,
			invert:     false,
			fixed:      true,
			lineNum:    false,
		},
	}
	for _, v := range l {
		res := grepLines(v.lines, v.pattern, v.before, v.after, v.context, v.ignoreCase, v.invert, v.fixed, v.lineNum)
		assert.Equal(t, res, v.linesRes)
	}
}
