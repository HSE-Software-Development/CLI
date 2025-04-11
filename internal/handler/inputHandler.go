package handler

import (
	"CLI/internal/environment"
	"CLI/internal/executor"
	"CLI/internal/parseline"
	"bufio"
	"fmt"
	"os"
	"runtime"
)

// TODO. InputHandler WILL store flags.
type InputHandler struct {
}

func New() *InputHandler {
	return &InputHandler{}
}



// Start: starts Read-Execute-Print Loop
func (handler *InputHandler) Start() {
	reader := bufio.NewReader(os.Stdin)
	env := environment.New()
	parser := parseline.New(env)
	executor := executor.New(env)

	for {
        fmt.Print("\n>>> ")
        input, _ := reader.ReadString('\n')

		pipeline, err := parser.ParsePipeline(cropLine(input))
		if err != nil {
			fmt.Print(err.Error())
			continue
		}
		result, err := executor.Execute(pipeline)
		if err != nil {
			fmt.Print(err.Error())
			continue
		}
		fmt.Print(result.String())
	}
}

func cropLine(input string) string {
	if runtime.GOOS == "windows" {
        return input[: len(input) - 2]
    } else {
		return input[: len(input) - 1]
	}
}