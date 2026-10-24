# go-foundation

A production-oriented Go foundation using gRPC, Protocol Buffers, Clean Architecture, PostgreSQL, local RBAC, audit trails, local JWT authentication, optional OIDC/Keycloak authentication, Prometheus, OpenTelemetry, Jaeger, Docker Compose, migrations, pagination, health checks, graceful shutdown, structured logging, and tests.

## Architecture

```text
Client
  -> gRPC / Protobuf or HTTP / JSON gateway
  -> request ID
  -> authentication (local JWT or optional OIDC)
  -> local RBAC permission lookup
  -> idempotency reservation / replay
  -> audit capture for mutations
  -> metrics / logging / recovery
  -> transport handler
  -> application use case
  -> domain
  -> repository port
  -> PostgreSQL
```

Authorization is always local. Enabling Keycloak/OIDC only changes how identity is established; roles and permissions remain in the starter database.

## Requirements

- Go 1.23+
- Docker + Docker Compose
- OpenSSL
- Optional local tooling: Buf, grpcurl, golang-migrate

Install developer tools:

```bash
make tools
```

## Quick start

The Compose configuration contains development credentials only. Change them before using this outside local development.

```bash
make docker-up
make logs
```

Services:

- gRPC: `localhost:50051`
- HTTP/JSON gateway: `localhost:8080`
- interactive API documentation: `localhost:8080/docs/`
- application metrics: `localhost:9090/metrics`
- Prometheus: `localhost:9091`
- Jaeger: `localhost:16686`
- Grafana: `localhost:3000` (`admin` / `admin`)
- PostgreSQL: `localhost:5432`

The development bootstrap administrator is:

```text
admin@example.com
ChangeMe123!
```

## Login

```bash
grpcurl -plaintext \
  -d '{"email":"admin@example.com","password":"ChangeMe123!"}' \
  localhost:50051 \
  auth.v1.AuthService/Login
```

Copy the returned token:

```bash
export ACCESS_TOKEN='...'
export REFRESH_TOKEN='...'
```

The same API is available over HTTP/JSON:

```bash
curl -X POST http://localhost:8080/v1/auth/login \
  -H 'Content-Type: application/json' \
  -d '{"email":"admin@example.com","password":"ChangeMe123!"}'
```

Inspect the current identity and effective RBAC:

```bash
grpcurl -plaintext \
  -H "Authorization: Bearer ${ACCESS_TOKEN}" \
  localhost:50051 \
  auth.v1.AuthService/Me
```

## Refresh tokens and sessions

Local login returns a short-lived JWT access token and an opaque refresh token. Refresh tokens are shown only in login and refresh responses; PostgreSQL stores only their SHA-256 hashes. A successful refresh rotates the token atomically, so the previous value cannot be used again. Sessions have an absolute lifetime controlled by `JWT_REFRESH_TOKEN_TTL` (default `720h`).

```bash
grpcurl -plaintext \
  -d "{\"refreshToken\":\"${REFRESH_TOKEN}\"}" \
  localhost:50051 \
  auth.v1.AuthService/Refresh
```

Replace both locally stored tokens with the values from every successful refresh response. To log out, revoke the session using its current refresh token:

```bash
grpcurl -plaintext \
  -d "{\"refreshToken\":\"${REFRESH_TOKEN}\"}" \
  localhost:50051 \
  auth.v1.AuthService/Logout
```

Authenticated users can inspect and revoke their own active sessions:

```bash
grpcurl -plaintext \
  -H "Authorization: Bearer ${ACCESS_TOKEN}" \
  localhost:50051 \
  auth.v1.AuthService/ListSessions

grpcurl -plaintext \
  -H "Authorization: Bearer ${ACCESS_TOKEN}" \
  -d '{"sessionId":"SESSION_UUID"}' \
  localhost:50051 \
  auth.v1.AuthService/RevokeSession
```

Revocation prevents further refreshes. Existing stateless access tokens remain valid until their short expiry, so keep `JWT_ACCESS_TOKEN_TTL` appropriately brief. Clients should keep refresh tokens in platform-secure storage and never local browser storage.

