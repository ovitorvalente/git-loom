package ui

import (
	"bytes"
	"strings"
	"testing"
)

func TestConfirmCommit(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name     string
		input    string
		expected bool
	}{
		{
			name:     "confirms with y",
			input:    "y\n",
			expected: true,
		},
		{
			name:     "confirms with yes",
			input:    "yes\n",
			expected: true,
		},
		{
			name:     "confirms empty input",
			input:    "\n",
			expected: true,
		},
		{
			name:     "rejects arbitrary input",
			input:    "nope\n",
			expected: false,
		},
	}

	for _, testCase := range testCases {
		testCase := testCase
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			output := &bytes.Buffer{}
			confirmed, err := ConfirmCommit(strings.NewReader(testCase.input), output, "criar commit?")
			if err != nil {
				t.Fatalf("expected nil error, got %v", err)
			}
			if confirmed != testCase.expected {
				t.Fatalf("expected %v, got %v", testCase.expected, confirmed)
			}
			outputStr := output.String()
			if !strings.Contains(outputStr, ">") {
				t.Fatalf("expected prompt to contain '>', got %q", outputStr)
			}
			if !strings.Contains(outputStr, "criar commit?") {
				t.Fatalf("expected prompt to contain 'criar commit?', got %q", outputStr)
			}
			if !strings.Contains(outputStr, "[Y/n]:") {
				t.Fatalf("expected prompt to contain '[Y/n]:', got %q", outputStr)
			}
		})
	}
}

func TestAskInput(t *testing.T) {
	t.Parallel()

	output := &bytes.Buffer{}
	answer, err := AskInput(strings.NewReader("nova mensagem\n"), output, "editar mensagem")
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if answer != "nova mensagem" {
		t.Fatalf("expected trimmed answer, got %q", answer)
	}
	if !strings.Contains(output.String(), "editar mensagem") {
		t.Fatalf("expected prompt output, got %q", output.String())
	}
}
