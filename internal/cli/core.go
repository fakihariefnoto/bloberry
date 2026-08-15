package cli

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

// Exit codes are API (cli/README.md §Exit codes). Stable once shipped.
const (
	ExitOK          = 0
	ExitFailed      = 1
	ExitInvocation  = 2
	ExitNotAuthed   = 3
	ExitForbidden   = 4
	ExitNotFound    = 5
	ExitQuota       = 6
	ExitPartial     = 7
	ExitConflict    = 8
	ExitBackendDown = 9
)

type ExitError struct {
	Code    int
	Message string
}

func (e *ExitError) Error() string { return e.Message }
func (e *ExitError) ExitCode() int { return e.Code }

func exit(code int, format string, args ...interface{}) error {
	return &ExitError{Code: code, Message: fmt.Sprintf(format, args...)}
}

// configPath returns the CLI config location per OS (XDG / APPDATA).
func configPath() string {
	if p := os.Getenv("BLOBERRY_CONFIG"); p != "" {
		return p
	}
	base := os.Getenv("XDG_CONFIG_HOME")
	if base == "" {
		base = "~/.config"
	}
	if runtime.GOOS == "windows" {
		if a := os.Getenv("APPDATA"); a != "" {
			base = a
		}
	}
	return filepath.Join(expand(base), "bloberry", "config.yaml")
}

func expand(p string) string {
	if len(p) > 0 && p[0] == '~' {
		home, err := os.UserHomeDir()
		if err == nil {
			return filepath.Join(home, p[1:])
		}
	}
	return p
}

// Config mirrors the CLI config file (cli/README.md §Config).
type Config struct {
	Server string `yaml:"server"`
	Tenant string `yaml:"tenant"`
	Output string `yaml:"output"`
	Color  string `yaml:"color"`
}

var errNoServer = errors.New("no server configured — run 'bloberry init'")
