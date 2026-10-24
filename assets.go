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
var Templates embed.FS
