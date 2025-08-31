
$componentFile = "component.json"
if (!(Test-Path $componentFile)) {
    Write-Error "The component.json file was not found"
    exit 1
}

$json = Get-Content $componentFile | ConvertFrom-Json
$name = $json.name
$version = $json.version
$type = $json.type
$tag = "$type/$name/v$version"

$existingTag = git tag --list $tag
if ($existingTag) {
    Write-Host "Tag $tag already exists."
    exit 1
}

Write-Host "Create a tag $tag for the component $name"
git tag -a $tag -m "Release $name $tag"
git push origin $tag

Write-Host "Tag $tag created and submitted to GitHub."
