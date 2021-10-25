package main

import (
	"fmt"
	"log"
	"net"
	"strings"
	"time"

	cssh "golang.org/x/crypto/ssh"
)

type NetScan struct {
	IPaddress string
}

func CiscoASA() {
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
	write("enable", sshIn)
	//write(DeviceUser, sshIn)
	if EnablePassword != "" {
		enablepassword := strings.TrimSpace(EnablePassword)
		write(enablepassword, sshIn)
	} else {
		write(DevicePass, sshIn)
	}
	write("terminal pager 0", sshIn)
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
func CiscoSwitch() {
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

func netscantest() {
	DeviceUser = strings.TrimSpace(UserName)
	//fmt.Println(DeviceUser)
	DevicePass = strings.TrimSpace(Password)
	//fmt.Println(DevicePass)
	DeviceConfigSplit := strings.Split(Config, "//")
	Deviceport := strings.TrimSpace(Port)
	DeviceAddress = strings.TrimSpace(Address)
	ip, ipnet, err := net.ParseCIDR(DeviceAddress)
	if err != nil {
		log.Fatal(err)
	}
	for ip := ip.Mask(ipnet.Mask); ipnet.Contains(ip); inc(ip) {
		DeviceIP := fmt.Sprint(ip)
		conn, err := cssh.Dial("tcp", DeviceIP+":"+Deviceport, &cssh.ClientConfig{
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
}
