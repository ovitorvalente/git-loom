package ui

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	"github.com/ovitorvalente/git-loom/internal/shared"
)

func ConfirmCommit(input io.Reader, output io.Writer, question string) (bool, error) {
	highlightedQuestion := colorizeLine(headerColor, question)
	promptPrefix := colorizeLine(borderColor, ">")
	if _, err := fmt.Fprintf(output, "\n%s %s %s", promptPrefix, highlightedQuestion, shared.MessageCommitPromptSuffix); err != nil {
		return false, err
	}

	reader := bufio.NewReader(input)
	answer, err := reader.ReadString('\n')
	if err != nil && err != io.EOF {
		return false, err
	}

	normalizedAnswer := strings.ToLower(strings.TrimSpace(answer))
	return normalizedAnswer == "" || normalizedAnswer == "y" || normalizedAnswer == "yes", nil
}
