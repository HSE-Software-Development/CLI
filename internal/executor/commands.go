package executor

import (
	"bytes"
	"os"
	"io"
	"fmt"
	"strings"
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
	content := b.String()
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

func exit(_ *bytes.Buffer) (*bytes.Buffer, error) {
	os.Exit(0)
	return nil, nil
}

func pwd(b *bytes.Buffer) (*bytes.Buffer, error) {
	dir, err := os.Getwd()
    if err != nil {
        return nil, err
    }
    b.Reset()
    b.WriteString(dir)
    return b, nil
}
func wc(b *bytes.Buffer) (*bytes.Buffer, error) {
	content := b.String()
	if len(content) == 0 {
		b.WriteString("0 0 0")
		return b, nil
	}
	file, err := os.Open(content)
	if err == nil {
		defer file.Close()
		data, err := io.ReadAll(file)
		if err != nil {
			return nil, err
		}
		content = string(data)
	} else {
		if content[0] == '"' || content[0] == '\'' {
			content = content[1:len(content) - 1]
			content = strings.ReplaceAll(content, "\\n", "\n")
		}
	}
	
    lines := strings.Count(content, "\n") 
    if len(content) > 0 && !strings.HasSuffix(content, "\n") {
        lines++ 
    }
    words := len(strings.Fields(content)) 
    characters := len(content)            
	result := fmt.Sprintf("%d %d %d", lines, words, characters)

    b.Reset()
    b.WriteString(result)
    return b, nil

}