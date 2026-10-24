# `__PROJECT_NAME__`

A Rust 2024 application foundation with a thin executable, reusable library, feature modules, typed errors, and no third-party dependencies.

## Foundry command scope

This is a generated application root. Use Cargo and the check scripts here for development, tests, and builds. Add features under `src/features/` using the structure below; Foundry has no `add` generators for this template. Go and Flutter scaffolding commands belong to their respective workspaces.

`foundry upgrade --check` compares the version in `.foundry.yaml` with the installed executable's bundled template; it changes no files. From another directory, use `foundry upgrade --project <root> --check`. Foundry does not run, build, or deploy this application.

## Run

Install Rust with rustup from https://rust-lang.org/tools/install/. The included toolchain file selects stable Rust with rustfmt and Clippy. The minimum supported compiler is Rust 1.85; CI checks both that version and stable.

```sh
cargo run --locked -- greet "Ada Lovelace"
cargo run --locked -- --help
```

The generated Cargo package and library are named `__RUST_CRATE__`. Foundry normalizes the project name to a valid crate identifier and prefixes reserved names or names beginning with a digit with `app_`.

## Structure

```text
src/
  main.rs                 Process startup, exit status, and error reporting
  lib.rs                  Reusable library root and public module tree
  cli.rs                  Argument parsing and output adapter
  features/
    greeting/
      mod.rs              Feature's public API
      domain.rs           Validated input and domain errors
      service.rs          Application operation and documentation example
tests/
  cli.rs                  Integration tests of the actual executable
.cargo/config.toml        Cargo convenience aliases
```

Replace the example greeting with your first feature. Keep business rules independent of command-line or transport details. Add another module under `features/` when it owns a separate capability. Share code only when multiple features need it. Grow into a Cargo workspace when independently reusable crates justify it.

The executable consumes the same public library API as other callers. Domain values enforce invariants at construction; typed errors propagate to the boundary. Command output uses an injected writer, making it testable without global output capture. External I/O and configuration belong at application boundaries. Introduce a CLI parser, async runtime, HTTP framework, or persistence library when the product requires them.

## Checks

Full gate on macOS/Linux (or a POSIX environment):

```sh
bash scripts/check.sh
```

Native Windows PowerShell:

```powershell
./scripts/check.ps1
```

The equivalent individual commands are:

```sh
cargo fmt --all -- --check
cargo lint
cargo test --locked
RUSTDOCFLAGS="${RUSTDOCFLAGS:-} -D warnings" cargo doc --locked --no-deps
cargo build --release --locked
```

`cargo test` runs unit, executable integration, and documentation tests. Tests cover Unicode input, validation, argument handling, output failures, and user-visible CLI behavior. `cargo lint` runs Clippy for every target with warnings treated as errors. Unsafe Rust is forbidden by the manifest. Release builds use thin LTO and remove debug information.

Commit `Cargo.lock` for this application. Review dependency and toolchain changes together with test results. CI checks Linux, macOS, Windows, and the minimum compiler version. For deployment, build a release executable on the target platform; configure additional cross-compilation targets as needed.

## Foundry delivery

Foundry embeds these source files directly. Generating this project needs no network access. Once its Rust toolchain is installed, the initial application can also build with `cargo build --locked --offline` because it has no external crates. Installing or updating Rust itself can require network access. Foundry does not invoke Cargo automatically.

## Command reference

Run Cargo commands from this generated project root. Both check scripts resolve the root from their own path and stop at the first failure.

| Command | What it does |
| --- | --- |
| `cargo run --locked -- greet World` | Build/debug-run the application; values after -- are application arguments |
| `cargo run --locked -- --help` | Show application help; different from cargo --help |
| `cargo fmt --all` / `cargo fmt --all -- --check` | Format/check Rust source |
| `cargo check --locked` | Type-check only; not the full quality gate |
| `cargo lint` | Alias for clippy --locked --all-targets -- -D warnings |
| `cargo test --locked` | Unit, executable integration, and documentation tests |
| `cargo doc --locked --no-deps` | Write API documentation under target/doc/ |
| `cargo build --release --locked` | Write optimized executable under target/release/; Windows adds .exe |
| `cargo clean` | Remove Cargo target artifacts, preserving Cargo.lock and source |
| `bash scripts/check.sh` | Format check → Clippy → all tests → docs (warnings denied) → release build |
| `./scripts/check.ps1` | Same full gate using native PowerShell |

Help: `bash scripts/check.sh --help` or `./scripts/check.ps1 -Help`. Scripts do not fix formatting or change the lockfile. Cargo's normal build/test/doc outputs are created under target/.

Use RUSTUP_TOOLCHAIN to select a preinstalled compiler (CI uses stable or 1.85.0). Set CARGO_NET_OFFLINE=true to disable crate downloads for the whole gate; toolchain installation is separate. Existing RUSTDOCFLAGS are preserved, with -D warnings added for documentation. The PowerShell script restores both the caller's working directory and doc flags.

Do not append --offline to cargo lint: the alias already contains --, so appended options go to the compiler. Instead use the offline environment variable or the expanded command:

```sh
CARGO_NET_OFFLINE=true bash scripts/check.sh
cargo clippy --locked --offline --all-targets -- -D warnings
```
