package handler

import (
	"CLI/internal/executor"
	"CLI/internal/parser_"
	"bufio"
	"fmt"
	"os"
)

type InputHandler struct {
}

func (handler *InputHandler) Start() {
	reader := bufio.NewReader(os.Stdin)
	exec := executor.NewExecutor()
	parser := parser_.Parser{}
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