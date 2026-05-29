# Awesome Autostart Manager - Windows Release Build Script
# This script compiles the project with optimal settings for production.
# It hides the console window, optimizes size, and packages the binary.

$ErrorActionPreference = "Stop"

# Define colors and formatting
function Write-Header ($text) {
    Write-Host "`n=== $text ===" -ForegroundColor Cyan
}

function Write-Success ($text) {
    Write-Host "✔ $text" -ForegroundColor Green
}

function Write-Info ($text) {
    Write-Host "ℹ $text" -ForegroundColor Yellow
}

function Write-ErrorMsg ($text) {
    Write-Host "✖ Error: $text" -ForegroundColor Red
}

Write-Header "Awesome Autostart Manager Build Script"
Write-Host "Preparing to compile a clean, optimized release binary..."

# 1. Check for Go
try {
    $goVersion = go version
    Write-Success "Go is installed: $goVersion"
} catch {
    Write-ErrorMsg "Go is not installed or not in PATH."
    Write-Info "Please install Go from https://go.dev/dl/ before building."
    Exit 1
}

# 2. Check for GCC (Required for CGO in Fyne)
try {
    $gccVersion = gcc --version | Select-Object -First 1
    Write-Success "GCC C Compiler is installed: $gccVersion"
} catch {
    Write-ErrorMsg "GCC (C Compiler) was not found in your PATH."
    Write-Info "Fyne GUI apps require a C compiler (CGO_ENABLED=1) to compile."
    Write-Info "Please install MSYS2 (https://www.msys2.org/) or TDM-GCC (https://jmeubank.github.io/tdm-gcc/)"
    Write-Info "and ensure 'gcc' is available in your environment variables/PATH."
    Exit 1
}

# 3. Clean previous builds
Write-Header "Cleaning Up"
if (Test-Path "aam.exe") {
    Remove-Item "aam.exe"
    Write-Success "Removed old aam.exe"
}
if (Test-Path "awesome-autostart-app-manager.zip") {
    Remove-Item "awesome-autostart-app-manager.zip"
    Write-Success "Removed old awesome-autostart-app-manager.zip"
}

# 4. Build Option A: Standard Optimized GUI Binary
Write-Header "Build Option A: Standard Optimized GUI Binary"
Write-Host "Compiling with options to hide the background terminal and optimize file size..."

$buildFlags = '-ldflags="-H=windowsgui -s -w"'
Write-Info "Command: go build $buildFlags -o aam.exe"

$env:CGO_ENABLED = "1"
Invoke-Expression "go build $buildFlags -o aam.exe"

if (Test-Path "aam.exe") {
    $sizeBytes = (Get-Item "aam.exe").Length
    $sizeMB = [math]::Round($sizeBytes / 1MB, 2)
    Write-Success "Successfully built optimized aam.exe! ($sizeMB MB)"
    Write-Info "Note: This binary is fully standalone. Copy it anywhere to run."
} else {
    Write-ErrorMsg "Failed to compile standard binary."
    Exit 1
}

# 5. Build Option B: Package with Fyne Tool (Adds Icon and App Manifest)
Write-Header "Build Option B: Package App with Embedded Icon"
Write-Host "This packages the executable, injecting the custom icon ('saves/awesome_logo.png') "
Write-Host "and creating an app manifest so it displays beautifully in Windows Explorer & the Taskbar."
Write-Info "Checking for Fyne CLI..."

$useFynePackage = $true
try {
    # We run the fyne CLI package tool via go run so the user doesn't have to install it globally
    Write-Host "Running Fyne CLI package tool..."
    go run fyne.io/fyne/v2/cmd/fyne@v2.3.5 package -os windows -icon saves/awesome_logo.png
} catch {
    Write-ErrorMsg "Fyne package execution failed: $_"
    Write-Info "Standard aam.exe is still available."
    $useFynePackage = $false
}

if ($useFynePackage -and (Test-Path "awesome-autostart-app-manager.zip")) {
    Write-Success "Successfully packaged application!"
    Write-Success "Packaged Zip: awesome-autostart-app-manager.zip"
    Write-Info "Extract this ZIP to find the fully branded executable with the embedded icon."
} else {
    Write-Info "Fyne package step was skipped or did not generate a zip."
}

Write-Header "Done"
Write-Host "All build steps finished! You are ready to distribute your application." -ForegroundColor Green
