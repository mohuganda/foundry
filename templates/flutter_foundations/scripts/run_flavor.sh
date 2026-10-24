#!/usr/bin/env bash
set -euo pipefail
case "${1:-}" in
  -h|--help) printf '%s\n' 'Usage: scripts/run_flavor.sh <app> [development|staging|production]' 'Runs an existing Flutter app with env/<environment>.json. Default: development.' 'The selected JSON must exist; staging/production examples must be copied and configured first.'; exit 0 ;;
esac
[[ $# -ge 1 && $# -le 2 ]] || { echo 'Usage: scripts/run_flavor.sh <app> [development|staging|production]' >&2; exit 2; }
app=$1
flavor=${2:-development}
[[ "$app" =~ ^[a-z][a-z0-9_]*$ ]] || { echo 'App must be a lowercase Dart package name.' >&2; exit 2; }
case "$flavor" in development|staging|production) ;; *) echo 'Environment must be development, staging, or production.' >&2; exit 2 ;; esac
project_root=$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)
cd "$project_root"
env_file="env/$flavor.json"
[[ -f "$env_file" ]] || { echo "Missing $env_file; configure it from the corresponding example." >&2; exit 1; }
[[ -d "apps/$app" ]] || { echo "Missing apps/$app; generate and register the app first." >&2; exit 1; }
cd "apps/$app"
fvm flutter run --dart-define-from-file="../../$env_file"
