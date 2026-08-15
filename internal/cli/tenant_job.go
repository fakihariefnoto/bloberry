package cli

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

func newTenantCmd() *cobra.Command {
	cmd := &cobra.Command{Use: "tenant", Short: "Tenant operations"}
	cmd.AddCommand(
		&cobra.Command{
			Use:   "list",
			Short: "List tenants the current user belongs to",
			RunE: func(cmd *cobra.Command, args []string) error {
				out := NewOutput(flagJSON, flagQuiet, flagNoColor)
				client, _, ok := configClient()
				if !ok {
					return exit(ExitNotAuthed, "not authenticated")
				}
				env, _, err := client.Do("GET", "/tenants", nil, "")
				if err != nil {
					return err
				}
				var ts []map[string]interface{}
				_ = env.JSON(&ts)
				return out.Data(ts)
			},
		},
		&cobra.Command{
			Use:   "use <slug>",
			Short: "Set the default tenant",
			Args:  cobra.ExactArgs(1),
			RunE: func(cmd *cobra.Command, args []string) error {
				cfg, err := loadConfig()
				if err != nil {
					return err
				}
				cfg.Tenant = args[0]
				return saveConfig(cfg)
			},
		},
		&cobra.Command{
			Use:   "usage",
			Short: "Show bytes, objects and estimated cost",
			RunE: func(cmd *cobra.Command, args []string) error {
				out := NewOutput(flagJSON, flagQuiet, flagNoColor)
				client, _, ok := configClient()
				if !ok {
					return exit(ExitNotAuthed, "not authenticated")
				}
				env, _, err := client.Do("GET", "/usage/me", nil, "")
				if err != nil {
					return err
				}
				var u map[string]interface{}
				_ = env.JSON(&u)
				return out.Data(u)
			},
		},
	)
	return cmd
}

func newJobCmd() *cobra.Command {
	cmd := &cobra.Command{Use: "job", Short: "Job operations"}
	cmd.AddCommand(
		&cobra.Command{
			Use:   "list",
			Short: "List jobs",
			RunE: func(cmd *cobra.Command, args []string) error {
				out := NewOutput(flagJSON, flagQuiet, flagNoColor)
				client, _, ok := configClient()
				if !ok {
					return exit(ExitNotAuthed, "not authenticated")
				}
				env, _, err := client.Do("GET", "/jobs", nil, "")
				if err != nil {
					return err
				}
				var jobs []map[string]interface{}
				_ = env.JSON(&jobs)
				return out.Data(jobs)
			},
		},
		&cobra.Command{
			Use:   "status <id>",
			Short: "Show a job's status",
			Args:  cobra.ExactArgs(1),
			RunE: func(cmd *cobra.Command, args []string) error {
				out := NewOutput(flagJSON, flagQuiet, flagNoColor)
				client, _, ok := configClient()
				if !ok {
					return exit(ExitNotAuthed, "not authenticated")
				}
				env, _, err := client.Do("GET", "/jobs/"+args[0], nil, "")
				if err != nil {
					return err
				}
				var j map[string]interface{}
				_ = env.JSON(&j)
				return out.Data(j)
			},
		},
		&cobra.Command{
			Use:   "watch <id>",
			Short: "Block until a job reaches a terminal state",
			Args:  cobra.ExactArgs(1),
			RunE: func(cmd *cobra.Command, args []string) error {
				client, _, ok := configClient()
				if !ok {
					return exit(ExitNotAuthed, "not authenticated")
				}
				for i := 0; i < 120; i++ {
					env, _, err := client.Do("GET", "/jobs/"+args[0], nil, "")
					if err != nil {
						return err
					}
					var j struct {
						State string `json:"state"`
					}
					_ = env.JSON(&j)
					if j.State == "succeeded" || j.State == "failed" {
						fmt.Fprintf(cmd.OutOrStdout(), "%s\n", j.State)
						return nil
					}
				}
				return exit(ExitFailed, "job did not finish in time")
			},
		},
	)
	return cmd
}

func newArchiveCmd() *cobra.Command {
	cmd := &cobra.Command{Use: "archive", Short: "Server-side archive operations"}
	cmd.AddCommand(
		&cobra.Command{
			Use:   "extract <path>",
			Short: "Extract an archive server-side into a folder",
			Args:  cobra.ExactArgs(1),
			RunE: func(cmd *cobra.Command, args []string) error {
				out := NewOutput(flagJSON, flagQuiet, flagNoColor)
				client, _, ok := configClient()
				if !ok {
					return exit(ExitNotAuthed, "not authenticated")
				}
				fileID := strings.TrimPrefix(args[0], "bloberry://")
				env, _, err := client.Do("POST", "/archives/extract", strings.NewReader(fmt.Sprintf(`{"file_id":%q,"target_folder_id":"root"}`, fileID)), "application/json")
				if err != nil {
					return err
				}
				var r struct {
					JobID string `json:"job_id"`
				}
				_ = env.JSON(&r)
				return out.Data(r.JobID)
			},
		},
		&cobra.Command{
			Use:   "bundle <paths...>",
			Short: "Create a bundle download of a folder",
			Args:  cobra.MinimumNArgs(1),
			RunE: func(cmd *cobra.Command, args []string) error {
				out := NewOutput(flagJSON, flagQuiet, flagNoColor)
				client, _, ok := configClient()
				if !ok {
					return exit(ExitNotAuthed, "not authenticated")
				}
				env, _, err := client.Do("POST", "/archives/bundle", strings.NewReader(`{"folder_id":"root"}`), "application/json")
				if err != nil {
					return err
				}
				var r struct {
					JobID string `json:"job_id"`
				}
				_ = env.JSON(&r)
				return out.Data(r.JobID)
			},
		},
	)
	return cmd
}

func newAdminCmd() *cobra.Command {
	cmd := &cobra.Command{Use: "admin", Short: "Platform admin operations"}
	backend := &cobra.Command{Use: "backend", Short: "Storage backend management"}
	backend.AddCommand(
		&cobra.Command{
			Use:   "list",
			Short: "List storage backends",
			RunE: func(cmd *cobra.Command, args []string) error {
				out := NewOutput(flagJSON, flagQuiet, flagNoColor)
				client, _, ok := configClient()
				if !ok {
					return exit(ExitNotAuthed, "not authenticated")
				}
				env, _, err := client.Do("GET", "/admin/backends", nil, "")
				if err != nil {
					return err
				}
				var bs []map[string]interface{}
				_ = env.JSON(&bs)
				return out.Data(bs)
			},
		},
		&cobra.Command{
			Use:   "add",
			Short: "Register a storage backend",
			RunE: func(cmd *cobra.Command, args []string) error {
				return exit(ExitInvocation, "use the web dashboard to register backends (credentials are secret-shaped)")
			},
		},
	)
	cmd.AddCommand(backend)
	return cmd
}
