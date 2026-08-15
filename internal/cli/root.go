package cli

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

// Build-time injected via -ldflags -X (cli/README.md §Output conventions).
var (
	Version   = "dev"
	Commit    = "none"
	BuildDate = "unknown"
)

var (
	flagConfig  string
	flagJSON    bool
	flagQuiet   bool
	flagVerbose bool
	flagNoColor bool
	flagTenant  string
)

var rootCmd = &cobra.Command{
	Use:   "bloberry",
	Short: "Bloberry — storage-agnostic object service CLI",
	Long:  "Bloberry CLI. Companion to a Bloberry server: move files, issue keys, manage shares, from a terminal or CI.",
	SilenceUsage:  true,
	SilenceErrors: true,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		// config + server must exist for commands that need it
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	},
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	pf := rootCmd.PersistentFlags()
	pf.StringVar(&flagConfig, "config", "", "config file path")
	pf.BoolVar(&flagJSON, "json", false, "emit structured JSON")
	pf.BoolVarP(&flagQuiet, "quiet", "q", false, "suppress non-essential output")
	pf.BoolVarP(&flagVerbose, "verbose", "v", false, "verbose output")
	pf.BoolVar(&flagNoColor, "no-color", false, "disable color")
	pf.StringVar(&flagTenant, "tenant", "", "tenant slug")

	rootCmd.AddCommand(
		newInitCmd(),
		newVersionCmd(),
		newConfigCmd(),
		newCompletionCmd(),
		newAuthCmd(),
		newLSCmd(),
		newStatCmd(),
		newCatCmd(),
		newCPSend(),
		newRMCmd(),
		newFolderCmd(),
		newShareCmd(),
		newKeyCmd(),
		newAppCmd(),
		newTenantCmd(),
		newJobCmd(),
		newArchiveCmd(),
		newAdminCmd(),
	)
}

// session returns the authenticated client for the current config + token.
func session() (*Client, *Config, error) {
	cfg, err := loadConfig()
	if err != nil {
		return nil, nil, err
	}
	if cfg.Server == "" {
		return nil, nil, errNoServer
	}
	envToken, envOK, err := getToken()
	if err != nil {
		return nil, nil, err
	}
	if !envOK {
		return nil, nil, exit(ExitNotAuthed, "not authenticated — run 'bloberry auth login' or set BLOBERRY_TOKEN")
	}
	c := NewClient(cfg.Server, envToken)
	// CLI presents as a non-web platform
	return c, cfg, nil
}

func loadConfig() (*Config, error) {
	path := flagConfig
	if path == "" {
		path = configPath()
	}
	cfg := &Config{Output: "table", Color: "auto"}
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return cfg, nil
		}
		return nil, err
	}
	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, fmt.Errorf("config %s: %w", path, err)
	}
	return cfg, nil
}

func saveConfig(cfg *Config) error {
	path := flagConfig
	if path == "" {
		path = configPath()
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return err
	}
	data, err := yaml.Marshal(cfg)
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o600)
}

func configClient() (*Client, *Config, bool) {
	cfg, err := loadConfig()
	if err != nil || cfg.Server == "" {
		return nil, cfg, false
	}
	envToken, envOK, _ := getToken()
	if !envOK {
		return nil, cfg, false
	}
	return NewClient(cfg.Server, envToken), cfg, true
}
