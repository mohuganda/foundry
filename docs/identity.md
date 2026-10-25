# Organisation and Ministry identity

All five foundations include a non-secret identity source at identity/identity.json.
The default profile is generic. Use moh-uganda for Ministry of Health / Republic of
Uganda names, the MoH short name, and https://health.go.ug. No approved visual assets
or official colour/font specifications are bundled. Neutral visual defaults remain
marked pending; supply approved material and record its review evidence yourself.

```sh
foundry identity example --profile moh-uganda --name "Programme Portal" > moh.json
# Edit application, branding, documents, links, navigation, and ownership.
foundry identity validate --file moh.json
foundry new programme-portal --template react --profile moh-uganda --identity moh.json
```

Partial override objects are supported. Asset paths in an input file are relative
to that file's folder. Generation copies images into identity/assets/ and writes
normalized references to the saved canonical JSON. SVG, PNG, JPEG, and ICO are
supported up to 5 MiB; image symlinks and active SVG content are rejected.

After editing the generated source:

```sh
foundry identity sync --project ./programme-portal --dry-run
foundry identity sync --project ./programme-portal
```

Sync derives public metadata, web design tokens/assets, Flutter package metadata,
Go identity APIs, and Rust constants/report helpers. It does not overwrite shared
UI components, wire authentication, configure native launcher icons, or delete
unused old asset files. Run the native formatter and quality gate after syncing.

The source contains organisation details, application identity, optional asset slots,
light/dark palettes, typography/spacing/borders, branding approval declarations,
document headers/footers, public links/navigation, and ownership references.
Ownership/operational references and private approval evidence are excluded from
public bindings. Keep the source itself free of credentials, patient data, private
keys, and other secrets: it is intended to be committed to the application repository.
Public application metadata is build-time data; environment labels do not configure
runtime credentials, endpoints, or deployment isolation.

Branding approval status is a declaration. Foundry cannot verify Ministry approval.

## Configuration sections

| Section | Configure | Consumption |
| --- | --- | --- |
| `organisation` | `name`, `country`, `short_name`, `website` | Ministry header, public links, reports |
| `branding.assets` | `logo_light`, `logo_dark`, `coat_of_arms`, `application_icon`, `favicon` | Optional images; web tab/touch icons; Flutter package assets |
| `branding.tokens` | Light/dark palettes, font family/size, small/medium/large spacing, radius, border width | Web CSS variables and Flutter theme |
| `application` | Name, description, department, programme, support email, environment | Product title, shared shell/footer, support links, nonproduction badge |
| `documents` | Title prefix, footer text, default version, `yyyy-MM-dd` date format | Report headers and export metadata |
| `links` / `navigation` | Privacy, accessibility, support URLs; labelled navigation destinations | Header navigation and footer links |
| `ownership` | Technical/service owner, repository URL, operational contact references | Repository metadata only, excluded from public bindings |
| `branding.approval` | Status, reference, approver, ISO approval date | Review declaration; reference and approver excluded from public bindings |

For example, save this partial override as `moh.json` beside any supplied assets:

```json
{
  "profile": "moh-uganda",
  "application": {
    "name": "Programme Portal",
    "description": "Programme services and reporting",
    "department": "Your department",
    "programme": "Your programme",
    "environment": "staging"
  },
  "documents": {
    "title_prefix": "Programme — ",
    "footer_text": "Ministry of Health, Republic of Uganda",
    "default_version": "draft",
    "date_format": "yyyy-MM-dd"
  },
  "navigation": [{ "label": "Home", "href": "/" }],
  "links": { "privacy": "/privacy", "accessibility": "/accessibility" },
  "ownership": {
    "technical_owner": "Team reference",
    "service_owner": "Programme owner reference",
    "operational_contacts": ["Operations directory reference"]
  }
}
```

The referenced routes/pages must be implemented by your application. Use a real
support email or HTTPS support link when available; no contact details are invented.
Environments are `development`, `test`, `staging`, or `production`. Only production
hides the environment badge. Omit unconfigured assets and links or leave them empty.
Use `foundry identity example` for the complete schema and both token palettes.

Colours use six-digit hex values. Palettes expose background, surface, text, muted,
primary, on-primary, focus, border, success, and error colours. Review contrast after
changing them. Web tokens map to CSS custom properties in
`src/shared/identity/tokens.css`; Flutter maps them to `identityTheme`. Font family
names configure selection, not font downloads: supply licensed font files and
register them separately in CSS or Flutter's pubspec. Native application/launcher
icons also require platform-specific setup. Web application icons map to the touch
icon link; Foundry does not generate a PWA manifest or resized icons.

An `approved` declaration requires a reference, approver, and `yyyy-MM-dd` date.
Until approved material exists, keep the default `pending` declaration. The JSON
file is limited to 1 MiB. Asset paths must be local relative paths without symlinks.
Passive SVG permits shapes/text but rejects scripts, event handlers, styles,
external references, animation, and URL-valued attributes; use PNG/JPEG if a complex
approved SVG cannot meet these restrictions. ICO is useful for favicons; prefer
SVG/PNG/JPEG for Flutter-rendered logos, whose platform codecs differ.

