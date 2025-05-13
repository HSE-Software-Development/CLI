package executor

import (
	"CLI/internal/environment"
	"CLI/internal/parseline"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// Executor stores a self-implemented functions and a reference to an object storing environment variables
type Executor struct {
	cmds commands
	env  environment.Env
	cwd  string
}

// Constructor: creates a new Executor and initializes commands
// Parameters:
// - env: environment.Env
func New(env environment.Env) *Executor {
	startDir, err := os.Getwd()
	if err != nil {
		startDir = "/"
	}
	return &Executor{
		env:  env,
		cmds: newCommands(),
		cwd:  startDir,
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

type grepOptions struct {
	ignoreCase   bool
	wordRegexp   bool
	afterContext int
	pattern      string
	files        []string
}

func (executor *Executor) execute(command parseline.Command, b *bytes.Buffer) (*bytes.Buffer, error) {
	var result *bytes.Buffer
	var res_error error
	if cmd, ok := executor.cmds[command.Name]; ok {
		if command.Name == "ls" || command.Name == "cd" {
			command.Args = append(command.Args, executor.cwd)
		}
		result, res_error = cmd(command, b)
		if command.Name == "cd" {
			executor.env.Set("OLD_PWD", executor.cwd)
			executor.cwd = result.String()
			executor.env.Set("PWD", executor.cwd)
			return &bytes.Buffer{}, res_error
		}
		return result, res_error
	} else if strings.ContainsRune(command.Name, '=') {
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
		res.Dir = executor.cwd
		output, err := res.Output()
		if err != nil {
			return nil, fmt.Errorf("command %s: %s", command.Name, err.Error())
		}
		return bytes.NewBufferString(string(output)), nil
	}
}
