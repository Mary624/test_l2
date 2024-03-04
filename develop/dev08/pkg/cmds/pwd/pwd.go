package pwd

import (
	"fmt"
	"unix-shell-utility/pkg/cmds"
	"unix-shell-utility/pkg/cmds/cd"
)

type Pwd struct {
	cd *cd.Cd
}

func New(cd *cd.Cd) *Pwd {
	return &Pwd{
		cd: cd,
	}
}

func (c *Pwd) Do(args ...string) (string, bool) {
	if len(args) != 0 {
		return cmds.ErrWrongArgs.Error(), false
	}

	return fmt.Sprintf("Path\n----\n%s", c.cd.WorkDir()), false
}
