# Backend Makefile — Go
#
# Reference template copied into apps/bloberry/backend/Makefile by the
# plan-app skill when Go is chosen. Assumes a cmd/server entrypoint
# (standard Go project layout) — adjust paths if the app differs.

.PHONY: build run debug test test-coverage lint security mocks swagger docker-build docker-up docker-down run-all tools clean help

APP_NAME := bloberry
VERSION  := $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
LDFLAGS  := -ldflags "-X main.version=$(VERSION)"

build: ## Build binary
	go build $(LDFLAGS) -o bin/$(APP_NAME) ./cmd/server

run: ## Run in dev mode
	go run ./cmd/server

debug: ## Run with delve debugger attached (dlv debug ./cmd/server)
	dlv debug ./cmd/server

test: ## Run tests
	go test ./... -v -count=1

test-coverage: ## Run tests with HTML coverage report
	go test ./... -coverprofile=coverage.out -covermode=atomic
	go tool cover -html=coverage.out -o coverage.html

lint: ## Run linter (includes gosec as one of the bundled linters)
	golangci-lint run ./...

security: ## Run gosec standalone (dedicated security scan, see templates/backend-go-defaults.md)
	gosec ./...

mocks: ## Generate mocks for handler/usecase/repository interfaces (see templates/backend-go-defaults.md naming convention)
	@find internal/ -name "handler.go" -o -name "usecase.go" -o -name "repository.go" | while read f; do \
		PKG=$$(dirname "$$f" | xargs basename); \
		DIR=$$(dirname "$$f"); \
		MOCK_DIR="$${DIR}/mocks"; \
		mkdir -p "$${MOCK_DIR}"; \
		echo "Generating mock for $$f -> $${MOCK_DIR}/$${PKG}_mock.go"; \
		mockgen -source="$$f" -destination="$${MOCK_DIR}/$${PKG}_mock.go" -package=mocks 2>/dev/null || true; \
	done

swagger: ## Generate OpenAPI spec from swag annotations (commit the output — see templates/backend-go-defaults.md)
	swag init -g cmd/server/main.go -o docs/swagger --parseDependency --parseInternal

docker-build: ## Build docker image
	docker build -t $(APP_NAME):$(VERSION) -f deployments/Dockerfile .

docker-up: ## Start docker compose stack (app + infra deps)
	docker-compose -f deployments/docker-compose.yml up -d

docker-down: ## Stop docker compose stack
	docker-compose -f deployments/docker-compose.yml down

run-all: docker-up run ## Start infra deps, then run the app

tools: ## Install dev tools (delve, golangci-lint, gosec, mockgen, swag)
	go install github.com/go-delve/delve/cmd/dlv@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install github.com/securego/gosec/v2/cmd/gosec@latest
	go install go.uber.org/mock/mockgen@latest
	go install github.com/swaggo/swag/cmd/swag@latest

clean: ## Clean build artifacts
	rm -rf bin/ coverage.out coverage.html

help: ## Show this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

.DEFAULT_GOAL := help
