#!/usr/bin/env sh
set -eu
case "${1:-}" in
  -h|--help) printf '%s\n' 'Usage: scripts/generate-dev-keys.sh' 'Creates missing development RSA keys in the project secrets/ directory; keeps existing keys.'; exit 0 ;;
esac
[ "$#" -eq 0 ] || { echo 'Unexpected argument; use --help.' >&2; exit 2; }
project_root=$(CDPATH= cd "$(dirname "$0")/.." && pwd)
cd "$project_root"
command -v openssl >/dev/null || { echo 'OpenSSL is required.' >&2; exit 127; }
mkdir -p secrets
if [ ! -f secrets/jwt_private.pem ]; then
  openssl genrsa -out secrets/jwt_private.pem 2048 >/dev/null 2>&1
  chmod 600 secrets/jwt_private.pem
fi
if [ ! -f secrets/jwt_public.pem ]; then
  openssl rsa -in secrets/jwt_private.pem -pubout -out secrets/jwt_public.pem >/dev/null 2>&1
fi
echo "JWT development keys are ready in ./secrets"
