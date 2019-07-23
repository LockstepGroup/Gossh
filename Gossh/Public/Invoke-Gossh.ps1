function Invoke-Gossh {
    [cmdletbinding()]
    Param (
        [Parameter(Mandatory = $true, Position = 0)]
        [string]$Hostname,

        [Parameter(Mandatory = $True, Position = 1)]
        [System.Management.Automation.PSCredential]
        [System.Management.Automation.Credential()]
        $Credential,

        [Parameter(Mandatory = $true, Position = 2)]
        [ValidateSet("ExtremeExos", "ExtremeEos", "CiscoASA", "PaloAlto", "HP", "CiscoSwitch", "Lab")]
        [string]$DeviceType,

        [Parameter(Mandatory = $true, Position = 3)]
        [string[]]$Command,

        [Parameter(Mandatory = $false)]
        [ValidateRange(0, 65535)]
        [int]$Port = 22
    )

    $VerbosePrefix = "Invoke-Gossh:"

    #region getGosshPath
    #############################################################

    switch -Regex (Get-OsVersion) {
        'MacOS' {
            $GosshBinaryName = 'gossh'
        }
        default {
            $GosshBinaryName = 'gossh.exe'
        }
    }
    $GosshPath = Join-Path -Path (Split-Path -Path $PSScriptRoot) -ChildPath "Bin/$GosshBinaryName"
    Write-Verbose "$VerbosePrefix GosshPath: $GosshPath"

    #############################################################
    #endregion getGosshPath

    #region invokeGossh
    #############################################################
    # \gossh.exe -host 1.1.1.1 -user admin -pass password -device lab -port 4001 -config "terminal pager 0/show run interface"

    $GosshCommand = $Command -join '/'
    $GosshUsername = $Credential.UserName
    $GosshPassword = $Credential.GetNetworkCredential().Password

    $GosshExpression = '. "' + $GosshPath + '"'
    $GosshExpression += ' -host ' + $Hostname
    $GosshExpression += ' -user ' + $GosshUsername
    $GosshExpression += ' -pass ' + $GosshPassword
    $GosshExpression += ' -device ' + $DeviceType
    $GosshExpression += ' -port ' + $Port
    $GosshExpression += ' -config "' + $GosshCommand + '"'

    Write-Verbose $GosshExpression

    $Results = Invoke-Expression -Command $GosshExpression

    #############################################################
    #endregion invokeGossh

    $Results
}