package handler

import (
	"CLI/internal/environment"
	"CLI/internal/executor"
	"CLI/internal/parseline"
	"bufio"
	"fmt"
	"os"
)

// TODO. InputHandler WILL store flags.
type InputHandler struct {
}



// Start: starts Read-Execute-Print Loop
func (handler *InputHandler) Start() {
	reader := bufio.NewReader(os.Stdin)
	env := environment.New()
	exec := executor.New(env)
	parser := parseline.New(env)
	for {
        fmt.Print("\n>>> ")
        input, _ := reader.ReadString('\n')
		cmd, b := parser.Parse(input)
        res, err := exec.Execute(cmd, b)
		if err == nil {
			fmt.Print(res.String())
		} else {
			fmt.Print(err)
		}
	}
}