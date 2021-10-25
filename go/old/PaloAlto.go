package main

import (
	"fmt"
	"strings"

	pan "github.com/zepryspet/GoPAN/utils"
	"golang.org/x/crypto/ssh"
)

func PaloAlto() {
	DeviceTypeTrim := strings.TrimSpace(Device)
	DeviceType := strings.ToLower(DeviceTypeTrim)
	DeviceAddress = strings.TrimSpace(Address)
	//fmt.Println(DeviceAddress)
	DeviceUser = strings.TrimSpace(UserName)
	//fmt.Println(DeviceUser)
	DevicePass = strings.TrimSpace(Password)
	//fmt.Println(DevicePass)
	DeviceConfigSplit := strings.Split(Config, "//")
	fmt.Println(DeviceType)
	for range DeviceConfigSplit {
		SplitConfig = DeviceConfigSplit[ConfigSplitCount]
		PanSSHConfig()
		ConfigSplitCount++
		//time.Sleep(5 * time.Millisecond)
	}
}

func PanSSHConfig() {
	deviceport := strings.TrimSpace(Port)

	// Start up ssh process
	//try username pass/first
	sshClt, err := ssh.Dial("tcp", DeviceAddress+":"+deviceport, &ssh.ClientConfig{
		User: DeviceUser,
		Auth: []ssh.AuthMethod{
			ssh.Password(DevicePass),
			//ssh.KeyboardInteractive(Challenge),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	})

	//if it fails try interactive keyboard auth type (show the banner and auto ack'd)
	if err != nil {
		sshClt, err = ssh.Dial("tcp", DeviceAddress+":"+deviceport, &ssh.ClientConfig{
			User: DeviceUser,
			Auth: []ssh.AuthMethod{
				//ssh.Password(pass),
				ssh.KeyboardInteractive(Challenge(DevicePass)),
			},
			HostKeyCallback: ssh.InsecureIgnoreHostKey(),
			BannerCallback:  ssh.BannerDisplayStderr(),
		})
		pan.Logerror(err, true)
	}

	session, err := sshClt.NewSession()
	pan.Logerror(err, true)
	sshOut, err := session.StdoutPipe()
	pan.Logerror(err, true)
	sshIn, err := session.StdinPipe()
	pan.Logerror(err, true)

	// Set up terminal modes
	modes := ssh.TerminalModes{
		ssh.ECHO:          0,     // disable echoing
		ssh.TTY_OP_ISPEED: 14400, // input speed = 14.4kbaud
		ssh.TTY_OP_OSPEED: 14400, // output speed = 14.4kbaud
	}
	// Request pseudo terminal
	if err := session.RequestPty("xterm", 40, 80, modes); err != nil {
		pan.Wlog("request for pseudo terminal failed: ", "error.txt", false)
		pan.Logerror(err, true)
	}
	// Start remote shell
	if err := session.Shell(); err != nil {
		pan.Wlog("request for remote shell failed: ", "error.txt", false)
		pan.Logerror(err, true)
	}
	//wait for banner
	// Once a Session is created, you can execute a single command on
	// the remote side using the Run method.

	cmdSend(sshOut, sshIn, SplitConfig, false, false, 20)
	session.Close()
}
