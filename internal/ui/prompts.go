package ui

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

func ConfirmCommit(input io.Reader, output io.Writer) (bool, error) {
	if _, err := fmt.Fprint(output, "create commit? [y/N]: "); err != nil {
		return false, err
	}

	reader := bufio.NewReader(input)
	answer, err := reader.ReadString('\n')
	if err != nil && err != io.EOF {
		return false, err
	}

	normalizedAnswer := strings.ToLower(strings.TrimSpace(answer))
	return normalizedAnswer == "y" || normalizedAnswer == "yes", nil
}
