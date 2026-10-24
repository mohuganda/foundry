#!/bin/sh
set -eu
case "${1:-}" in
  -h|--help)
    echo 'Usage: scripts/homebrew-formula.sh <version> <SHA256SUMS>'
    echo 'Print a Homebrew formula for a stable release using its four macOS/Linux archive checksums.'
    exit 0 ;;
esac
[ "$#" -eq 2 ] || { echo 'Expected version and checksum manifest; use --help.' >&2; exit 2; }
version=${1#v}
printf '%s\n' "$version" | awk '/^[0-9]+\.[0-9]+\.[0-9]+$/ { valid=1 } END { exit !valid }' || {
  echo 'Expected a stable version such as 1.3.0.' >&2; exit 2;
}
manifest=$2
checksum() {
  digest=$(awk -v file="foundry_${version}_$1_$2.tar.gz" '$2 == file { value=$1; count++ } END { if (count != 1) exit 1; print value }' "$manifest") || {
    echo "Missing or duplicate checksum for $1/$2." >&2; exit 1;
  }
  [ "${#digest}" -eq 64 ] || { echo 'Invalid SHA-256 length.' >&2; exit 1; }
  case "$digest" in *[!0-9a-fA-F]*) echo 'Invalid SHA-256 digest.' >&2; exit 1 ;; esac
  printf '%s' "$digest"
}
# Validate every required digest before emitting a formula.
darwin_arm64=$(checksum darwin arm64)
darwin_amd64=$(checksum darwin amd64)
linux_arm64=$(checksum linux arm64)
linux_amd64=$(checksum linux amd64)
cat <<FORMULA
class Foundry < Formula
  desc "Create and extend projects from reusable foundations"
  homepage "https://github.com/mohuganda/foundry"
  version "$version"

  on_macos do
    on_arm do
      url "https://github.com/mohuganda/foundry/releases/download/v$version/foundry_${version}_darwin_arm64.tar.gz"
      sha256 "$darwin_arm64"
    end
    on_intel do
      url "https://github.com/mohuganda/foundry/releases/download/v$version/foundry_${version}_darwin_amd64.tar.gz"
      sha256 "$darwin_amd64"
    end
  end

  on_linux do
    on_arm do
      url "https://github.com/mohuganda/foundry/releases/download/v$version/foundry_${version}_linux_arm64.tar.gz"
      sha256 "$linux_arm64"
    end
    on_intel do
      url "https://github.com/mohuganda/foundry/releases/download/v$version/foundry_${version}_linux_amd64.tar.gz"
      sha256 "$linux_amd64"
    end
  end

  def install
    bin.install "foundry"
    doc.install "README.md", "CHANGELOG.md", "docs/commands.md", "docs/workspaces.md"
  end

  test do
    assert_match "foundry #{version}", shell_output("#{bin}/foundry version")
    system bin/"foundry", "new", "brew-smoke", "--template", "rust", "--no-hooks"
    assert_path_exists testpath/"brew-smoke/Cargo.toml"
  end
end
FORMULA
