# Foundry

Foundry is a single Go CLI for creating and extending projects from reusable foundations.

## Highlights

- Go, Flutter, Vanilla JavaScript, React, and Rust templates bundled directly as source files in one executable
- offline project generation with no template downloads or packaging step
- `foundry template list` to inspect the bundled templates
- `foundry add service` and `foundry add migration` for Go projects
- `foundry add feature` for Flutter, Vanilla JavaScript, React, and Rust; `foundry add package` for Flutter

## Bundled templates

- `go-foundations` (`go`, `grpc`, `backend`) — Go, gRPC, Protobuf, Clean Architecture, PostgreSQL, local JWT, optional OIDC, local RBAC, OpenTelemetry, Prometheus, Docker, migrations and tests.
- `flutter-foundations` (`flutter`, `mobile`, `app`) — Flutter, Riverpod, go_router, Dio, Melos, Mason-ready feature generation and reusable foundation packages.
- `vanilla-js-foundations` (`vanilla-js`, `vanilla`, `javascript`, `js`) — JavaScript ES modules, Vite, feature-first DOM components, ESLint, Prettier, Vitest, and Testing Library.
- `react-foundations` (`react`) — React with JavaScript, Vite, feature-first components, Hooks linting, ESLint, Prettier, Vitest, and Testing Library.
- `rust-foundations` (`rust`, `rs`) — Rust 2024, a thin CLI and reusable library, feature modules, validated domain types, typed errors, unit/integration/doc tests, rustfmt, and Clippy.

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

### Homebrew (macOS and Linux)

```sh
brew tap mohuganda/foundry https://github.com/mohuganda/foundry
brew install mohuganda/foundry/foundry
foundry version
```

The tap uses this repository's `Formula/` directory and installs verified release binaries for Apple Silicon/Intel macOS and arm64/amd64 Linux. Upgrade with `brew update && brew upgrade mohuganda/foundry/foundry`.

### Docker (Linux containers)

```sh
docker run --rm ghcr.io/mohuganda/foundry:1.3.0 version
# From the parent directory where the generated project should be written:
docker run --rm --user "$(id -u):$(id -g)" \
  --mount "type=bind,source=$(pwd),target=/workspace" \
  ghcr.io/mohuganda/foundry:1.3.0 new my-api --template go --no-hooks
```

Images support linux/amd64 and linux/arm64. Pin a version for reproducibility; `latest` follows the latest stable release. The image contains the CLI and embedded templates. Install language SDKs and run generated-project development commands on the host or in a separate development container. Use `--no-hooks` because generation hooks need SDKs that this image does not include. See [distribution details](docs/distribution.md) for Windows mounts, interactive use, and publishing.

### Go source installation

Install directly from source with Go 1.23 or newer:

```bash
go install github.com/jabahum/foundry/cmd/foundry@latest
```

For the complete CLI, script, template, and release reference, see [Commands and workflows](docs/commands.md). Each generated template also documents its own commands in its README.

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

## Command categories and workspaces

| Category | Commands | Where to run |
| --- | --- | --- |
| Discover | `help`, `list`, `template list`, `info`, `version` | Any directory |
| Diagnose | `doctor [--project <root>]` | Detect the current/parent workspace or select an explicit root |
| Create | `foundry`, `new` | Parent folder of the new project; or specify `--destination` |
| Extend Go | `add service`, `repository`, `usecase`, `proto`, `migration` | Generated Go root, or select it with `--project` |
| Extend Rust | `add feature` | Generated Rust root, or select it with `--project`; register the module afterward |
| Extend web | `add feature` | Generated Vanilla JS/React root, or select it with `--project` |
| Extend Flutter | `add feature`, `add package`, `mason` | Generated Flutter workspace root; `add` accepts `--project` |
| Maintain | `upgrade --check` | Any generated root, or select it with `--project` |
| Develop and validate | Make, Melos, npm, Cargo | The relevant generated workspace root |
| Contribute and release Foundry | Root `make` targets and release scripts | Foundry's own repository |

Start with the [workspace guide](docs/workspaces.md) for the support matrix and end-to-end examples for Go, Flutter, Vanilla JS, React, and Rust. Use the [command reference](docs/commands.md) for flags, scripts, prerequisites, outputs, and side effects. `foundry help` shows the same CLI categories.

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

Flutter workspace and app:

```bash
foundry new mediguide --template flutter --no-hooks
cd mediguide
bash scripts/setup.sh
fvm dart pub global run mason_cli:mason make flutter_app \
  --name mediguide --organization ug.go.health
# Add apps/mediguide to workspace: in pubspec.yaml.
fvm dart pub get
```

Configure native platform folders, permissions, and signing in the app before a platform build. The template README explains the remaining app setup.

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

Rust application:

```bash
foundry new command-tool --template rust
cd command-tool
cargo run --locked -- greet "Ada Lovelace"
cargo test --locked
cargo lint
```

The Rust starter has no third-party dependencies. Its toolchain file selects stable Rust with rustfmt and Clippy, and its minimum supported version is Rust 1.85. Domain rules live in feature modules, with CLI input/output at the boundary. Foundry does not run Cargo automatically.

