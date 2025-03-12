package main

import (
	"CLI/pkg/executor"
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

func main() {
	sigChan := make(chan os.Signal, 1)
    signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

    go func() {
        sig := <-sigChan 
        fmt.Println("\nexit:", sig)
        os.Exit(0) 
    }()

	reader := bufio.NewReader(os.Stdin)
	exec := executor.NewExecutor()

	for {
        fmt.Print("\n>>> ")
        input, _ := reader.ReadString('\n')

		input = input[: len(input) - 1]
		words := strings.SplitN(input, " ", 2)
		var b *bytes.Buffer
		
		if len(words) > 1 {
			b = bytes.NewBuffer([]byte(words[1]))
		} else {
			b = bytes.NewBuffer(make([]byte, 0))
		}
        res, err := exec.Execute(words[0], b)
		if err == nil {
			fmt.Print(res.String())
		} else {
			fmt.Print(err)
		}
	}


}