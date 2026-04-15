package ui

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	"github.com/ovitorvalente/git-loom/internal/shared"
)

func ConfirmCommit(input io.Reader, output io.Writer, question string) (bool, error) {
	promptPrefix := colorizeText(statusPromptColor, ">")
	highlightedQuestion := colorizeText(defaultColor, question)
	suffix := colorizeText(emphasisColor, shared.MessageCommitPromptSuffix)
	if _, err := fmt.Fprintf(output, "\n%s %s %s", promptPrefix, highlightedQuestion, suffix); err != nil {
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

func AskInput(input io.Reader, output io.Writer, question string) (string, error) {
	promptPrefix := colorizeText(statusPromptColor, ">")
	highlightedQuestion := colorizeText(defaultColor, question)
	suffix := colorizeText(emphasisColor, ": ")
	if _, err := fmt.Fprintf(output, "\n%s %s%s", promptPrefix, highlightedQuestion, suffix); err != nil {
		return "", err
	}

	reader := bufio.NewReader(input)
	answer, err := reader.ReadString('\n')
	if err != nil && err != io.EOF {
		return "", err
	}

	return strings.TrimSpace(answer), nil
}
