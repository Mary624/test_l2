package executor

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"
	"unix-shell-utility/pkg/cmds/cd"
	"unix-shell-utility/pkg/cmds/echo"
	"unix-shell-utility/pkg/cmds/forkexec"
	"unix-shell-utility/pkg/cmds/kill"
	"unix-shell-utility/pkg/cmds/ps"
	"unix-shell-utility/pkg/cmds/pwd"
)

var (
	ErrWrongNameOfCmd = errors.New("wrong cmd")
)

type Cmd interface {
	Do(...string) (string, bool)
}

type Cd interface {
	WorkDir() string
}

type Executor struct {
	Cmds map[string]Cmd
	Cd   Cd
}

func New() (*Executor, error) {
	c1, err := cd.New()
	if err != nil {
		return nil, err
	}

	c2 := pwd.New(c1)
	c3 := echo.New()
	c4 := ps.New()
	c5 := kill.New()
	c6 := forkexec.New(c1)

	cmds := make(map[string]Cmd, 6)
	cmds["cd"] = c1
	cmds["pwd"] = c2
	cmds["echo"] = c3
	cmds["ps"] = c4
	cmds["kill"] = c5
	cmds["exec"] = c6

	return &Executor{
		Cmds: cmds,
		Cd:   c1,
	}, nil
}

func (e *Executor) Run() {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Printf("PS %s>", e.Cd.WorkDir())
	for scanner.Scan() {
		argsStr := scanner.Text()
		if argsStr == "" {
			fmt.Println()
			continue
		}
		if argsStr == `\quit` {
			break
		}
		var out string
		// если комманд несколько, выполняем их все и передаем аргументы дальше
		if strings.Contains(argsStr, "|") {
			out = e.executePipe(argsStr)
		} else {
			out, _ = e.execute(e.getArgs(argsStr))
		}

		if out != "" {
			fmt.Println(out)
		}
		fmt.Printf("PS %s>", e.Cd.WorkDir())
	}
}

// получаем аргументы из строки
func (e *Executor) getArgs(argsStr string) []string {
	if strings.Contains(argsStr, "\n") {
		return strings.Split(argsStr, "\n")
	}
	res := strings.Split(argsStr, " ")
	args := make([]string, 0, len(res))
	for _, arg := range res {
		if arg != "" && arg != "|" {
			args = append(args, arg)
		}
	}
	return args
}

func (e *Executor) executePipe(argsStr string) string {
	cmds := strings.Split(argsStr, "|")
	var out string
	var isErr bool
	var newArgs []string
	for i, cmd := range cmds {
		args := e.getArgs(cmd)
		if i != 0 {
			args = append(args[:1], newArgs...)
		}
		out, isErr = e.execute(args)
		if isErr {
			break
		}
		newArgs = e.getArgs(out)
	}
	return out
}

// выполняем
func (e *Executor) execute(args []string) (string, bool) {
	c, ok := e.Cmds[args[0]]
	if !ok {
		fmt.Printf("%s\n", ErrWrongNameOfCmd.Error())
		return "", true
	}
	out, isErr := c.Do(args[1:]...)
	return out, isErr
}
