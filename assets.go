package foundry

import "embed"

// Templates contains the bundled project foundations as source files.
//
//go:embed all:templates/go-foundations/.gitignore
//go:embed all:templates/go-foundations/Dockerfile
//go:embed all:templates/go-foundations/FILES.txt
//go:embed all:templates/go-foundations/Makefile
//go:embed all:templates/go-foundations/README.md
//go:embed all:templates/go-foundations/api
//go:embed all:templates/go-foundations/buf.gen.yaml
//go:embed all:templates/go-foundations/buf.lock
//go:embed all:templates/go-foundations/buf.yaml
//go:embed all:templates/go-foundations/cmd
//go:embed all:templates/go-foundations/compose.yaml
//go:embed all:templates/go-foundations/configs
//go:embed all:templates/go-foundations/deploy
//go:embed all:templates/go-foundations/gen
//go:embed all:templates/go-foundations/go.mod.tmpl
//go:embed all:templates/go-foundations/go.sum
//go:embed all:templates/go-foundations/internal
//go:embed all:templates/go-foundations/migrations
//go:embed all:templates/go-foundations/scripts
//go:embed all:templates/go-foundations/secrets
//go:embed all:templates/flutter_foundations/.env.example
//go:embed all:templates/flutter_foundations/.fvmrc
//go:embed all:templates/flutter_foundations/.gitignore
//go:embed all:templates/flutter_foundations/README.md
//go:embed all:templates/flutter_foundations/VERSION
//go:embed all:templates/flutter_foundations/analysis_options.yaml
//go:embed all:templates/flutter_foundations/apps
//go:embed all:templates/flutter_foundations/bricks
//go:embed all:templates/flutter_foundations/env
//go:embed all:templates/flutter_foundations/mason.yaml
//go:embed templates/flutter_foundations/packages
//go:embed all:templates/flutter_foundations/pubspec.lock
//go:embed all:templates/flutter_foundations/pubspec.yaml
//go:embed all:templates/flutter_foundations/scripts
//go:embed templates/vanilla-js-foundations/src templates/vanilla-js-foundations/index.html templates/vanilla-js-foundations/package*.json
//go:embed templates/vanilla-js-foundations/*.config.js templates/vanilla-js-foundations/README.md all:templates/vanilla-js-foundations/public
//go:embed templates/vanilla-js-foundations/.editorconfig templates/vanilla-js-foundations/.env.example templates/vanilla-js-foundations/.gitignore
//go:embed templates/vanilla-js-foundations/.nvmrc templates/vanilla-js-foundations/.prettierignore templates/vanilla-js-foundations/.prettierrc.json
//go:embed templates/react-foundations/src templates/react-foundations/index.html templates/react-foundations/package*.json
//go:embed templates/react-foundations/*.config.js templates/react-foundations/README.md all:templates/react-foundations/public
//go:embed templates/react-foundations/.editorconfig templates/react-foundations/.env.example templates/react-foundations/.gitignore
//go:embed templates/react-foundations/.nvmrc templates/react-foundations/.prettierignore templates/react-foundations/.prettierrc.json
var Templates embed.FS
