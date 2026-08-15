# Desktop Makefile — Wails (Go + web frontend)
#
# Reference template copied into apps/bloberry/desktop/Makefile by the
# plan-app skill when Wails is chosen. Assumes the standard `wails init`
# layout (Go backend at the repo root, frontend/ holding the web UI).

.PHONY: install dev build build-debug package-mac package-win package-linux test lint format clean help

install: ## Install Go + frontend deps
	go mod download
	cd frontend && pnpm install

dev: ## Run in dev mode (hot-reloads frontend + rebuilds Go on change)
	wails dev

build: ## Production build for the current platform (see templates/desktop-defaults.md → Packaging for signing)
	wails build

build-debug: ## Debug build with dev tools enabled
	wails build -debug

# Installer targets — Wails builds binaries, not installers, so each OS needs an
# extra tool: NSIS on Windows, create-dmg on macOS, nfpm on Linux. Keep these in
# sync with desktop/README.md's packaging table and delete unshipped OSes.
# See templates/desktop-packaging.md. (Wails v3 uses `wails3 task <os>:package`
# instead — replace these targets wholesale if the app pins v3.)

package-mac: ## macOS universal .app, then wrap into a .dmg (needs create-dmg; codesign+notarize separately)
	wails build -platform darwin/universal
	create-dmg --volname "bloberry" build/bin/bloberry.dmg build/bin/bloberry.app

package-win: ## Windows NSIS installer (requires NSIS installed: winget install NSIS / scoop install nsis)
	@command -v makensis >/dev/null || { echo "makensis not found — the -nsis flag will silently emit only a bare .exe"; exit 1; }
	wails build -platform windows/amd64 -nsis

package-linux: ## Linux binary + .deb/.rpm via nfpm (nfpm.yaml must install the .desktop entry, icon and MIME files)
	wails build -platform linux/amd64
	nfpm pkg --packager deb --target build/bin/
	nfpm pkg --packager rpm --target build/bin/

test: ## Run Go-side tests
	go test ./... -v -count=1

lint: ## Lint Go + frontend
	golangci-lint run ./...
	cd frontend && pnpm lint

format: ## Format Go + frontend
	gofmt -w .
	cd frontend && pnpm exec prettier --write .

clean: ## Clean build artifacts
	rm -rf build/bin

help: ## Show this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

.DEFAULT_GOAL := help
