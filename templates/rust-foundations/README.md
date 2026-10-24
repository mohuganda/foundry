# `__PROJECT_NAME__`

A Rust 2024 application foundation with a thin executable, reusable library, feature modules, typed errors, and no third-party dependencies.

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

```sh
cargo fmt --all -- --check
cargo test --locked
cargo lint
cargo doc --locked --no-deps
cargo build --release --locked
```

`cargo test` runs unit, executable integration, and documentation tests. Tests cover Unicode input, validation, argument handling, output failures, and user-visible CLI behavior. `cargo lint` runs Clippy for every target with warnings treated as errors. Unsafe Rust is forbidden by the manifest. Release builds use thin LTO and remove debug information.

Commit `Cargo.lock` for this application. Review dependency and toolchain changes together with test results. CI checks Linux, macOS, Windows, and the minimum compiler version. For deployment, build a release executable on the target platform; configure additional cross-compilation targets as needed.

## Foundry delivery

Foundry embeds these source files directly. Generating this project needs no network access. Once its Rust toolchain is installed, the initial application can also build with `cargo build --locked --offline` because it has no external crates. Installing or updating Rust itself can require network access. Foundry does not invoke Cargo automatically.
