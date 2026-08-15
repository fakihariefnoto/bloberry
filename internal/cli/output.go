package cli

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

// Output follows the stdout-is-data convention (cli/README.md §Output).
// Decoration and errors go to stderr; data goes to stdout.
type Output struct {
	JSON bool
	Quiet bool
	NoColor bool
	Stdout io.Writer
	Stderr io.Writer
}

func NewOutput(jsonMode, quiet, noColor bool) *Output {
	return &Output{
		JSON: jsonMode, Quiet: quiet, NoColor: noColor || os.Getenv("NO_COLOR") != "" || os.Getenv("TERM") == "dumb",
		Stdout: os.Stdout, Stderr: os.Stderr,
	}
}

func (o *Output) Data(v interface{}) error {
	if o.JSON {
		enc := json.NewEncoder(o.Stdout)
		enc.SetIndent("", "  ")
		return enc.Encode(v)
	}
	// default text: print scalar, or JSON for structured
	switch t := v.(type) {
	case string:
		fmt.Fprintln(o.Stdout, t)
	default:
		enc := json.NewEncoder(o.Stdout)
		return enc.Encode(v)
	}
	return nil
}

func (o *Output) Info(format string, args ...interface{}) {
	if !o.Quiet {
		fmt.Fprintf(o.Stderr, format+"\n", args...)
	}
}

func (o *Output) Warn(format string, args ...interface{}) {
	fmt.Fprintf(o.Stderr, "warning: "+format+"\n", args...)
}
