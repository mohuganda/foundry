#!/bin/sh
set -eu
case "${1:-}" in
  -h|--help)
    printf '%s\n' 'Usage: scripts/check.sh' \
      'Checks formatting, Clippy, all tests, documentation, and the release build in order.' \
      'Stops at the first failure. Set RUSTUP_TOOLCHAIN to select a compiler or CARGO_NET_OFFLINE=true to disable downloads.'
    exit 0 ;;
esac
[ "$#" -eq 0 ] || { echo 'Unexpected argument; use --help.' >&2; exit 2; }
project_root=$(CDPATH= cd "$(dirname "$0")/.." && pwd)
cd "$project_root"
cargo fmt --all -- --check
cargo lint
cargo test --locked
RUSTDOCFLAGS="${RUSTDOCFLAGS:-} -D warnings" cargo doc --locked --no-deps
cargo build --release --locked
