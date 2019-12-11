#requires -module Corkscrew
[Cmdletbinding()]
Param (
)

switch -Regex (Get-OsVersion) {
    'Windows' {
        $PathSeparator = ';'
    }
    'MacOSX' {
        if (-not (Test-Path -Path $env:GOROOT)) {
            $env:GOROOT = "/usr/local/opt/go/libexec"
        }
        $PathSeparator = ':'
    }
}

$env:GOPATH = Join-Path -Path (Resolve-Path -Path ./) -ChildPath "go"
$env:PATH = $env:PATH + $PathSeparator + (Join-Path -Path $env:GOPATH -ChildPath 'bin')
$env:PATH = $env:PATH + $PathSeparator + (Join-Path -Path $env:GOROOT -ChildPath 'bin')
$OutputPath = Join-Path -Path (Resolve-Path -Path ./) -ChildPath "Gossh/Bin"

Push-Location -Path $env:GOPATH

Write-Verbose "Getting go dependencies"
# load go dependencies
go get "github.com/zepryspet/GoPAN/utils"
go get "golang.org/x/crypto/ssh"
go get "github.com/sfreiberg/simplessh"

Write-Debug "dependencies loaded, continuing to build"

Write-Verbose "Building for macos"
$env:GOARCH = 'amd64'
$env:GOOS = 'darwin'

go build -o "$OutputPath/gossh"

Write-Verbose "Building for windows"
$env:GOARCH = 'amd64'
$env:GOOS = 'windows'

go build -o "$OutputPath/gossh.exe"

Pop-Location