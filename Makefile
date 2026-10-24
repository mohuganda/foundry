BINARY := bin/foundry
VERSION ?= dev
LDFLAGS := -s -w -X github.com/jabahum/foundry/internal/cli.Version=$(VERSION)
.DEFAULT_GOAL := help

.PHONY: help build test format format-check lint check ci smoke install release clean fmt fmt-check vet release-snapshot

help: ## Show commands; see docs/commands.md for prerequisites and side effects.
	@awk 'BEGIN {FS = ":.*## "; print "Foundry commands:"} /^[a-zA-Z_-]+:.*## / {printf "  make %-18s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

build: ## Build bin/foundry for this machine (VERSION defaults to dev).
	mkdir -p bin
	CGO_ENABLED=0 go build -trimpath -ldflags "$(LDFLAGS)" -o $(BINARY) ./cmd/foundry

test: ## Run Foundry tests, including embedded generation for all five templates.
	go test ./...

format: ## Format Foundry's Go source; does not format generated template projects.
	gofmt -w assets.go cmd internal

format-check: ## Check Foundry formatting without rewriting source.
	test -z "$$(gofmt -l assets.go cmd internal)"

lint: ## Run go vet on Foundry packages.
	go vet ./...

check: format-check lint test build ## Check formatting, lint, tests, and the native build.

ci: check smoke ## Run Foundry checks and the isolated generation/add smoke test.

smoke: build ## Generate all five templates and exercise Go add commands in a temporary folder.
	./scripts/smoke.sh

install: ## Install Foundry into GOBIN or GOPATH/bin.
	go install -trimpath -ldflags "$(LDFLAGS)" ./cmd/foundry

release: ## Build six executable archives and checksums; does not publish (VERSION=dev).
	./scripts/build-release.sh "$(VERSION)"

clean: ## Remove local bin/ and dist/ artifacts.
	rm -rf bin dist

# Compatibility aliases. Prefer the canonical names above in new automation.
fmt: format
fmt-check: format-check
vet: lint
release-snapshot: release
