package executor

import (
	"bytes"
	"os/exec"
)

type Executor struct {
	cmds commands
}

func NewExecutor() *Executor {
	return &Executor{
		cmds: newCommands(),
	}
}
func (executor *Executor) Execute(command string, b *bytes.Buffer) (*bytes.Buffer, error) {
	if cmd, ok := executor.cmds[command]; ok {
		return cmd(b)
		
	} else {
		var res *exec.Cmd 
		if len(b.String()) > 0 {
			res = exec.Command(command, b.String())
		} else {
			res = exec.Command(command)
		}
		output, err := res.Output()
		if err != nil {
			return nil, err
		}
		b.Reset()
		b.WriteString(string(output))
		return b, nil
	}
}