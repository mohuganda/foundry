#!/bin/sh
set -eu

case "${1:-}" in
  -h|--help)
    printf '%s\n' 'Usage: scripts/build-release.sh [version]' \
      'Builds static Foundry executables for Linux/macOS/Windows on amd64/arm64.' \
      'Writes six executable archives and SHA256SUMS to FOUNDRY_DIST_DIR (default: repository dist/).' \
      'Version defaults to dev. Does not test, create a tag, or publish a release.'
    exit 0 ;;
esac
[ "$#" -le 1 ] || { echo 'Expected at most one version; use --help.' >&2; exit 2; }
version=${1:-dev}
case "$version" in
  ""|[._-]*|*[!0-9A-Za-z._-]*) echo "invalid version: $version" >&2; exit 2 ;;
esac
repo_root=$(CDPATH= cd "$(dirname "$0")/.." && pwd)
cd "$repo_root"
for tool in go tar zip; do
  command -v "$tool" >/dev/null || { echo "Required tool not found: $tool" >&2; exit 127; }
done
if ! command -v sha256sum >/dev/null && ! command -v shasum >/dev/null; then
  echo 'sha256sum or shasum is required.' >&2; exit 127
fi
output_dir=${FOUNDRY_DIST_DIR:-dist}
mkdir -p "$output_dir"
output_dir=$(CDPATH= cd "$output_dir" && pwd)
scratch=$(mktemp -d "$output_dir/.foundry-release.XXXXXX")
trap 'rm -rf "$scratch"' EXIT HUP INT TERM
ldflags="-s -w -X github.com/jabahum/foundry/internal/cli.Version=${version}"
set --

for target in linux/amd64 linux/arm64 darwin/amd64 darwin/arm64 windows/amd64 windows/arm64; do
  goos=${target%/*}
  goarch=${target#*/}
  archive="foundry_${version}_${goos}_${goarch}"
  stage="$scratch/$archive"
  binary=foundry
  [ "$goos" != windows ] || binary=foundry.exe
  mkdir -p "$stage"
  CGO_ENABLED=0 GOOS="$goos" GOARCH="$goarch" \
    go build -trimpath -ldflags "$ldflags" -o "$stage/$binary" ./cmd/foundry
  cp README.md CHANGELOG.md "$stage/"
  cp -R docs "$stage/"
  if [ "$goos" = windows ]; then
    filename="$archive.zip"
    (cd "$scratch" && zip -qr "$filename" "$archive")
  else
    filename="$archive.tar.gz"
    tar -C "$scratch" -czf "$scratch/$filename" "$archive"
  fi
  set -- "$@" "$filename"
  rm -r "$stage"
done

# Hash only the six archives from this invocation, even if dist contains older releases.
(
  cd "$scratch"
  if command -v sha256sum >/dev/null; then
    sha256sum "$@"
  else
    shasum -a 256 "$@"
  fi
) > "$scratch/SHA256SUMS"
for filename in "$@"; do mv "$scratch/$filename" "$output_dir/$filename"; done
mv "$scratch/SHA256SUMS" "$output_dir/SHA256SUMS"
printf 'Release artifacts: %s\n' "$output_dir"