## Idempotency

Authenticated mutation RPCs accept an optional `idempotency-key` gRPC metadata value or `Idempotency-Key` HTTP header. The key must contain 8-255 letters, digits, `.`, `_`, `:`, or `-`. Keys are scoped to the authenticated user and RPC method, and only their SHA-256 hashes are stored.

```bash
curl -X POST http://localhost:8080/v1/users \
  -H "Authorization: Bearer ${ACCESS_TOKEN}" \
  -H "Idempotency-Key: $(uuidgen)" \
  -H 'Content-Type: application/json' \
  -d '{"name":"Jeremy","email":"jeremy@example.com","password":"StrongPass123!"}'
```

The first request reserves the key and stores its protobuf response. Repeating the same method, key, actor, and payload replays that response without executing the handler or writing a second audit event. HTTP responses expose `Idempotency-Replayed: true|false`; gRPC clients receive the equivalent response metadata.

- Reusing a key with a different payload returns `IDEMPOTENCY_KEY_REUSED` (`AlreadyExists`, HTTP 409).
- Concurrent requests using the same key return `IDEMPOTENCY_REQUEST_IN_PROGRESS` with retry metadata.
- Transient server failures release the reservation so the request can be retried.
- Completed records expire after `IDEMPOTENCY_TTL` (default `24h`) and are cleaned incrementally.
- Abandoned in-progress reservations can be taken over after `IDEMPOTENCY_LOCK_TIMEOUT` (default `30s`); configure it above the maximum expected mutation duration.

Idempotency covers user creation, RBAC relationship mutations, and authenticated session revocation. Login, refresh, and logout are intentionally excluded because safely replaying them would require persisting raw authentication tokens. Set `IDEMPOTENCY_ENABLED=false` to bypass the feature. This mechanism prevents ordinary duplicate execution but is not a substitute for a database transaction spanning both the business mutation and idempotency record.

## Create a user

```bash
grpcurl -plaintext \
  -H "Authorization: Bearer ${ACCESS_TOKEN}" \
  -d '{"name":"Jeremy","email":"jeremy@example.com","password":"StrongPass123!"}' \
  localhost:50051 \
  user.v1.UserService/CreateUser
```

## List users with pagination

```bash
grpcurl -plaintext \
  -H "Authorization: Bearer ${ACCESS_TOKEN}" \
  -d '{"pageSize":20,"search":"jeremy"}' \
  localhost:50051 \
  user.v1.UserService/ListUsers
```

Pass the returned `nextPageToken` into the next request.

## Manage RBAC

List roles:

```bash
grpcurl -plaintext \
  -H "Authorization: Bearer ${ACCESS_TOKEN}" \
  localhost:50051 \
  rbac.v1.RBACService/ListRoles
```

List permissions:

```bash
grpcurl -plaintext \
  -H "Authorization: Bearer ${ACCESS_TOKEN}" \
  localhost:50051 \
  rbac.v1.RBACService/ListPermissions
```

Assign a role using IDs returned from the APIs:

```bash
grpcurl -plaintext \
  -H "Authorization: Bearer ${ACCESS_TOKEN}" \
  -d '{"userId":"USER_UUID","roleId":"ROLE_UUID"}' \
  localhost:50051 \
  rbac.v1.RBACService/AssignRole
```

The gRPC authorization policy checks permissions, never hard-coded role names.

## Audit trail

Authorized user creation and RBAC mutations are recorded in the append-only `audit_events` table. Each event contains the actor, action, target, request ID, RPC method, gRPC outcome, structured error reason, client address, user agent, timestamp, and explicitly allow-listed metadata. Complete request bodies are never serialized, so passwords and access tokens cannot enter the audit trail.

Only callers with `audit:read` can query the trail; the default `super_admin` role receives that permission. Results use cursor pagination and can be filtered by actor, action, resource type, or resource ID:

```bash
grpcurl -plaintext \
  -H "Authorization: Bearer ${ACCESS_TOKEN}" \
  -d '{"pageSize":50,"action":"role.assign"}' \
  localhost:50051 \
  audit.v1.AuditService/ListAuditEvents
```

