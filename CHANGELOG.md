# Changelog

## Unreleased

- Add a Homebrew tap backed by this repository's Formula directory and verified macOS/Linux release archives.
- Add a non-root Docker image with all five embedded foundations, offline generation tests, and multi-platform GHCR publishing after stable releases.
- Automate Homebrew version/checksum updates after successful release publication and document distribution usage and maintenance.

## 1.3.0 - 2026-10-24

### Compatibility notes

- Template ZIP delivery, external registries, downloads, caches, and their CLI commands are retired. Generate projects from the five source foundations bundled in the executable; `--version` only asserts the bundled template version.
- Existing generated projects remain independent copies. `foundry upgrade --check` compares recorded and bundled versions without modifying project files.

### Changes

- Categorize CLI help by purpose and workspace; document root selection, supported generators, and complete workflows for all five templates.

- Align quality command meanings across the CLI repository and all five foundations; retain existing Make/Melos aliases.
- Add comprehensive command/workflow documentation, script help, CWD-independent helpers, and native PowerShell Rust checks.
- Make Go contract generation explicit, isolate smoke tests, and checksum only the current release's six executable archives.
- Correct use-case argument parsing so a trailing --project flag is honored; return success for command help.

- Add embedded Rust 2024 foundations with a thin CLI, reusable library, feature modules, typed errors, unit/integration/doc tests, rustfmt, and Clippy.
- Validate generated Rust projects across Linux, macOS, Windows, and the minimum Rust version.

- Complete the project metadata package rename from `internal/projectconfig` to `internal/config`.

- Add embedded Vanilla JavaScript and React foundations with Vite, feature-first structure, reproducible npm installs, linting, formatting, behavior tests, and production builds.
- Expose both web templates in the CLI menu, aliases, generation metadata, and cross-version Node CI.

- Deliver all templates as embedded source files and restore `.tmpl` filenames during generation.
- Remove template ZIP packaging, downloads, external registries, and cache commands.
- Resolve template version checks and upgrade reports from the bundled registry only.

## 1.2.0 - 2026-08-30

- Remove checked-in build artifacts and an orphan generated application file.
- Add static release builds for Linux, macOS, and Windows on amd64 and arm64.
- Add tag-driven GitHub Releases with archives and SHA-256 checksums.
- Validate formatting and cross-compilation in CI.
- Accept the documented directory-first syntax for `foundry template pack`.
- Add `FOUNDRY_CACHE_DIR` for deterministic CI and container cache placement.
- Added `.foundry.yaml` project metadata to generated projects.
- Added safe, read-only `foundry upgrade --check` foundation version reporting.
- Added `foundry add repository` for Go repository port/PostgreSQL adapter scaffolding.
- Added `foundry add usecase` for Go application use-case scaffolding.
- Added `foundry add proto` for standalone protobuf service contracts.
- Added `foundry mason` passthrough for Flutter Mason bricks.
- Added `foundry template pack` for release ZIP creation, SHA-256 generation, and registry-entry output.
- Preserved bundled/offline templates and v1.1 registry/cache behavior.

## 1.1.0

- Added external template registry, verified downloads and offline cache.
- Added initial Go and Flutter component generators.
