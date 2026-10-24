#!/bin/sh
set -eu
case "${1:-}" in
  -h|--help) echo 'Usage: scripts/test-docker.sh [image]'; echo 'Check CLI help and offline generation of all five templates in a mounted temporary directory.'; exit 0 ;;
esac
[ "$#" -le 1 ] || { echo 'Expected at most one image; use --help.' >&2; exit 2; }
image=${1:-foundry:test}
workspace=$(mktemp -d)
trap 'rm -rf "$workspace"' EXIT HUP INT TERM
run() {
  docker run --rm --network none --user "$(id -u):$(id -g)" \
    --mount "type=bind,source=$workspace,target=/workspace" "$image" "$@"
}
run version
run help
for template in go flutter vanilla-js react rust; do
  run new "$template" --template "$template" --destination "$template" --no-hooks
  test -f "$workspace/$template/.foundry.yaml"
done
test -f "$workspace/go/go.mod"
test -f "$workspace/flutter/pubspec.yaml"
test -f "$workspace/vanilla-js/package.json"
test -f "$workspace/react/package.json"
test -f "$workspace/rust/Cargo.toml"
test -x "$workspace/rust/scripts/check.sh"
echo 'Docker offline generation checks passed for all five templates.'
