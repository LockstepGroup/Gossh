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
        [ValidateSet("ExtremeExos", "ExtremeEos", "CiscoASA", "PaloAlto", "HpSwitch", "HpRouter", "HpAruba", "CiscoSwitch", "GetConsole", "SonicWall", "CiscoSB")]
        [string]$DeviceType,

        [Parameter(Mandatory = $true, Position = 3)]
        [AllowEmptyString()]
        [string[]]$Command,

        [Parameter(Mandatory = $false)]
        [ValidateRange(0, 65535)]
        [int]$Port = 22,

        [Parameter(Mandatory = $false)]
        [System.Management.Automation.PSCredential]
        [System.Management.Automation.Credential()]
        $EnableCredential
    )

    $VerbosePrefix = "Invoke-Gossh:"

    #region getGosshPath
    #############################################################

    switch -Regex (Get-OsVersion) {
        'MacOS' {
            $GosshBinaryName = 'gossh'
            $GosshPath = Join-Path -Path (Split-Path -Path $PSScriptRoot) -ChildPath "Bin/$GosshBinaryName"
            # check to see if executable
            $NixPath = $GosshPath -replace '([\ \(\)])', '\$1'
            $ExecutableCheckCommand = 'bash -c "if [ -x ' + $NixPath + ' ]; then echo true; else echo false; fi"'
            $IsExecutable = Invoke-Expression -Command $ExecutableCheckCommand
            if ($IsExecutable -notmatch 'true') {
                $ExecutableCheckCommand = 'bash -c "chmod +x ' + $NixPath + '"'
                $MakeExecutable = Invoke-Expression -Command $ExecutableCheckCommand
            }
        }
        default {
            $GosshBinaryName = 'gossh.exe'
            $GosshPath = Join-Path -Path (Split-Path -Path $PSScriptRoot) -ChildPath "Bin/$GosshBinaryName"
        }
    }

    Write-Verbose "$VerbosePrefix GosshPath: $GosshPath"

    #############################################################
    #endregion getGosshPath

    #region invokeGossh
    #############################################################
    # \gossh.exe -host 1.1.1.1 -user admin -pass password -device lab -port 4001 -config "terminal pager 0/show run interface"

    $GosshCommand = $Command -join '//'
    $GosshUsername = $Credential.UserName
    $GosshPassword = $Credential.GetNetworkCredential().Password

    $GosshExpression = '. "' + $GosshPath + '"'
    $GosshExpression += ' -host ' + $Hostname
    $GosshExpression += ' -user ' + $GosshUsername
    <#     if ('' -eq $GosshPassword) {
        Write-Verbose "$Password is null"
        $GosshExpression += ' -pass ' + "'" + '\r\n' + "'"
    } else {
        $GosshExpression += ' -pass ' + "'" + $GosshPassword + "'"
    } #>
    $GosshExpression += ' -pass ' + "'" + $GosshPassword + "'"
    $GosshExpression += ' -device ' + $DeviceType
    $GosshExpression += ' -port ' + $Port
    $GosshExpression += ' -command "' + $GosshCommand + '"'

    if ($EnableCredential) {
        #$EnableCredential = New-Object System.Management.Automation.PSCredential ('test', $EnablePassword)
        $EnablePassword = $EnableCredential.GetNetworkCredential().Password
        $GosshExpression += " -enable '" + $EnablePassword + "'"
    }

    # required to make error variable work
    $GosshExpression += " 2>''"

    Write-Verbose $GosshExpression

    $Results = Invoke-Expression -Command $GosshExpression -ErrorVariable GosshError

    if ($GosshError) {
        Throw $GosshError
    }

    #############################################################
    #endregion invokeGossh

    #region cleanup
    #############################################################

    Remove-Variable -Name 'GosshCommand', 'GosshUsername', 'GosshPassword'

    #############################################################
    #endregion cleanup

    $Results
}