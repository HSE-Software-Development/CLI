package executor

import (
	"CLI/internal/parseline"
	"bytes"
	"fmt"
	"os"
	"strings"
	"errors"
	"strconv"
	"regexp"

	"github.com/spf13/pflag"
)

type commands map[string]func(parseline.Command, *bytes.Buffer) (*bytes.Buffer, error)


func newCommands() commands {
	cmds := make(commands)

	// Here you can add new command in CLI
	// cmd["name_command"] = name_command
	// below you need to implement a command with the following signature:
	// func(parseline.Command, *bytes.Buffer) (*bytes.Buffer, error)
	cmds["cat"] = cat
	cmds["echo"] = echo
	cmds["exit"] = exit
	cmds["pwd"] = pwd
	cmds["wc"] = wc
	cmds["grep"] = grep


	return cmds
}

//Here you need implement a command with the following signature:
// func name_command(parseline.Command, *bytes.Buffer) (*bytes.Buffer, error) {}

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
	var input string
	if len(cmd.Args) == 0 {
		input = b.String()
	} else if len(cmd.Args) > 0 {
		data, err := os.ReadFile(cmd.Args[0])
		if err != nil {
			return nil, fmt.Errorf("wc: %w", err)
		}
		input = string(data)
	} else {
		return nil, errors.New("wc: no input provided")
	}
	lines := strings.Count(input, "\n")
	words := len(strings.Fields(input))
	chars := len(input)

	result := fmt.Sprintf("%d %d %d", lines, words, chars)
	return bytes.NewBufferString(result), nil
}

func grep(cmd parseline.Command, input *bytes.Buffer) (*bytes.Buffer, error) {
	var (
		caseInsensitive bool
		wordRegexp      bool
		afterContext    int
	)
	

	flagSet := pflag.NewFlagSet("grep", pflag.ContinueOnError)
	flagSet.BoolVarP(&caseInsensitive, "ignore-case", "i", false, "Case-insensitive search")
	flagSet.BoolVarP(&wordRegexp, "word-regexp", "w", false, "Match whole word")
	flagSet.IntVarP(&afterContext, "after-context", "A", 0, "Number of trailing context lines to print")

	if err := flagSet.Parse(cmd.Args); err != nil {
		return nil, err
	}

	args := flagSet.Args()
	if len(args) < 1 {
		return nil, errors.New("pattern is required")
	}
	pattern := args[0]

	var reBuilder strings.Builder
	if caseInsensitive {
		reBuilder.WriteString("(?i)")
	}
	if wordRegexp {
		reBuilder.WriteString(`\b`)
	}
	reBuilder.WriteString(pattern)
	if wordRegexp {
		reBuilder.WriteString(`\b`)
	}

	re, err := regexp.Compile(reBuilder.String())
	if err != nil {
		return nil, fmt.Errorf("invalid regex pattern: %v", err)
	}
	
	inputData := input.String()
	if inputData == "" {
		data, err := os.ReadFile(cmd.Args[len(cmd.Args) - 1])
		if err != nil {
			inputData = ""
		}
		inputData = string(data)
	}
	lines := strings.Split(inputData, "\n")
	printed := make(map[int]struct{})

	for i, line := range lines {
		if re.MatchString(line) {
			end := i + afterContext
			if end >= len(lines) {
				end = len(lines) - 1
			}
			for j := i; j <= end; j++ {
				printed[j] = struct{}{}
			}
		}
	}

	var resultBuffer bytes.Buffer
	for i := 0; i < len(lines); i++ {
		if _, ok := printed[i]; ok {
			resultBuffer.WriteString(lines[i])
			resultBuffer.WriteByte('\n')
		}
	}

	return &resultBuffer, nil
}