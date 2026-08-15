package main

import (
	"fmt"
	"os"

	"github.com/fakihariefnoto/bloberry/internal/cli"
)

func main() {
	if err := cli.Execute(); err != nil {
		if code, ok := err.(interface{ ExitCode() int }); ok {
			os.Exit(code.ExitCode())
		}
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
