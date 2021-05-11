package main

import (
	"fmt"
	"strings"
	"time"

	cssh "golang.org/x/crypto/ssh"
)

func HpSwitch() {
	DeviceTypeTrim := strings.TrimSpace(Device)
	DeviceType := strings.ToLower(DeviceTypeTrim)
	DeviceAddress = strings.TrimSpace(Address)
	//fmt.Println(DeviceAddress)
	DeviceUser = strings.TrimSpace(UserName)
	//fmt.Println(DeviceUser)
	DevicePass = strings.TrimSpace(Password)
	//fmt.Println(DevicePass)
	DeviceConfigSplit := strings.Split(Config, "//")
	Deviceport := strings.TrimSpace(Port)
	fmt.Println(DeviceType)
	conn, err := cssh.Dial("tcp", DeviceAddress+":"+Deviceport, &cssh.ClientConfig{
		User:            DeviceUser,
		Auth:            []cssh.AuthMethod{cssh.Password(DevicePass)},
		HostKeyCallback: cssh.InsecureIgnoreHostKey(),
		Timeout:         5 * time.Second,
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
	write("\r\n", sshIn)
	write("no page", sshIn)
	//write(DeviceUser, sshIn)
	//write(DevicePass, sshIn)
	//write("terminal pager 0", sshIn)
	for range DeviceConfigSplit {
		fmt.Println(DeviceConfigSplit[ConfigSplitCount])
		// create a new connection
		write(DeviceConfigSplit[ConfigSplitCount], sshIn)
		ConfigSplitCount++
	}
	write("logout", sshIn)
	write("y", sshIn)
	test := strings.Replace(readBuffForString2(sshOut), "[24;1H[2K[24;1H[1;24r[24;1H", "", -1)
	fmt.Println(test)
	session.Close()
	conn.Close()
}
