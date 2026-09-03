package main

import (
	"os"

	"app-do-breno-agent/internal/cli"
)

func main() {
	os.Exit(cli.Execute())
}
