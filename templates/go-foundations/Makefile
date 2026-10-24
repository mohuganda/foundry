DATABASE_URL ?= postgres://grpc:grpc@localhost:5432/go_foundation?sslmode=disable
BINARY ?= bin/server
.DEFAULT_GOAL := help

.PHONY: help tools codegen contracts-check format format-check lint test build dev check ci clean keys docker-up docker-down logs migrate-up migrate-down migrate-create tidy proto fmt run

help: ## Show project commands and their purpose.
	@awk 'BEGIN {FS = ":.*## "; print "Go foundation commands:"} /^[a-zA-Z_-]+:.*## / {printf "  make %-18s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

tools: ## Install Buf, protobuf generators, gateway generators, and grpcurl into GOBIN/GOPATH/bin.
	go install github.com/bufbuild/buf/cmd/buf@v1.47.2
	go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.36.1
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.5.1
	go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@v2.23.0
	go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@v2.23.0
	go install github.com/fullstorydev/grpcurl/cmd/grpcurl@latest

codegen: ## Regenerate protobuf, gRPC, gateway, and OpenAPI outputs; requires tools.
	buf generate

contracts-check: ## Lint protobuf sources without regenerating them; requires Buf.
	buf lint

format: ## Format application and generated Go source.
	gofmt -w ./cmd ./internal ./gen

format-check: ## Check Go formatting without rewriting source.
	test -z "$$(gofmt -l ./cmd ./internal ./gen)"

lint: ## Run go vet against the checked-in generated contracts.
	go vet ./...

test: ## Run all Go tests with the race detector; does not regenerate contracts.
	go test -race ./...

build: ## Build bin/server using existing generated contracts.
	mkdir -p bin
	go build -trimpath -o $(BINARY) ./cmd/server

dev: ## Run the server; requires configured PostgreSQL and development keys.
	go run ./cmd/server

check: format-check lint test build ## Check formatting, vet, tests, and build without changing contracts.

ci: ## Regenerate contracts, lint them, then run check; requires generation tools.
	$(MAKE) codegen
	$(MAKE) contracts-check
	$(MAKE) check

clean: ## Remove local binaries and coverage output; preserve generated contracts.
	rm -rf bin coverage.out

keys: ## Create missing development JWT keys in secrets/; preserve existing keys.
	./scripts/generate-dev-keys.sh

docker-up: keys ## Build and start Compose services, apply migrations, and initialize the API.
	docker compose up --build -d

docker-down: ## Stop Compose services; preserve the named PostgreSQL volume.
	docker compose down

logs: ## Follow the API container's logs until interrupted.
	docker compose logs -f api

migrate-up: ## Apply pending database migrations using DATABASE_URL.
	migrate -path migrations -database "$(DATABASE_URL)" up

migrate-down: ## Roll back one migration using DATABASE_URL; changes the database.
	migrate -path migrations -database "$(DATABASE_URL)" down 1

migrate-create: ## Create paired SQL migration files (name=create_something).
	@test -n "$(name)" || (echo "usage: make migrate-create name=create_something" && exit 1)
	migrate create -ext sql -dir migrations -seq $(name)

tidy: ## Update go.mod/go.sum to match imports; may download dependencies.
	go mod tidy

# Compatibility aliases.
proto: codegen
fmt: format
run: dev
