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
