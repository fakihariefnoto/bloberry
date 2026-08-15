# Web Makefile — Vue / Next / React / Svelte / htmx
#
# Reference template copied into apps/bloberry/web/Makefile by the
# plan-app skill. Framework-agnostic: it just wraps the package.json
# scripts every framework choice defines (dev/build/start/lint/test).
#
# Runtime / package manager: Bun (see ../TRD.md Stack decisions). The commands below
# are already converted from the template default of Node + pnpm.

.PHONY: install dev build start debug test lint format clean help

install: ## Install dependencies
	bun install

dev: ## Run dev server
	bun run dev

debug: ## Run dev server with the Bun inspector attached (breakpoint debugging)
	bun --inspect run dev

build: ## Production build
	bun run build

start: ## Run the production build
	bun run start

test: ## Run tests
	bun run test

lint: ## Lint
	bun run lint

format: ## Format with Prettier
	bunx prettier --write .

clean: ## Clean build artifacts and caches
	rm -rf .next dist build .svelte-kit node_modules/.cache

help: ## Show this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

.DEFAULT_GOAL := help
