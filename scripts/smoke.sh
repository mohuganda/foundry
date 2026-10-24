#!/bin/sh
set -eu

case "${1:-}" in
  -h|--help)
    printf '%s\n' 'Usage: [FOUNDRY_BIN=/absolute/path/to/foundry] scripts/smoke.sh' \
      'Generates all five templates without hooks, checks their manifests, and exercises Go add commands.' \
      'Uses a temporary folder which is removed on success or failure. Build Foundry first with make build.'
    exit 0 ;;
esac
[ "$#" -eq 0 ] || { echo 'Unexpected argument; use --help.' >&2; exit 2; }
repo_root=$(CDPATH= cd "$(dirname "$0")/.." && pwd)
binary=${FOUNDRY_BIN:-"$repo_root/bin/foundry"}
[ -x "$binary" ] || { echo 'Foundry executable not found; run make build or set FOUNDRY_BIN.' >&2; exit 1; }
scratch=$(mktemp -d "${TMPDIR:-/tmp}/foundry-smoke.XXXXXX")
trap 'rm -rf "$scratch"' EXIT HUP INT TERM

"$binary" version
"$binary" template list
for template in go flutter vanilla-js react rust; do
  project="$scratch/$template"
  if [ "$template" = go ]; then
    "$binary" new smoke --template go --module smoke --destination "$project" --no-hooks
    grep -q '^module smoke$' "$project/go.mod"
  else
    "$binary" new smoke --template "$template" --destination "$project" --no-hooks
  fi
  test -f "$project/.foundry.yaml"
  case "$template" in
    flutter) test -f "$project/pubspec.yaml" ;;
    vanilla-js|react) grep -q '"name": "smoke"' "$project/package.json" ;;
    rust) grep -q '^name = "smoke"$' "$project/Cargo.toml" ;;
  esac
done
"$binary" add migration create_smoke --project "$scratch/go"
set -- "$scratch/go"/migrations/*_create_smoke.up.sql
test -f "$1"
"$binary" add service smoke --project "$scratch/go"
test -f "$scratch/go/api/proto/smoke/v1/smoke.proto"
printf '%s\n' 'Smoke checks passed for all five templates.'
