VERSION ?= dev
GIT_COMMIT ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo unknown)
BUILD_DATE ?= $(shell date -u +%Y-%m-%dT%H:%M:%SZ 2>/dev/null || echo unknown)
LDFLAGS = -X github.com/ovitorvalente/git-loom/internal/cli.Version=$(VERSION) -X github.com/ovitorvalente/git-loom/internal/cli.GitCommit=$(GIT_COMMIT) -X github.com/ovitorvalente/git-loom/internal/cli.BuildDate=$(BUILD_DATE)

test:
	go test ./...

vet:
	go vet ./...

build:
	go build -ldflags "$(LDFLAGS)" ./cmd/gitloom
