package executor

import (
	"CLI/internal/environment"
	"CLI/internal/parseline"
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

// Executor stores a self-implemented functions and a reference to an object storing environment variables
type Executor struct {
	cmds commands
	env environment.Env
}

// Constructor: creates a new Executor and initializes commands
// Parameters:
// - env: environment.Env
func New(env environment.Env) *Executor {
	return &Executor{
		env: env,
		cmds: newCommands(),
	}
}

func (executor *Executor) execute(command parseline.Command, b *bytes.Buffer) (*bytes.Buffer, error) {
	if cmd, ok := executor.cmds[command.Name]; ok {
		return cmd(command, b)
	} else if strings.ContainsRune(command.Name, '=' ){
		split := strings.Split(command.Name, "=")
		if len(split) != 2 {
			return nil, fmt.Errorf("command %s: '=' is incorrect symbol for variable or value", command.Name)
		} 
		if len(split[0]) == 0 || len(split[1]) == 0 {
			return nil, fmt.Errorf("command %s: incorrect lenght of variable or value", command.Name)
		}   
		executor.env.Set(split[0], split[1])
		return bytes.NewBufferString(""), nil

	} else {
		var res *exec.Cmd
		content := b.String() 
		if len(content) > 0 {
			res = exec.Command(command.Name, append(command.Args, content)...)
		} else {
			res = exec.Command(command.Name, command.Args...)
		}
		output, err := res.Output()
		if err != nil {
			return nil, fmt.Errorf("command %s: %s", command.Name, err.Error())
		}
		return bytes.NewBufferString(string(output)), nil
	}
}

// Execute: execute commands, and returns resulting buffer.
// Parameters:
// - commands: []parseline.Command
// Returns:
// - buffer: resulting buffer.
// - err: error of execute.
func (executor *Executor) Execute(commands []parseline.Command) (*bytes.Buffer, error) {
	var err error
	buffer := bytes.NewBufferString("")
	
	for _, cmd := range commands {
		buffer, err = executor.execute(cmd, buffer)
		if err != nil {
			return nil, err
		}
	}
	return buffer, nil
} 