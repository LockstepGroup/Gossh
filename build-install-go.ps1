<#
.Description
Download and install Golang on Windows. It sets the GOROOT environment
variable and adds GOROOT\bin to the PATH environment variable.
#>
Param(
    [Parameter(Mandatory = $false, Position = 0)]
    [string]$Version = '1.13.4'
)

$GoRoot = "C:\go$version"
$DownloadDir = $env:TEMP

try {
    $GoVersion = Invoke-Expression -Command "go version"
    if ($GoVersion -match $Version) {
        Write-Verbose "Desired Go version ($Version) already installed, exiting"
        exit
    }
} catch {
    Write-Verbose "Go not in `$env:PATH, checking manual path"
    try {
        $env:GOROOT = $GoRoot
        $env:PATH = $env:PATH + ';' + (Join-Path -Path $env:GOROOT -ChildPath 'bin')
        $GoVersion = Invoke-Expression -Command "go version"
        if ($GoVersion -match $Version) {
            Write-Verbose "Desired Go version ($Version) already installed, adding environment variables and continuing"

            exit
        }
    } catch {
        Write-Verbose "Go not installed, continuing"
    }
}
Write-Debug "waiting to continue"

$url32 = 'https://storage.googleapis.com/golang/go' + $version + '.windows-386.zip'
$url64 = 'https://storage.googleapis.com/golang/go' + $version + '.windows-amd64.zip'


# Determine type of system
if ($ENV:PROCESSOR_ARCHITECTURE -eq "AMD64") {
    $url = $url64
} else {
    $url = $url32
}

<# if (Test-Path "$GoRoot\bin\go.exe") {
    Write-Host "Go is installed to $GoRoot"
    exit
} #>

Write-Verbose "Downloading $url"
$zip = "$DownloadDir\golang-$version.zip"
if (!(Test-Path "$zip")) {
    $downloader = new-object System.Net.WebClient
    $downloader.DownloadFile($url, $zip)
}

Write-Verbose "Extracting $zip to $GoRoot"
if (Test-Path "$DownloadDir\go") {
    rm -Force -Recurse -Path "$DownloadDir\go"
}
Add-Type -AssemblyName System.IO.Compression.FileSystem
[System.IO.Compression.ZipFile]::ExtractToDirectory("$zip", $DownloadDir)
mv "$DownloadDir\go" $GoRoot

Write-Verbose "Setting GOROOT for session"
$env:GOROOT = $GoRoot
<#
[System.Environment]::SetEnvironmentVariable("GOROOT", "$goroot", "Machine")
$p = [System.Environment]::GetEnvironmentVariable("PATH", "Machine")
$p = "$goroot\bin;$p"
[System.Environment]::SetEnvironmentVariable("PATH", "$p", "Machine") #>