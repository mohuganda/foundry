# Flutter Foundation Workspace v3

A reusable production-oriented Flutter monorepo foundation built around **Melos + Mason**, with Riverpod, go_router, Dio, Freezed/JSON, Drift-ready offline boundaries, secure storage, typed environments, design-system primitives, observability, connectivity awareness, files/permissions, forms, notifications, scanning, localization and reusable test helpers.

The workspace is pinned to Flutter 3.44.8 through FVM. Install [FVM](https://fvm.app/) before running the setup commands; use `fvm flutter` and `fvm dart` rather than system SDK commands.

## Foundry command scope

Run `foundry add feature <name>` and `foundry add package <name>` from this workspace root, or append `--project <root>` from another directory. Feature generation selects automatically when exactly one app under `apps/` has a pubspec. With multiple apps, run `foundry add feature <name> --app <app-directory>`; omitting the selection lists candidates and stops before writing. Add `--dry-run` to preview the files. The FVM Mason feature command also accepts explicit `--app <name>`. Resolve workspace dependencies after adding/registering a package and wire the generated feature into the application.

`foundry mason <brick> [arguments...]` runs Mason from PATH in the current directory; use the explicit FVM commands below to guarantee the workspace SDK. Run workspace commands at this root and platform builds inside the selected app. `foundry upgrade --check` only compares template versions. Go service/repository/usecase/proto/migration generators belong to Go projects.

## Foundation packages

```text
packages/
  foundation_core/             Result, Failure, pagination
  foundation_network/          Dio + bearer auth + serialized token refresh
  foundation_auth/             Session/token refresh boundary
  foundation_storage/          secure and ordinary key/value storage
  foundation_testing/          Riverpod/widget helpers + Mocktail
  foundation_ui/               async states, feedback, responsive layout, media widgets
  foundation_observability/    AppLogger/Talker + CrashReporter/Sentry adapters
  foundation_connectivity/     connectivity service + Riverpod providers
  foundation_permissions/      centralized runtime permission service
  foundation_files/            documents, images, sharing, paths and external URLs
  foundation_forms/            standardized FormBuilder wrappers and validators
  foundation_notifications/    local notification service
  foundation_scanner/          QR/barcode scanner wrapper
  foundation_sync/             connectivity-aware sync job coordinator
  foundation_localization/     shared localization delegates/config
  foundation_device/           app version/build and device diagnostics
  foundation_content/          Markdown, HTML and PDF rendering
  foundation_deeplinks/        app/universal/deep-link stream wrapper
```

Third-party libraries stay behind these packages where practical, so product features depend on foundation APIs rather than plugin-specific APIs.

## Mason generators

```text
flutter_app            create a new application
feature                create a clean feature module
flutter_package        create a reusable workspace package
api_resource           create API resource code
riverpod_controller    create Riverpod controller/provider code
drift_table            create a Drift table
test_suite             create feature test scaffolding
page                   create a presentation page
form                   create a validated form
model                  create a Freezed JSON model
repository             create repository contract + implementation
scanner                create a scanner page
notification           create notification adapter
offline_feature        create a SyncJob boundary
```

## Setup

```bash
fvm install
fvm dart pub global activate mason_cli
fvm dart pub get
fvm dart pub global run mason_cli:mason get
```

## Example workflow

```bash
fvm dart pub global run mason_cli:mason make flutter_app --name mediguide --description "Clinical Guidelines"
# Add `apps/mediguide` under `workspace:` in the root pubspec.yaml.
fvm dart pub global run mason_cli:mason make feature --app mediguide --name guidelines
fvm dart pub global run mason_cli:mason make model --app mediguide --feature guidelines --name guideline
fvm dart pub global run mason_cli:mason make repository --app mediguide --feature guidelines --name guidelines
fvm dart pub global run mason_cli:mason make form --app mediguide --feature guidelines --name guideline_search
fvm dart pub global run mason_cli:mason make page --app mediguide --feature guidelines --name guidelines
fvm dart pub global run mason_cli:mason make offline_feature --app mediguide --feature guidelines --name guidelines
fvm dart pub global run mason_cli:mason make scanner --app mediguide --feature guidelines --name guideline_qr
fvm dart pub global run mason_cli:mason make test_suite --app mediguide --feature guidelines

fvm dart pub get
fvm dart run melos run codegen
fvm dart run melos run check
fvm dart run melos run bricks:validate
```

## UX conventions

- Use `AsyncContent` for loading/error/empty/data presentation.
- Use `BuildContext.success/error/warning/info` instead of ad-hoc SnackBars.
- Use `ResponsiveLayout` and shared breakpoints for phone/tablet/desktop behavior.
- Surface offline state explicitly with `AppOfflineBanner`; never infer API reachability from connectivity type alone.
- Use `foundation_forms` wrappers to keep form appearance and validation consistent.
- Keep camera/files/permissions/scanning behind their foundation services.
- Use `foundation_content` for Markdown/HTML/PDF rendering and `foundation_deeplinks` for URL-driven navigation.
- Use `foundation_device` for support screens that report app version/build and device diagnostics.
- All custom widgets should preserve semantic labels, text scaling and minimum touch targets.

## Observability

`TalkerAppLogger` is the developer/local logging implementation. `CrashReporter` is an abstraction with `SentryCrashReporter` and `NoopCrashReporter`, allowing products to switch reporting providers without coupling feature code to Sentry.

## Offline sync

`foundation_sync` provides a connectivity-aware `SyncEngine` and `SyncJob` contract. Domain features own their own queue/storage schema and conflict policy. This avoids a universal database schema that cannot represent every application's data rules.

## Environments

```bash
./scripts/run_flavor.sh mediguide development
./scripts/run_flavor.sh mediguide staging
```

Do not place secrets in committed Dart-define JSON files. Use secure runtime retrieval or CI secret injection for sensitive values.

## Quality

```bash
fvm dart run melos run format:check
fvm dart run melos run analyze
fvm dart run melos run test
fvm dart run melos run codegen
fvm dart run melos run check
```

## Platform setup still required per app

Some plugins require native configuration (camera/photo/notification permissions, notification icons/channels, background execution, iOS usage descriptions, etc.). Keep that configuration in the generated app's Android/iOS projects and document any product-specific requirements there.

## Command reference

Run Melos commands from this workspace root through FVM: `fvm dart run melos run <task>`. The scripts below resolve their own project root and also accept `--help` without performing setup or validation.

| Task | Work and side effects |
| --- | --- |
| `bootstrap`, `pub:get` | Resolve workspace packages; writes tool caches and may update pubspec.lock |
| `format` / `format:check` | Rewrite/check Dart under apps/ and packages/ |
| `analyze`, `lint` | Analyze every workspace package; lint is an alias |
| `test` | Run Flutter tests in Flutter packages with test/ and Dart tests in non-Flutter packages with test/ |
| `test:coverage` | Flutter coverage only, in Flutter packages with tests; not Dart-only coverage |
| `codegen` / `codegen:watch` | Run build_runner in packages declaring it; writes generated source; watch stays running |
| `check` | Format check → lint → test; no code generation or native app build |
| `quality` | Compatibility alias for check |
| `ci` | Codegen → check → bricks:validate; may download dependencies and change generated source |
| `bricks:list` | List bricks through FVM-activated Mason |
| `bricks:validate` | Render/generate/analyze/test every brick using the script below |
| `outdated` | Report newer packages without upgrading constraints |
| `clean` | Run flutter clean in Flutter packages; preserve source, lockfiles, Dart-only caches and installed SDKs |

- `bash scripts/setup.sh`: require FVM, install the configured SDK, activate Mason CLI in Dart's global package cache, resolve packages, and fetch bricks. It may use the network. Equivalent manual steps are in Setup above.
- `bash scripts/run_flavor.sh <app> [development|staging|production]`: default to development; require apps/<app> and env/<environment>.json; run the app using Dart defines. This selects JSON configuration, not a native `--flavor` build variant. Copy/configure the ignored staging/production JSON from its example first.
- `bash scripts/validate_bricks.sh`: require fvm install and FVM-activated Mason; fetch bricks, render every generator in a temporary workspace, resolve dependencies, format-check, generate code, analyze and test the sample app/package, and remove the temporary workspace on exit. SDK/pub/Mason caches remain.

Generate/register an app before running it: create it with FVM's Mason CLI, add apps/<name> to workspace: in pubspec.yaml, and run fvm dart pub get. Configure native platform folders, permissions and signing inside that app. Build with `fvm flutter build <platform>` from apps/<app>. Workspace check deliberately does not guess a platform or build every app.

Mason commands in these instructions use `fvm dart pub global run mason_cli:mason` so they use the selected SDK. `foundry mason` is an optional wrapper around `mason make` from PATH; use it only when that global executable is configured for the workspace SDK. Shell scripts need Bash; FVM/Dart commands can also be run natively on Windows.

## Organisation identity

Generate with `foundry new portal --template flutter --profile moh-uganda` to use
Ministry organisation details. Add `--identity ./moh.json` for programme-specific
settings. Every generated project contains `identity/identity.json` and
`identity/README.md`. Edit the JSON, preview with `foundry identity sync --dry-run`
from the generated root, then run `foundry identity sync` and the native formatter
and quality gate. Shared layouts and report helpers consume the generated public
metadata. Ownership references stay in the canonical source; never store secrets
there. No approved logos, colours, or fonts are supplied.

## Configurable authentication

OIDC is optional. Create/edit settings with `foundry auth example > oidc.json`,
validate using `foundry auth validate --file oidc.json --template flutter`, then pass
`--auth oidc --auth-config oidc.json` during creation. Generated `auth/README.md`
covers enabling/disabling authentication, callback setup, provider compatibility,
and local Keycloak tests. Edit `auth/auth.json` and run `foundry auth configure`
from the root to update bindings, then format, test, and rebuild.
