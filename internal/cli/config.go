package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

func newConfigCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config",
		Short: "Manage CLI configuration",
	}
	cmd.AddCommand(
		&cobra.Command{
			Use:   "get <key>",
			Short: "Print a config value",
			Args:  cobra.ExactArgs(1),
			RunE: func(cmd *cobra.Command, args []string) error {
				out := NewOutput(flagJSON, flagQuiet, flagNoColor)
				cfg, err := loadConfig()
				if err != nil {
					return err
				}
				v := ""
				switch args[0] {
				case "server":
					v = cfg.Server
				case "tenant":
					v = cfg.Tenant
				case "output":
					v = cfg.Output
				case "color":
					v = cfg.Color
				default:
					return exit(ExitInvocation, "unknown key %q (server|tenant|output|color)", args[0])
				}
				return out.Data(v)
			},
		},
		&cobra.Command{
			Use:   "set <key> <value>",
			Short: "Set a config value",
			Args:  cobra.ExactArgs(2),
			RunE: func(cmd *cobra.Command, args []string) error {
				cfg, err := loadConfig()
				if err != nil {
					return err
				}
				switch args[0] {
				case "server":
					cfg.Server = args[1]
				case "tenant":
					cfg.Tenant = args[1]
				case "output":
					cfg.Output = args[1]
				case "color":
					cfg.Color = args[1]
				default:
					return exit(ExitInvocation, "unknown key %q", args[0])
				}
				return saveConfig(cfg)
			},
		},
		&cobra.Command{
			Use:   "path",
			Short: "Print the config file path",
			RunE: func(cmd *cobra.Command, args []string) error {
				fmt.Println(configPath())
				return nil
			},
		},
	)
	return cmd
}

func newCompletionCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:       "completion <bash|zsh|fish|powershell>",
		Short:     "Generate a shell completion script",
		Args:      cobra.ExactArgs(1),
		ValidArgs: []string{"bash", "zsh", "fish", "powershell"},
		RunE: func(cmd *cobra.Command, args []string) error {
			switch args[0] {
			case "bash":
				return rootCmd.GenBashCompletion(cmd.OutOrStdout())
			case "zsh":
				return rootCmd.GenZshCompletion(cmd.OutOrStdout())
			case "fish":
				return rootCmd.GenFishCompletion(cmd.OutOrStdout(), true)
			case "powershell":
				return rootCmd.GenPowerShellCompletion(cmd.OutOrStdout())
			}
			return nil
		},
	}
	return cmd
}
