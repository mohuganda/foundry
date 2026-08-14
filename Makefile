BINARY := bin/foundry

.PHONY: build test fmt vet smoke clean install

build:
	mkdir -p bin
	go build -trimpath -o $(BINARY) ./cmd/foundry

test:
	go test ./...

fmt:
	gofmt -w cmd internal

vet:
	go vet ./...

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
	go install ./cmd/foundry

clean:
	rm -rf bin dist .tmp-smoke