The equivalent HTTP endpoint is `GET /v1/audit-events`. Audit persistence failures are logged with the request ID but do not rewrite an already-completed business operation into a client-visible failure.

## Default permissions

- `users:read`
- `users:create`
- `users:update`
- `users:disable`
- `users:delete`
- `roles:read`
- `roles:manage`
- `roles:assign`
- `audit:read`

Default roles are `super_admin`, `user_admin`, and `viewer`.

## Optional OIDC / Keycloak

Local authentication is enabled by default. To additionally accept OIDC access tokens:

```dotenv
OIDC_ENABLED=true
OIDC_ISSUER_URL=https://auth.example.org/realms/example
OIDC_CLIENT_ID=grpc-api
```

External identities must be linked to a local user so local RBAC remains authoritative. The relevant fields are `users.auth_provider='oidc'` and `users.external_subject=<OIDC sub>`.

Example linking an existing user:

```sql
UPDATE users
SET auth_provider = 'oidc',
    external_subject = 'OIDC-SUBJECT-HERE',
    password_hash = NULL
WHERE email = 'person@example.com';
```

If you want both local and external login for the same human, model identities in a separate `user_identities` table instead of replacing these fields; that is intentionally left as the first extension point for applications that need account linking.

## Protobuf

```bash
make proto
```

Contracts are versioned under:

```text
api/proto/auth/v1
api/proto/audit/v1
api/proto/user/v1
api/proto/rbac/v1
```

Generated Go and gRPC-Gateway code goes under `gen/proto`. The generated OpenAPI document is embedded into the server from `internal/transport/http/docs/go-foundation.swagger.json`; generated files are ignored by Git.

HTTP routes are declared with `google.api.http` annotations in the protobuf service definitions. Requests sent to the gateway on port `8080` pass through the existing gRPC authentication and authorization interceptors. For example:

```bash
curl http://localhost:8080/v1/users \
  -H "Authorization: Bearer ${ACCESS_TOKEN}"
```

Swagger UI and its pinned assets are embedded in the application binary and served at `/docs/`; the underlying specification is available at `/openapi.json`. Documentation metadata and security declarations live in `api/openapi.yaml`, keeping presentation-only OpenAPI options out of the protobuf contracts. Set `DOCS_ENABLED=false` to disable both documentation endpoints, such as in a production deployment where the API surface should not be published.

## Request validation and errors

Request validation rules are declared alongside fields in the protobuf contracts using Protovalidate. A unary interceptor applies them consistently to native gRPC and HTTP/JSON requests before authorization and handler execution.

Errors use `google.rpc.Status` with structured details:

- `google.rpc.ErrorInfo` provides a stable reason, the `go-foundation` domain, and request ID metadata.
- `google.rpc.BadRequest` identifies every invalid field and its validation failure.
- `google.rpc.ResourceInfo` identifies resources involved in not-found and conflict errors.

The gRPC-Gateway renders the same details as JSON, so clients can branch on `ErrorInfo.reason` instead of matching human-readable messages. Unexpected Go errors are converted to a generic internal error and are never returned verbatim.

## Database migrations

```bash
make migrate-up
make migrate-down
make migrate-create name=add_something
```

## Testing

```bash
make test
```

Unit tests use the in-memory repository where appropriate. PostgreSQL integration tests can be added under `tests/integration` without changing application code.

## Contract quality

Generated protobuf, gRPC-Gateway, and OpenAPI outputs are included in this template. Run the contract and application checks locally with:

```bash
buf lint
buf generate
go mod tidy
go test -race ./...
go vet ./...
docker build -t go-foundation:local .
```

After regenerating contracts, review the generated changes before committing them to your application's repository.

## Foundry

This template is maintained in Foundry's `templates/go-foundations` folder. Its Go module uses `go-foundations` as the placeholder that Foundry replaces with the application's module path during generation.

Foundry embeds this folder directly. Go source files use `.go.tmpl` and the module manifest uses `go.mod.tmpl` so the template is separate from Foundry's own Go module. Generation removes the `.tmpl` suffix and substitutes the application's module path.

