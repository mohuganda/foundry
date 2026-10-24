#!/usr/bin/env bash
set -euo pipefail
case "${1:-}" in
  -h|--help) printf '%s\n' 'Usage: scripts/validate_bricks.sh' 'Fetches bricks, renders examples in a temporary workspace, resolves dependencies, generates code, formats, analyzes, and tests.' 'Requires fvm install and Mason CLI activated through FVM. Removes the temporary workspace on exit.'; exit 0 ;;
esac
[[ $# -eq 0 ]] || { echo 'Unexpected argument; use --help.' >&2; exit 2; }

repo_root=$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)
dart_bin="$repo_root/.fvm/flutter_sdk/bin/dart"
flutter_bin="$repo_root/.fvm/flutter_sdk/bin/flutter"
validation_dir=$(mktemp -d "${TMPDIR:-/tmp}/flutter-foundation-bricks.XXXXXX")

cleanup() {
  rm -rf "$validation_dir"
}
trap cleanup EXIT

[[ -x "$dart_bin" && -x "$flutter_bin" ]] || {
  echo "Run 'fvm install' before validating bricks."
  exit 1
}

export PATH="$repo_root/.fvm/flutter_sdk/bin:$PATH"
"$dart_bin" pub global run mason_cli:mason --version >/dev/null || {
  echo "Activate Mason with 'fvm dart pub global activate mason_cli'." >&2
  exit 1
}
mason() { "$dart_bin" pub global run mason_cli:mason "$@"; }

cd "$repo_root"
mason get

cat >"$validation_dir/pubspec.yaml" <<'YAML'
name: brick_validation_workspace
publish_to: none

environment:
  sdk: ">=3.12.0 <4.0.0"

workspace:
  - apps/sample_app
  - packages/*

dev_dependencies:
  flutter_lints: ^6.0.0
YAML

cat >"$validation_dir/analysis_options.yaml" <<'YAML'
include: package:flutter_lints/flutter.yaml

analyzer:
  exclude:
    - "**/*.g.dart"
    - "**/*.freezed.dart"
  errors:
    invalid_annotation_target: ignore

linter:
  rules:
    always_use_package_imports: true
    avoid_print: true
    directives_ordering: true
YAML

mason make flutter_app --quiet --output-dir "$validation_dir" \
  --name sample_app --description "Sample application" --organization com.example
mason make flutter_package --quiet --output-dir "$validation_dir" \
  --name sample_package --description "Sample package"
mason make feature --quiet --output-dir "$validation_dir" \
  --app sample_app --name catalog
mason make api_resource --quiet --output-dir "$validation_dir" \
  --app sample_app --name product
mason make riverpod_controller --quiet --output-dir "$validation_dir" \
  --app sample_app --feature catalog --name catalog
mason make drift_table --quiet --output-dir "$validation_dir" \
  --app sample_app --name product
mason make page --quiet --output-dir "$validation_dir" \
  --app sample_app --feature catalog --name overview
mason make form --quiet --output-dir "$validation_dir" \
  --app sample_app --feature catalog --name product
mason make model --quiet --output-dir "$validation_dir" \
  --app sample_app --feature catalog --name product
mason make repository --quiet --output-dir "$validation_dir" \
  --app sample_app --feature catalog --name product
mason make scanner --quiet --output-dir "$validation_dir" \
  --app sample_app --feature catalog --name product
mason make notification --quiet --output-dir "$validation_dir" \
  --app sample_app --feature catalog --name product
mason make offline_feature --quiet --output-dir "$validation_dir" \
  --app sample_app --feature catalog --name product
mason make test_suite --quiet --output-dir "$validation_dir" \
  --app sample_app --feature catalog

mkdir -p "$validation_dir/packages"
for package_dir in "$repo_root"/packages/foundation_*; do
  cp -R "$package_dir" "$validation_dir/packages/"
done

cd "$validation_dir"
"$dart_bin" pub get
"$dart_bin" format --output=none --set-exit-if-changed \
  apps/sample_app/lib apps/sample_app/test packages/sample_package

cd "$validation_dir/apps/sample_app"
"$dart_bin" run build_runner build
"$flutter_bin" analyze
"$flutter_bin" test

cd "$validation_dir/packages/sample_package"
"$flutter_bin" analyze
"$flutter_bin" test
