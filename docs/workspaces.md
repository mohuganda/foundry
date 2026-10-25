# Commands by workspace

Use Foundry to create projects, add supported scaffolds, and inspect template versions. Use each generated project's native tools to install dependencies, run the application, validate changes, and build outputs. Foundry does not provide `run`, `test`, `build`, or `deploy` subcommands.

## Identify your workspace

| Workspace | Root markers | Scope |
| --- | --- | --- |
| Foundry source repository | `assets.go`, `cmd/foundry/`, `templates/` | Develop the generator and embedded templates; produce the Foundry executable |
| Generated Go | `.foundry.yaml`, `go.mod`, `cmd/server/`, `Makefile` | One backend application |
| Generated Flutter | `.foundry.yaml`, root `pubspec.yaml` with `workspace:`, `apps/`, `packages/`, `mason.yaml` | Multiple apps and shared packages managed together |
| Generated Vanilla JS | `.foundry.yaml`, `package.json`, `src/main.js` | One Vite browser application |
| Generated React | `.foundry.yaml`, `package.json`, `src/main.jsx` | One React/Vite browser application |
| Generated Rust | `.foundry.yaml`, `Cargo.toml`, `rust-toolchain.toml` | One application package with a library and executable; not yet a multi-crate Cargo workspace |

The source folders under Foundry's `templates/` are inputs to generation. Generate a project into a separate directory before running its development workflow. In particular, Go templates still have `.tmpl` suffixes and Rust contains crate-name tokens.

Generated `.foundry.yaml` records the template ID/version and project settings. Keep it for `foundry upgrade --check`. It does not select a directory automatically or make Foundry search ancestor folders. Run commands at the appropriate root, or use `--project` where supported. Paths are relative to your current directory.

## Foundry command categories

| Category | Commands | Working directory and effect |
| --- | --- | --- |
| Discovery | `foundry help`, `foundry list`, `foundry template list`, `foundry info <template>`, `foundry version` | Anywhere; print help or bundled metadata |
| Environment inspection | `foundry doctor` | Search the current folder and parents for a workspace, check relevant tools/versions, and print fixes; `--project <root>` selects an exact root |
| Creation | `foundry` or `foundry new <name> --template <alias>` | Parent directory; create a project folder. `--destination <path>` overrides its location |
| Go scaffolding | `foundry add service`, `repository`, `usecase`, `proto`, `migration` | Generated Go root, or `--project <root>`; create files for you to implement and wire |
| Flutter scaffolding | `foundry add feature`, `foundry add package` | Generated Flutter workspace root, or `--project <root>`; create feature/package files |
| Mason generation | `foundry mason <brick> [arguments...]` | Flutter workspace root; execute `mason make` from `PATH`, forwarding arguments. No Foundry `--project` option |
| Maintenance | `foundry upgrade --check` | Any generated root, or `--project <root>`; compare recorded and bundled versions, without modifying files |

