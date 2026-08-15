package main

import (
	"embed"
	"fmt"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/wailsapp/wails/v3/pkg/application"
	"gopkg.in/yaml.v3"
)

//go:embed all:web
var webFS embed.FS

var app *application.App

// DesktopService exposes native-only capabilities to the embedded web UI:
// file dialogs, folder drag-drop target, and the runtime API base (desktop/README.md §Sub-choices).
type DesktopService struct{}

func (s *DesktopService) GetServerURL() string {
	return serverURL()
}

func (s *DesktopService) SetServerURL(url string) error {
	cfg, err := loadConfig()
	if err != nil {
		return err
	}
	cfg.Server = strings.TrimSuffix(url, "/")
	return saveConfig(cfg)
}

func (s *DesktopService) OpenFileDialog() (string, error) {
	dialog := app.Dialog.OpenFile()
	return dialog.PromptForSingleSelection()
}

func (s *DesktopService) Version() string {
	return "0.1.0-desktop"
}

func main() {
	web, err := webFS.Open("web")
	if err != nil {
		// web/dist not built — fail loudly, not a blank window.
		fmt.Fprintln(os.Stderr, "desktop: web/dist missing — run `make web` first")
		os.Exit(1)
	}
	web.Close()

	var assets fs.FS
	sub, err := fs.Sub(webFS, "web")
	if err != nil {
		fmt.Fprintln(os.Stderr, "desktop: bad embed:", err)
		os.Exit(1)
	}
	assets = sub

	app = application.New(application.Options{
		Name:        "Bloberry",
		Description: "Bloberry desktop — storage-agnostic object service",
		Assets: application.AssetOptions{
			Handler: application.BundledAssetFileServer(assets),
			Middleware: application.ChainMiddleware(
				injectAPIConfig,
			),
		},
		Services: []application.Service{
			application.NewService(&DesktopService{}),
		},
	})

	// Native chrome: the one surface with no web equivalent is the first-run
	// server-URL screen, handled by the web app reading __BLOBERRY_API_BASE__.
	app.Window.NewWithOptions(application.WebviewWindowOptions{
		Title:    "Bloberry",
		Width:    1280,
		Height:   800,
		MinWidth: 960,
		MinHeight: 600,
		URL:      "/",
	})

	// System tray: uploads and sync continue after the window closes
	// (desktop/README.md §Window, menu, shortcuts — tray = yes).
	tray := app.SystemTray.New()
	tray.SetLabel("Bloberry")
	tray.SetTooltip("Bloberry — click to open")

	app.Run()
}

// injectAPIConfig sets window.__BLOBERRY_API_BASE__ before the web app boots,
// so the SPA's runtime API-base hook points at the configured server
// (web/src/lib/api.ts, web/components.md).
func injectAPIConfig(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "index.html") || r.URL.Path == "/" {
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			fmt.Fprintf(w, `<script>window.__BLOBERRY_API_BASE__=%q;</script>`, serverURL())
			// After the injected script, serve the real index.html.
			next.ServeHTTP(w, r)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// --- config (mirror of the CLI's, stored per-user) ---

func configPath() string {
	if p := os.Getenv("BLOBERRY_CONFIG"); p != "" {
		return p
	}
	home, err := os.UserHomeDir()
	if err != nil {
		home = "."
	}
	return filepath.Join(home, ".config", "bloberry", "desktop.yaml")
}

type config struct {
	Server string `yaml:"server"`
}

func loadConfig() (config, error) {
	var cfg config
	data, err := os.ReadFile(configPath())
	if err != nil {
		if os.IsNotExist(err) {
			return cfg, nil
		}
		return cfg, err
	}
	_ = yamlUnmarshal(data, &cfg)
	return cfg, nil
}

func saveConfig(cfg config) error {
	if err := os.MkdirAll(filepath.Dir(configPath()), 0o700); err != nil {
		return err
	}
	return yamlWrite(cfg)
}

func serverURL() string {
	cfg, err := loadConfig()
	if err != nil || cfg.Server == "" {
		return "http://localhost:8080"
	}
	return cfg.Server
}

func yamlUnmarshal(data []byte, v interface{}) error { return yaml.Unmarshal(data, v) }
func yamlWrite(v interface{}) error {
	data, err := yaml.Marshal(v)
	if err != nil {
		return err
	}
	return os.WriteFile(configPath(), data, 0o600)
}
