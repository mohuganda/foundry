# Distribution

Foundry ships as a standalone executable with five embedded source foundations. GitHub Releases, Homebrew, and Docker distribute the same CLI. Template downloads and separate template packages are unnecessary.

## Homebrew

The existing Foundry repository doubles as a custom tap; the explicit repository URL is required when first adding it:

```sh
brew tap mohuganda/foundry https://github.com/mohuganda/foundry
brew install mohuganda/foundry/foundry
foundry version
```

The formula selects the macOS/Linux archive for arm64 or amd64, verifies its SHA-256, and installs the executable and documentation. No Go SDK is needed. The fully qualified formula name avoids collisions with other tools named Foundry.

```sh
brew update
brew upgrade mohuganda/foundry/foundry
brew test mohuganda/foundry/foundry
brew uninstall mohuganda/foundry/foundry
```

Homebrew's formula test checks the version and generates a Rust project without hooks. Generated files stay independent of the installed formula; upgrading/uninstalling the CLI does not upgrade/delete your applications.

## Docker

Images are published to `ghcr.io/mohuganda/foundry` for Linux amd64 and arm64. Docker Desktop runs these Linux images on macOS and Windows. Version tags such as `1.3.0` identify a release; `latest` follows the latest stable GitHub release. Use a pinned version or image digest for repeatable automation.

```sh
docker run --rm ghcr.io/mohuganda/foundry:1.3.0 help
docker run --rm ghcr.io/mohuganda/foundry:1.3.0 list
```

### Generate into the host directory

On macOS/Linux, run from the parent directory where the new project belongs:

```sh
docker run --rm --user "$(id -u):$(id -g)" \
  --mount "type=bind,source=$(pwd),target=/workspace" \
  ghcr.io/mohuganda/foundry:1.3.0 \
  new customer-mobile --template flutter --no-hooks
```

`/workspace` is the image's working directory. The bind mount persists generated output on the host; `--rm` removes the container afterward. Supplying your UID/GID gives the files your host ownership. The image defaults to non-root UID/GID 65532; on Linux, grant that user access or use the UID/GID override above. The host directory must be writable by the selected user.

Native Windows PowerShell with Docker Desktop in Linux-container mode:

```powershell
docker run --rm --mount "type=bind,source=$($PWD.Path),target=/workspace" `
  ghcr.io/mohuganda/foundry:1.3.0 new dashboard --template react --no-hooks
```

Docker Desktop handles host file access on Windows. For interactive creation on macOS/Linux, add `-it` and run the executable without arguments, overriding the image's default help command:

```sh
docker run --rm -it --user "$(id -u):$(id -g)" \
  --mount "type=bind,source=$(pwd),target=/workspace" \
  --entrypoint "" ghcr.io/mohuganda/foundry:1.3.0 /usr/local/bin/foundry
```

The image's default command is `help`; the example explicitly runs the executable without arguments to select interactive creation.

### Extend or inspect a mounted project

```sh
# Run from the parent containing my-api/:
docker run --rm --user "$(id -u):$(id -g)" \
  --mount "type=bind,source=$(pwd),target=/workspace" \
  ghcr.io/mohuganda/foundry:1.3.0 add service orders --project my-api
```

`--project` and `--destination` refer to container paths. A destination outside the mounted directory will not persist after the container is removed. The image supports generation, built-in Go/Flutter/web/Rust scaffolding, and read-only template version checks. It contains no shell, Go/Flutter/Node/Rust SDK, Mason executable, database, or Docker daemon. Use `--no-hooks` when generating; develop and validate the result using the host's native tools. `doctor` inside the image reports container tooling, rather than the host's environment.

### Build and check locally

From the Foundry repository with Docker running:

```sh
docker build --build-arg VERSION=dev -t foundry:test .
./scripts/test-docker.sh foundry:test
```

The multi-stage Dockerfile compiles a static CLI and copies it into a minimal scratch image. Build context exclusions omit local dependencies and build outputs. `test-docker.sh` runs version/help and generates all five templates into a mounted temporary directory with networking disabled, checking manifests and executable script permissions. It cleans up the temporary directory on exit. Docker is required; generated-project SDKs are not.

## Release automation

1. Merge validated code, create a stable `vX.Y.Z` tag, and push it.
2. `Release` checks Foundry, builds six executable archives plus `SHA256SUMS`, and publishes the GitHub release.
3. `Distribute release` runs after that successful tag-triggered workflow. It verifies a published stable release exists, builds/tests its CLI image, then publishes Linux amd64/arm64 images to GHCR.
4. For the latest stable release, it also updates `Formula/foundry.rb` on the default branch using that release's checksums. Historical backfills publish only their versioned image and preserve the latest formula/image alias.

The distribution workflow uses the repository's `GITHUB_TOKEN` with `packages: write` and `contents: write`. No cross-repository token or separate tap repository is needed. The default branch must permit its formula-update commit; if branch rules prohibit it, apply the generated formula through a normal reviewed PR instead. Formula-only automation commits include `[skip ci]`; CI already tests the generator and installation of the checked-in formula.

On first GHCR publication, check the package's settings and ensure visibility is **Public** for anonymous pulls. GitHub initially creates packages as private; the image's source label links it to this repository. Repository/workflow policies may restrict package creation or token permissions.

To backfill an already published stable release, or retry a failed distribution run:

```sh
gh workflow run distribute.yml --repo mohuganda/foundry -f version=v1.3.0
```

Manual `Release` workflow runs only validate archives; they do not start automatic distribution. Stable tags are supported; prerelease tags are not distributed by this workflow.

### Regenerate a formula manually

Download the release checksum manifest, then run from the Foundry repository:

```sh
gh release download v1.3.0 --repo mohuganda/foundry \
  --pattern SHA256SUMS --dir /tmp/foundry-homebrew
./scripts/homebrew-formula.sh 1.3.0 /tmp/foundry-homebrew/SHA256SUMS > Formula/foundry.rb
ruby -c Formula/foundry.rb
```

The generator accepts a stable version with or without its `v` prefix, requires one valid SHA-256 per supported macOS/Linux archive, and writes the formula to stdout. Review the diff and commit it through the usual workflow. Archive URLs and hashes always refer to an existing release; building local snapshots does not change the tap.
