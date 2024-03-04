package forkexec

import (
	"os"
	"unix-shell-utility/pkg/cmds"
	"unix-shell-utility/pkg/cmds/cd"
)

type Cd interface {
	WorkDir() string
}

type ForkExec struct {
	cd *cd.Cd
}

func New(cd *cd.Cd) *ForkExec {
	return &ForkExec{
		cd: cd,
	}
}

func (c *ForkExec) Do(args ...string) (string, bool) {
	var err error
	var procAttr os.ProcAttr
	procAttr.Files = []*os.File{os.Stdin,
		os.Stdout, os.Stderr}
	procAttr.Dir = c.cd.WorkDir()
	procAttr.Env = os.Environ()
	if len(args) == 1 {
		_, err = os.StartProcess(args[0], nil, &procAttr)
	} else if len(args) > 1 {
		_, err = os.StartProcess(args[0], args[1:], &procAttr)
	} else {
		return cmds.ErrWrongArgs.Error(), true
	}
	if err != nil {
		return err.Error(), true
	}

	return "", false
}
