package cli

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/spf13/cobra"
)

func httpNewRequest(method, url string, body interface{ Read([]byte) (int, error) }) (*http.Request, error) {
	if r, ok := body.(*strings.Reader); ok {
		return http.NewRequest(method, url, r)
	}
	return http.NewRequest(method, url, body.(interface{ Read([]byte) (int, error) }))
}

func newFolderCmd() *cobra.Command {
	cmd := &cobra.Command{Use: "folder", Short: "Folder operations"}
	cmd.AddCommand(
		&cobra.Command{
			Use:   "create <path>",
			Short: "Create a folder",
			Args:  cobra.ExactArgs(1),
			RunE: func(cmd *cobra.Command, args []string) error {
				client, _, ok := configClient()
				if !ok {
					return exit(ExitNotAuthed, "not authenticated")
				}
				name := strings.Trim(args[0], "/")
				parent := "root"
				if i := strings.LastIndex(name, "/"); i >= 0 {
					parent = name[:i]
					name = name[i+1:]
				}
				body := fmt.Sprintf(`{"name":%q,"parent_id":%q}`, name, parent)
				if _, _, err := client.Do("POST", "/folders", strings.NewReader(body), "application/json"); err != nil {
					return err
				}
				return nil
			},
		},
		&cobra.Command{
			Use:   "tree [path]",
			Short: "Show the folder hierarchy as a tree",
			Args:  cobra.MaximumNArgs(1),
			RunE: func(cmd *cobra.Command, args []string) error {
				out := NewOutput(flagJSON, flagQuiet, flagNoColor)
				client, _, ok := configClient()
				if !ok {
					return exit(ExitNotAuthed, "not authenticated")
				}
				path := "root"
				if len(args) > 0 {
					path = strings.TrimPrefix(args[0], "bloberry://")
				}
				env, _, err := client.Do("GET", "/folders/"+path+"/tree", nil, "")
				if err != nil {
					return err
				}
				var tree struct {
					Folders []struct{ Name string `json:"name"` } `json:"folders"`
				}
				_ = env.JSON(&tree)
				return out.Data(tree)
			},
		},
	)
	return cmd
}

func newShareCmd() *cobra.Command {
	var ttl int
	cmd := &cobra.Command{Use: "share", Short: "Create and manage share links"}
	cmd.PersistentFlags().IntVar(&ttl, "ttl", 3600, "link TTL in seconds")
	cmd.AddCommand(
		&cobra.Command{
			Use:   "link <path>",
			Short: "Create a signed share URL",
			Args:  cobra.ExactArgs(1),
			RunE: func(cmd *cobra.Command, args []string) error {
				out := NewOutput(flagJSON, flagQuiet, flagNoColor)
				client, _, ok := configClient()
				if !ok {
					return exit(ExitNotAuthed, "not authenticated")
				}
				objectID := strings.TrimPrefix(args[0], "bloberry://")
				body := fmt.Sprintf(`{"object_id":%q,"ttl":%d}`, objectID, ttl)
				env, _, err := client.Do("POST", "/shares", strings.NewReader(body), "application/json")
				if err != nil {
					return err
				}
				var link struct {
					URL string `json:"url"`
				}
				_ = env.JSON(&link)
				return out.Data(link.URL)
			},
		},
		&cobra.Command{
			Use:   "short <path>",
			Short: "Create a short URL",
			Args:  cobra.ExactArgs(1),
			RunE: func(cmd *cobra.Command, args []string) error {
				out := NewOutput(flagJSON, flagQuiet, flagNoColor)
				client, _, ok := configClient()
				if !ok {
					return exit(ExitNotAuthed, "not authenticated")
				}
				body := fmt.Sprintf(`{"object_id":%q}`, strings.TrimPrefix(args[0], "bloberry://"))
				env, _, err := client.Do("POST", "/shares/short", strings.NewReader(body), "application/json")
				if err != nil {
					return err
				}
				var link struct {
					URL string `json:"url"`
				}
				_ = env.JSON(&link)
				return out.Data(link.URL)
			},
		},
		&cobra.Command{
			Use:   "list",
			Short: "List the current tenant's share links",
			RunE: func(cmd *cobra.Command, args []string) error {
				out := NewOutput(flagJSON, flagQuiet, flagNoColor)
				client, _, ok := configClient()
				if !ok {
					return exit(ExitNotAuthed, "not authenticated")
				}
				env, _, err := client.Do("GET", "/shares", nil, "")
				if err != nil {
					return err
				}
				var links []map[string]interface{}
				_ = env.JSON(&links)
				return out.Data(links)
			},
		},
		&cobra.Command{
			Use:   "revoke <id>",
			Short: "Revoke a share link",
			Args:  cobra.ExactArgs(1),
			RunE: func(cmd *cobra.Command, args []string) error {
				client, _, ok := configClient()
				if !ok {
					return exit(ExitNotAuthed, "not authenticated")
				}
				_, _, err := client.Do("DELETE", "/shares/"+args[0], nil, "")
				return err
			},
		},
	)
	return cmd
}

