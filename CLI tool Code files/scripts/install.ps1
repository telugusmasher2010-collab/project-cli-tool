#requires -Version 5.1
param(
    [string]$Version = "latest"
)

$ErrorActionPreference = "Stop"

$Repo = "telugusmasher2010-collab/project-cli-tool"
$InstallDir = Join-Path $env:LOCALAPPDATA "proj-init"
$BinPath = Join-Path $InstallDir "proj-init.exe"

function Write-Info([string]$Msg) { Write-Host $Msg -ForegroundColor Green }
function Write-Err([string]$Msg)  { Write-Host $Msg -ForegroundColor Red; exit 1 }

# --- detect OS / arch -------------------------------------------------------
$Arch = $env:PROCESSOR_ARCHITECTURE
switch ($Arch) {
    "AMD64" { $Arch = "amd64" }
    "ARM64" { $Arch = "arm64" }
    default { Write-Err "Unsupported architecture: $Arch" }
}

# --- resolve release asset --------------------------------------------------
if ($Version -eq "latest") {
    $Release = Invoke-RestMethod -Uri "https://api.github.com/repos/$Repo/releases/latest"
    $Version = $Release.tag_name
}

$File = "proj-init_$($Version.TrimStart('v'))_windows_$Arch.zip"
$Url = "https://github.com/$Repo/releases/download/$Version/$File"

$Tmp = Join-Path $env:TEMP ("proj-init-" + [Guid]::NewGuid().ToString("N"))
New-Item -ItemType Directory -Path $Tmp -Force | Out-Null

try {
    Write-Info "Downloading $Url"
    Invoke-WebRequest -Uri $Url -OutFile (Join-Path $Tmp "proj-init.zip")

    Expand-Archive -Path (Join-Path $Tmp "proj-init.zip") -DestinationPath $Tmp -Force

    $NewBin = Get-ChildItem -Path $Tmp -Filter "proj-init.exe" -Recurse | Select-Object -First 1
    if (-not $NewBin) { Write-Err "proj-init.exe not found in archive" }

    New-Item -ItemType Directory -Path $InstallDir -Force | Out-Null
    Copy-Item -Path $NewBin.FullName -Destination $BinPath -Force

    Write-Info "proj-init $Version installed to $BinPath"
} finally {
    Remove-Item -Path $Tmp -Recurse -Force -ErrorAction SilentlyContinue
}

# --- add to user PATH if missing --------------------------------------------
$UserPath = [Environment]::GetEnvironmentVariable("Path", "User")
if ($UserPath -notlike "*$InstallDir*") {
    $NewPath = if ([string]::IsNullOrEmpty($UserPath)) { $InstallDir } else { "$UserPath;$InstallDir" }
    [Environment]::SetEnvironmentVariable("Path", $NewPath, "User")
    Write-Info "Added $InstallDir to your user PATH. Restart your terminal to use 'proj-init'."
} else {
    Write-Info "Run 'proj-init --help' to get started."
}