Edit these source templates and rebuild Foundry to deliver changes. To run this foundation's checks, generate a project first and run the commands from its folder.

## Health and operations

Check standard gRPC health:

```bash
grpcurl -plaintext localhost:50051 grpc.health.v1.Health/Check
```

HTTP liveness for the operations listener:

```bash
curl http://localhost:9090/livez
```

The server marks gRPC health `NOT_SERVING` before graceful shutdown.

## Security notes

- Passwords use Argon2id.
- Local JWTs use RS256.
- Development RSA keys are generated by `make keys` and are not committed.
- Never commit production private keys.
- Tokens do not contain RBAC permissions; effective permissions are read from the database so revocation takes effect immediately.
- Do not log access tokens, passwords, authorization headers, or private keys.
- The bootstrap admin environment variables are for initialization and local development; use secret management in production.

## Suggested production extensions

This starter intentionally gives you the core service platform. Depending on the application, typical next extensions are refresh-token/session rotation, a separate `user_identities` table for multi-provider account linking, Redis permission caching with invalidation, rate limiting, database integration tests via Testcontainers, TLS/mTLS, and generated API documentation.

## Command reference

Run these commands from this generated project's root. `make` and `make help` list targets. Make helpers need Make, a POSIX shell, and awk; native Go commands work on Windows too.

| Command | What it does |
| --- | --- |
| `make format` / `make format-check` | Rewrite/check Go under cmd/, internal/, and gen/ |
| `make lint` | Run go vet with existing generated contracts |
| `make test` | Run race-enabled tests without generating contracts; requires a C toolchain for the race detector |
| `make build` | Build bin/server; override BINARY if needed |
| `make dev` | Start the Go process using exported environment values; needs PostgreSQL and configured JWT keys |
| `make check` | Format check → lint → race tests → build; does not regenerate tracked contracts |
| `make tools` | Install pinned Buf/protobuf/gateway tooling and current grpcurl; modifies Go's tool bin and may download modules |
| `make codegen` | Regenerate protobuf, gRPC, gateway and OpenAPI outputs; requires tools on PATH; modifies generated files |
| `make contracts-check` | Buf lint without code generation |
| `make ci` | Sequential codegen → contract lint → check; review generated changes |
| `make keys` | Run scripts/generate-dev-keys.sh; create only missing development RSA key files under secrets/ |
| `make docker-up` | Ensure keys, build/start Compose stack, run migrations and API bootstrap |
| `make docker-down` | Stop/remove Compose services; retain the named database volume |
| `make logs` | Follow API container logs until interrupted |
| `make migrate-up` | Apply pending migrations to DATABASE_URL using the local migrate CLI |
| `make migrate-down` | Roll back one migration in that database |
| `make migrate-create name=...` | Create paired SQL files using the migrate CLI |
| `make tidy` | Update go.mod/go.sum to match imports; may download dependencies |
| `make clean` | Remove bin/ and coverage.out; preserve contracts and secrets/ |

Compatibility aliases: `fmt` → `format`, `proto` → `codegen`, `run` → `dev`. Tests and dev no longer regenerate bindings as a side effect. After changing a contract, run codegen explicitly, review its output, then check. After generation with `--no-hooks`, run format once if module substitution changed import order.

The server reads process environment variables; it does not load configs/.env.example automatically. Export DATABASE_URL and overrides before make dev. The Make DATABASE_URL default is for migration commands, not automatic application configuration. Compose supplies its own container environment. Local migration commands additionally require the migrate CLI; make tools does not install it.

The key script accepts `--help`, resolves its project root regardless of caller directory, preserves existing keys, and sets the private key to mode 600 when creating it. OpenSSL is required. Database migrations and Docker startup are explicit actions; check/test do not run them.

On Windows without Make, use `gofmt -l ./cmd ./internal ./gen`, `go vet ./...`, `go test -race ./...`, and `go build -trimpath -o bin/server.exe ./cmd/server`. Use a POSIX environment for shell helpers, or configure development keys and services through your own platform tooling.
