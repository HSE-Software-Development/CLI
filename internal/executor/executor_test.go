package executor

import (
	"CLI/internal/environment"
	"CLI/internal/parseline"
	"os"
	"testing"
)

func TestPipeline_EchoToWc(t *testing.T) {
	env := environment.New()
	executor := New(env)
	commands := []parseline.Command{
		{Name: "echo", Args: []string{"hello\nworld\n"}},
		{Name: "wc", Args: []string{}},
	}

	result, err := executor.Execute(commands)
	if err != nil {
		t.Fatalf("ExecutePipeline failed: %v", err)
	}

	expected := "2 2 12"
	if result.String() != expected {
		t.Errorf("Expected %q, got %q", expected, result.String())
	}
}

func TestPipeline_CatToWc(t *testing.T) {
	env := environment.New()
	executor := New(env)
	executor.cwd = t.TempDir()
	tmpfile, err := os.CreateTemp(executor.cwd, "testfile")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())
	tmpfile.WriteString("hello\nworld\n")

	commands := []parseline.Command{
		{Name: "cat", Args: []string{tmpfile.Name()}},
		{Name: "wc", Args: []string{}},
	}

	result, err := executor.Execute(commands)
	if err != nil {
		t.Fatalf("ExecutePipeline failed: %v", err)
	}

	expected := "2 2 12"
	if result.String() != expected {
		t.Errorf("Expected %q, got %q", expected, result.String())
	}
}

func TestPipeline_EchoCatWc(t *testing.T) {
	env := environment.New()
	executor := New(env)
	commands := []parseline.Command{
		{Name: "echo", Args: []string{"hello\nworld\n"}},
		{Name: "cat", Args: []string{}},
		{Name: "wc", Args: []string{}},
	}

	result, err := executor.Execute(commands)
	if err != nil {
		t.Fatalf("ExecutePipeline failed: %v", err)
	}

	expected := "2 2 12"
	if result.String() != expected {
		t.Errorf("Expected %q, got %q", expected, result.String())
	}
}
