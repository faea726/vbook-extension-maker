# Build script for vbook-cli cross-platform distribution
# This script builds the CLI tool for multiple platforms

param(
    [string]$Version = "1.0.0",
    [string]$OutputDir = "dist"
)

Write-Host "Building vbook-cli v$Version for multiple platforms..." -ForegroundColor Green

# Create output directory
if (Test-Path $OutputDir) {
    Remove-Item -Recurse -Force $OutputDir
}
New-Item -ItemType Directory -Path $OutputDir | Out-Null

# Set common build variables
$env:CGO_ENABLED = "0"
$BuildFlags = @("-ldflags", "-s -w -X main.version=$Version")

# Define target platforms
$Platforms = @(
    @{ OS = "windows"; Arch = "amd64"; Ext = ".exe" },
    @{ OS = "windows"; Arch = "arm64"; Ext = ".exe" },
    @{ OS = "linux"; Arch = "amd64"; Ext = "" },
    @{ OS = "linux"; Arch = "arm64"; Ext = "" },
    @{ OS = "darwin"; Arch = "amd64"; Ext = "" },
    @{ OS = "darwin"; Arch = "arm64"; Ext = "" }
)

foreach ($Platform in $Platforms) {
    $OutputName = "vbook-cli-$($Platform.OS)-$($Platform.Arch)$($Platform.Ext)"
    $OutputPath = Join-Path $OutputDir $OutputName
    
    Write-Host "Building for $($Platform.OS)/$($Platform.Arch)..." -ForegroundColor Yellow
    
    $env:GOOS = $Platform.OS
    $env:GOARCH = $Platform.Arch
    
    $BuildArgs = @("build") + $BuildFlags + @("-o", $OutputPath, ".")
    
    try {
        & go @BuildArgs
        if ($LASTEXITCODE -eq 0) {
            $FileSize = (Get-Item $OutputPath).Length
            Write-Host "  Success: Built $OutputName ($([math]::Round($FileSize / 1024 / 1024, 2)) MB)" -ForegroundColor Green
        } else {
            Write-Host "  Error: Failed to build $OutputName" -ForegroundColor Red
        }
    } catch {
        Write-Host "  Error: Error building $OutputName - $_" -ForegroundColor Red
    }
}

Write-Host "`nBuild completed! Binaries are in the '$OutputDir' directory." -ForegroundColor Green

# List all built files
Write-Host "`nBuilt files:" -ForegroundColor Cyan
Get-ChildItem $OutputDir | ForEach-Object {
    $Size = [math]::Round($_.Length / 1024 / 1024, 2)
    Write-Host "  $($_.Name) ($Size MB)" -ForegroundColor White
}