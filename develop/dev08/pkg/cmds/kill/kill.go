package kill

import (
	"os"
	"strconv"
	"unix-shell-utility/pkg/cmds"
)

type Kill struct {
}

func New() *Kill {
	return &Kill{}
}

func (c *Kill) Do(args ...string) (string, bool) {
	if len(args) != 1 || args[0] == "" {
		return cmds.ErrWrongArgs.Error(), true
	}

	err := c.kill(args[0])
	if err != nil {
		return err.Error(), true
	}
	return "", false
}

func (c *Kill) kill(arg string) error {
	pid, err := strconv.Atoi(arg)
	if err != nil {
		return err
	}

	proc, err := os.FindProcess(pid)
	if err != nil {
		return err
	}
	err = proc.Kill()
	if err != nil {
		return err
	}
	return nil
}
