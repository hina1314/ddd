param(
    [Parameter(Mandatory = $true)]
    [ValidateSet(3001, 3002)]
    [int]$TargetPort,

    [string]$NginxExe = 'D:\dev\EServer\core\childApp\server\nginx\nginx.exe'
)

$ErrorActionPreference = 'Stop'

$repoRoot = (Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$configPath = Join-Path $repoRoot 'ops\nginx\nginx.lab.conf'
$prefixPath = Join-Path $repoRoot '.local\nginx-lab'
$prefixArg = $prefixPath.Replace('\', '/').TrimEnd('/') + '/'
$configArg = $configPath.Replace('\', '/')

if (-not (Test-Path -LiteralPath $NginxExe -PathType Leaf)) {
    throw "找不到 Nginx: $NginxExe"
}
if (-not (Test-Path -LiteralPath $configPath -PathType Leaf)) {
    throw "找不到实验配置: $configPath"
}
if (-not (Test-Path -LiteralPath (Join-Path $prefixPath 'logs\nginx.pid') -PathType Leaf)) {
    throw '实验 Nginx 尚未运行，找不到 PID 文件。'
}

function Get-ReadyState {
    param([int]$Port)

    $output = & curl.exe --noproxy '*' --silent --include --max-time 3 "http://127.0.0.1:$Port/readyz" 2>$null
    if ($LASTEXITCODE -ne 0) {
        return [pscustomobject]@{ Status = 0; Upstream = '' }
    }

    $response = $output -join "`n"
    $status = [regex]::Match($response, '(?m)^HTTP/\S+\s+(\d{3})\b')
    $upstream = [regex]::Match($response, '(?im)^X-Lab-Upstream:\s*([^\r\n]+)')

    return [pscustomobject]@{
        Status   = if ($status.Success) { [int]$status.Groups[1].Value } else { 0 }
        Upstream = if ($upstream.Success) { $upstream.Groups[1].Value.Trim() } else { '' }
    }
}

function Invoke-Nginx {
    param([string[]]$Arguments)

    & $NginxExe @Arguments -p $prefixArg -c $configArg
    if ($LASTEXITCODE -ne 0) {
        throw "Nginx 命令失败: $($Arguments -join ' ') (exit $LASTEXITCODE)"
    }
}

function Wait-ForRoute {
    param([int]$ExpectedPort)

    for ($attempt = 0; $attempt -lt 30; $attempt++) {
        $state = Get-ReadyState -Port 8080
        if ($state.Status -eq 200 -and $state.Upstream -eq "127.0.0.1:$ExpectedPort") {
            return $true
        }
        Start-Sleep -Milliseconds 200
    }
    return $false
}

$original = [System.IO.File]::ReadAllText($configPath)
$routePattern = [regex]::new('proxy_pass\s+http://127\.0\.0\.1:(3001|3002);')
$routes = $routePattern.Matches($original)
if ($routes.Count -ne 1) {
    throw '实验配置必须恰好包含一条指向 3001 或 3002 的 proxy_pass。'
}

$currentPort = [int]$routes[0].Groups[1].Value
$currentState = Get-ReadyState -Port 8080
if ($currentState.Status -eq 0 -or $currentState.Upstream -ne "127.0.0.1:$currentPort") {
    throw "当前代理状态与配置不一致：配置=$currentPort，HTTP=$($currentState.Status)，上游=$($currentState.Upstream)"
}
if ($currentState.Status -ne 200) {
    Write-Warning "当前实例不健康（HTTP $($currentState.Status)）；只要目标实例健康，仍允许紧急回滚。"
}

$targetState = Get-ReadyState -Port $TargetPort
if ($targetState.Status -ne 200) {
    throw "目标实例 127.0.0.1:$TargetPort 尚未 ready，停止切流。"
}

if ($currentPort -eq $TargetPort) {
    Write-Host "当前已指向健康的 127.0.0.1:$TargetPort，无需切流。"
    exit 0
}

$backupDir = Join-Path $prefixPath 'backups'
New-Item -ItemType Directory -Path $backupDir -Force | Out-Null
$backupPath = Join-Path $backupDir ("nginx.lab.conf.$(Get-Date -Format 'yyyyMMddHHmmssfff').bak")
[System.IO.File]::Copy($configPath, $backupPath)

$updated = $routePattern.Replace($original, "proxy_pass http://127.0.0.1:$TargetPort;", 1)
$utf8 = [System.Text.UTF8Encoding]::new($false)

try {
    [System.IO.File]::WriteAllText($configPath, $updated, $utf8)
    Invoke-Nginx -Arguments @('-t')
    Invoke-Nginx -Arguments @('-s', 'reload')

    if (-not (Wait-ForRoute -ExpectedPort $TargetPort)) {
        throw "切流后 8080 未通过验收：目标端口 $TargetPort。"
    }

    Write-Host "切流成功：8080 -> $TargetPort；原配置备份在 $backupPath"
} catch {
    $failure = $_.Exception.Message
    [System.IO.File]::WriteAllText($configPath, $original, $utf8)

    try {
        Invoke-Nginx -Arguments @('-t')
        Invoke-Nginx -Arguments @('-s', 'reload')
        if (-not (Wait-ForRoute -ExpectedPort $currentPort)) {
            throw "回滚后 8080 未恢复到 $currentPort。"
        }
    } catch {
        throw "切流失败：$failure；自动回滚也失败：$($_.Exception.Message)。备份：$backupPath"
    }

    throw "切流失败：$failure；已恢复 8080 -> $currentPort。备份：$backupPath"
}