## Shared layouts and APIs

### Vanilla JavaScript and React

The generated home page already uses the identity shell. Application identity is
read from the generated public config, rather than a separate product-name setting
in environment configuration. `applyDocumentIdentity` sets the document title,
description, favicon, and touch icon at startup. Vite's configured base path is
respected for copied assets.

Vanilla exports from `src/shared/identity/dom.js` include
`createIdentityShell(content)`, `createIdentityState(kind, options)`, and
`createIdentityReportHeader(title, options)`. Each returns an element and a destroy
function; state kinds are `loading`, `error`, and `sign-in`. Wire `onAction` to a real
retry/sign-in handler and call destroy when replacing a view.

React exports `IdentityShell`, `IdentityLoading`, `IdentityError`, `IdentitySignIn`,
and `IdentityReportHeader` from `src/shared/identity/Identity.jsx`:

```jsx
<IdentityShell>
  <IdentitySignIn onSignIn={startSignIn} />
</IdentityShell>
```

The error component accepts `message` and `onRetry`. The loading component accepts
`message`. Sign-in/retry buttons remain disabled until a callback is supplied.
These are presentation components; connect the application's authentication and
routing. Neither starter implements a new authentication provider.

Report headers accept `version` and `issuedOn`; pass the report date explicitly:

```jsx
<article>
  <IdentityReportHeader title="Quarterly report" version="1.0" issuedOn="2026-10-25" />
  <p>Report content</p>
</article>
```

`reportMetadata(title, { version, issuedOn })` from `index.js` also supplies a plain
metadata object and export title for a separate export implementation. Print styles
hide the application header/navigation/footer, retain document branding, and use a
light palette. PDF/CSV/export engines are application responsibilities.

### Flutter workspaces

Mason-generated apps depend on `foundation_identity`. Their `AppRoot` uses the
shared light/dark theme and wraps router content in `IdentityShell`. Internal paths
go through the app router; HTTPS/email links use the platform URL launcher with
safe failure feedback. Add routes before configuring their navigation links.

The package exports `appIdentity`, `identityTheme`, `IdentityShell`,
`IdentitySignIn`, `IdentityLoading`, `IdentityError`, and `IdentityReportHeader`.
Pass identity to the sign-in/report widgets, `onSignIn` to sign-in, and `onRetry` to
errors. For a workspace with several apps, override application name/description
at the app boundary while retaining shared organisation and branding:

```dart
AppRoot(identity: appIdentity.withApplication(name: 'Programme Mobile'))
```

The reusable shell accepts `onOpenLink` for other hosts. Reports can use
`IdentityReportHeader(identity: appIdentity, title: 'Quarterly report',
version: '1.0', issuedOn: '2026-10-25')` or `appIdentity.reportHeader(...)` for plain
text metadata. Flutter print/PDF integration requires an application export library;
the web print stylesheet does not apply to Flutter.

### Go backends

The generated `identity` package embeds public metadata. `identity.Load()` returns
the typed profile and an error; the server uses application/environment as logging
context. `identity.JSON()` returns a copy suitable for an explicitly implemented
public metadata endpoint. Foundry does not add an unauthenticated endpoint.

`identity.ReportHeader("Quarterly report", "1.0", "2026-10-25")` returns structured
header fields and an error. Escape these values for the destination document
format. Go backends expose metadata rather than browser layouts.

### Rust applications

`src/identity.rs` exports application/organisation constants, `PUBLIC_JSON`, and
`report_header(title, version, issued_on)`. `cargo run --locked -- identity` prints
public identity JSON. The report helper returns plain text; integrate it with your
chosen export format. Rust's CLI template does not include a graphical UI.

## Keeping generated bindings current

1. Edit `identity/identity.json`, or copy approved assets under `identity/assets/`
   and update their references relative to the identity folder.
2. Run `foundry identity validate --file identity/identity.json`.
3. Preview `foundry identity sync --dry-run`, then run `foundry identity sync` from
   the generated root (or supply `--project /path/to/root`).
4. Run the native formatter and quality gate, review the diff, and rebuild/redeploy.

Sync rewrites canonical normalized JSON, `identity/public.json`, and the platform's
marked generated bindings/assets. Shared editable layout components are created
at initial generation and left alone during sync. Existing projects need those
components installed or a fresh foundation before these layouts are available;
sync is not a template upgrade command. It never deletes stale assets: remove
unused source/public/package copies explicitly after checking references.

Public bindings omit ownership, operational contacts, and approval reference/name.
The repository source is still readable by repository collaborators. Shipping the
entire source tree in a container or static directory can expose it; build/deploy
only the intended application output. No arbitrary ownership data should be sent
to browser clients, logs, reports, or public endpoints.
