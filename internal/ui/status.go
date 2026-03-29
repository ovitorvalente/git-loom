package ui

import (
	"fmt"
	"io"
)

func PrintStatus(w io.Writer, message string) {
	fmt.Fprintf(w, "\r%s\n", colorizeLine(accentColor, "  "+message))
}

func PrintStatusDone(w io.Writer, message string) {
	fmt.Fprintf(w, "\r%s\n", colorizeLine(successColor, "  "+message))
}
