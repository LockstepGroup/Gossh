build:
	go build -pkgdir gossh/ -a -o '/Users/nmuragian/OneDrive - Lockstep Technology Group/Compiled/GoSsh'
compile:
	echo "Compiling for Windows and MacOS"
	GOOS=darwin GOARCH=amd64 go build -pkgdir gossh/ -a -o $$HOME'/OneDrive - Lockstep Technology Group/Compiled/GoSsh'
	GOOS=windows GOARCH=amd64 go build -pkgdir gossh/ -a -o $$HOME'/OneDrive - Lockstep Technology Group/Compiled/GoSsh'
PSSwitch:
	echo "Compiling for PSSwitch"
	rm $$HOME'/OneDrive - Lockstep Technology Group/code/git/Powershell/PSSwitch/private/bin/gossh.exe'
	rm $$HOME'/OneDrive - Lockstep Technology Group/code/git/Powershell/PSSwitch/private/bin/gossh'
	ls $$HOME'/OneDrive - Lockstep Technology Group/code/git/Powershell/PSSwitch/private/bin/'
	GOOS=darwin GOARCH=amd64 go build -pkgdir gossh/ -a -o $$HOME'/OneDrive - Lockstep Technology Group/code/git/Powershell/PSSwitch/private/bin'
	GOOS=windows GOARCH=amd64 go build -pkgdir gossh/ -a -o $$HOME'/OneDrive - Lockstep Technology Group/code/git/Powershell/PSSwitch/private/bin'
	GOOS=darwin GOARCH=amd64 go build -pkgdir gossh/ -a -o $$HOME'/OneDrive - Lockstep Technology Group/Compiled/GoSsh'
	GOOS=windows GOARCH=amd64 go build -pkgdir gossh/ -a -o $$HOME'/OneDrive - Lockstep Technology Group/Compiled/GoSsh'
	ls $$HOME'/OneDrive - Lockstep Technology Group/code/git/Powershell/PSSwitch/private/bin/'