function Invoke-Gossh {
    [cmdletbinding()]
    Param (
        [Parameter(Mandatory = $true, Position = 0)]
        [string]$Hostname,

        [Parameter(Mandatory = $True, Position = 1)]
        [System.Management.Automation.PSCredential]
        [System.Management.Automation.Credential()]
        $Credential,

        [Parameter(Mandatory = $true, Position = 3)]
        [AllowEmptyString()]
        [string[]]$Command,

        [Parameter(Mandatory = $false)]
        [ValidateRange(0, 65535)]
        [int]$Port = 22,

        [Parameter(Mandatory = $false)]
        [System.Management.Automation.PSCredential]
        [System.Management.Automation.Credential()]
        $EnableCredential,

        [Parameter(Mandatory = $false)]
        [bool]$ConfigFile = $false
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
    # \gossh.exe -h 1.1.1.1 -u admin -p password -P 4001 -C "terminal pager 0/show run interface" -f=false -t 35

    $GosshCommand = $Command -join '||'
    $GosshUsername = $Credential.UserName
    $GosshPassword = $Credential.GetNetworkCredential().Password

    $GosshExpression = '. "' + $GosshPath + '"'
    $GosshExpression += ' -h ' + $Hostname
    $GosshExpression += ' -u ' + $GosshUsername
    $GosshExpression += ' -p ' + "'" + $GosshPassword + "'"
    $GosshExpression += ' -P ' + $Port
    $GosshExpression += ' -C "' + $GosshCommand + '"'
    $GosshExpression += ' -f=' + $ConfigFile

    if ($EnableCredential) {
        #$EnableCredential = New-Object System.Management.Automation.PSCredential ('test', $EnablePassword)
        $EnablePassword = $EnableCredential.GetNetworkCredential().Password
        $GosshExpression += " -e '" + $EnablePassword + "'"
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