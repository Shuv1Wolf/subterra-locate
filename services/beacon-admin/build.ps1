#!/usr/bin/env pwsh

Set-StrictMode -Version latest
$ErrorActionPreference = "Stop"

# Generate image and container names using the data in the "component.json" file
$component = Get-Content -Path "component.json" | ConvertFrom-Json
$buildImage="$($component.registry)/$($component.name):$($component.version)-$($component.build)-build"
$container=$component.name

# Copy private keys to access git repo
if (-not (Test-Path -Path "$PSScriptRoot/docker/id_rsa")) {
    if (-not [string]::IsNullOrEmpty($env:GIT_PRIVATE_KEY)) {
        Write-Host "Creating docker/id_rsa from environment variable..."
        Set-Content -Path "$PSScriptRoot/docker/id_rsa" -Value $($env:GIT_PRIVATE_KEY).Replace("\n", "`n")
    }
    elseif (Test-Path -Path "~/.ssh/id_rsa") {
        Write-Host "Copying ~/.ssh/id_rsa to docker..."
        Copy-Item -Path "~/.ssh/id_rsa" -Destination "docker"
    }
    else {
        Write-Host "Missing ~/.ssh/id_rsa file..."
        Set-Content -Path "$PSScriptRoot/docker/id_rsa" -Value ""
    }
}

# Remove build files
if (Test-Path "./dist") {
    Remove-Item -Recurse -Force -Path "./dist/*"
} else {
    New-Item -ItemType Directory -Force -Path "./dist"
}

# Build docker image
docker build -f docker/Dockerfile.build -t $buildImage .

# Create and copy compiled files, then destroy the container
docker create --name $container $buildImage
docker cp "$($container):/app/main" "$PSScriptRoot/dist/main"
docker rm $container

if (-not(Test-Path "./dist")) {
    Write-Error "dist folder doesn't exist in root dir. Build failed. Watch logs above."
    exit 1
}
