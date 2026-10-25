# Command and script reference

This guide covers the Foundry executable, development of Foundry itself, and commands inside generated projects. Run `foundry help` for CLI syntax, `make help` in either Make-based project for targets, and a shell script with `--help` for its inputs. Rust's PowerShell script uses `-Help`.

See [Distribution](distribution.md) for Homebrew installation, Docker usage, and release publication automation.

## Choose a command category

The [workspace guide](workspaces.md) explains which commands apply to each workspace and gives creation-to-validation recipes. This page is the detailed reference.

| Category | Scope | Reference |
| --- | --- | --- |
| Discovery | Any directory: help, templates, version, tool availability | [Inspect and generate](#inspect-and-generate) |
| Creation | Parent folder: interactive generation or `new` | [Creation flags and delivery](#inspect-and-generate) |
| Go scaffolding | Generated Go root: service, repository, use case, proto, migration | [Extend a project](#extend-a-project) |
| Flutter scaffolding | Workspace root: feature, package, Mason bricks | [Extend a project](#extend-a-project) |
| Maintenance | Generated roots: read-only template version comparison | [Maintenance and errors](#maintenance-and-errors) |
| Workspace development | Native tools for setup, run, format, lint, test, build, clean | Sections 5–8 below |
| Foundry development and distribution | This repository's Make targets, smoke tests, executable release archives | [Develop and release Foundry](#4-develop-and-release-foundry) |

## 1. Shared command meanings

Use each ecosystem's native command. The names differ where the ecosystem already has a standard; their responsibilities are consistent.

| Task | Foundry repository | Generated Go | Flutter workspace | Vanilla JS / React | Rust project |
| --- | --- | --- | --- | --- | --- |
| Set up | Go 1.23+, then `make build` | Go 1.23+, `go mod download`; optional `make tools` | `bash scripts/setup.sh` | `npm ci` | Install Rust via rustup |
| Format source | `make format` | `make format` | `fvm dart run melos run format` | `npm run format` | `cargo fmt --all` |
| Check formatting | `make format-check` | `make format-check` | `fvm dart run melos run format:check` | `npm run format:check` | `cargo fmt --all -- --check` |
| Lint / analyze | `make lint` | `make lint` | `fvm dart run melos run lint` | `npm run lint` | `cargo lint` |
| Run tests | `make test` | `make test` | `fvm dart run melos run test` | `npm test` | `cargo test --locked` |
| Full quality gate | `make check` | `make check` | `fvm dart run melos run check` | `npm run check` | `bash scripts/check.sh` or `./scripts/check.ps1` |
| Build | `make build` | `make build` | `fvm flutter build <platform>` from the app folder | `npm run build` | `cargo build --release --locked` |
| Run locally | `./bin/foundry help` | `make dev` | `bash scripts/run_flavor.sh <app> [environment]` | `npm run dev` | `cargo run --locked -- greet World` |
| Clean build outputs | `make clean` | `make clean` | `fvm dart run melos run clean` | `npm run clean` | `cargo clean` |

Formatting intentionally rewrites source. Formatting checks and linting do not fix source. Quality gates stop on failure and do not intentionally regenerate tracked source; they do create ordinary compiler, test, documentation, or build outputs. Fix a failed step and rerun the gate. `make ci` in Go and `melos run ci` in Flutter additionally regenerate code; review those changes before committing.

Flutter has no generic workspace build because each app needs a platform target, native project files, and platform SDKs. Its `check` covers workspace formatting, analysis, and tests. Build the selected app/platform separately.

**Rust distinction:** native `cargo check` only type-checks code. It is not the full quality gate. Use the supplied check script for formatting, Clippy, unit/integration/doc tests, documentation, and release build.

## 2. Prerequisites and working directories

| Work | Required tooling |
| --- | --- |
| Foundry source | Go 1.23+; Make and a POSIX shell for Make targets; `awk` for `make help` |
| Foundry executable | No language runtime; generated projects require their own SDKs |
| Release archives | Go, POSIX shell, `tar`, `zip`, and either `sha256sum` or `shasum` |
| Go project checks | Go and a C compiler supported by Go's race detector; project dependencies |
| Go code generation | Buf and the generators installed by `make tools`, on `PATH`; registry access when dependencies are uncached |
| Go runtime / database | PostgreSQL, development JWT keys when local auth is enabled, exported environment configuration; Docker Compose for the included stack; `migrate` for local migration commands; OpenSSL for key generation |
| Flutter | FVM, its pinned Flutter SDK, resolved packages; Mason CLI activated with the same SDK; native SDKs/devices for app builds/runs |
| Web templates | Node 22.18+ on 22.x, Node 24, or a newer supported even-numbered release; npm |
| Rust | Rust via rustup with rustfmt, Clippy, and the platform linker; stable selected by `rust-toolchain.toml`; Rust 1.85 minimum |

Run Make, npm, Melos, and Cargo commands from the corresponding project/workspace root unless a row says otherwise. The provided shell scripts resolve their own repository/project root and can be invoked by path from another directory. Rust's PowerShell script and web cleanup script do the same.

Commands that install SDKs, packages, tools, containers, or uncached Go dependencies may use the network. Offline **generation** does not mean a machine without SDKs/dependencies can build the result. Installing dependencies is explicit except for the legacy best-effort Go/Flutter generation hooks described below.

### Windows

The Foundry executable, npm commands, and Cargo commands work natively. Make and Bash helpers require a POSIX environment such as Git Bash or WSL with the listed tools. Rust also has a native PowerShell quality gate:

```powershell
./scripts/check.ps1
```

Native alternatives for Foundry source checks:

```powershell
$files = gofmt -l assets.go cmd internal
if ($files) { throw "Unformatted files: $files" }
go vet ./...
if ($LASTEXITCODE) { exit $LASTEXITCODE }
go test ./...
if ($LASTEXITCODE) { exit $LASTEXITCODE }
go build -trimpath -o bin/foundry.exe ./cmd/foundry
```

For a generated Go project, run `gofmt -l ./cmd ./internal ./gen`, `go vet ./...`, `go test -race ./...`, and `go build -trimpath -o bin/server.exe ./cmd/server`. The race detector requires a supported C toolchain. A bare `go build` does not start a server or database.

## 3. Foundry CLI

### Inspect and generate

| Command | Behavior |
| --- | --- |
| `foundry` | Ask for a project name and one of five templates, then generate it |
| `foundry help`, `--help`, `-h` | Show commands without generating anything |
| `foundry version`, `--version`, `-v` | Print the Foundry executable version |
| `foundry list` | Print bundled IDs, versions, and descriptions |
| `foundry template list` | Same bundled listing; `templates` is a command alias and `foundry template` also lists |
| `foundry info <template>` | Show its canonical ID, version, aliases, and description |
| `foundry doctor` | Detect the workspace and check relevant tool/SDK versions; print fixes without installing anything |
| `foundry new <name> --template <template> [flags]` | Copy the selected bundled foundation and write `.foundry.yaml` |

| Template ID | Accepted short aliases |
| --- | --- |
| `go-foundations` | `go`, `grpc`, `backend` |
| `flutter-foundations` | `flutter`, `mobile`, `app` |
| `vanilla-js-foundations` | `vanilla-js`, `vanilla`, `javascript`, `js` |
| `react-foundations` | `react` |
| `rust-foundations` | `rust`, `rs` |

Put the project name before flags: `foundry new portal --template react --destination ./apps/portal`. Flags-first form is also supported when the project name is last. Unexpected leftover arguments are rejected rather than silently ignored.

| Flag | Default and effect |
| --- | --- |
| `--template` | Required ID or alias |
| `--version` | Bundled version by default; if specified, it must equal that bundled version. It does not download another release |
| `--destination` | Project name as output path; an explicit path changes the output location |
| `--module` | Go module override; otherwise the normalized project name is used |
| `--org` | Flutter organization placeholder; native Flutter scaffolding defaults to `com.example` if omitted |
| `--package` | Flutter package placeholder/metadata token; Mason app variables are configured separately when generating a Mason app |
| `--no-hooks` | False; when true, skip optional pre/post generation tools |
| `--dry-run` | False; validate and list generated paths without changing the destination; hooks are skipped |
| `--force` | False; allow writing into a non-empty destination. Matching template files can be overwritten; unrelated files are retained |

Generation copies embedded source files, removes file suffixes such as `.go.tmpl`, substitutes project tokens, and restores executable shell-script permissions. JavaScript manifests, HTML, configuration and lockfiles receive the project name. Rust manifests, source, tests and lockfiles receive a valid crate identifier; reserved names and digit-leading names are prefixed with `app_`.

The optional Go hook attempts `go fmt ./...`. Legacy Flutter hooks attempt native scaffolding with `flutter` from `PATH`, rename the scaffolded app, and bootstrap with `melos` from `PATH` when available. Hook errors are warnings and may still return generation success. For the FVM/Mason workflow, generate Flutter with `--no-hooks`, then run its setup and app-generation steps explicitly. JS/React/Rust do not install dependencies or run package-manager hooks.

`.foundry.yaml` records project name, template ID/version/kind, and supplied module/org/package settings. Template versions and the Foundry binary version are separate: updating a binary does not rewrite an existing project.

### Extend a project

All `add` commands default to the current folder; `--project <path>` chooses another root. They create source scaffolds and print the remaining implementation/wiring steps. They refuse to overwrite existing files. Workspace type and all planned collisions are checked before writes. `--dry-run` validates the same plan and lists paths without changing the project. Handled write failures restore overwritten files and remove files/directories created by the operation; rollback failures are reported.

| Command | Outputs / follow-up |
| --- | --- |
| `foundry add service <name>` | Go proto, domain/repository port, application service, gRPC handler; implement persistence, run `make codegen`, wire the handler |
| `foundry add grpc-service <name>` | Alias for `add service` |
| `foundry add repository <name>` | Go repository port and PostgreSQL adapter skeleton; define the entity, implement SQL, wire it |
| `foundry add usecase <domain> [name]` | Go application use case; omitted name defaults to `execute`; supply the domain/entity/repository and implement behavior |
| `foundry add proto <name>` | Go protobuf service contract; run `make codegen` and implement/register the service |
| `foundry add migration <name>` | Next numbered SQL up/down pair; edit SQL and apply it explicitly |
| `foundry add feature <name>` | Flutter feature domain/application/data/presentation skeleton in the detected app; wire routing and add tests |
| `foundry add package <name>` | Flutter package skeleton under `packages/`; resolve workspace dependencies afterward |

Examples:

```sh
foundry add service orders --project ./backend
foundry add usecase orders create --project ./backend
foundry add usecase orders --project ./backend cancel
foundry add migration add_order_status --project ./backend
```

Flutter `add feature` selects the first app (alphabetical directory order) under apps/ with a pubspec.yaml. For explicit app selection in a multi-app workspace, use the Mason feature brick with `--app <name>`.

There are no `add` generators for JavaScript/React/Rust yet. Create feature modules using each template's documented structure.

`foundry mason <brick> [arguments...]` runs **`mason make` from PATH** with inherited input/output and arguments. It requires Mason installed and the correct workspace as the current directory. For FVM-managed workspaces the explicit equivalent is `fvm dart pub global run mason_cli:mason make <brick> ...`, which guarantees the selected SDK.

### Maintenance and errors

`foundry upgrade --project <root> --check` reads `.foundry.yaml` and compares the recorded template version to the bundled registry. It reports an available update and changes no files. `--check` defaults to true; disabling it is rejected because automatic upgrades are not implemented.

Exit codes: `0` for success/help, `1` for operational failures, and `2` for invalid usage/flags. Optional generation-hook warnings do not necessarily change the exit code. Use `foundry new --help`, `foundry add --help`, `foundry template --help`, `foundry upgrade --help`, or `foundry mason --help` for command-specific help.

Template pack/update/cache delivery commands and `--offline` are retired. Templates always ship as embedded files. `npm ci` is a dependency installation command; it does not mean `npm run check`.

## 4. Develop and release Foundry

All Make targets below run from the Foundry repository root.

| Command | What runs / outputs |
| --- | --- |
| `make`, `make help` | List canonical targets and descriptions |
| `make format` | Rewrite `assets.go`, `cmd/`, and `internal/` with gofmt |
| `make format-check` | Fail if those Go sources need formatting; does not scan node_modules or template implementations |
| `make lint` | `go vet ./...` |
| `make test` | `go test ./...`, including generation tests for all five embedded templates |
| `make build` | Static native `bin/foundry`; injects `VERSION`, default `dev` |
| `make check` | Formatting check, lint, tests, native build |
| `make smoke` | Build then run `scripts/smoke.sh` |
| `make ci` | `check` plus `smoke`; template SDK checks are separate CI jobs |
| `make install VERSION=...` | Build/install to `GOBIN`, or Go's `GOPATH/bin` default |
| `make release VERSION=...` | Run release helper; create six executable archives plus checksum manifest; no publishing |
| `make clean` | Remove repository `bin/` and `dist/`; leaves sources, SDKs, and package caches |

Compatibility aliases remain: `fmt` → `format`, `fmt-check` → `format-check`, `vet` → `lint`, `release-snapshot` → `release`.

### `scripts/smoke.sh`

Needs a built Foundry executable. It prints version/list, generates Go/Flutter/Vanilla JS/React/Rust into a fresh temporary directory with hooks disabled, checks their manifests and metadata, then adds a Go migration and service. It does not install template dependencies or run generated app tests. The temporary directory is removed on success, failure, or interruption; it does not reuse or delete a fixed `.tmp-smoke` folder.

| Input | Effect |
| --- | --- |
| `FOUNDRY_BIN` | Optional **absolute** path to a compiled Foundry; default repository `bin/foundry` |
| `TMPDIR` | Parent for temporary validation workspace; system `/tmp` fallback |
| `--help` | Print usage, perform no builds or generation |

### `scripts/build-release.sh [version]`

The script resolves the repository root, validates a plain version label, and builds with `CGO_ENABLED=0`, `-trimpath`, and linker flags that inject the version and strip Go debug symbols. Version labels may contain letters, digits, `.`, `_`, and `-`, and cannot start with `.`, `_`, or `-`; default is `dev`.

It stages Linux/macOS/Windows × amd64/arm64 builds in its own temporary directory. Linux/macOS archives use `.tar.gz`; Windows uses `.zip` and contains `foundry.exe`. Each archive includes the executable, README, changelog, and this `docs/` directory. On successful completion it moves all six archives and `SHA256SUMS` to the output folder and removes staging. A failed build cleans its staging folder without publishing staged partial artifacts.

`FOUNDRY_DIST_DIR` chooses the output folder, default `dist/`. Relative paths are resolved from the repository root even when the script is called elsewhere. Existing archives with the same version/target names are replaced. Older versions may coexist, but `SHA256SUMS` lists **only the six archives from the latest successful invocation**.

```sh
make ci
make release VERSION=1.2.0
(cd dist && shasum -a 256 -c SHA256SUMS)
# Linux alternative: (cd dist && sha256sum -c SHA256SUMS)
```

These archives deliver the **Foundry executable**; they are not template ZIP delivery. Running this script does not test, tag, or publish. The release workflow runs `make check`, packages on a `v*` tag (stripping `v` from the version label), uploads artifacts, and publishes a GitHub Release only for tags. A manual workflow run builds/uploads artifacts without publishing. Template SDK checks run in the separate CI workflow, so check their status before tagging.

## 5. Generated Go commands

Run inside a **generated** Go project. The template source folder has `.go.tmpl`/`go.mod.tmpl` files and is not itself the runnable module.

| Command | Responsibility / side effects |
| --- | --- |
| `make help` | List targets |
| `make tools` | Install pinned Buf/protobuf/gateway tooling and current grpcurl; may download and change GOBIN/GOPATH/bin; does not install the local `migrate` CLI |
| `make format`, `make format-check` | Format/check `cmd`, `internal`, and generated Go under `gen` |
| `make lint` | Vet using existing generated bindings |
| `make test` | Race-enabled tests using existing bindings; no regeneration or database migration |
| `make build` | Compile to `bin/server` (`BINARY` can override output path) |
| `make dev` | `go run ./cmd/server`; stays running until stopped; requires configuration and PostgreSQL |
| `make check` | Formatting check, vet, race tests, build; may fetch uncached Go modules, but no Buf tooling or live database is needed for current tests |
| `make codegen` | Regenerate Go/gRPC/Gateway/OpenAPI files with Buf; changes generated files |
| `make contracts-check` | Buf protobuf lint |
| `make ci` | Sequentially codegen → contract lint → check; requires tool installation; review generated changes |
| `make keys` | Create missing dev RSA key pair in `secrets/`; preserves existing key files |
| `make docker-up` | Ensure keys, build/start Compose stack; Compose runs migrations and initializes configured bootstrap data |
| `make docker-down` | Stop/remove Compose containers/networks, preserve named database volume |
| `make logs` | Follow API logs; interrupt to stop following |
| `make migrate-up` | Apply all pending SQL migrations to `DATABASE_URL` |
| `make migrate-down` | Roll back **one** migration; changes the selected database |
| `make migrate-create name=...` | Local `migrate` CLI creates numbered up/down files |
| `make tidy` | Rewrite module manifests to match imports |
| `make clean` | Remove `bin/` and `coverage.out`, preserve source/generated contracts/keys |

Aliases: `fmt` → `format`, `proto` → `codegen`, `run` → `dev`. Tests and dev no longer regenerate contracts implicitly.

For local execution, export `DATABASE_URL` (and any overrides) into the process environment. The server does **not** automatically load `configs/.env.example` or `.env`. `make migrate-up DATABASE_URL='...'` overrides only the Make migration connection; it does not export application configuration for other commands. The Compose file supplies its own container configuration.

`make docker-up` is the easiest complete local stack. To run the Go process separately, start/configure PostgreSQL, apply migrations, generate keys, export configuration, then `make dev`. A project generated with `--no-hooks` may need `make format` after module-name substitution before its first formatting check. After contract changes run `make codegen`, review outputs, then `make check`.

## 6. Flutter commands and scripts

Use FVM for all Dart/Flutter/Melos/Mason invocations. From the workspace root:

```sh
bash scripts/setup.sh
fvm dart pub global run mason_cli:mason make flutter_app --name my_app --description "My application"
# Add apps/my_app to workspace: in pubspec.yaml.
fvm dart pub get
fvm dart run melos run check
```

Mason generates application source, not a complete native platform deployment. Create/configure the desired Flutter native platforms in the app folder and follow the generated README before running or building it.

| `fvm dart run melos run <task>` | Effect |
| --- | --- |
| `bootstrap`, `pub:get` | Resolve workspace packages with `fvm dart pub get`; writes tool caches and may update pubspec.lock |
| `format` / `format:check` | Format/check `apps/` and `packages/` |
| `analyze` / `lint` | Analyze each workspace package; `lint` aliases `analyze` |
| `test` | Run flutter tests in Flutter packages with test folders and dart tests in non-Flutter packages with test folders; skips packages without tests |
| `test:coverage` | Flutter test coverage for Flutter packages with tests; does not produce Dart-only package coverage |
| `codegen` / `codegen:watch` | Run build_runner in packages declaring it; watch keeps running; writes generated source |
| `check` | Formatting check → lint → test; no code generation or app build |
| `quality` | Compatibility alias for `check` |
| `ci` | Codegen → check → all-brick validation; can fetch dependencies and generate source |
| `bricks:list` | List configured bricks through FVM's Mason CLI |
| `bricks:validate` | Run the brick validation script below |
| `outdated` | Report newer dependency versions; does not upgrade constraints |
| `clean` | `flutter clean` in Flutter workspace packages; does not remove source, lockfiles, all Dart caches, or globally installed SDKs |

| Script | Inputs, work, and outputs |
| --- | --- |
| `scripts/setup.sh` | No arguments. Resolve workspace root; require FVM; install configured Flutter SDK; activate Mason CLI globally in Dart's package cache; resolve packages; fetch local bricks. Requires network if uncached |
| `scripts/run_flavor.sh <app> [environment]` | Default `development`; also `staging` / `production`. Require an existing lowercase app directory and `env/<environment>.json`; run `fvm flutter run` from that app with Dart defines. It does not select a native `--flavor` build variant |
| `scripts/validate_bricks.sh` | Require installed SDK and FVM-activated Mason. Fetch bricks; render every brick in a temporary workspace; resolve packages; format-check; generate code; analyze and test generated app/package; remove the workspace on exit. SDK/pub/Mason caches can remain |

Copy staging/production example JSON to its corresponding ignored JSON file before using those environments. Values embedded in a client app are public; keep server secrets outside them. SDKs, emulators/devices, native permissions, signing, and app registration remain app-specific prerequisites. All three scripts accept `--help` before running setup or validation.

## 7. Vanilla JavaScript and React commands

Both templates expose the same scripts. From a generated project:

| Command | Behavior / output |
| --- | --- |
| `nvm install` | Optional convenience for nvm users; select the Node major in `.nvmrc` |
| `npm ci` | Install exactly the lockfile dependency tree; replaces node_modules; fails if dependency manifests disagree; does not change the lockfile |
| `npm run dev` | Vite development server, stays running; terminal shows URL/port |
| `npm run build` | Production static assets under `dist/` |
| `npm run preview` | Local preview of the already-built dist; run build first; not a production hosting service |
| `npm run format` / `format:check` | Rewrite/check files using Prettier; build/dependency outputs and lockfile are excluded |
| `npm run lint` | ESLint with warnings treated as failures |
| `npm test` / `npm run test:watch` | One Vitest pass / interactive watch process |
| `npm run check` | Formatting check → lint → tests → production build |
| `npm run clean` | Run scripts/clean.js; remove only this project's dist/, coverage/, and .vitest/; preserve node_modules, source, and lockfile |

`src/shared/config/env.js` reads public `VITE_` variables with defaults. Copy `.env.example` to `.env.local` for overrides. Vite applies them at startup/build time; rebuilding is required to change deployed values. `npm install` is for intentional dependency updates; commit its resulting package-lock.json and rerun check. The source templates use a project placeholder in both manifests; generated projects receive the normalized name.

## 8. Rust commands and scripts

Run in a generated Rust project. `rust-toolchain.toml` selects stable and requests rustfmt/Clippy. `Cargo.lock` is committed, and the initial project has no third-party crates.

| Command | Behavior / output |
| --- | --- |
| `cargo run --locked -- greet World` | Build debug executable if necessary and run greeting; arguments after `--` go to the app |
| `cargo run --locked -- --help` | App help (different from `cargo --help`) |
| `cargo fmt --all` / `cargo fmt --all -- --check` | Format/check package Rust source |
| `cargo check --locked` | Fast type-check only; not the full quality gate |
| `cargo lint` | .cargo/config.toml alias for Clippy across all targets with locked dependencies and warnings denied |
| `cargo test --locked` | Unit, real executable integration, and documentation tests |
| `cargo doc --locked --no-deps` | Generate target/doc/; the check scripts also deny documentation warnings |
| `cargo build --release --locked` | Optimized executable under target/release/; Windows adds .exe |
| `cargo clean` | Remove Cargo target artifacts including docs and builds; preserve Cargo.lock/source |
| `bash scripts/check.sh` | POSIX full gate below |
| `./scripts/check.ps1` | Native PowerShell full gate below |

Both check scripts resolve the project root, stop at the first failure, and run **format check → cargo lint → cargo test → documentation (warnings denied) → release build**. They do not fix source or update the lockfile. The PowerShell script restores the caller's location and RUSTDOCFLAGS on exit. Help uses `--help` for the shell script and `-Help` for PowerShell.

| Environment variable | Use |
| --- | --- |
| `RUSTUP_TOOLCHAIN` | Select a preinstalled toolchain; CI uses stable or 1.85.0 |
| `CARGO_NET_OFFLINE=true` | Disable crate downloads for the full gate; toolchain installation is a separate concern |
| `RUSTDOCFLAGS` | Existing doc flags are retained; the gate adds `-D warnings` for doc generation |

Do not append `--offline` to `cargo lint`: the alias already contains `--` and later arguments would be forwarded to Clippy/rustc. Use the offline environment variable, or write the expanded command with Cargo options before `--`:

```sh
CARGO_NET_OFFLINE=true bash scripts/check.sh
cargo clippy --locked --offline --all-targets -- -D warnings
```

## 9. How CI maps to local work

- Foundry OS jobs run native Go tests, vet, and build on Linux/macOS/Windows; Linux also runs the isolated smoke helper.
- Formatting calls `make format-check` and cross-compilation calls `make release VERSION=ci`.
- Web jobs generate fresh projects from the bundled binary, install with npm ci, and run the same npm check gate on Node 22 and 24.
- Rust jobs generate fresh projects and invoke the same POSIX/PowerShell check scripts. Stable runs on all three OSes; Rust 1.85 runs on Linux. Cargo is offline after toolchain installation.
- Flutter CI-equivalent local validation is `fvm dart run melos run ci` after setup. Native app builds require a chosen platform. There is no automated Flutter SDK job in the current repository CI workflow.
- Go-foundation full local validation is `make ci` in a generated project after installing code-generation tools. Foundry Go tests check generation, not all backend integration behavior. There is no separate automated generated-backend test job in the current workflow.

For template changes, generate a fresh project and run its native gate. Passing Foundry's Go tests alone does not prove a template's runtime code builds. No release script automatically merges a PR or upgrades existing generated projects.

## Version reporting

`foundry version` reports an injected release version when supplied by the release build. For `go install ...@<version>`, it reads the main module version from Go build metadata. Local builds without a module version report `dev`, with the VCS revision and dirty marker when available. Make's explicit `VERSION=dev` remains `dev`; no version is inferred from the bundled template versions.

## Workspace-aware doctor

Run `foundry doctor` from a generated root or subdirectory, or use `foundry doctor --project ./backend` to select an exact root. Foundry metadata takes precedence over native manifests; inconsistent/unsupported metadata produces an error. Without metadata, Go/Cargo/npm manifests and Flutter workspace markers identify the project. An explicit path must identify a supported root.

Go checks the compiler requirement in go.mod; Rust checks Cargo.toml's minimum compiler plus Cargo, Clippy, and rustfmt. Web projects check the foundation's supported Node versions and npm. Flutter checks FVM and compares the cached selected SDK version with .fvmrc without invoking Flutter or installing SDKs. Version probes run from the workspace root with a five-second timeout. Optional Go contract/container tooling is reported separately.

Exit 1 means a required check or workspace detection failed; optional tooling failures do not fail the command. Outside a recognized workspace, doctor shows optional tool availability and returns success. It reports tool readiness, not whether dependencies are installed or application tests pass. Use the project's native quality gate afterward.

## Scaffold safety and previews

```sh
foundry new backend --template go --dry-run
foundry add service orders --project ./backend --dry-run
```

`new` renders into a temporary staging directory, then validates and applies only generated paths. Dry runs remove their staging directory and never run hooks or change the destination. Existing non-empty destinations require `--force`; metadata for another template kind is rejected even with force. Unrelated destination files are not substituted or formatted. SDK/cache output produced by legacy staging hooks is excluded from the resulting project.

`add` requires the matching Go/Flutter workspace and native root manifest. It never overwrites existing files. Both creation and extension reject symlinks at the selected root or inside planned output paths, escaping paths, and file/directory conflicts. Handled write failures trigger rollback, including original contents and permissions for overwritten files. This is not a crash-safe transaction or a lock against concurrent filesystem changes: avoid concurrent edits to the same planned paths, and inspect any reported rollback errors. External generation hooks can affect SDK/package caches outside the destination; use `--no-hooks` for an explicit setup workflow.
