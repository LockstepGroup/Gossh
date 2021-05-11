package main

import (
	"flag"
	"strings"
)

var DeviceAddress = ""
var DeviceUser = ""
var DevicePass = ""
var DeviceConfig = ""
var ConfigSplitCount = 0
var HostSplitCount = 0
var seconds int
var flag1 bool
var flag2 bool
var Config string
var SplitConfig string
var Port string
var Address string
var UserName string
var Password string
var Device string
var EnablePassword string

func main() {
	flag.StringVar(&Address, "host", "", "the hostname or ip address of the devices that you are attempting to connect to")
	flag.StringVar(&UserName, "user", "", "username used to connect to the device")
	flag.StringVar(&Password, "pass", "", "password used to connect to the device")
	flag.StringVar(&Config, "command", "", "commands that you would like to run on the host use a // to seperate diffrent commands")
	flag.StringVar(&Device, "device", "", "Current Device Types \"ExtremeExos\",\"ExtremeEos\",\"CiscoASA\",\"PaloAlto\",\"HP\",\"CiscoSwitch\",\"GetConsole\",\"SonicWall\"")
	flag.StringVar(&Port, "port", "", "specify the port that the ssh connection will be preformed on")
	flag.StringVar(&EnablePassword, "enable", "", "specify the enable password")
	flag.Parse()
	DeviceTypeTrim := strings.TrimSpace(Device)
	DeviceType := strings.ToLower(DeviceTypeTrim)
	if strings.Contains(DeviceType, "extremeeos") {
		ExtremeEos()
	}
	if strings.Contains(DeviceType, "extremeexos") {
		ExtremeExos()
	}
	if strings.Contains(DeviceType, "ciscoasa") {
		CiscoASA()
	}
	if strings.Contains(DeviceType, "paloalto") {
		PaloAlto()
	}
	if strings.Contains(DeviceType, "hp") {
		HpSwitch()
	}
	if strings.Contains(DeviceType, "ciscoswitch") {
		CiscoSwitch()
	}
	if strings.Contains(DeviceType, "getconsole") {
		GetConsole()
	}
	if DeviceType == "netscantest" {
		netscantest()
	}
	if DeviceType == "sonicwall" {
		SonicWall()
	}
	if DeviceType == "esxi" {
		ESXI()
	}
}
