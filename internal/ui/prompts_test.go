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
			name:     "rejects empty input",
			input:    "\n",
			expected: false,
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
			confirmed, err := ConfirmCommit(strings.NewReader(testCase.input), output)
			if err != nil {
				t.Fatalf("expected nil error, got %v", err)
			}
			if confirmed != testCase.expected {
				t.Fatalf("expected %v, got %v", testCase.expected, confirmed)
			}
			if !strings.Contains(output.String(), "create commit? [y/N]: ") {
				t.Fatalf("unexpected prompt output: %q", output.String())
			}
		})
	}
}
