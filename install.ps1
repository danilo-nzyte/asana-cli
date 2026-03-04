$ErrorActionPreference = "Stop"

$Repo = "danilo-nzyte/asana-cli"
$InstallDir = "$HOME\.local\bin"
$SkillDir = "$HOME\.claude\skills\asana"

# Detect architecture
$Arch = if ([Environment]::Is64BitOperatingSystem) {
    if ($env:PROCESSOR_ARCHITECTURE -eq "ARM64" -or $env:PROCESSOR_IDENTIFIER -match "ARM") {
        "arm64"
    } else {
        "amd64"
    }
} else {
    Write-Error "Unsupported: 32-bit systems are not supported."
    exit 1
}

$Archive = "asana-cli_windows_${Arch}.zip"

# Get latest release tag
Write-Host "==> Fetching latest release..."
$Release = Invoke-RestMethod -Uri "https://api.github.com/repos/$Repo/releases/latest"
$Tag = $Release.tag_name
Write-Host "    Latest release: $Tag"

# Download
$Url = "https://github.com/$Repo/releases/download/$Tag/$Archive"
$TmpDir = New-Item -ItemType Directory -Path (Join-Path $env:TEMP "asana-cli-install-$(Get-Random)")

Write-Host "==> Downloading $Archive..."
Invoke-WebRequest -Uri $Url -OutFile (Join-Path $TmpDir $Archive)

Write-Host "==> Extracting..."
Expand-Archive -Path (Join-Path $TmpDir $Archive) -DestinationPath $TmpDir -Force

# Install binary
Write-Host "==> Installing to $InstallDir..."
New-Item -ItemType Directory -Path $InstallDir -Force | Out-Null
Copy-Item (Join-Path $TmpDir "asana-cli.exe") (Join-Path $InstallDir "asana-cli.exe") -Force

# Add to PATH if not already there
$UserPath = [Environment]::GetEnvironmentVariable("Path", "User")
if ($UserPath -notlike "*$InstallDir*") {
    Write-Host "==> Adding $InstallDir to user PATH..."
    [Environment]::SetEnvironmentVariable("Path", "$UserPath;$InstallDir", "User")
    $env:Path = "$env:Path;$InstallDir"
}

# Install Claude Code skill
Write-Host "==> Installing Claude Code skill..."
New-Item -ItemType Directory -Path $SkillDir -Force | Out-Null
$SkillUrl = "https://raw.githubusercontent.com/$Repo/$Tag/skill/SKILL.md"
Invoke-WebRequest -Uri $SkillUrl -OutFile (Join-Path $SkillDir "SKILL.md")
Write-Host "    Skill installed to $SkillDir"

# Cleanup
Remove-Item -Recurse -Force $TmpDir

Write-Host ""
Write-Host "==> Installed asana-cli $Tag to $InstallDir\asana-cli.exe"
Write-Host ""
Write-Host "Next steps:"
Write-Host "  1. Restart your terminal (for PATH changes)"
Write-Host "  2. Set environment variables (ASANA_CLIENT_ID, ASANA_CLIENT_SECRET, ASANA_WORKSPACE_ID)"
Write-Host "  3. Run: asana-cli auth login"
Write-Host "  4. Verify: asana-cli auth status"
Write-Host ""
Write-Host "See https://github.com/$Repo#authentication-setup for details."
