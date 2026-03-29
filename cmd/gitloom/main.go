package main

import (
	"fmt"
	"os"

	"github.com/ovitorvalente/git-loom/internal/cli"
)

func main() {
	if err := cli.Execute(); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
