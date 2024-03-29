package main

import (
	"flag"
)

var Banner = false
var bannerchecked = false
var returnstring = ""

func main() {
	var HostName string
	var UserName string
	var Password string
	var EnablePassword string
	var SSH bool
	var File bool
	var Port string
	var CLIArguments string
	var TimeoutValue int
	//setting CLI Parameters//
	flag.StringVar(&HostName, "h", "", "specify the host you are connecting to with either a hostname or ip address")
	flag.StringVar(&UserName, "u", "", "specify device username")
	flag.StringVar(&Password, "p", "", "specify device password")
	flag.StringVar(&EnablePassword, "e", "", "specify device enable password")
	flag.BoolVar(&SSH, "c", true, "specify connection type default is SSH")
	flag.StringVar(&Port, "P", "22", "use this to specify port default is 22")
	flag.StringVar(&CLIArguments, "C", "", "specify CLI commands to run on host device")
	flag.BoolVar(&File, "f", false, "set to true of you are sending a config file")
	flag.IntVar(&TimeoutValue, "t", 35, "change default timeout value")
	flag.Parse()
	switch File {
	case false:
		switch SSH {
		case true:
			InvokeSSH(HostName, Port, UserName, Password, StringToArray(CLIArguments), TimeoutValue)
		case false:
			//InvokeTelnet(HostName, UserName, Password, EnablePassword, Port, StringToArray1(CLIArguments))
		}
	case true:
		switch SSH {
		case true:
			InvokeSSH(HostName, Port, UserName, Password, readLines(CLIArguments), TimeoutValue)
		case false:
			//InvokeTelnet(HostName, UserName, Password, EnablePassword, Port, StringToArray1(CLIArguments))
		}
	}
}