func newKeyCmd() *cobra.Command {
	var folder, permission, expires, app string
	cmd := &cobra.Command{Use: "key", Short: "Application access keys"}
	cmd.PersistentFlags().StringVar(&folder, "folder", "", "scope to a folder subtree")
	cmd.PersistentFlags().StringVar(&permission, "permission", "read", "permission(s) — comma-separated")
	cmd.PersistentFlags().StringVar(&expires, "expires", "", "expiry (RFC3339)")
	cmd.PersistentFlags().StringVar(&app, "app", "", "application ID")
	cmd.AddCommand(
		&cobra.Command{
			Use:   "create",
			Short: "Create a scoped access key (secret shown once)",
			RunE: func(cmd *cobra.Command, args []string) error {
				out := NewOutput(flagJSON, flagQuiet, flagNoColor)
				client, _, ok := configClient()
				if !ok {
					return exit(ExitNotAuthed, "not authenticated")
				}
				if app == "" {
					return exit(ExitInvocation, "pass --app <application-id> to create a key")
				}
				perms := strings.Split(permission, ",")
				var scope []string
				if folder != "" {
					scope = strings.Split(folder, ",")
				}
				payload := map[string]interface{}{"permissions": perms, "scope_folder_ids": scope}
				if expires != "" {
					payload["expires_at"] = expires
				}
				bodyBytes, _ := json.Marshal(payload)
				env, _, err := client.Do("POST", "/applications/"+app+"/keys", strings.NewReader(string(bodyBytes)), "application/json")
				if err != nil {
					return err
				}
				var key map[string]interface{}
				_ = env.JSON(&key)
				if flagJSON {
					return out.Data(key)
				}
				fmt.Fprintf(out.Stdout, "secret: %s\n", key["secret"])
				fmt.Fprintf(out.Stderr, "this secret is shown once\n")
				return nil
			},
		},
		&cobra.Command{
			Use:   "list",
			Short: "List access keys",
			RunE: func(cmd *cobra.Command, args []string) error {
				out := NewOutput(flagJSON, flagQuiet, flagNoColor)
				client, _, ok := configClient()
				if !ok {
					return exit(ExitNotAuthed, "not authenticated")
				}
				env, _, err := client.Do("GET", "/applications", nil, "")
				if err != nil {
					return err
				}
				var apps []struct {
					ID   string `json:"id"`
					Name string `json:"name"`
				}
				_ = env.JSON(&apps)
				return out.Data(apps)
			},
		},
		&cobra.Command{
			Use:   "revoke <id>",
			Short: "Revoke an access key",
			Args:  cobra.ExactArgs(1),
			RunE: func(cmd *cobra.Command, args []string) error {
				client, _, ok := configClient()
				if !ok {
					return exit(ExitNotAuthed, "not authenticated")
				}
				_, _, err := client.Do("DELETE", "/applications/x/keys/"+args[0], nil, "")
				return err
			},
		},
	)
	return cmd
}

func newAppCmd() *cobra.Command {
	cmd := &cobra.Command{Use: "app", Short: "Applications (non-human principals)"}
	cmd.AddCommand(
		&cobra.Command{
			Use:   "create <name>",
			Short: "Register an application",
			Args:  cobra.ExactArgs(1),
			RunE: func(cmd *cobra.Command, args []string) error {
				out := NewOutput(flagJSON, flagQuiet, flagNoColor)
				client, _, ok := configClient()
				if !ok {
					return exit(ExitNotAuthed, "not authenticated")
				}
				env, _, err := client.Do("POST", "/applications", strings.NewReader(fmt.Sprintf(`{"name":%q}`, args[0])), "application/json")
				if err != nil {
					return err
				}
				var app map[string]interface{}
				_ = env.JSON(&app)
				return out.Data(app)
			},
		},
		&cobra.Command{
			Use:   "list",
			Short: "List applications",
			RunE: func(cmd *cobra.Command, args []string) error {
				out := NewOutput(flagJSON, flagQuiet, flagNoColor)
				client, _, ok := configClient()
				if !ok {
					return exit(ExitNotAuthed, "not authenticated")
				}
				env, _, err := client.Do("GET", "/applications", nil, "")
				if err != nil {
					return err
				}
				var apps []map[string]interface{}
				_ = env.JSON(&apps)
				return out.Data(apps)
			},
		},
		&cobra.Command{
			Use:   "delete <id>",
			Short: "Delete an application",
			Args:  cobra.ExactArgs(1),
			RunE: func(cmd *cobra.Command, args []string) error {
				client, _, ok := configClient()
				if !ok {
					return exit(ExitNotAuthed, "not authenticated")
				}
				_, _, err := client.Do("DELETE", "/applications/"+args[0], nil, "")
				return err
			},
		},
	)
	return cmd
}
