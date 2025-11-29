#!/usr/bin/env pwsh

Set-StrictMode -Version latest
$ErrorActionPreference = "Stop"

# Get component metadata and set necessary variables
$component = Get-Content -Path "$PSScriptRoot/component.json" | ConvertFrom-Json
$protosImage = "$($component.registry)/$($component.name):$($component.version)-$($component.build)-protos"
$container=$component.name

# Remove documentation files
if (Test-Path -Path "$PSScriptRoot/protos") {
    Remove-Item -Recurse -Force -Path "$PSScriptRoot/protos/*.go"
}

# Build docker image
docker build -f "$PSScriptRoot/docker/Dockerfile.proto" -t $protosImage .

# Create and copy compiled files, then destroy
docker create --name $container $protosImage
docker cp "$($container):/app/protos" "$PSScriptRoot/"
docker rm $container

# Verify that protos folder was indeed created after generating proto files
if (-not (Test-Path "$PSScriptRoot/protos")) {
    Write-Error "protos folder doesn't exist in src dir. Build failed. See logs above for more information."
}
