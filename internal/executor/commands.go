package executor

import (
	"CLI/internal/parseline"
	"bytes"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/spf13/pflag"
)

type commands map[string]func(parseline.Command, *bytes.Buffer, *Executor) (*bytes.Buffer, error)

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
	cmds["cd"] = cd	cmds["ls"] = ls


	return cmds
}

//Here you need implement a command with the following signature:
// func name_command(parseline.Command, *bytes.Buffer) (*bytes.Buffer, error) {}

func cat(cmd parseline.Command, b *bytes.Buffer, executor *Executor) (*bytes.Buffer, error) {
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

func echo(cmd parseline.Command, b *bytes.Buffer, executor *Executor) (*bytes.Buffer, error) {
	if len(cmd.Args) == 0 {
		return bytes.NewBufferString(""), nil
	}

	content := strings.Join(cmd.Args, " ")
	if content[0] == '"' {
		content = content[1 : len(content)-1]
	}
	if content[0] == '\'' {
		content = content[1 : len(content)-1]
		content = strings.ReplaceAll(content, "\\n", "\n")
	}
	b.Reset()
	b.WriteString(content)

	return b, nil
}

func exit(cmd parseline.Command, b *bytes.Buffer, executor *Executor) (*bytes.Buffer, error) {
	code := 0
	if len(cmd.Args) > 0 {
		if c, err := strconv.Atoi(cmd.Args[0]); err == nil {
			code = c
		}
	}
	os.Exit(code)
	return nil, nil
}

func pwd(cmd parseline.Command, _ *bytes.Buffer, executor *Executor) (*bytes.Buffer, error) {
	dir, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("pwd: %w", err)
	}
	output := bytes.NewBufferString(dir)
	return output, nil
}

func wc(cmd parseline.Command, b *bytes.Buffer, executor *Executor) (*bytes.Buffer, error) {
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

func grep(cmd parseline.Command, input *bytes.Buffer, executor *Executor) (*bytes.Buffer, error) {
	opts, err := parseArgs(cmd.Args)
	if err != nil {
		return nil, err
	}

	re, err := compileRegex(opts)
	if err != nil {
		return nil, fmt.Errorf("regex error: %v", err)
	}

	data, err := readData(opts, input)
	if err != nil {
		return nil, err
	}

	result := processData(data, re, opts.afterContext)
	return result, nil
}

func parseArgs(args []string) (*grepOptions, error) {
	opts := &grepOptions{}
	fs := pflag.NewFlagSet("grep", pflag.ContinueOnError)

	fs.IntVarP(&opts.afterContext, "after-context", "A", 0, "Lines after match")
	fs.BoolVarP(&opts.ignoreCase, "ignore-case", "i", false, "Ignore case")
	fs.BoolVarP(&opts.wordRegexp, "word-regexp", "w", false, "Whole word match")

	if err := fs.Parse(args); err != nil {
		return nil, err
	}

	remaining := fs.Args()
	if len(remaining) < 1 {
		return nil, fmt.Errorf("pattern required")
	}

	opts.pattern = strings.Trim(remaining[0], `"'`)
	if len(remaining) > 1 {
		opts.files = remaining[1:]
	}

	return opts, nil
}

func compileRegex(opts *grepOptions) (*regexp.Regexp, error) {
	pattern := opts.pattern

	if opts.wordRegexp {
		pattern = fmt.Sprintf(`\b%s\b`, pattern)
	}

	if opts.ignoreCase {
		pattern = "(?i)" + pattern
	}

	return regexp.Compile(pattern)
}

func readData(opts *grepOptions, input *bytes.Buffer) ([]byte, error) {
	if len(opts.files) > 0 {
		var data []byte
		for _, f := range opts.files {
			fileData, err := os.ReadFile(f)
			if err != nil {
				return nil, fmt.Errorf("error reading %s: %v", f, err)
			}
			data = append(data, fileData...)
		}
		return data, nil
	}
	return input.Bytes(), nil
}

func processData(data []byte, re *regexp.Regexp, context int) *bytes.Buffer {
	lines := bytes.Split(data, []byte{'\n'})
	var output bytes.Buffer
	remaining := 0

	for _, line := range lines {
		if remaining > 0 {
			output.Write(line)
			output.WriteByte('\n')
			remaining--
		}

		if re.Match(line) {
			output.Write(line)
			output.WriteByte('\n')
			remaining = context
		}
	}

	return &output
}

func cd(cmd parseline.Command, b *bytes.Buffer, executor *Executor) (*bytes.Buffer, error) {
	var path string
	if len(cmd.Args) == 0 {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return nil, fmt.Errorf("cd: cannot get home directory: %w", err)
		}
		path = homeDir
	} else {
		path = cmd.Args[0]
	}

	absPath, err := filepath.Abs(filepath.Join(executor.cwd, path))
	if err != nil {
		return nil, fmt.Errorf("cd: cannot resolve path: %w", err)
	}

	info, err := os.Stat(absPath)
	if err != nil || !info.IsDir() {
		return nil, fmt.Errorf("cd: %s: not a directory", path)
	}

	executor.cwd = absPath
	return b, nil
}


// Takes first provided argument as path to a directory and lists all files in it.
func ls(cmd parseline.Command, buffer *bytes.Buffer) (*bytes.Buffer, error) {
	var path string
	if (len(cmd.Args) == 0) {
		wd, err := os.Getwd()
		if err != nil {
			return nil, fmt.Errorf("pwd: %w", err)
		}
		path = wd
	} else {
		path = cmd.Args[0]
	}
	entries, err := os.ReadDir(regexp.QuoteMeta(path))
	if err != nil {
		return nil, fmt.Errorf("ls: %w", err)
	}

	for _, entry := range entries {
		buffer.WriteString(entry.Name() + "\n")
	}

	return buffer, nil
}
