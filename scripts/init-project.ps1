param(
    [Parameter(Mandatory = $true)]
    [ValidatePattern('^[A-Za-z0-9._~-]+(?:/[A-Za-z0-9._~-]+)+$')]
    [string]$Module,

    [Parameter(Mandatory = $true)]
    [ValidatePattern('^[a-z0-9][a-z0-9-]*$')]
    [string]$AppName
)

$ErrorActionPreference = 'Stop'
$repoRoot = (Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$goMod = Join-Path $repoRoot 'go.mod'
$currentModule = ((Get-Content -LiteralPath $goMod -TotalCount 1) -replace '^module\s+', '').Trim()

Get-ChildItem -LiteralPath $repoRoot -Recurse -File -Filter '*.go' |
    Where-Object { $_.FullName -notmatch '[\\/]\.git[\\/]' } |
    ForEach-Object {
        $content = [IO.File]::ReadAllText($_.FullName)
        $updated = $content.Replace('"' + $currentModule + '/', '"' + $Module + '/')
        if ($updated -ne $content) {
            [IO.File]::WriteAllText($_.FullName, $updated, [Text.UTF8Encoding]::new($false))
        }
    }

Push-Location $repoRoot
try {
    go mod edit "-module=$Module"
    foreach ($path in @('app.env.example', 'ops/docker/app.env.example')) {
        $fullPath = Join-Path $repoRoot $path
        $content = [IO.File]::ReadAllText($fullPath)
        $content = [regex]::Replace($content, '(?m)^APP_NAME=.*$', "APP_NAME=$AppName")
        $content = [regex]::Replace($content, '(?m)^TOKEN_ISSUER=.*$', "TOKEN_ISSUER=$AppName")
        [IO.File]::WriteAllText($fullPath, $content, [Text.UTF8Encoding]::new($false))
    }
    go mod tidy
    go run github.com/google/wire/cmd/wire@v0.6.0 ./internal/di
    go test ./...
} finally {
    Pop-Location
}

Write-Host "Initialized $AppName with module $Module"
