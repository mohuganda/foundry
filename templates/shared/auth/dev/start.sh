#!/usr/bin/env bash
set -euo pipefail
root=$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)
case "${1:-}" in
  -h|--help) echo 'Usage: bash auth/dev/start.sh'; echo 'Creates disposable local TLS certificates, starts Keycloak, waits for discovery. Requires Docker Compose, OpenSSL and curl. Stops with docker compose -f auth/dev/compose.yaml down.'; exit 0 ;;
esac
[[ $# -eq 0 ]] || { echo 'Unexpected argument.' >&2; exit 2; }
mkdir -p "$root/certs"
if [[ ! -f "$root/certs/server.crt" || ! -f "$root/certs/server.key" ]] || ! openssl x509 -checkend 3600 -noout -in "$root/certs/server.crt"; then
  openssl req -x509 -newkey rsa:2048 -nodes -days 7 -subj /CN=localhost \
    -addext 'subjectAltName=DNS:localhost,IP:127.0.0.1' \
    -keyout "$root/certs/server.key" -out "$root/certs/server.crt"
  # The non-root container must read this disposable test key. Never use it outside this fixture.
  chmod 644 "$root/certs/server.key"
fi
docker compose -f "$root/compose.yaml" up -d
for attempt in {1..90}; do
  if curl --fail --silent --cacert "$root/certs/server.crt" https://localhost:8443/realms/foundry/.well-known/openid-configuration >/dev/null; then
    echo 'Local Keycloak ready at https://localhost:8443/realms/foundry'
    exit 0
  fi
  sleep 2
done
docker compose -f "$root/compose.yaml" logs >&2
exit 1
