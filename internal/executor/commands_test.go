package executor

import (
	"CLI/internal/environment"
	"CLI/internal/parseline"
	"bytes"
	"os"
	"path/filepath"
	"slices"
	"strings"
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
			got, err := cat(tt.cmd, tt.input, nil)
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
			got, err := echo(tt.cmd, tt.input, nil)
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
	got, err := pwd(cmd, nil, nil)
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
			got, err := wc(tt.cmd, tt.input, nil)
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

	output, err := grep(cmd, input, nil)
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

	output, err := grep(cmd, input, nil)
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

	output, err := grep(cmd, input, nil)
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

	output, err := grep(cmd, input, nil)
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

	output, err := grep(cmd, input, nil)
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

	output, err := grep(cmd, input, nil)
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

	_, err := grep(cmd, input, nil)
	if err == nil {
		t.Error("Expected error for invalid regex, got nil")
	}
}

func TestLsEmpty(t *testing.T) {
	tempDir := t.TempDir()

	cmd := parseline.Command{
		Name: "ls",
		Args: []string{tempDir},
	}
	var testInput []byte
	input := bytes.NewBuffer(testInput)
	expected := ``

	output, err := ls(cmd, input)
	if err != nil {
		t.Fatalf("ls failed: %v", err)
	}
	if output.String() != expected {
		t.Errorf("Expected:\n%s\nGot:\n%s", expected, output.String())
	}
}

func TestLsNonEmpty(t *testing.T) {
	tempDir := t.TempDir()
	var filePaths []string
	for i := 0; i < 5; i++ {
		file, err := os.CreateTemp(tempDir, "testfile-*.tmp")
		if err != nil {
			t.Fatalf("Failed to create temporary file: %v", err)
		}
		defer file.Close()
		filePaths = append(filePaths, filepath.Base(file.Name()))
	}

	cmd := parseline.Command{
		Name: "ls",
		Args: []string{tempDir},
	}
	var testInput []byte
	input := bytes.NewBuffer(testInput)
	slices.Sort(filePaths)
	expected := strings.Join(filePaths, "\n") + "\n"

	output, err := ls(cmd, input)
	if err != nil {
		t.Fatalf("ls failed: %v", err)
	}
	if output.String() != expected {
		t.Errorf("Expected:\n%s\nGot:\n%s", expected, output.String())
	}
}

func TestLsCurrentDir(t *testing.T) {
	cmd := parseline.Command{
		Name: "ls",
	}
	var testInput []byte
	input := bytes.NewBuffer(testInput)

	output, err := ls(cmd, input)
	if err != nil {
		t.Fatalf("ls failed: %v", err)
	}
	output_tokens := strings.Split(output.String(), "\n")
	t.Log(output_tokens)
	if !slices.Contains(output_tokens, "commands_test.go") {
		t.Errorf("ls without argument does not work")
	}
}
func TestCdCasualUsage(t *testing.T) {
	tmpDir := t.TempDir()
	subDir := filepath.Join(tmpDir, "subdir")
	if err := os.Mkdir(subDir, 0755); err != nil {
		t.Fatalf("failed to create subdir: %v", err)
	}

	env := environment.New()
	exec := New(env)
	exec.cwd = tmpDir

	if exec.cwd != tmpDir {
		t.Fatalf("expected cwd to be %s, got %s", tmpDir, exec.cwd)
	}

	cmd := parseline.Command{Name: "cd", Args: []string{"subdir"}}
	_, err := cd(cmd, &bytes.Buffer{}, exec)
	if err != nil {
		t.Fatalf("cd command failed: %v", err)
	}

	if exec.cwd != subDir {
		t.Errorf("expected cwd to be %s after cd, got %s", subDir, exec.cwd)
	}
}

func TestCdRelativeAndDotNotation(t *testing.T) {
	tmpRoot := t.TempDir()

	subdir1 := filepath.Join(tmpRoot, "dir1")
	subdir2 := filepath.Join(subdir1, "dir2")
	if err := os.MkdirAll(subdir2, 0755); err != nil {
		t.Fatalf("could not create dirs: %v", err)
	}

	env := environment.New()
	exec := New(env)
	exec.cwd = tmpRoot

	tests := []struct {
		input   string
		wantDir string
	}{
		{"dir1", subdir1},
		{".", subdir1},
		{"dir2", subdir2},
		{"../..", tmpRoot},
	}

	for _, tt := range tests {
		cmd := parseline.Command{Name: "cd", Args: []string{tt.input}}
		_, err := cd(cmd, &bytes.Buffer{}, exec)
		if err != nil {
			t.Errorf("cd %q failed: %v", tt.input, err)
			continue
		}
		if exec.cwd != tt.wantDir {
			t.Errorf("after cd %q: expected %q, got %q", tt.input, tt.wantDir, exec.cwd)
		}
	}
}

func TestCdToHome(t *testing.T) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		t.Skip("can't test home dir; os.UserHomeDir() failed")
	}

	env := environment.New()
	exec := New(env)
	exec.cwd = "/"

	cmd := parseline.Command{Name: "cd"} // no args
	_, err = cd(cmd, &bytes.Buffer{}, exec)
	if err != nil {
		t.Fatalf("cd to home failed: %v", err)
	}
	if exec.cwd != homeDir {
		t.Errorf("expected home dir %q, got %q", homeDir, exec.cwd)
	}
}

func TestCdToInvalidPath(t *testing.T) {
	env := environment.New()
	exec := New(env)
	exec.cwd = t.TempDir()

	cmd := parseline.Command{Name: "cd", Args: []string{"not-a-dir"}}
	_, err := cd(cmd, &bytes.Buffer{}, exec)
	if err == nil {
		t.Fatal("expected error when cd into nonexistent dir, got none")
	}
}
