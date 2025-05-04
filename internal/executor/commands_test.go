package executor

import (
	"CLI/internal/parseline"
	"bytes"
	"os"
	"testing"
)

func TestCat(t *testing.T) {
	tmpfile, err := os.CreateTemp("", "testfile")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())
	tmpfile.WriteString("test data\n")

	tests := []struct {
		name    string
		cmd     parseline.Command
		input   *bytes.Buffer
		want    string
		wantErr bool
	}{
		{
			name:    "Read file",
			cmd:     parseline.Command{Name: "cat", Args: []string{tmpfile.Name()}},
			input:   nil,
			want:    "test data\n",
			wantErr: false,
		},
		{
			name:    "No file",
			cmd:     parseline.Command{Name: "cat", Args: []string{"nonexistent.txt"}},
			input:   nil,
			want:    "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := cat(tt.cmd, tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("cat() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != nil && got.String() != tt.want {
				t.Errorf("cat() = %v, want %v", got.String(), tt.want)
			}
		})
	}
}

func TestEcho(t *testing.T) {
	tests := []struct {
		name    string
		cmd     parseline.Command
		input   *bytes.Buffer
		want    string
		wantErr bool
	}{
		{
			name:    "Simple echo",
			cmd:     parseline.Command{Name: "echo", Args: []string{"hello"}},
			input:   bytes.NewBufferString(""),
			want:    "hello",
			wantErr: false,
		},
		{
			name:    "Multiple args",
			cmd:     parseline.Command{Name: "echo", Args: []string{`"Hello\nWorld"`}},
			input:   bytes.NewBufferString(""),
			want:    "Hello\\nWorld",
			wantErr: false,
		},
		{
			name:    "Multiple args2",
			cmd:     parseline.Command{Name: "echo", Args: []string{`'Hello\nWorld'`}},
			input:   bytes.NewBufferString(""),
			want:    "Hello\nWorld",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := echo(tt.cmd, tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("echo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got.String() != tt.want {
				t.Errorf("echo() = %v, want %v", got.String(), tt.want)
			}
		})
	}
}

func TestPwd(t *testing.T) {
	want, _ := os.Getwd()
	cmd := parseline.Command{Name: "pwd", Args: nil}
	got, err := pwd(cmd, nil)
	if err != nil {
		t.Errorf("pwd() error = %v", err)
	}
	if got.String() != want {
		t.Errorf("pwd() = %v, want %v", got.String(), want)
	}
}

func TestWc(t *testing.T) {
	tests := []struct {
		name    string
		cmd     parseline.Command
		input   *bytes.Buffer
		want    string
		wantErr bool
	}{
		{
			name:    "Count lines/words/chars",
			cmd:     parseline.Command{Name: "wc", Args: nil},
			input:   bytes.NewBufferString("hello world\n"),
			want:    "1 2 12",
			wantErr: false,
		},
		{
			name:    "Read from file",
			cmd:     parseline.Command{Name: "wc", Args: []string{"testfile.txt"}},
			input:   nil,
			want:    "",
			wantErr: true, // Файла нет
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := wc(tt.cmd, tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("wc() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != nil && got.String() != tt.want {
				t.Errorf("wc() = %v, want %v", got.String(), tt.want)
			}
		})
	}
}

var testInput = []byte(`Lorem ipsum dolor sit amet,
consectetur adipiscing elit.
Lorem IPSUM dolor sit amet,
Praesent non Word_WORD word.
Word at start.
end Word
`)

func TestGrepBasic(t *testing.T) {
	cmd := parseline.Command{
		Name: "grep",
		Args: []string{"Lorem"},
	}
	input := bytes.NewBuffer(testInput)
	expected := `Lorem ipsum dolor sit amet,
Lorem IPSUM dolor sit amet,
`

	output, err := grep(cmd, input)
	if err != nil {
		t.Fatalf("Grep failed: %v", err)
	}
	if output.String() != expected {
		t.Errorf("Expected:\n%s\nGot:\n%s", expected, output.String())
	}
}

func TestGrepCaseInsensitive(t *testing.T) {
	cmd := parseline.Command{
		Name: "grep",
		Args: []string{"-i", "ipsum"},
	}
	input := bytes.NewBuffer(testInput)
	expected := `Lorem ipsum dolor sit amet,
Lorem IPSUM dolor sit amet,
`

	output, err := grep(cmd, input)
	if err != nil {
		t.Fatalf("Grep failed: %v", err)
	}
	if output.String() != expected {
		t.Errorf("Expected:\n%s\nGot:\n%s", expected, output.String())
	}
}

func TestGrepWholeWord(t *testing.T) {
	cmd := parseline.Command{
		Name: "grep",
		Args: []string{"-w", "Word"},
	}
	input := bytes.NewBuffer(testInput)
	expected := `Word at start.
end Word
`

	output, err := grep(cmd, input)
	if err != nil {
		t.Fatalf("Grep failed: %v", err)
	}
	if output.String() != expected {
		t.Errorf("Expected:\n%s\nGot:\n%s", expected, output.String())
	}
}

func TestGrepAfterContext(t *testing.T) {
	cmd := parseline.Command{
		Name: "grep",
		Args: []string{"-A", "1", "elit"},
	}
	input := bytes.NewBuffer(testInput)
	expected := `consectetur adipiscing elit.
Lorem IPSUM dolor sit amet,
`

	output, err := grep(cmd, input)
	if err != nil {
		t.Fatalf("Grep failed: %v", err)
	}
	if output.String() != expected {
		t.Errorf("Expected:\n%s\nGot:\n%s", expected, output.String())
	}
}

func TestGrepOverlappingContext(t *testing.T) {
	cmd := parseline.Command{
		Name: "grep",
		Args: []string{"-A", "2", "Lorem"},
	}
	input := bytes.NewBuffer(testInput)
	expected := `Lorem ipsum dolor sit amet,
consectetur adipiscing elit.
Lorem IPSUM dolor sit amet,
Lorem IPSUM dolor sit amet,
Praesent non Word_WORD word.
Word at start.
`

	output, err := grep(cmd, input)
	if err != nil {
		t.Fatalf("Grep failed: %v", err)
	}
	if output.String() != expected {
		t.Errorf("Expected:\n%s\nGot:\n%s", expected, output.String())
	}
}

func TestGrepEmptyInput(t *testing.T) {
	cmd := parseline.Command{
		Name: "grep",
		Args: []string{"pattern"},
	}
	input := bytes.NewBuffer([]byte{})
	expected := ""

	output, err := grep(cmd, input)
	if err != nil {
		t.Fatalf("Grep failed: %v", err)
	}
	if output.String() != expected {
		t.Errorf("Expected empty output, got: %s", output.String())
	}
}

func TestGrepInvalidRegex(t *testing.T) {
	cmd := parseline.Command{
		Name: "grep",
		Args: []string{"*invalid["},
	}
	input := bytes.NewBuffer(testInput)

	_, err := grep(cmd, input)
	if err == nil {
		t.Error("Expected error for invalid regex, got nil")
	}
}

