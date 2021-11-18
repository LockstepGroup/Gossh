package main

import (
	"fmt"
	"time"

	"golang.org/x/crypto/ssh"
)
func InvokeSSH(HostName string, Port string, UserName string, Password string, Commands []string, TimeoutValue int, EnablePassword string) {
	if Password == "blank" {
		Password = ""
	}
	conn, err := ssh.Dial("tcp", HostName+":"+Port, &ssh.ClientConfig{
		User:    UserName,
		Timeout: 5 * time.Second,
		Auth: []ssh.AuthMethod{ssh.Password(Password)},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	})
	if err != nil {
		fmt.Println(err)
	}
	session, err := conn.NewSession()
	handleError(err)
	sshOut, err := session.StdoutPipe()
	handleError(err)
	sshIn, err := session.StdinPipe()

	err = session.Shell()
	handleError(err)
	defer session.Close()
	defer conn.Close()
	if EnablePassword != ""{
		writer(sshIn,"en")
		writer(sshIn,EnablePassword)
	}
	for _, CLICommands := range Commands {
			write(CLICommands, sshIn)
			readBuff(sshOut, TimeoutValue,HostName)
			writeblank(sshIn)
			readBuffblank(sshOut, TimeoutValue,HostName)
	}
	fmt.Println(returnstring)
}
