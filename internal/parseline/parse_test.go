package parseline

import (
	"CLI/internal/environment"
	"testing"

	"github.com/stretchr/testify/assert"
)


func TestSubstitution(t *testing.T) {
	parser := newTestParser()
	tests := []struct{
		name string
		input string
		expected string
		wantErr error
	}{
		{
			name:     "Simple variable ($VAR)",
			input:    "User: $USER",
			expected: "User: alice",
			wantErr:  nil,
		},
		{
			name:     "Braced variable (${VAR})",
			input:    "Home: ${HOME}",
			expected: "Home: /home/alice",
			wantErr:  nil,
		},
		
		{
			name:     "Multiple variables",
			input:    "App: $APP_NAME v$VERSION",
			expected: "App: testapp v1.0",
			wantErr:  nil,
		},
		
		{
			name:     "Undefined variable",
			input:    "Path: $TTTT",
			expected: "Path: ",
			wantErr:  nil,
		},
		
		{
			name:     "Empty braces (${})",
			input:    "Test: ${}",
			expected: "Test: ${}",
			wantErr:  nil,
		},
		{
			name:     "Dollar sign only ($)",
			input:    "Just $",
			expected: "Just $",
			wantErr:  nil,
		},
		{
			name:     "Mixed content",
			input:    "Run $APP_NAME in ${HOME} with DEBUG=$DEBUG",
			expected: "Run testapp in /home/alice with DEBUG=true",
			wantErr:  nil,
		},
		
		{
			name:     "Unclosed braces (${VAR)",
			input:    "Error: ${USER",
			expected: "",
			wantErr:  nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := parser.substitution(tt.input)
			if err == nil {
				assert.Equal(t, tt.expected, result)
			} 
		})
	}




}

func TestParsePipeline(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []Command
		wantErr  bool
	}{
		{
			name:  "Single command",
			input: "ls -l",
			expected: []Command{
				{Name: "ls", Args: []string{"-l"}},
			},
			wantErr: false,
		},
		{
			name:  "Pipe with two commands",
			input: "cat file.txt | grep 'hello'",
			expected: []Command{
				{Name: "cat", Args: []string{"file.txt"}},
				{Name: "grep", Args: []string{"'hello'"}},
			},
			wantErr: false,
		},
		{
			name:  "Quoted arguments",
			input: `echo "Hello, World!" | awk '{print 1}'`,
			expected: []Command{
				{Name: "echo", Args: []string{`"Hello, World!"`}},
				{Name: "awk", Args: []string{`'{print 1}'`}},
			},
			wantErr: false,
		},
		{
			name:  "Escaped pipe inside quotes",
			input: `echo "Hello | World" | sed 's/|/PIPE/g'`,
			expected: []Command{
				{Name: "echo", Args: []string{`"Hello | World"`}},
				{Name: "sed", Args: []string{`'s/|/PIPE/g'`}},
			},
			wantErr: false,
		},
		{
			name:  "Escaped",
			input: `echo "Hello"\n`,
			expected: []Command{
				{Name: "echo", Args: []string{`"Hello"\n`}},
			},
			wantErr: false,
		},
		{
			name:     "Unclosed quotes (error)",
			input:    `echo "Hello`,
			expected: nil,
			wantErr:  true,
		},
		{
			name:     "Enter variable",
			input:    `echo 123 | x=y`,
			expected: []Command{
				{Name: "echo", Args: []string{"123"}},
				{Name: "x=y", Args: []string{}},
			},
			wantErr:  false,
		},
	}
	parser := newTestParser()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parser.ParsePipeline(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParsePipeline() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !compareCommands(got, tt.expected) {
				t.Errorf("ParsePipeline() = %v, want %v", got, tt.expected)
			}
		})
	}
}

// compareCommands сравнивает два списка команд.
func compareCommands(a, b []Command) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i].Name != b[i].Name || !compareStringSlices(a[i].Args, b[i].Args) {
			return false
		}
	}
	return true
}

// compareStringSlices сравнивает два слайса строк.
func compareStringSlices(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

// создает Parser и Env с переменными для тестирования
func newTestParser() *Parser {
	env := environment.New()
	parser := New(env)
	variables := map[string]string {
		"USER":     "alice",
		"HOME":     "/home/alice",
		"VERSION":  "1.0",
		"APP_NAME": "testapp",
		"DEBUG":    "true",
	}
	for k, v := range variables {
		parser.env.Set(k, v)
	}

	return parser
}