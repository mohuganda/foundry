#!/usr/bin/env bash
set -euo pipefail
case "${1:-}" in
  -h|--help) printf '%s\n' 'Usage: scripts/setup.sh' 'Installs the pinned Flutter SDK, activates Mason CLI, resolves workspace packages, and fetches bricks.'; exit 0 ;;
esac
[[ $# -eq 0 ]] || { echo 'Unexpected argument; use --help.' >&2; exit 2; }
project_root=$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)
cd "$project_root"

command -v fvm >/dev/null || { echo "FVM is required: https://fvm.app"; exit 1; }

fvm install
fvm dart pub global activate mason_cli
fvm dart pub get
fvm dart pub global run mason_cli:mason get

echo "Flutter Foundation workspace is ready."
echo "Create an app: fvm dart pub global run mason_cli:mason make flutter_app --name my_app"
