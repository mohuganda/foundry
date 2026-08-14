BINARY := bin/foundry
VERSION ?= dev
LDFLAGS := -s -w -X github.com/jabahum/foundry/internal/cli.Version=$(VERSION)

.PHONY: build test fmt fmt-check vet check smoke install release-snapshot clean

build:
	mkdir -p bin
	CGO_ENABLED=0 go build -trimpath -ldflags "$(LDFLAGS)" -o $(BINARY) ./cmd/foundry

test:
	go test ./...

fmt:
	gofmt -w cmd internal

fmt-check:
	test -z "$$(gofmt -l .)"

vet:
	go vet ./...

check: fmt-check test vet

smoke: build
	./$(BINARY) version
	./$(BINARY) list
	rm -rf .tmp-smoke
	./$(BINARY) new .tmp-smoke --template go --module example.com/smoke --no-hooks
	test -f .tmp-smoke/go.mod
	grep -q 'module example.com/smoke' .tmp-smoke/go.mod
	./$(BINARY) add migration create_smoke --project .tmp-smoke
	test -f .tmp-smoke/migrations/000004_create_smoke.up.sql
	./$(BINARY) add service smoke --project .tmp-smoke
	test -f .tmp-smoke/api/proto/smoke/v1/smoke.proto
	rm -rf .tmp-smoke

install:
	go install -trimpath -ldflags "$(LDFLAGS)" ./cmd/foundry

release-snapshot:
	./scripts/build-release.sh "$(VERSION)"

clean:
	rm -rf bin dist .tmp-smoke
