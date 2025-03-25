package parseline

import (
	"CLI/internal/environment"
	"bytes"
	"runtime"
	"strings"
)

// TODO. Parser WILL store flags for parsing.
type Parser struct {
	env environment.Env
}


func New(env environment.Env) *Parser {
	return &Parser{
		env: env,
	}
}

// Parse: parses the received string into a command and buffer
// Parameters:
// - input: string 
// Returns:
// - cmd: name of command.
// - buffer: args.
func (p * Parser) Parse(input string) (string, *bytes.Buffer) {
	if runtime.GOOS == "windows" {
        input = input[: len(input) - 2]
    } else {
		input = input[: len(input) - 1]
	}
	words := strings.SplitN(input, " ", 2)
	
	var b *bytes.Buffer
	
	if len(words) > 1 {
		b = bytes.NewBuffer([]byte(words[1]))
	} else {
		b = bytes.NewBuffer(make([]byte, 0))
	}
	return words[0], b
}