package cli

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

// init — one-command first-run (cli/README.md): server → reachability → tenant → auth.
func newInitCmd() *cobra.Command {
	var server, tenant, token string
	cmd := &cobra.Command{
		Use:   "init",
		Short: "First-run setup: server, tenant, auth",
		RunE: func(cmd *cobra.Command, args []string) error {
			out := NewOutput(flagJSON, flagQuiet, flagNoColor)
			cfg, err := loadConfig()
			if err != nil {
				return err
			}
			// 1. server
			if server == "" {
				if cfg.Server != "" {
					server = cfg.Server
				} else {
					return exit(ExitInvocation, "no server given — pass --server <url>")
				}
			}
			cfg.Server = strings.TrimSuffix(server, "/")

			// 2. reachability
			probe := NewClient(cfg.Server, "")
			if _, _, err := probe.Do("GET", "/health", nil, ""); err != nil {
				return err
			}
			// 3. tenant
			if tenant != "" {
				cfg.Tenant = tenant
			}
			// 4. auth: non-interactive via --token, else keychain
			if token != "" {
				if err := setRefreshToken(token); err != nil {
					return err
				}
			}
			if err := saveConfig(cfg); err != nil {
				return err
			}
			out.Info("configured server %s", cfg.Server)
			if token != "" {
				out.Info("token stored in OS keychain")
			} else {
				out.Info("run 'bloberry auth login' to authenticate")
			}
			out.Data(map[string]string{"server": cfg.Server, "tenant": cfg.Tenant})
			return nil
		},
	}
	cmd.Flags().StringVar(&server, "server", "", "server URL")
	cmd.Flags().StringVar(&tenant, "tenant", "", "tenant slug")
	cmd.Flags().StringVar(&token, "token", "", "refresh token (never written to the config file)")
	return cmd
}

func newVersionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print version, commit and build date",
		RunE: func(cmd *cobra.Command, args []string) error {
			out := NewOutput(flagJSON, flagQuiet, flagNoColor)
			if flagJSON {
				return out.Data(map[string]string{"version": Version, "commit": Commit, "build_date": BuildDate})
			}
			fmt.Fprintf(out.Stdout, "bloberry %s (%s, %s)\n", Version, Commit, BuildDate)
			return nil
		},
	}
}
