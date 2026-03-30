VERSION ?= dev
GIT_COMMIT ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo unknown)
BUILD_DATE ?= $(shell date -u +%Y-%m-%dT%H:%M:%SZ 2>/dev/null || echo unknown)
DIST_DIR ?= dist
BINARY_NAME ?= gitloom
PLATFORMS ?= darwin/amd64 darwin/arm64 linux/amd64 linux/arm64 windows/amd64 windows/arm64
GO_CACHE_DIR ?= .cache/go-build
LDFLAGS = -X github.com/ovitorvalente/git-loom/internal/cli.Version=$(VERSION) -X github.com/ovitorvalente/git-loom/internal/cli.GitCommit=$(GIT_COMMIT) -X github.com/ovitorvalente/git-loom/internal/cli.BuildDate=$(BUILD_DATE)

.PHONY: test vet build clean dist dist-checksums release-artifacts install uninstall

INSTALL_VERSION ?= latest
INSTALL_DIR ?=
INSTALL_MODE ?= auto

install:
	@if [ -n "$(INSTALL_DIR)" ]; then \
		curl -sSL https://raw.githubusercontent.com/ovitorvalente/git-loom/main/scripts/install.sh | bash -s -- -v $(INSTALL_VERSION) -d "$(INSTALL_DIR)" -m $(INSTALL_MODE); \
	else \
		curl -sSL https://raw.githubusercontent.com/ovitorvalente/git-loom/main/scripts/install.sh | bash -s -- -v $(INSTALL_VERSION) -m $(INSTALL_MODE); \
	fi

uninstall:
	@curl -sSL https://raw.githubusercontent.com/ovitorvalente/git-loom/main/scripts/install.sh | bash -s -- -u

install-local: build
	@mkdir -p ~/.local/bin 2>/dev/null || true
	@cp $(BINARY_NAME) ~/.local/bin/
	@chmod +x ~/.local/bin/$(BINARY_NAME)
	@echo "Installed to ~/.local/bin/$(BINARY_NAME)"

test:
	go test ./...

vet:
	go vet ./...

build:
	go build -buildvcs=false -ldflags "$(LDFLAGS)" -o $(BINARY_NAME) ./cmd/gitloom

clean:
	rm -rf $(DIST_DIR) $(BINARY_NAME) $(GO_CACHE_DIR)

dist: clean
	mkdir -p $(DIST_DIR)
	mkdir -p $(GO_CACHE_DIR)
	set -eu; \
	for platform in $(PLATFORMS); do \
		os=$${platform%/*}; \
		arch=$${platform#*/}; \
		artifact="$(BINARY_NAME)_$(VERSION)_$${os}_$${arch}"; \
		bin_name="$(BINARY_NAME)"; \
		if [ "$$os" = "windows" ]; then \
			bin_name="$(BINARY_NAME).exe"; \
		fi; \
		build_dir="$(DIST_DIR)/$$artifact"; \
		mkdir -p "$$build_dir"; \
		echo "building $$os/$$arch"; \
		GOCACHE="$(CURDIR)/$(GO_CACHE_DIR)" GOOS=$$os GOARCH=$$arch CGO_ENABLED=0 go build -buildvcs=false -trimpath -ldflags "$(LDFLAGS)" -o "$$build_dir/$$bin_name" ./cmd/gitloom; \
		cp README.md LICENSE "$$build_dir/"; \
		if [ "$$os" = "windows" ]; then \
			( cd "$(DIST_DIR)" && zip -qr "$$artifact.zip" "$$artifact" ); \
		else \
			tar -C "$(DIST_DIR)" -czf "$(DIST_DIR)/$$artifact.tar.gz" "$$artifact"; \
		fi; \
		rm -rf "$$build_dir"; \
	done
	$(MAKE) dist-checksums

dist-checksums:
	set -eu; \
	cd $(DIST_DIR); \
	if command -v sha256sum >/dev/null 2>&1; then \
		sha256sum * > checksums.txt; \
	else \
		shasum -a 256 * > checksums.txt; \
	fi

release-artifacts: dist
