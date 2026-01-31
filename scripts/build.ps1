# PowerShell build script for Windows
# Cross-platform compilation for InstaCli

$VERSION = "1.0.0"
$BUILD_DIR = ".\dist"
$APP_NAME = "instacli"

Write-Host "🔨 Building InstaCli v$VERSION..." -ForegroundColor Cyan
Write-Host ""

# Clean and create build directory
if (Test-Path $BUILD_DIR) {
    Remove-Item -Recurse -Force $BUILD_DIR
}
New-Item -ItemType Directory -Path $BUILD_DIR | Out-Null

# Platform configurations
$platforms = @(
    @{GOOS="linux"; GOARCH="amd64"; ext=""},
    @{GOOS="linux"; GOARCH="arm64"; ext=""},
    @{GOOS="darwin"; GOARCH="amd64"; ext=""},
    @{GOOS="darwin"; GOARCH="arm64"; ext=""},
    @{GOOS="windows"; GOARCH="amd64"; ext=".exe"}
)

foreach ($platform in $platforms) {
    $env:GOOS = $platform.GOOS
    $env:GOARCH = $platform.GOARCH
    
    $outputName = "$APP_NAME-$($platform.GOOS)-$($platform.GOARCH)$($platform.ext)"
    $outputPath = Join-Path $BUILD_DIR $outputName
    
    Write-Host "📦 Building for $($platform.GOOS)/$($platform.GOARCH)..." -ForegroundColor Yellow
    
    go build -ldflags="-s -w" -o $outputPath ./cmd/instacli
    
    if ($LASTEXITCODE -eq 0) {
        Write-Host "   ✓ $outputName" -ForegroundColor Green
    } else {
        Write-Host "   ✗ Failed!" -ForegroundColor Red
    }
}

# Reset environment
$env:GOOS = ""
$env:GOARCH = ""

Write-Host ""
Write-Host "✅ Build complete! Binaries are in $BUILD_DIR/" -ForegroundColor Green
Write-Host ""
Get-ChildItem $BUILD_DIR | Format-Table Name, Length
