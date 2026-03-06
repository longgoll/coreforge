# ──────────────────────────────────────────────────────────────
# Forge CLI — Cross-compile Build Script (Windows PowerShell)
# ──────────────────────────────────────────────────────────────
# Usage: .\build.ps1 [-All] [-Version "1.0.0"]

param(
    [switch]$All,
    [string]$Version = "dev"
)

$ErrorActionPreference = "Stop"
$APP_NAME = "forge"
$BUILD_DIR = "dist"
$BUILD_DATE = (Get-Date -Format "yyyy-MM-ddTHH:mm:ssZ")

# Try to get version from git tags
if ($Version -eq "dev") {
    try {
        $Version = (git describe --tags --always --dirty 2>$null)
        if (-not $Version) { $Version = "dev" }
    } catch {
        $Version = "dev"
    }
}

$LDFLAGS = "-s -w -X 'github.com/longgoll/forge-cli/cmd.Version=$Version' -X 'github.com/longgoll/forge-cli/cmd.BuildDate=$BUILD_DATE'"

Write-Host ""
Write-Host "  ⚡ Forge CLI Build Script" -ForegroundColor Cyan
Write-Host "  ─────────────────────────────────────────" -ForegroundColor DarkGray
Write-Host "  Version:    $Version" -ForegroundColor DarkGray
Write-Host "  Build Date: $BUILD_DATE" -ForegroundColor DarkGray
Write-Host ""

if ($All) {
    # ── Cross-compile for all platforms ──
    if (Test-Path $BUILD_DIR) { Remove-Item -Recurse -Force $BUILD_DIR }
    New-Item -ItemType Directory -Path $BUILD_DIR -Force | Out-Null

    $targets = @(
        @{ GOOS="windows"; GOARCH="amd64"; ext=".exe" },
        @{ GOOS="windows"; GOARCH="arm64"; ext=".exe" },
        @{ GOOS="linux";   GOARCH="amd64"; ext="" },
        @{ GOOS="linux";   GOARCH="arm64"; ext="" },
        @{ GOOS="darwin";  GOARCH="amd64"; ext="" },
        @{ GOOS="darwin";  GOARCH="arm64"; ext="" }
    )

    foreach ($target in $targets) {
        $output = "$BUILD_DIR/$APP_NAME-$($target.GOOS)-$($target.GOARCH)$($target.ext)"
        Write-Host "  → $($target.GOOS)/$($target.GOARCH)" -ForegroundColor Yellow -NoNewline

        $env:GOOS = $target.GOOS
        $env:GOARCH = $target.GOARCH
        go build -ldflags $LDFLAGS -o $output .

        if ($LASTEXITCODE -eq 0) {
            $size = [math]::Round((Get-Item $output).Length / 1MB, 1)
            Write-Host "  ✓ ${size}MB" -ForegroundColor Green
        } else {
            Write-Host "  ✗ FAILED" -ForegroundColor Red
        }
    }

    # Reset environment
    Remove-Item Env:GOOS -ErrorAction SilentlyContinue
    Remove-Item Env:GOARCH -ErrorAction SilentlyContinue

    Write-Host ""
    Write-Host "  ✓ All builds complete! Check $BUILD_DIR/" -ForegroundColor Green

    # Generate checksums
    Write-Host ""
    Write-Host "  SHA256 Checksums:" -ForegroundColor Cyan
    Get-ChildItem $BUILD_DIR -File | Where-Object { $_.Name -ne "checksums.txt" } | ForEach-Object {
        $hash = (Get-FileHash $_.FullName -Algorithm SHA256).Hash.ToLower()
        $line = "$hash  $($_.Name)"
        Write-Host "  $line" -ForegroundColor DarkGray
        $line
    } | Out-File -FilePath "$BUILD_DIR/checksums.txt" -Encoding utf8

} else {
    # ── Build for current OS only ──
    Write-Host "  → Building for current platform..." -ForegroundColor Yellow
    go build -ldflags $LDFLAGS -o "$APP_NAME.exe" .

    if ($LASTEXITCODE -eq 0) {
        $size = [math]::Round((Get-Item "$APP_NAME.exe").Length / 1MB, 1)
        Write-Host "  ✓ $APP_NAME.exe (${size}MB)" -ForegroundColor Green
    } else {
        Write-Host "  ✗ Build failed!" -ForegroundColor Red
        exit 1
    }
}

Write-Host ""
