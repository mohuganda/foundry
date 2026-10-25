#!/bin/sh
set -eu

case "${1:-}" in
  -h|--help)
    printf '%s\n' 'Usage: [FOUNDRY_BIN=/absolute/path/to/foundry] scripts/smoke.sh' \
      'Generates all five templates without hooks, checks manifests and exercises Go plus web/Rust feature generators.' \
      'Uses a temporary folder which is removed on success or failure. Build Foundry first with make build.'
    exit 0 ;;
esac
[ "$#" -eq 0 ] || { echo 'Unexpected argument; use --help.' >&2; exit 2; }
repo_root=$(CDPATH= cd "$(dirname "$0")/.." && pwd)
binary=${FOUNDRY_BIN:-"$repo_root/bin/foundry"}
[ -x "$binary" ] || { echo 'Foundry executable not found; run make build or set FOUNDRY_BIN.' >&2; exit 1; }
scratch=$(mktemp -d "${TMPDIR:-/tmp}/foundry-smoke.XXXXXX")
trap 'rm -rf "$scratch"' EXIT HUP INT TERM

"$binary" auth example > "$scratch/oidc.json"
"$binary" version
"$binary" template list
for template in go flutter vanilla-js react rust; do
  project="$scratch/$template"
  if [ "$template" = go ]; then
    "$binary" new smoke --template go --module smoke --profile moh-uganda --destination "$project" --no-hooks
    grep -q '^module smoke$' "$project/go.mod"
  else
    "$binary" new smoke --template "$template" --profile moh-uganda --destination "$project" --no-hooks
  fi
  test -f "$project/.foundry.yaml"
  if [ "$template" != rust ]; then
    "$binary" auth validate --file "$scratch/oidc.json" --template "$template"
    "$binary" auth configure --project "$project" --file "$scratch/oidc.json" --dry-run
    "$binary" auth configure --project "$project" --file "$scratch/oidc.json"
  fi
  "$binary" identity validate --file "$project/identity/identity.json"
  "$binary" identity sync --project "$project" --dry-run
  grep -q '"name": "Ministry of Health"' "$project/identity/public.json"
  if grep -q '"ownership"' "$project/identity/public.json"; then
    echo 'Ownership leaked into public identity.' >&2
    exit 1
  fi
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
for template in vanilla-js react rust; do
  "$binary" add feature order-history --project "$scratch/$template" --dry-run
  test ! -e "$scratch/$template/src/features/order_history"
  "$binary" add feature order-history --project "$scratch/$template"
  test -f "$scratch/$template/src/features/order_history/README.md"
done
printf '%s\n' 'Smoke checks passed for all five templates.'
