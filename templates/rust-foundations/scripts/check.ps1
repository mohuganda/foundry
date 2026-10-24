[CmdletBinding()]
param([switch]$Help)

$ErrorActionPreference = 'Stop'
if ($Help) {
    Write-Output 'Usage: scripts/check.ps1 [-Help]'
    Write-Output 'Checks formatting, Clippy, all tests, documentation, and the release build in order.'
    Write-Output 'Set RUSTUP_TOOLCHAIN to select a compiler or CARGO_NET_OFFLINE=true to disable downloads.'
    exit 0
}

function Invoke-Cargo {
    param([string[]]$CargoArgs)
    & cargo @CargoArgs
    if ($LASTEXITCODE -ne 0) {
        throw "cargo $($CargoArgs -join ' ') failed with exit code $LASTEXITCODE"
    }
}

$previousRustdocFlags = $env:RUSTDOCFLAGS
Push-Location (Split-Path -Parent $PSScriptRoot)
try {
    Invoke-Cargo -CargoArgs @('fmt', '--all', '--', '--check')
    Invoke-Cargo -CargoArgs @('lint')
    Invoke-Cargo -CargoArgs @('test', '--locked')
    $env:RUSTDOCFLAGS = "$previousRustdocFlags -D warnings"
    Invoke-Cargo -CargoArgs @('doc', '--locked', '--no-deps')
    Invoke-Cargo -CargoArgs @('build', '--release', '--locked')
}
finally {
    $env:RUSTDOCFLAGS = $previousRustdocFlags
    Pop-Location
}
