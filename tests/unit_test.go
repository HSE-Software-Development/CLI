package unit_tests

import (
	"CLI/internal/executor"
	"testing"
	"bytes"
	"os"
)

func TestPwdCommand(t *testing.T) {
	executor := executor.NewExecutor()
	input := bytes.NewBufferString("")
	output, err := executor.Execute("pwd", input)
	
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	expected, _ := os.Getwd()
	if output.String() != expected {
		t.Errorf("Expected '%s', got '%s'", expected, output.String())
	}
}

func TestWcCommand(t *testing.T) {
    executor := executor.NewExecutor()
    input := bytes.NewBufferString("")
    output, err := executor.Execute("wc", input)
    if err != nil {
        t.Fatalf("Expected no error, got %v", err)
    }
    expected := "0 0 0"
    if output.String() != expected {
        t.Errorf("Expected '%s', got '%s'", expected, output.String())
    }

    input = bytes.NewBufferString("Hello world\nThis is a test")
    output, err = executor.Execute("wc", input)
    if err != nil {
        t.Fatalf("Expected no error, got %v", err)
    }
    expected = "2 6 26" 
    if output.String() != expected {
        t.Errorf("Expected '%s', got '%s'", expected, output.String())
    }
}

func TestExitCommand(t *testing.T) {
    executor := executor.NewExecutor()

    // Сценарий: Вызов exit
    input := bytes.NewBufferString("")
    defer func() {
        if r := recover(); r == nil {
            t.Errorf("Expected program to exit, but it didn't")
        }
    }()
    executor.Execute("exit", input)
}

func TestEchoCommand(t *testing.T) {
    executor := executor.NewExecutor()

    input := bytes.NewBufferString("Hello, echo!")
    output, err := executor.Execute("echo", input)
    if err != nil {
        t.Fatalf("Expected no error, got %v", err)
    }
    expected := "Hello, echo!"
    if output.String() != expected {
        t.Errorf("Expected '%s', got '%s'", expected, output.String())
    }
}

func TestEchoCommand2(t *testing.T) {
    executor := executor.NewExecutor()

    input := bytes.NewBufferString("'Hello, \necho!'")
    output, err := executor.Execute("echo", input)
    if err != nil {
        t.Fatalf("Expected no error, got %v", err)
    }
    expected := "Hello, \necho!"
    if output.String() != expected {
        t.Errorf("Expected '%s', got '%s'", expected, output.String())
    }
}