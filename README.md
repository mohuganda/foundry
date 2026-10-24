# Foundry

Foundry is a single Go CLI for creating and extending projects from reusable foundations.

## Highlights

- external, versioned template registry support
- SHA-256 verified template ZIP downloads
- local template cache and offline generation
- `foundry template update` / `foundry template list`
- `foundry add service` and `foundry add migration` for Go projects
- `foundry add feature` and `foundry add package` for Flutter workspaces
- bundled templates remain available without a registry or network connection

## Bundled templates

- `go-foundations` (`go`, `grpc`, `backend`) — Go, gRPC, Protobuf, Clean Architecture, PostgreSQL, local JWT, optional OIDC, local RBAC, OpenTelemetry, Prometheus, Docker, migrations and tests.
- `flutter-foundations` (`flutter`, `mobile`, `app`) — Flutter, Riverpod, go_router, Dio, Melos, Mason-ready feature generation and reusable foundation packages.

## Install

Foundry is distributed as a static executable and does not need a runtime. Release
archives are built for:

| Operating system | Architectures |
| --- | --- |
| Linux | amd64, arm64 |
| macOS | Intel (amd64), Apple Silicon (arm64) |
| Windows | amd64, arm64 |

Download the archive for your platform from GitHub Releases, verify it against
`SHA256SUMS`, extract it, and place `foundry` (or `foundry.exe`) on your `PATH`.

Install directly from source with Go 1.23 or newer:

```bash
go install github.com/jabahum/foundry/cmd/foundry@latest
```

For local development:

```bash
make build
./bin/foundry version
```

On Windows, where `make` may not be installed, use the Go command directly:

```powershell
go build -trimpath -o foundry.exe ./cmd/foundry
.\foundry.exe version
```

## Create projects

Interactive:

```bash
foundry
```

Go backend:

```bash
foundry new users-service \
  --template go \
  --module github.com/your-org/users-service
```

Flutter app:

```bash
foundry new mediguide \
  --template flutter \
  --org ug.go.health
```

Select an external/cached version:

```bash
foundry new users-service \
  --template go \
  --version 1.2.0 \
  --module github.com/your-org/users-service
```

Force offline operation:

```bash
foundry new users-service --template go --version 1.2.0 --offline
```

## Extend an existing project

Go:

```bash
cd users-service
foundry add service notifications
foundry add migration create_notifications
```

A service scaffold creates a Protobuf contract plus domain, application and gRPC transport skeletons. Foundry intentionally prints the remaining persistence and server-wiring steps rather than silently editing your composition root.

Flutter:

```bash
cd mediguide
foundry add feature emergency_guidelines
foundry add package foundation_analytics
```

A Flutter feature is created as a small clean slice with `domain`, `application`, `data`, and `presentation` directories.

## External template registry

Foundry does not require a registry. The built-in templates always remain available. A registry lets templates evolve independently from the Foundry binary.

Configure it once:

```bash
export FOUNDRY_REGISTRY_URL="https://raw.githubusercontent.com/your-org/foundry-templates/main/registry.json"
foundry template update
foundry template list
```

or provide the URL explicitly:

```bash
foundry template update --registry https://example.org/foundry/registry.json
```

Registry format:

```json
{
  "schema": 1,
  "updated_at": "2026-08-14T12:00:00Z",
  "templates": [
    {
      "id": "go-foundations",
      "name": "Go Foundations",
      "version": "1.2.0",
      "description": "Production Go gRPC foundation",
      "kind": "go",
      "aliases": ["go", "grpc", "backend"],
      "url": "https://github.com/your-org/foundry-templates/releases/download/go-foundations-v1.2.0/go-foundations.zip",
      "sha256": "<sha256-of-the-zip>"
    }
  ]
}
```

The `url` can be a GitHub Release asset or any direct HTTP(S) ZIP. Foundry verifies the downloaded archive against `sha256` before adding it to the cache.

Cache layout is platform-specific under the operating system's standard user-cache directory:

```text
foundry/
├── registry.json
└── templates/
    └── go-foundations/
        └── 1.2.0/
            └── template.zip
```

Set `FOUNDRY_CACHE_DIR` to use an explicit cache location, which is useful in
CI, containers, and other isolated environments.

Once downloaded:

```bash
foundry new demo --template go --version 1.2.0 --offline
```

works without a network connection.

## Commands

```text
foundry
foundry new <name> --template <go|flutter>
foundry add service <name>
foundry add migration <name>
foundry add feature <name>
foundry add package <name>
foundry template list
foundry template update --registry <url>
foundry list
foundry info <template>
foundry doctor
foundry cache list
foundry cache clean
foundry version
```

## Project-generation flags

```text
--module       Go module path
--org          Flutter organization identifier
--package      Flutter package token
--version      external/cached template version
--offline      do not download missing templates
--destination  output path
--force        generate into a non-empty directory
--no-hooks     skip post-generation hooks
```

## Architecture

```text
Foundry CLI
   |
   +-- built-in registry
   +-- external registry
   |      `-- SHA-256 verified releases
   |
   +-- template cache
   |
   +-- project generator
   |      +-- Go foundation
   |      `-- Flutter foundation
   |
   +-- add generators
   |      +-- Go service
   |      +-- Go migration
   |      +-- Flutter feature
   |      `-- Flutter package
   |
   +-- doctor
   `-- cache
```

Authentication/RBAC choices belong to the generated Go foundation, not to Foundry itself. Flutter's platform folders are generated using the installed Flutter SDK before the foundation overlay is applied when Flutter is available.

## Development

```bash
make fmt
make check
make build
make smoke
```

`go test ./...` includes tests for registry fetching/caching/checksum validation and the Go/Flutter add generators.

### Releases

Create all supported archives and checksums locally:

```bash
make clean
make release-snapshot VERSION=1.2.0
```

Pushing a tag such as `v1.2.0` runs tests, produces static Linux, macOS, and
Windows archives, injects the tag into `foundry version`, generates SHA-256
checksums, and publishes a GitHub Release. The release workflow can also be run
manually to validate packaging without publishing a release.

## Foundry v1.2 project maintenance

Every newly generated project now includes `.foundry.yaml`, recording the template, template version, project kind and generation metadata. This lets Foundry reason about the project's foundation without owning or overwriting application code.

```bash
foundry upgrade --check
```

The v1.2 upgrade command is intentionally read-only: it compares the project template version with bundled/cached registry releases and reports whether an update is available.

Additional Go generators:

```bash
foundry add repository orders
foundry add usecase orders create
foundry add proto billing
```

Flutter Mason bricks can be invoked consistently through Foundry:

```bash
foundry mason feature --name profile
```

Template maintainers can create versioned release archives and SHA-256 checksums:

```bash
foundry template pack ./my-template \
  --id my-template \
  --version 1.0.0 \
  --kind go \
  --out my-template-1.0.0.zip
```

If `--url` is supplied, Foundry also prints a registry JSON entry containing the generated SHA-256 digest.
# foundry
