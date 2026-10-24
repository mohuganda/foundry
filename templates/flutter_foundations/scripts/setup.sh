#!/usr/bin/env bash
set -euo pipefail

command -v fvm >/dev/null || { echo "FVM is required: https://fvm.app"; exit 1; }

fvm install
fvm dart pub global activate mason_cli
fvm dart pub get
fvm dart pub global run mason_cli:mason get

echo "Flutter Foundation workspace is ready."
echo "Create an app: mason make flutter_app --name my_app"
