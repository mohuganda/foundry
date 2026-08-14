#!/bin/sh
set -eu

version="${1:-dev}"
output_dir="${FOUNDRY_DIST_DIR:-dist}"
ldflags="-s -w -X github.com/jabahum/foundry/internal/cli.Version=${version}"

case "$version" in
  ""|[._-]*|*[!0-9A-Za-z._-]*)
    echo "invalid version: $version" >&2
    exit 2
    ;;
esac

mkdir -p "$output_dir"

for target in \
  linux/amd64 linux/arm64 \
  darwin/amd64 darwin/arm64 \
  windows/amd64 windows/arm64
do
  goos=${target%/*}
  goarch=${target#*/}
  archive="foundry_${version}_${goos}_${goarch}"
  stage="$output_dir/$archive"
  binary="foundry"

  if [ "$goos" = "windows" ]; then
    binary="foundry.exe"
  fi

  mkdir -p "$stage"
  CGO_ENABLED=0 GOOS="$goos" GOARCH="$goarch" \
    go build -trimpath -ldflags "$ldflags" -o "$stage/$binary" ./cmd/foundry
  cp README.md CHANGELOG.md "$stage/"

  if [ "$goos" = "windows" ]; then
    (cd "$output_dir" && zip -qr "$archive.zip" "$archive")
  else
    tar -C "$output_dir" -czf "$output_dir/$archive.tar.gz" "$archive"
  fi
  rm -r "$stage"
done

(
  cd "$output_dir"
  if command -v sha256sum >/dev/null 2>&1; then
    sha256sum ./*.tar.gz ./*.zip > SHA256SUMS
  else
    shasum -a 256 ./*.tar.gz ./*.zip > SHA256SUMS
  fi
)
