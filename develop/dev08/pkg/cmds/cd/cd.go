package cd

import (
	"errors"
	"os"
	"path/filepath"
	"regexp"
	"unix-shell-utility/pkg/cmds"
)

var (
	ErrIsFile = errors.New("it's file")
)

type Cd struct {
	workdir string
}

func New() (*Cd, error) {
	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	return &Cd{
		workdir: wd,
	}, nil
}

func (c *Cd) Do(args ...string) (string, bool) {
	if len(args) != 1 || args[0] == "" {
		return cmds.ErrWrongArgs.Error(), true
	}
	err := c.changeWd(args[0])
	if err != nil {
		return err.Error(), true
	}
	return "", false
}

func (c *Cd) WorkDir() string {
	return c.workdir
}

func (c *Cd) changeWd(newWd string) error {
	r := regexp.MustCompile(`^[A-Z]:[\/\\]`)
	if r.MatchString(newWd) {
		info, err := os.Stat(newWd)
		if os.IsNotExist(err) || err != nil {
			return err
		}
		if !info.IsDir() {
			return ErrIsFile
		}
		c.workdir = newWd
		return nil
	}
	fullPath := filepath.Join(c.workdir, newWd)
	info, err := os.Stat(fullPath)
	if os.IsNotExist(err) || err != nil {
		return err
	}
	if !info.IsDir() {
		return ErrIsFile
	}
	c.workdir = fullPath
	return nil
}
