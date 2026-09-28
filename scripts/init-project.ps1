param(
    [Parameter(Mandatory = $true)]
    [ValidatePattern('^[A-Za-z0-9._~-]+(?:/[A-Za-z0-9._~-]+)+$')]
    [string]$Module,

    [Parameter(Mandatory = $true)]
    [ValidatePattern('^[a-z0-9][a-z0-9-]*$')]
    [string]$AppName,

    # Omit to retain the kit version pinned by the template.
    [Alias('FrameworkVersion')]
    [ValidatePattern('^v\d+\.\d+\.\d+(?:-[0-9A-Za-z.-]+)?$')]
    [string]$KitVersion
)

$ErrorActionPreference = 'Stop'
$repoRoot = (Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$goMod = Join-Path $repoRoot 'go.mod'
$currentModule = ((Get-Content -LiteralPath $goMod -TotalCount 1) -replace '^module\s+', '').Trim()
$kitModule = 'github.com/hina1314/kit'
$importPattern = '"' + [regex]::Escape($currentModule) + '/'

Get-ChildItem -LiteralPath $repoRoot -Recurse -File -Filter '*.go' |
    Where-Object { $_.FullName -notmatch '[\\/](\.git|\.local|build|dist)[\\/]' } |
    ForEach-Object {
        $content = [IO.File]::ReadAllText($_.FullName)
        $updated = [regex]::Replace($content, $importPattern, '"' + $Module + '/')
        if ($updated -ne $content) {
            [IO.File]::WriteAllText($_.FullName, $updated, [Text.UTF8Encoding]::new($false))
        }
    }

Push-Location $repoRoot
try {
    go mod edit "-module=$Module"
    if ($LASTEXITCODE -ne 0) { throw 'Failed to update application module path.' }
    go mod edit "-dropreplace=$kitModule"
    if ($LASTEXITCODE -ne 0) { throw 'Failed to remove the local kit replacement.' }
    if ($KitVersion) {
        go get "$kitModule@$KitVersion"
        if ($LASTEXITCODE -ne 0) { throw 'Failed to download the requested kit version.' }
    }
    foreach ($path in @('app.env.example', 'ops/docker/app.env.example')) {
        $fullPath = Join-Path $repoRoot $path
        $content = [IO.File]::ReadAllText($fullPath)
        $content = [regex]::Replace($content, '(?m)^APP_NAME=.*$', "APP_NAME=$AppName")
        $content = [regex]::Replace($content, '(?m)^TOKEN_ISSUER=.*$', "TOKEN_ISSUER=$AppName")
        [IO.File]::WriteAllText($fullPath, $content, [Text.UTF8Encoding]::new($false))
    }
    go mod tidy
    if ($LASTEXITCODE -ne 0) { throw 'Failed to tidy application dependencies.' }
    go run github.com/google/wire/cmd/wire ./internal/di
    if ($LASTEXITCODE -ne 0) { throw 'Failed to generate dependency injection.' }
    go test ./...
    if ($LASTEXITCODE -ne 0) { throw 'Application tests failed.' }
} finally {
    Pop-Location
}

Write-Host "Initialized $AppName with module $Module"
Write-Host 'Using github.com/hina1314/kit at the version pinned in go.mod.'
