# Application identity

Edit identity/identity.json, then run foundry identity sync from the project root.
Use --dry-run to preview generated bindings. The source is non-secret configuration.
Do not put passwords, tokens, patient information, or private keys in this file.

The moh-uganda profile supplies organisation names and the Ministry website only.
No approved official logos, coat of arms, fonts, or colour specifications ship here.
Neutral theme values are examples; approval remains pending until your organisation
supplies approved assets and records the reference, approver, and date.

All asset paths in the saved identity JSON are relative to identity/. Asset images
are copied into public web or Flutter package assets during sync. Only passive SVG,
PNG, JPEG, and ICO files up to 5 MiB are supported. SVG styles, scripts, embedded HTML,
external references, event handlers, and symlink inputs are rejected. Native mobile
launcher icons and signing still require platform setup; an application-icon slot
does not change Android/iOS launcher resources automatically.

Ownership and operational-contact references remain in the source JSON and are
excluded from public browser/mobile/Rust metadata and the Go public API. They must
still contain references or business contacts only: this source file is versioned.
Report dates/versions are explicit application inputs, separate from build dates.
Links accept HTTPS, application /paths, and #fragments. Empty optional links/assets
are omitted by UI components. Sign-in layouts require your actual auth callback.

Sync updates generated configuration, tokens, bindings, and copied assets. It does
not rewrite your shared UI components or remove old unused image files. Remove
obsolete copied assets explicitly before deployment. After syncing, run the native
formatter and quality gate. Branding approval metadata records your declaration;
Foundry cannot verify Ministry authorisation or approve assets on your behalf.
