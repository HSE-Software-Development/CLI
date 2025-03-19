package parser_

import (
	"bytes"
	
	"strings"
)

type Parser struct {

}

func (p * Parser) Parse(input string) (string, *bytes.Buffer) {
	input = input[: len(input) - 1]
	words := strings.SplitN(input, " ", 2)
	
	var b *bytes.Buffer
	
	if len(words) > 1 {
		b = bytes.NewBuffer([]byte(words[1]))
	} else {
		b = bytes.NewBuffer(make([]byte, 0))
	}
	return words[0], b
}