`foundry new --help` lists creation flags. `foundry add --help` groups generators by ecosystem. The [full reference](commands.md#3-foundry-cli) describes aliases, hooks, overwrite behavior, errors, and follow-up work.

### Generator support matrix

| Operation | Go | Flutter | Vanilla JS | React | Rust |
| --- | --- | --- | --- | --- | --- |
| Create from bundled source | Yes | Yes | Yes | Yes | Yes |
| Check recorded template version | Yes | Yes | Yes | Yes | Yes |
| Add backend service/repository/use case/proto/migration | Yes | — | — | — | — |
| Add workspace feature/package | — | Yes | — | — | — |
| Generate Mason bricks | — | Yes, with Mason installed | — | — | — |
| Run/test/build | Make / Go | FVM / Melos / Flutter | npm | npm | Cargo / check scripts |

A dash means the operation has no generator for that template. Foundry does not enforce every workspace/type combination before writing files: choose the matching command explicitly. For web and Rust projects, add modules using the structure documented in their generated README.

## Go workflow

Run creation from a parent folder, then work from the generated Go root:

```sh
foundry new orders-api --template go --module your.org/orders-api --no-hooks
cd orders-api
go mod download
make tools
make format
make ci
```

`make tools` installs contract tooling; ensure Go's binary directory is on `PATH`. `make ci` regenerates contracts, lints them, then checks formatting, Go analysis, race-enabled tests, and the server build. Review generated changes. Go and a race-detector-compatible C toolchain are required.

Add a capability from that same root:

```sh
foundry add service orders
foundry add usecase orders create
foundry add migration create_orders
# Implement the generated skeletons, SQL, and application wiring.
make format
make ci
```

Generators create scaffolds; they do not implement persistence, register every handler, or apply migrations. `make docker-up` starts the development stack, including its database/migration/bootstrap workflow; `make logs` follows API logs and `make docker-down` stops the stack while retaining database volumes. `make dev` runs only the local server and requires exported environment configuration, a database, and any configured auth keys. See the generated README before running it.

For a repository containing several projects, select the backend explicitly:

```sh
# From the parent containing orders-api/:
foundry add repository orders --project ./orders-api
foundry upgrade --project ./orders-api --check
```

`--project` selects scaffolding/metadata roots; it does not change where later `make` commands run.

## Flutter workflow

Generate the workspace with hooks disabled so setup uses its pinned FVM SDK:

```sh
foundry new customer-mobile --template flutter --no-hooks
cd customer-mobile
bash scripts/setup.sh
fvm dart pub global run mason_cli:mason make flutter_app \
  --name customer --organization org.acme
# Edit root pubspec.yaml: add apps/customer under workspace:.
fvm dart pub get
fvm dart pub global run mason_cli:mason make feature --app customer --name orders
# Wire the feature's routing and application behavior.
fvm dart run melos run format
fvm dart run melos run ci
```

Setup installs/resolves workspace tooling and dependencies. Register each new app/package as required by the root Dart workspace before resolving packages. Run Melos from the workspace root, where it can operate across apps and shared packages. `ci` runs code generation, the quality gate, and brick validation; `check` runs formatting checks, analysis, and tests without code generation or brick validation.

For a workspace with one app, `foundry add feature orders` is a lightweight scaffold option; `foundry add package order_models` creates a shared package skeleton. Resolve dependencies after adding a package and implement/wire the generated code. `add feature` chooses the first app alphabetically under `apps/` with a `pubspec.yaml`. In a multi-app workspace, use the Mason feature brick with an explicit `--app` as above.

`foundry mason feature --app customer --name orders` is an optional shorthand when the `mason` executable on `PATH` uses the correct SDK. The explicit FVM command above guarantees the workspace SDK. These are alternative ways to invoke the brick; run it once per feature.

To run an app after configuring its native platform folders/device:

```sh
# From the workspace root; development is the default environment:
bash scripts/run_flavor.sh customer development
# Or build a selected platform from the app root:
cd apps/customer
fvm flutter build apk
```

The run script selects Dart environment configuration, not a native platform flavor. Native SDKs, permissions, signing, and platform-specific configuration remain per-app work. There is no single platform build for the whole workspace. Flutter scaffolding commands belong at the workspace root, not inside `apps/customer`.

## Vanilla JavaScript and React workflows

Both templates expose the same npm scripts. Choose one creation command:

```sh
foundry new storefront --template vanilla-js
# For React instead: foundry new storefront --template react
cd storefront
npm ci
npm run dev
```

Stop the development server with Ctrl+C. Add product capabilities under `src/features/`, compose screens under `src/pages/`, and keep reusable code in `src/shared/`. Foundry has no web `add feature` command.

Before committing or handing off a build:

```sh
npm run format
npm run check
# Preview the dist/ output locally:
npm run preview
```

`check` runs formatting verification, ESLint, tests, and the production build. `preview` stays running until stopped; it serves the existing `dist/` locally. Deploy `dist/` through your chosen static hosting workflow. Foundry does not publish it. Keep `package-lock.json` committed; use `npm ci` for reproducible installs. Each app has its own package root; run npm there, even when several apps share a parent repository.

## Rust workflow

```sh
foundry new report-tool --template rust
cd report-tool
cargo run --locked -- greet World
# Add modules under src/features/ and expose them through the library.
cargo fmt --all
bash scripts/check.sh
```

On native Windows, use `./scripts/check.ps1` for the full gate. The scripts check formatting, run Clippy, unit/integration/doc tests, generate API documentation with warnings denied, and build a release executable under `target/release/`.

Use `cargo test --locked` for a focused test run. `cargo check --locked` only type-checks; it is not equivalent to the full gate. Foundry has no Rust scaffolding commands and does not install Rust or invoke Cargo during generation. Start with the included package; introduce a multi-crate Cargo workspace when your application needs it.

## Maintenance across projects

```sh
# From the directory containing these generated projects:
foundry upgrade --project ./orders-api --check
foundry upgrade --project ./customer-mobile --check
foundry upgrade --project ./storefront --check
foundry upgrade --project ./report-tool --check
```

This compares `.foundry.yaml` with the templates bundled in the installed executable. It does not fetch remote templates, merge changes, update dependencies, or rewrite application code. Automatic upgrades are not implemented. Install a newer Foundry executable to inspect its bundled versions; review and apply template changes manually to an existing application.

## Developing Foundry itself

From the Foundry repository root:

```sh
make help
make format
make ci
```

Here `make ci` validates the CLI, generation tests, native executable build, and smoke generation of all five templates. Native SDK checks for the generated projects run separately in CI. The generated Go project's `make ci` has a different scope: contracts and that application's quality gate.

Edit `templates/` when changing future generated output, rebuild Foundry, and generate a fresh project to validate the change. Existing projects are independent copies and are not edited by rebuilding Foundry.

`make release VERSION=<version>` builds executable distribution archives and checksums. It does not publish a release or create template delivery archives. See the [release and script reference](commands.md#4-develop-and-release-foundry) for prerequisites and outputs.
