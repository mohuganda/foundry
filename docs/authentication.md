# Configurable authentication

Authentication defaults to `none`. Go retains local JWT authentication; web/Flutter
pages remain available without an OIDC gate. Enable OIDC using public configuration:

```sh
foundry auth example > oidc.json
# Edit issuer, audience, and platform public-client callback settings.
foundry auth validate --file oidc.json --template react
foundry auth configure --project . --file oidc.json --dry-run
foundry auth configure --project . --file oidc.json
```

Creation accepts `--auth oidc --auth-config oidc.json`. Configuration has no secret
fields and is separate from Ministry branding. Edit auth/auth.json and run configure
without --file to synchronize bindings. Set mode to none and configure to disable
OIDC. Format, test and rebuild after changes. Earlier foundations need the documented
migration; configure does not silently retrofit application code or dependencies.

Web uses Authorization Code + PKCE through oidc-client-ts. Tokens stay in memory;
redirect transaction state uses sessionStorage. Reload requires another sign-in
redirect (the provider may reuse SSO). Register exact HTTPS callback/logout URLs at
the same application origin, on different paths. The host must serve the SPA on
both paths, including any Vite base prefix. The authenticatedRequest helper restricts
bearer tokens to VITE_API_BASE_URL, rejects traversal and redirects, and requires
HTTPS outside loopback. Use getAccessToken only for trusted API destinations.
Client gates are presentation; servers must authorize every protected operation.

Flutter uses AppAuth with OS-browser sign-in on Android/iOS/macOS and issuer/client
namespaced secure storage. Register native callback schemes/links before running.
This adapter explicitly rejects Flutter web/Linux/Windows OIDC. Pass a configured
OidcAuth to AppRoot for per-app settings. Bootstrap uses it as the API token provider
when enabled; legacy AuthService remains available for local backend sessions.
Refresh is deduplicated, rotated tokens are saved, failed refresh clears the session,
and logout clears local credentials before contacting the provider.

Go reads generated authconfig defaults. Runtime overrides: OIDC_ENABLED,
OIDC_ISSUER_URL, OIDC_AUDIENCE, OIDC_TOKEN_PROFILE, AUTH_LOCAL_ENABLED. JSON OIDC mode
disables local JWT by default; explicitly enable local auth for dual acceptance.
The API audience must differ from frontend client IDs. token_profile=keycloak requires
typ=Bearer in signed claims; rfc9068 requires an at+jwt header. Both verify signature,
issuer, audience, expiry and subject. Opaque tokens and other access-token profiles
require a separate verifier; never send ID tokens to APIs. Local RBAC remains
authoritative. Apply migration 000007 and link users using auth_provider=oidc,
external_issuer=<exact issuer>, external_subject=<verified sub>. Existing OIDC users
need explicit issuer backfill. Never auto-link by email. One configured issuer per
backend is supported; simultaneous providers need a separate identity-link model.

Local provider fixture (Docker Compose, OpenSSL, curl):

```sh
bash auth/dev/start.sh
# Trust auth/dev/certs/server.crt only for local development.
# Go projects can run the real-provider test:
FOUNDRY_KEYCLOAK_TEST_CA="$PWD/auth/dev/certs/server.crt" \
  go test ./internal/infrastructure/auth/oidc -run TestKeycloakIntegration -v
# Stop the fixture:
docker compose -f auth/dev/compose.yaml down
```

Fixture issuer: https://localhost:8443/realms/foundry. API audience: programme-api.
Public clients: programme-web and programme-mobile. Web callbacks use
https://localhost:5173/auth/callback and /auth/signed-out; mobile callbacks use
ug.go.health.programme://auth/callback and /auth/signed-out. Configure HTTPS for web
development. The default example.org JSON values must be changed for this fixture.
Admin: admin / local-development-only; override KEYCLOAK_ADMIN_PASSWORD before first
start. Create realm users, then link subjects and grant local application permissions.
The foundry-test service client uses local-test-only solely for the integration test.
The start-dev database, seven-day self-signed certificate and credentials are disposable
test settings. Generated certs/keys are Git-ignored; start renews near-expiry certs.

Rust CLI authentication is not implemented; OIDC generation fails before writing.
Non-OIDC providers/protocols need additional adapters. Never put credentials, tokens
or patient data in public config. MFA policy, provider/client provisioning, deployment
TLS and application authorization remain explicit application/deployment work.

## Native setup and application APIs

For Flutter, first create native platform folders for the selected Mason app. AppAuth
12.1 requires Android API 24+, iOS 13+, or macOS 10.15+. Set the Android minimum SDK
accordingly and add the callback scheme to `android/app/build.gradle.kts`:

