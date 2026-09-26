param(
    [Parameter(Mandatory = $true)]
    [string]$ArchivePath,

    [Parameter(Mandatory = $true)]
    [string]$ChecksumsPath,

    [Parameter(Mandatory = $true)]
    [string]$ConfigPath,

    [Parameter(Mandatory = $true)]
    [ValidateSet(3001, 3002)]
    [int]$TargetPort
)

$ErrorActionPreference = 'Stop'

$archive = (Resolve-Path -LiteralPath $ArchivePath -ErrorAction Stop).Path
$checksums = (Resolve-Path -LiteralPath $ChecksumsPath -ErrorAction Stop).Path
$config = (Resolve-Path -LiteralPath $ConfigPath -ErrorAction Stop).Path
$archiveName = [System.IO.Path]::GetFileName($archive)

if ($archiveName -notmatch '^study-api_(v\d+\.\d+\.\d+)_windows_amd64\.zip$') {
    throw "不是预期的 Windows Release 包：$archiveName"
}
$version = $Matches[1]

$checksumLine = @(Get-Content -LiteralPath $checksums | Where-Object {
    $_ -match ('^[0-9a-fA-F]{64}\s+\*?' + [regex]::Escape($archiveName) + '$')
})
if ($checksumLine.Count -ne 1) {
    throw "checksums.txt 中必须恰好有一条 $archiveName 的校验记录。"
}
$expectedHash = $checksumLine[0].Substring(0, 64).ToLowerInvariant()
$actualHash = (Get-FileHash -LiteralPath $archive -Algorithm SHA256).Hash.ToLowerInvariant()
if ($actualHash -ne $expectedHash) {
    throw "SHA256 校验失败：$archiveName；不会启动候选实例。"
}

if (Get-NetTCPConnection -State Listen -LocalPort $TargetPort -ErrorAction SilentlyContinue) {
    throw "端口 $TargetPort 已被占用；先确认它不是当前服务，不会覆盖或停止现有进程。"
}

$repoRoot = (Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$releaseRoot = Join-Path $repoRoot '.local\releases'
$stageDir = Join-Path $releaseRoot "$version-$TargetPort"
if (Test-Path -LiteralPath $stageDir) {
    throw "候选目录已经存在，不会覆盖：$stageDir"
}

$configText = [System.IO.File]::ReadAllText($config)
$addressPattern = [regex]::new('(?m)^SERVER_ADDRESS\s*=.*$')
if ($addressPattern.Matches($configText).Count -ne 1) {
    throw '配置文件必须恰好有一条 SERVER_ADDRESS=...；不会猜测监听地址。'
}
$configText = $addressPattern.Replace($configText, "SERVER_ADDRESS=127.0.0.1:$TargetPort", 1)

New-Item -ItemType Directory -Path $stageDir -Force | Out-Null
Expand-Archive -LiteralPath $archive -DestinationPath $stageDir
$appDir = Join-Path $stageDir "study-api_${version}_windows_amd64"
$exePath = Join-Path $appDir 'study-api.exe'
if (-not (Test-Path -LiteralPath $exePath -PathType Leaf)) {
    throw "包中未找到可执行文件：$exePath"
}

[System.IO.File]::WriteAllText(
    (Join-Path $appDir 'app.env'),
    $configText,
    [System.Text.UTF8Encoding]::new($false)
)

# 防止当前 PowerShell 的同名环境变量覆盖候选目录中的 app.env。
$previousAddress = [Environment]::GetEnvironmentVariable('SERVER_ADDRESS', 'Process')
try {
    [Environment]::SetEnvironmentVariable('SERVER_ADDRESS', "127.0.0.1:$TargetPort", 'Process')
    $process = Start-Process -FilePath $exePath -WorkingDirectory $appDir `
        -WindowStyle Hidden -PassThru `
        -RedirectStandardOutput (Join-Path $stageDir 'stdout.log') `
        -RedirectStandardError (Join-Path $stageDir 'stderr.log')
} finally {
    [Environment]::SetEnvironmentVariable('SERVER_ADDRESS', $previousAddress, 'Process')
}

$process.Id | Set-Content -LiteralPath (Join-Path $stageDir 'pid.txt')
for ($attempt = 0; $attempt -lt 30; $attempt++) {
    if ($process.HasExited) {
        throw "候选进程已退出（exit $($process.ExitCode)）；查看 $stageDir\stderr.log 和 stdout.log。"
    }
    $httpCode = & curl.exe --noproxy '*' --silent --output NUL --write-out '%{http_code}' `
        --max-time 2 "http://127.0.0.1:$TargetPort/readyz" 2>$null
    if ($LASTEXITCODE -eq 0 -and $httpCode -eq '200') {
        $listener = Get-NetTCPConnection -State Listen -LocalPort $TargetPort -ErrorAction SilentlyContinue |
            Where-Object { $_.OwningProcess -eq $process.Id }
        if (-not $listener) {
            throw "端口 $TargetPort 返回了 200，但监听进程不是候选 PID $($process.Id)；不会切流。"
        }
        Write-Host "候选实例已就绪：$version -> 127.0.0.1:$TargetPort (PID $($process.Id))"
        Write-Host "目录：$appDir"
        Write-Host '当前只完成候选实例准备；8080 尚未切流。'
        exit 0
    }
    Start-Sleep -Seconds 1
}

throw "候选实例在 30 秒内未 ready（PID $($process.Id)）；查看 $stageDir\stderr.log 和 stdout.log。8080 未切流。"
