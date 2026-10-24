#!/usr/bin/env bash
set -euo pipefail
APP=${1:?Usage: ./scripts/run_flavor.sh <app> <development|staging|production>}
FLAVOR=${2:-development}
ENV_FILE="env/${FLAVOR}.json"
[[ -f "$ENV_FILE" ]] || { echo "Missing $ENV_FILE"; exit 1; }
(cd "apps/$APP" && fvm flutter run --dart-define-from-file="../../$ENV_FILE")