```kotlin
android {
    defaultConfig {
        minSdk = 24
        manifestPlaceholders["appAuthRedirectScheme"] = "ug.go.health.programme"
    }
}
```

For iOS/macOS, register the scheme under CFBundleURLTypes in the app's Info.plist:

```xml
<key>CFBundleURLTypes</key>
<array><dict>
  <key>CFBundleURLSchemes</key>
  <array><string>ug.go.health.programme</string></array>
</dict></array>
```

Keep the entire callback URI registered at the provider. Review platform keychain
and secure-storage entitlements; macOS sandboxed apps also need outgoing networking.
Use claimed HTTPS app/universal links where your deployment supports them; custom
schemes must be unique to the application. Native setup is explicit because Foundry
creates a multi-app Dart workspace before platform projects exist.

The foundation_auth package exports OidcConfig, OidcAuth and OidcStatus. To customize
one app, pass OidcAuth(config: OidcConfig(...)) into AppRoot and use the same instance
as the API client's AuthTokenProvider. Bootstrap already does this for generated
settings and exposes oidcAuthProvider. In a multi-app workspace use separate client
registrations/callbacks for each app; shared branding does not imply shared clients.
The caller owns injected controllers and must dispose them when their host is removed.

Browser callers use auth.signIn(), auth.signOut(), auth.getSession(), and
await auth.getAccessToken() from src/shared/auth/client.js. The session snapshot
contains status and public user information, not tokens. subscribe returns an
unsubscribe function. React uses useSyncExternalStore; Vanilla JS tears down old views
when session state changes. The optional authenticatedRequest helper is available:

```js
import { authenticatedRequest } from './shared/auth/request.js';
const reports = await authenticatedRequest('/reports', {
  query: { period: 'quarter' },
  signal: abortController.signal,
});
```

Set VITE_API_BASE_URL to the trusted API origin/base path before building. The helper
refreshes near-expiry sessions, never retries a write after a 401, and leaves API
permission errors to the caller. Configure CORS for the exact web origin. Generic
requestJSON remains available for public requests.

## Local browser development

After starting the Keycloak fixture, change the example JSON issuer and both web
callback origins to the localhost values above. Trust the fixture certificate in
the development OS/browser; never disable certificate verification globally. Both
web templates support explicit development TLS:

```sh
DEV_TLS_CERT="$PWD/auth/dev/certs/server.crt" \
DEV_TLS_KEY="$PWD/auth/dev/certs/server.key" npm run dev
```

The server binds localhost:5173 and fails if the port is occupied, so callback URLs
stay predictable. These Node-only environment settings are not VITE_ client data.
The test certificate is only for loopback development. A physical mobile device
needs a reachable HTTPS issuer with a certificate it trusts and matching client
registrations; localhost on a device refers to that device.

## Migration and provider compatibility

New foundations include the integration in disabled mode. For an earlier generated
project, use a fresh foundation as a comparison: add the adapter/dependencies and
lockfiles, application gate, API token-provider wiring and canonical/generated config
before invoking configure. Configure detects missing integration markers and fails
with guidance; it does not rewrite existing application code. Back up and apply the
Go migration before deploying new repository code. Backfill trusted external issuers;
blank legacy issuers cannot authenticate. Rollback can fail if distinct issuers now
share a subject, preserving the old uniqueness constraint rather than merging users.
Replace legacy OIDC_CLIENT_ID with OIDC_AUDIENCE for the API verifier.

The generic protocol adapter supports providers exposing OIDC discovery and public
Authorization Code + PKCE. Backend access tokens must use either the Keycloak signed
Bearer claim profile or RFC9068 JWT profile. Configure a distinct resource audience
at the provider (the fixture includes an audience mapper). Opaque-token introspection,
SAML, provider-specific token shapes, device authorization and confidential browser
BFF sessions are not implemented. Logout/refresh capabilities depend on provider/client
configuration. No provider list implies those capabilities were tested individually.

CI tests disabled defaults and configuration safety in Foundry, enabled browser gates
and adapter lifecycle, native Flutter storage/refresh races and Mason app generation,
Go signature/issuer/audience/expiry/type/key rotation, issuer-bound local users, and
real Keycloak discovery/client-credentials access tokens using the disposable fixture.
Browser and mobile interactive login use library-backed adapters with mocked provider
boundaries in tests; automated device/browser sign-in against Keycloak is not included.

See the [shared login design guide](login-design.md) for consistent client layouts and the matching Keycloak theme.
