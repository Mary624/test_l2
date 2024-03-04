package echo

import "strings"

type Echo struct {
}

func New() *Echo {
	return &Echo{}
}

func (c *Echo) Do(args ...string) (string, bool) {
	var b strings.Builder

	for i, arg := range args {
		b.WriteString(arg)
		if i != len(args)-1 {
			b.WriteRune('\n')
		}
	}
	return b.String(), false
}
