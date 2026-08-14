# Changelog

## Unreleased

- Remove checked-in build artifacts and an orphan generated application file.
- Add static release builds for Linux, macOS, and Windows on amd64 and arm64.
- Add tag-driven GitHub Releases with archives and SHA-256 checksums.
- Validate formatting and cross-compilation in CI.
- Accept the documented directory-first syntax for `foundry template pack`.

## 1.2.0

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
