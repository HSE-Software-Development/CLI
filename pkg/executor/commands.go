package executor

import (
	"bytes"
	"os"
	"io"
)

type commands map[string]func(*bytes.Buffer) (*bytes.Buffer, error)


func newCommands() commands {
	cmds := make(commands)
	cmds["cat"] = cat
	cmds["echo"] = echo
	cmds["exit"] = exit
	cmds["pwd"] = pwd
	cmds["wc"] = wc
	return cmds
}

func cat(b *bytes.Buffer) (*bytes.Buffer, error) {
	file, err := os.Open(b.String())
	if err != nil {
		return nil, err
	}
	defer file.Close()
	b.Reset()
	_, err = io.Copy(b, file)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func echo(b *bytes.Buffer) (*bytes.Buffer, error) {
	return b, nil
}

func exit(_ *bytes.Buffer) (*bytes.Buffer, error) {
	os.Exit(0)
	return nil, nil
}

func pwd(b *bytes.Buffer) (*bytes.Buffer, error) {
	return nil, nil
}
func wc(b *bytes.Buffer) (*bytes.Buffer, error) {
	return nil, nil
}