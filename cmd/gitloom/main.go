package main

import (
	"log"

	"github.com/ovitorvalente/git-loom/internal/cli"
)

func main() {
	if err := cli.Execute(); err != nil {
		log.Fatal(err)
	}
}
