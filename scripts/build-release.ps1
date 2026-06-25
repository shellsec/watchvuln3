# 交叉编译 Windows / Linux / macOS 发布包（纯 Go，无需 CGO）
# 用法: .\scripts\build-release.ps1
#       .\scripts\build-release.ps1 -Version v3.1.0

param(
    [string]$Version = "v3.1.0"
)

$ErrorActionPreference = "Stop"
$root = Split-Path -Parent $PSScriptRoot
$dist = Join-Path $root "dist"
New-Item -ItemType Directory -Force -Path $dist | Out-Null

$ldflags = "-s -w -X main.Version=$Version"

$targets = @(
    @{ GOOS = "windows"; GOARCH = "amd64"; Out = "watchvuln-windows-amd64.exe" },
    @{ GOOS = "linux";   GOARCH = "amd64"; Out = "watchvuln-linux-amd64" },
    @{ GOOS = "linux";   GOARCH = "arm64"; Out = "watchvuln-linux-arm64" },
    @{ GOOS = "darwin";  GOARCH = "amd64"; Out = "watchvuln-darwin-amd64" },
    @{ GOOS = "darwin";  GOARCH = "arm64"; Out = "watchvuln-darwin-arm64" }
)

Push-Location $root
try {
    foreach ($t in $targets) {
        $out = Join-Path $dist $t.Out
        Write-Host "==> $($t.GOOS)/$($t.GOARCH) -> dist/$($t.Out)"
        $env:GOOS = $t.GOOS
        $env:GOARCH = $t.GOARCH
        $env:CGO_ENABLED = "0"
        go build -trimpath -ldflags $ldflags -o $out .
        if ($LASTEXITCODE -ne 0) { exit $LASTEXITCODE }
    }
    Write-Host ""
    Write-Host "Done. Artifacts in dist/:"
    Get-ChildItem $dist | ForEach-Object { Write-Host ("  {0,12}  {1}" -f ($_.Length), $_.Name) }
}
finally {
    Remove-Item Env:GOOS, Env:GOARCH, Env:CGO_ENABLED -ErrorAction SilentlyContinue
    Pop-Location
}
