package cli

import (
	"encoding/json"

	"github.com/spf13/cobra"
)

type authResp struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`
}

func newAuthCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "auth",
		Short: "Authenticate with a Bloberry server",
	}
	cmd.AddCommand(
		&cobra.Command{
			Use:   "login",
			Short: "Log in and store the session in the OS keychain",
			RunE: func(cmd *cobra.Command, args []string) error {
				cfg, err := loadConfig()
				if err != nil {
					return err
				}
				if cfg.Server == "" {
					return errNoServer
				}
				return exit(ExitInvocation, "interactive login needs a browser device flow; use 'bloberry init --token <refresh-token>' or set BLOBERRY_TOKEN")
			},
		},
		&cobra.Command{
			Use:   "logout",
			Short: "Revoke the session server-side and clear the local token",
			RunE: func(cmd *cobra.Command, args []string) error {
				out := NewOutput(flagJSON, flagQuiet, flagNoColor)
				client, _, ok := configClient()
				if !ok {
					return exit(ExitNotAuthed, "not authenticated")
				}
				if _, _, err := client.Do("POST", "/auth/logout", nil, ""); err != nil {
					return err
				}
				_ = deleteRefreshToken()
				out.Info("logged out (session revoked server-side)")
				return nil
			},
		},
		&cobra.Command{
			Use:   "status",
			Short: "Show who you are, which tenant, and the auth source",
			RunE: func(cmd *cobra.Command, args []string) error {
				out := NewOutput(flagJSON, flagQuiet, flagNoColor)
				client, cfg, ok := configClient()
				if !ok {
					return exit(ExitNotAuthed, "not authenticated — run 'bloberry auth login' or set BLOBERRY_TOKEN")
				}
				env, _, err := client.Do("GET", "/users/me", nil, "")
				if err != nil {
					return err
				}
				var me map[string]interface{}
				_ = env.JSON(&me)
				_, envOK, _ := getToken()
				source := "keychain"
				if envOK {
					source = "BLOBERRY_TOKEN"
				}
				if flagJSON {
					return out.Data(map[string]interface{}{"user": me, "tenant": cfg.Tenant, "source": source})
				}
				out.Data(map[string]interface{}{"user": me["email"], "tenant": cfg.Tenant, "source": source})
				return nil
			},
		},
	)
	return cmd
}

var _ = json.Marshal
