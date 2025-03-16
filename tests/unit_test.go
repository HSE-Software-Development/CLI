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