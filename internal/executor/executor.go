package executor

import (
	"bytes"
	"os/exec"
)
// Executor stores a self-implemented functions. 
type Executor struct {
	cmds commands
}

// NewExecutor: create a new Executor
func NewExecutor() *Executor {
	return &Executor{
		cmds: newCommands(),
	}
}


// Execute: gets command and buffer, and returns resulting buffer.
// Parameters:
// - command: string
// - b: buffer with args 
// Returns:
// - buffer: resulting buffer.
// - err: error of execute.
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