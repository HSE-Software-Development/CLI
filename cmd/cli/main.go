package main

import (
	"CLI/internal/handler"
	"fmt"
	"os"
	"os/signal"
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
	handler := handler.InputHandler{}
	handler.Start()
}