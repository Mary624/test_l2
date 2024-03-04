package ps

import (
	"fmt"
	"strings"

	ps "github.com/mitchellh/go-ps"
)

type Ps struct {
}

func New() *Ps {
	return &Ps{}
}

func (c *Ps) Do(args ...string) (string, bool) {
	procs, err := ps.Processes()
	if err != nil {
		return err.Error(), true
	}
	var b strings.Builder
	for i, proc := range procs {
		b.WriteString(fmt.Sprintf("Pid: %d\tPPid: %d\tExecutable: %s", proc.Pid(), proc.PPid(), proc.Executable()))
		if i != len(procs)-1 {
			b.WriteRune('\n')
		}
	}
	return b.String(), false
}
