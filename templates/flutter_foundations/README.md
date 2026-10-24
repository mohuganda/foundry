# Flutter Foundation Workspace v3

A reusable production-oriented Flutter monorepo foundation built around **Melos + Mason**, with Riverpod, go_router, Dio, Freezed/JSON, Drift-ready offline boundaries, secure storage, typed environments, design-system primitives, observability, connectivity awareness, files/permissions, forms, notifications, scanning, localization and reusable test helpers.

The workspace is pinned to Flutter 3.44.8 through FVM. Install [FVM](https://fvm.app/) before running the setup commands; use `fvm flutter` and `fvm dart` rather than system SDK commands.

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
mason make flutter_app --name mediguide --description "Clinical Guidelines"
# Add `apps/mediguide` under `workspace:` in the root pubspec.yaml.
mason make feature --app mediguide --name guidelines
mason make model --app mediguide --feature guidelines --name guideline
mason make repository --app mediguide --feature guidelines --name guidelines
mason make form --app mediguide --feature guidelines --name guideline_search
mason make page --app mediguide --feature guidelines --name guidelines
mason make offline_feature --app mediguide --feature guidelines --name guidelines
mason make scanner --app mediguide --feature guidelines --name guideline_qr
mason make test_suite --app mediguide --feature guidelines

fvm dart pub get
fvm dart run melos run codegen
fvm dart run melos run quality
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
fvm flutter run --project-dir apps/mediguide --dart-define-from-file=env/staging.json
```

Do not place secrets in committed Dart-define JSON files. Use secure runtime retrieval or CI secret injection for sensitive values.

## Quality

```bash
fvm dart run melos run format:check
fvm dart run melos run analyze
fvm dart run melos run test
fvm dart run melos run codegen
fvm dart run melos run quality
```

## Platform setup still required per app

Some plugins require native configuration (camera/photo/notification permissions, notification icons/channels, background execution, iOS usage descriptions, etc.). Keep that configuration in the generated app's Android/iOS projects and document any product-specific requirements there.
