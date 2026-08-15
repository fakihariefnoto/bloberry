package cli

import (
	"fmt"
	"io"
	"net/url"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

// parsePath splits a bloberry:// path into (remote, remainder).
func parseRemote(p string) (string, bool) {
	return strings.TrimPrefix(p, "bloberry://"), strings.HasPrefix(p, "bloberry://")
}

func newLSCmd() *cobra.Command {
	var long bool
	cmd := &cobra.Command{
		Use:   "ls [path]",
		Short: "List a folder's contents",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			out := NewOutput(flagJSON, flagQuiet, flagNoColor)
			client, _, ok := configClient()
			if !ok {
				return exit(ExitNotAuthed, "not authenticated")
			}
			path := "root"
			if len(args) > 0 {
				rp, isRemote := parseRemote(args[0])
				if !isRemote {
					// local path → resolve tenant root folder if it's a bare path
					rp = strings.Trim(args[0], "/")
				}
				if rp != "" {
					path = rp
				}
			}
			// Resolve folder by name: v1 resolves "root" or the tenant root.
			env, _, err := client.Do("GET", "/folders/"+url.PathEscape(path)+"/children", nil, "")
			if err != nil {
				return err
			}
			var listing struct {
				Folders []struct {
					ID   string `json:"id"`
					Name string `json:"name"`
				} `json:"folders"`
				Objects []struct {
					ID         string `json:"id"`
					Name       string `json:"name"`
					SizeBytes  int64  `json:"size_bytes"`
					Visibility string `json:"visibility"`
				} `json:"objects"`
			}
			if err := env.JSON(&listing); err != nil {
				return err
			}
			if flagJSON {
				return out.Data(listing)
			}
			for _, f := range listing.Folders {
				fmt.Fprintf(out.Stdout, "%-40s  %s\n", f.Name+"/", "folder")
			}
			for _, o := range listing.Objects {
				if long {
					fmt.Fprintf(out.Stdout, "%-40s  %10d  %s\n", o.Name, o.SizeBytes, o.Visibility)
				} else {
					fmt.Fprintln(out.Stdout, o.Name)
				}
			}
			return nil
		},
	}
	cmd.Flags().BoolVarP(&long, "long", "l", false, "long listing")
	cmd.Flags().Bool("recursive", false, "recurse into folders")
	return cmd
}

func newStatCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "stat <path>",
		Short: "Show metadata for an object or folder",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			out := NewOutput(flagJSON, flagQuiet, flagNoColor)
			client, _, ok := configClient()
			if !ok {
				return exit(ExitNotAuthed, "not authenticated")
			}
			rp, _ := parseRemote(args[0])
			// v1: stat by object id if the path looks like an id, else 404.
			env, _, err := client.Do("GET", "/objects/"+url.PathEscape(rp), nil, "")
			if err != nil {
				return err
			}
			var obj map[string]interface{}
			if err := env.JSON(&obj); err != nil {
				return err
			}
			return out.Data(obj)
		},
	}
}

func newCatCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "cat <path>",
		Short: "Stream an object to stdout",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, _, ok := configClient()
			if !ok {
				return exit(ExitNotAuthed, "not authenticated")
			}
			rp, _ := parseRemote(args[0])
			resp, err := client.HTTP.Get(client.BaseURL + "/objects/" + url.PathEscape(rp) + "/download")
			if err != nil {
				return err
			}
			defer resp.Body.Close()
			if resp.StatusCode >= 400 {
				return exit(ExitNotFound, "object not found")
			}
			_, err = io.Copy(os.Stdout, resp.Body)
			return err
		},
	}
}

func newCPSend() *cobra.Command {
	var recursive, yes bool
	cmd := &cobra.Command{
		Use:   "cp <src> <dst>",
		Short: "Copy objects or folders (up, down, or remote→remote)",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			out := NewOutput(flagJSON, flagQuiet, flagNoColor)
			client, _, ok := configClient()
			if !ok {
				return exit(ExitNotAuthed, "not authenticated")
			}
			src, srcRemote := parseRemote(args[0])
			dst, dstRemote := parseRemote(args[1])

			// v1 supports upload (local → remote) via presigned PUT.
			if !srcRemote && dstRemote {
				// determine target folder + name
				parts := strings.Split(strings.Trim(dst, "/"), "/")
				name := parts[len(parts)-1]
				folder := "root"
				if len(parts) > 1 {
					folder = strings.Join(parts[:len(parts)-1], "/")
				}
				// presign
				fi, err := os.Stat(src)
				if err != nil {
					return err
				}
				var env *Envelope
				var status int
				var serr error
				if recursive && fi.IsDir() {
					return exit(ExitInvocation, "directory copy via -r not yet wired; upload files individually")
				}
				body := fmt.Sprintf(`{"folder_id":%q,"name":%q,"size":%d}`, folder, name, fi.Size())
				env, status, serr = client.Do("POST", "/objects/presign-put", strings.NewReader(body), "application/json")
				if serr != nil {
					return serr
				}
				_ = status
				var presign struct {
					FileID    string            `json:"file_id"`
					UploadURL string            `json:"upload_url"`
					Headers   map[string]string `json:"headers"`
				}
				if err := env.JSON(&presign); err != nil {
					return err
				}
				// PUT bytes
				f, err := os.Open(src)
				if err != nil {
					return err
				}
				defer f.Close()
				req, err := httpNewRequest("PUT", presign.UploadURL, f)
				if err != nil {
					return err
				}
				for k, v := range presign.Headers {
					req.Header.Set(k, v)
				}
				resp, err := client.HTTP.Do(req)
				if err != nil {
					return exit(ExitBackendDown, "upload failed: %v", err)
				}
				resp.Body.Close()
				if resp.StatusCode >= 300 {
					return exit(ExitFailed, "upload failed with status %d", resp.StatusCode)
				}
				// complete
				if _, _, err := client.Do("POST", "/objects/"+presign.FileID+"/complete", strings.NewReader(`{"etag":"x"}`), "application/json"); err != nil {
					return err
				}
				out.Info("copied %s → %s", args[0], args[1])
				return nil
			}
			return exit(ExitInvocation, "unsupported copy direction (v1 supports local → bloberry://)")
		},
	}
	cmd.Flags().BoolVarP(&recursive, "recursive", "r", false, "copy recursively")
	cmd.Flags().BoolVarP(&yes, "yes", "y", false, "skip confirmation")
	return cmd
}

func newRMCmd() *cobra.Command {
	var recursive, yes bool
	cmd := &cobra.Command{
		Use:   "rm <path>",
		Short: "Delete an object or folder",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, _, ok := configClient()
			if !ok {
				return exit(ExitNotAuthed, "not authenticated")
			}
			rp, _ := parseRemote(args[0])
			if !yes {
				return exit(ExitInvocation, "refusing without confirmation; pass --yes (typed-name confirmation is a web-dashboard action)")
			}
			_, _, err := client.Do("DELETE", "/objects/"+url.PathEscape(rp), nil, "")
			return err
		},
	}
	cmd.Flags().BoolVarP(&recursive, "recursive", "r", false, "delete recursively")
	cmd.Flags().BoolVarP(&yes, "yes", "y", false, "skip confirmation")
	return cmd
}
