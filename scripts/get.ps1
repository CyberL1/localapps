#!/usr/bin/env pwsh

$ErrorActionPreference = 'Stop'

$Dir = "${Home}\localapps\bin"
$Zip = "$Dir\localapps.zip"
$Exe = "$Dir\localapps.exe"
$OldExe = "$env:Temp\localapps-old.exe"

$Target = "windows-amd64"

$DownloadUrl = "https://github.com/CyberL1/localapps/releases/latest/download/localapps-${Target}.zip"

if (!(Test-Path $Dir)) {
  New-Item $Dir -ItemType Directory | Out-Null
}

curl.exe -Lo $Zip $DownloadUrl

if (Test-Path $Exe) {
  Move-Item -Path $Exe -Destination $OldExe -Force
}

Expand-Archive -LiteralPath $Zip -DestinationPath $BinPath -Force
Remove-Item $Zip

$User = [System.EnvironmentVariableTarget]::User
$Path = [System.Environment]::GetEnvironmentVariable('Path', $User)

if (!(";${Path};".ToLower() -like "*;${Dir};*".ToLower())) {
  [System.Environment]::SetEnvironmentVariable('Path', "${Path};${Dir}", $User)
  $Env:Path += ";${Dir}"
}

Write-Output "localapps CLI was installed to $Exe"
Write-Output "Run 'localapps up' to get started"