All templates are available offline. For the FVM/Mason workflow, generate Flutter with `--no-hooks` and run `bash scripts/setup.sh` explicitly. Optional generation hooks may invoke installed SDKs and their package managers; use `--no-hooks` to skip them.

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
foundry add feature emergency_guidelines --app mediguide
foundry add package foundation_analytics
```

A Flutter feature is created as a small clean slice with `domain`, `application`, `data`, and `presentation` directories.

Web (Vanilla JavaScript or React):

```sh
foundry add feature order-history --project ./storefront --dry-run
foundry add feature order-history --project ./storefront
```

Foundry selects the UI/test scaffold from the workspace kind. Import the generated feature from a page, implement its behavior, and run `npm run check`.

Rust:

```sh
foundry add feature order-history --project ./command-tool
# Add pub mod order_history; to src/features/mod.rs, then wire its API.
```

The Rust scaffold includes a domain type, typed errors, an application operation, and unit tests. Tests compile and run after the module is registered. Run the project's full quality gate afterward.

## Template delivery

All templates are maintained under `templates/` and embedded directly into the Foundry executable. Generation copies their files, restores template filenames, and substitutes project settings.

| Template | Source folder | Delivery |
| --- | --- | --- |
| Go Foundations | `templates/go-foundations/` | Embedded source files |
| Flutter Foundations | `templates/flutter_foundations/` | Embedded source files |
| Vanilla JavaScript Foundations | `templates/vanilla-js-foundations/` | Embedded source files |
| React Foundations | `templates/react-foundations/` | Embedded source files |
| Rust Foundations | `templates/rust-foundations/` | Embedded source files |

Go source files use `.go.tmpl` and its module manifest uses `go.mod.tmpl` to keep the template separate from Foundry's Go module. Generated projects receive ordinary `.go` files and `go.mod`.

Template changes ship with a new Foundry build/release. There is no external template registry, template archive download, or template cache. `foundry template list` shows the available versions; `--version` can assert the bundled version but cannot select an external release.

To update a template, edit its source files and rebuild Foundry with `make build`.

## Commands and workflows

[The command reference](docs/commands.md) documents every CLI command and flag, each template's native commands, setup requirements, generated outputs, side effects, Windows alternatives, and CI/release behavior.

```sh
foundry help
foundry new --help
foundry template list
make help
```

The shared quality workflow is format → lint/analyze → test → build where a generic build exists:

| Project | Full quality gate |
| --- | --- |
| Foundry / generated Go | `make check` |
| Flutter workspace | `fvm dart run melos run check` (build the chosen app/platform separately) |
| Vanilla JS / React | `npm run check` |
| Rust | `bash scripts/check.sh` or `./scripts/check.ps1` |

Formatting checks do not fix source; regeneration is explicit. Use `make format`, `npm run format`, `fvm dart run melos run format`, or `cargo fmt --all` to fix formatting. Existing `fmt`, `fmt-check`, `vet`, `proto`, `run`, `quality`, and `release-snapshot` aliases remain available in their relevant projects.

## Architecture

```text
Foundry CLI
   |
   +-- built-in registry
   +-- project generator
   |      +-- Go foundation
   |      +-- Flutter foundation
   |      +-- Vanilla JavaScript foundation
   |      +-- React foundation
   |      `-- Rust foundation
   |
   +-- add generators
   |      +-- Go service
   |      +-- Go migration
   |      +-- Flutter feature
   |      `-- Flutter package
   |
   `-- doctor
```

Authentication/RBAC choices belong to the generated Go foundation, not to Foundry itself. Legacy Flutter hooks use the installed Flutter executable before the foundation overlay, only when hooks are enabled. The recommended FVM-first setup is explicit after `--no-hooks`.

## Development

```bash
make format
make ci
```

`go test ./...` checks generation of all five bundled templates, project setting substitution, CLI behavior, and the Go/Flutter/web/Rust add generators. Web-template CI generates fresh projects and feature scaffolds, then runs `npm ci` and `npm run check` on Node 22 and 24. Rust-template CI generates and registers a feature module, then invokes the same POSIX/PowerShell check scripts for formatting, tests, Clippy, documentation, and release builds on stable Rust across Linux/macOS/Windows and on Rust 1.85 for minimum-version compatibility.

### Releases

Create all supported archives and checksums locally:

```bash
make ci
make release VERSION=1.2.0
```

Pushing a tag such as `v1.2.0` runs Foundry checks, produces static Linux, macOS, and
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

With a Mason executable configured on PATH, Foundry forwards brick arguments to `mason make`:

```bash
foundry mason feature --app mediguide --name profile
```

## Organisation identity

All five templates support `--profile moh-uganda` and an optional `--identity file.json`.
Edit the generated `identity/identity.json`, preview with
`foundry identity sync --project ./portal --dry-run`, then sync and run the native
formatter/quality gate. See the [identity guide](docs/identity.md) for configuration, shared
layouts, report branding, public metadata, and ownership references. Official
logos/fonts/colours are not bundled; visual defaults are neutral and pending approval.

## Optional authentication

Go, React, Vanilla JS, and native Flutter support configurable OIDC authentication.
Generate an example with `foundry auth example > oidc.json`, edit and validate it, then
create with `--auth oidc --auth-config oidc.json`. Authentication remains disabled by
default in frontends; Go retains local JWT. See the [authentication guide](docs/authentication.md) for
provider configuration, Keycloak development setup, and platform limitations.
