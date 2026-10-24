# Foundry

Foundry is a single Go CLI for creating and extending projects from reusable foundations.

## Highlights

- Go, Flutter, Vanilla JavaScript, and React templates bundled directly as source files in one executable
- offline project generation with no template downloads or packaging step
- `foundry template list` to inspect the bundled templates
- `foundry add service` and `foundry add migration` for Go projects
- `foundry add feature` and `foundry add package` for Flutter workspaces

## Bundled templates

- `go-foundations` (`go`, `grpc`, `backend`) — Go, gRPC, Protobuf, Clean Architecture, PostgreSQL, local JWT, optional OIDC, local RBAC, OpenTelemetry, Prometheus, Docker, migrations and tests.
- `flutter-foundations` (`flutter`, `mobile`, `app`) — Flutter, Riverpod, go_router, Dio, Melos, Mason-ready feature generation and reusable foundation packages.
- `vanilla-js-foundations` (`vanilla-js`, `vanilla`, `javascript`, `js`) — JavaScript ES modules, Vite, feature-first DOM components, ESLint, Prettier, Vitest, and Testing Library.
- `react-foundations` (`react`) — React with JavaScript, Vite, feature-first components, Hooks linting, ESLint, Prettier, Vitest, and Testing Library.

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

Vanilla JavaScript:

```bash
foundry new website --template vanilla-js
cd website
npm ci
npm run dev
```

React (JavaScript):

```bash
foundry new dashboard --template react
cd dashboard
npm ci
npm run dev
```

The web templates use Node.js 22.18+ on the 22.x line, Node 24, or a newer supported even-numbered release. They include a feature-first `src/app`, `src/pages`, `src/features`, and `src/shared` structure, colocated tests, public environment configuration, reusable UI/styles, and an HTTP helper. Run `npm run check` for formatting, lint, tests, and a production build. Dependency installation is explicit; Foundry does not run npm automatically.

All templates are available offline. Optional generation hooks may invoke installed SDKs and their package managers; use `--no-hooks` to skip them.

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

## Template delivery

All templates are maintained under `templates/` and embedded directly into the Foundry executable. Generation copies their files, restores template filenames, and substitutes project settings.

| Template | Source folder | Delivery |
| --- | --- | --- |
| Go Foundations | `templates/go-foundations/` | Embedded source files |
| Flutter Foundations | `templates/flutter_foundations/` | Embedded source files |
| Vanilla JavaScript Foundations | `templates/vanilla-js-foundations/` | Embedded source files |
| React Foundations | `templates/react-foundations/` | Embedded source files |

Go source files use `.go.tmpl` and its module manifest uses `go.mod.tmpl` to keep the template separate from Foundry's Go module. Generated projects receive ordinary `.go` files and `go.mod`.

Template changes ship with a new Foundry build/release. There is no external template registry, template archive download, or template cache. `foundry template list` shows the available versions; `--version` can assert the bundled version but cannot select an external release.

To update a template, edit its source files and rebuild Foundry with `make build`.

## Commands

```text
foundry
foundry new <name> --template <go|flutter|vanilla-js|react>
foundry add service <name>
foundry add migration <name>
foundry add feature <name>
foundry add package <name>
foundry template list
foundry list
foundry info <template>
foundry doctor
foundry version
```

## Project-generation flags

```text
--module       Go module path
--org          Flutter organization identifier
--package      Flutter package token
--version      require this bundled template version
--destination  output path
--force        generate into a non-empty directory
--no-hooks     skip post-generation hooks
```

## Architecture

```text
Foundry CLI
   |
   +-- built-in registry
   +-- project generator
   |      +-- Go foundation
   |      +-- Flutter foundation
   |      +-- Vanilla JavaScript foundation
   |      `-- React foundation
   |
   +-- add generators
   |      +-- Go service
   |      +-- Go migration
   |      +-- Flutter feature
   |      `-- Flutter package
   |
   `-- doctor
```

Authentication/RBAC choices belong to the generated Go foundation, not to Foundry itself. Flutter's platform folders are generated using the installed Flutter SDK before the foundation overlay is applied when Flutter is available.

## Development

```bash
make fmt
make check
make build
make smoke
```

`go test ./...` checks generation of all four bundled templates, project setting substitution, CLI behavior, and the Go/Flutter add generators. Web-template CI generates fresh projects and runs `npm ci` and `npm run check` on Node 22 and 24.

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

The v1.2 upgrade command is intentionally read-only: it compares the project template version with the bundled template version and reports whether an update is available.

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
