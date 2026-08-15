.PHONY: all build generate openapi web test lint security mocks tools run dev clean

ORG        := github.com/fakihariefnoto/bloberry
BIN        := bin
SERVER     := $(BIN)/bloberry-server
CLI        := $(BIN)/bloberry

all: build

# Build order (architecture.md §7 edge 2): openapi → web build → go build.
build: generate web go-build

generate: openapi
	$(MAKE) -C . generate-server

openapi:
	@test -f api/openapi.yaml || (echo "api/openapi.yaml missing"; exit 1)

generate-server: api/openapi.yaml
	oapi-codegen -config api/oapi-server.yaml api/openapi.yaml
	oapi-codegen -config api/oapi-go-sdk.yaml api/openapi.yaml
	goimports -w internal/platform/api/server.gen.go sdk/go/client.gen.go

web:
	$(MAKE) -C web build
	rm -rf internal/platform/web/static
	cp -r web/dist internal/platform/web/static

go-build:
	go build -o $(SERVER) ./cmd/bloberry-server
	go build -o $(CLI) ./cmd/bloberry

# Dev infra + app (01-setup smoke run).
dev:
	docker compose -f deploy/docker-compose.yml up -d mongo redis
	go run ./cmd/bloberry-server

test:
	go test ./...

lint:
	golangci-lint run ./...

security:
	gosec ./...

mocks:
	mockgen -destination internal/authz/mocks/mocks.go github.com/fakihariefnoto/bloberry/internal/authz,internal/user Repository,Usecase

tools:
	go install github.com/go-delve/delve/cmd/dlv@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install github.com/securego/gosec/v2/cmd/gosec@latest
	go install go.uber.org/mock/mockgen@latest
	go install github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@latest

clean:
	rm -rf $(BIN)
	$(MAKE) -C web clean
