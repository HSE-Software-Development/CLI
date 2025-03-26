package parseline

import (
	"CLI/internal/environment"
	"fmt"
	"strings"
	"errors"
	"unicode"
)

// TODO. Parser WILL store flags for parsing.
type Parser struct {
	env environment.Env
}

// Constructor of parser            
// Parameters: env environment.Env
func New(env environment.Env) *Parser {
	return &Parser{
		env: env,
	}
}

type Command struct {
	Name string   
	Args []string 
}

// ParsePipeline: parses the received string into pipeline
// Parameters: 
// - input: string 
// Returns: 
// - []Command: pipeline
// - error:
func (parser * Parser)ParsePipeline(input string) ([]Command, error) {
	var commands []Command
	var currentCmd strings.Builder
	var inSingleQuote, inDoubleQuote, escaped bool
	var currentArg strings.Builder
	var args []string
	expectingCommand := true

	input, err := parser.substitution(input)
	if err != nil {
		return nil, err
	}

	for _, r := range input {
		switch {
		case escaped:
			currentArg.WriteRune(r)
			escaped = false
		case r == '\\':
			escaped = true
			currentArg.WriteRune(r)
		case r == '\'' && !inDoubleQuote:
			inSingleQuote = !inSingleQuote
			currentArg.WriteRune(r)
		case r == '"' && !inSingleQuote:
			inDoubleQuote = !inDoubleQuote
			currentArg.WriteRune(r)
		case r == '|' && !inSingleQuote && !inDoubleQuote:
			if currentArg.Len() > 0 {
				args = append(args, currentArg.String())
				currentArg.Reset()
			}
			if len(args) > 0 || currentCmd.Len() > 0 {
				if currentCmd.Len() == 0 && len(args) > 0 {
					currentCmd.WriteString(args[0])
					args = args[1:]
				}
				commands = append(commands, Command{
					Name: currentCmd.String(),
					Args: args,
				})
				currentCmd.Reset()
				args = nil
				expectingCommand = true
			}
		case unicode.IsSpace(r) && !inSingleQuote && !inDoubleQuote:
			// Конец аргумента (если не в кавычках)
			if currentArg.Len() > 0 {
				args = append(args, currentArg.String())
				currentArg.Reset()
			}
		default:
			currentArg.WriteRune(r)
			if expectingCommand && currentCmd.Len() == 0 && currentArg.Len() > 0 {
				expectingCommand = false
			}
		}
	}

	if currentArg.Len() > 0 {
		args = append(args, currentArg.String())
	}

	if currentCmd.Len() == 0 && len(args) > 0 {
		currentCmd.WriteString(args[0])
		args = args[1:]
	}

	if currentCmd.Len() > 0 || len(args) > 0 {
		commands = append(commands, Command{
			Name: currentCmd.String(),
			Args: args,
		})
	}

	if inSingleQuote || inDoubleQuote {
		return nil, errors.New("unclosed quotes in input")
	}

	if escaped {
		return nil, errors.New("unfinished escape sequence")
	}

	return commands, nil
}

func (parser *Parser) substitution(s string) (string, error) {
	var result strings.Builder
	i := 0
	n := len(s)

	for i < n {
		if s[i] == '$' {
			i++
			if i < n && s[i] == '{' {
				i++ 
				varNameStart := i

				for i < n && s[i] != '}' {
					i++
				}

				if i >= n {
					return "", errors.New("unclosed ${...} variable")
				}

				varName := s[varNameStart:i]
				i++ 

				value, err := parser.env.Get(varName)
				if err != nil {
					return "", fmt.Errorf("error getting variable %s: %w", varName, err)
				}
				result.WriteString(value)
			} else {
				varNameStart := i

				for i < n && (isAlphaNum(s[i]) || s[i] == '_') {
					i++
				}

				varName := s[varNameStart:i]

				if varName == "" {
					result.WriteByte('$')
				} else {

					value, err := parser.env.Get(varName)
					if err != nil {
						return "", fmt.Errorf("error getting variable %s: %w", varName, err)
					}
					
					result.WriteString(value)
				}
			}
		} else {
			result.WriteByte(s[i])
			i++
		}
	}

	return result.String(), nil
}

func isAlphaNum(c byte) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9')
}