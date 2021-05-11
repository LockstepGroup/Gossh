package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/sfreiberg/simplessh"
	cssh "golang.org/x/crypto/ssh"
)

func ExtremeEos() {
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
	for range DeviceConfigSplit {
		fmt.Println(DeviceConfigSplit[ConfigSplitCount])
		// create a new connection
		write(DeviceConfigSplit[ConfigSplitCount], sshIn)
		ConfigSplitCount++
	}
	write("exit", sshIn)
	readBuffForString1(sshOut)
	session.Close()
	conn.Close()
}
func ExtremeExos() {
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
	fmt.Println(DeviceAddress + ":" + Deviceport)
	client, err := simplessh.ConnectWithPassword(DeviceAddress+":"+Deviceport, DeviceUser, DevicePass)
	if err != nil {
		panic(err)
	}
	if strings.Contains(Config, "//") {
		for range DeviceConfigSplit {
			fmt.Println(DeviceConfigSplit[ConfigSplitCount])
			output, err := client.Exec(DeviceConfigSplit[ConfigSplitCount])
			if err != nil {
				panic(err)
			}
			fmt.Printf("%s\n", output)
			ConfigSplitCount++
		}
	} else {
		output, err := client.Exec(Config)
		if err != nil {
			panic(err)
		}
		fmt.Printf("%s\n", output)
	}
	client.Close()
}
