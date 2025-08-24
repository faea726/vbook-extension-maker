# build.ps1
$ErrorActionPreference = "Stop"

$AppName = "vbook-cli"
$OutputDir = "build"

# List of target platforms (OS/ARCH)
$Platforms = @(
    "windows/amd64",
    "windows/386",
    "linux/amd64",
    "linux/386",
    "linux/arm64",
    "darwin/amd64",
    "darwin/arm64",
    "android/arm64"   # ✅ Android (aarch64)
)

# Clean & recreate build directory
if (Test-Path $OutputDir) {
    Remove-Item $OutputDir -Recurse -Force
}
New-Item -ItemType Directory -Path $OutputDir | Out-Null

foreach ($Platform in $Platforms) {
    $parts = $Platform.Split("/")
    $GOOS = $parts[0]
    $GOARCH = $parts[1]

    $OutputName = "$AppName-$GOOS-$GOARCH"
    if ($GOOS -eq "windows") {
        $OutputName += ".exe"
    }

    Write-Host "Building $OutputName ..."

    $env:GOOS = $GOOS
    $env:GOARCH = $GOARCH
    $env:CGO_ENABLED = "0"   # ✅ Disable CGO for cross-build portability

    go build -o "$OutputDir/$OutputName" .
}

Write-Host "[OK] Build finished. Binaries are in '$OutputDir/'"
