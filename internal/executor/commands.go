package executor

import (
	"CLI/internal/parseline"
	"bytes"
	"fmt"
	"os"
	"strings"
	"errors"
	"strconv"
)

type commands map[string]func(parseline.Command, *bytes.Buffer) (*bytes.Buffer, error)


func newCommands() commands {
	cmds := make(commands)
	cmds["cat"] = cat
	cmds["echo"] = echo
	cmds["exit"] = exit
	cmds["pwd"] = pwd
	cmds["wc"] = wc
	return cmds
}

func cat(cmd parseline.Command, b *bytes.Buffer) (*bytes.Buffer, error) {
	output := bytes.NewBuffer(nil)
	if len(cmd.Args) == 0 {
		if b != nil {
			_, err := output.Write(b.Bytes())
			return output, err
		}
		return nil, errors.New("no input provided")
	}

	for _, filename := range cmd.Args {
		data, err := os.ReadFile(filename)
		if err != nil {
			return nil, fmt.Errorf("cat: %w", err)
		}
		output.Write(data)
	}

	return output, nil
}

func echo(cmd parseline.Command, b *bytes.Buffer) (*bytes.Buffer, error) {
	if len(cmd.Args) == 0 {
		return bytes.NewBufferString(""), nil
	}

	content := strings.Join(cmd.Args, " ")
	if content[0] == '"' {
		content = content[1:len(content) - 1]
	}
	if content[0] == '\'' {
		content = content[1:len(content) - 1]
		content = strings.ReplaceAll(content, "\\n", "\n")
	}
	b.Reset()
	b.WriteString(content)
	b.WriteByte('\n')
	return b, nil
}

func exit(cmd parseline.Command, b *bytes.Buffer) (*bytes.Buffer, error) {
	code := 0
	if len(cmd.Args) > 0 {
		if c, err := strconv.Atoi(cmd.Args[0]); err == nil {
			code = c
		}
	}
	os.Exit(code)
	return nil, nil 
}

func pwd(cmd parseline.Command, _ *bytes.Buffer) (*bytes.Buffer, error) {
	dir, err := os.Getwd()
    if err != nil {
		return nil, fmt.Errorf("pwd: %w", err)
	}
    output := bytes.NewBufferString(dir)
	return output, nil
}

func wc(cmd parseline.Command, b *bytes.Buffer) (*bytes.Buffer, error) {
	var input []byte
	if b != nil {
		input = b.Bytes()
	} else if len(cmd.Args) > 0 {
		data, err := os.ReadFile(cmd.Args[0])
		if err != nil {
			return nil, fmt.Errorf("wc: %w", err)
		}
		input = data
	} else {
		return nil, errors.New("wc: no input provided")
	}

	lines := bytes.Count(input, []byte{'\n'})
	words := len(bytes.Fields(input))
	chars := len(input)

	result := fmt.Sprintf("%d %d %d", lines, words, chars)
	return bytes.NewBufferString(result), nil
}