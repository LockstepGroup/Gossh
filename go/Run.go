package main

import (
	"fmt"
	"github.com/TwiN/go-color"
	"github.com/jessevdk/go-flags"
	"os"
)

var  returnStrings string
var options Options
var verbose bool
var secondaryLogin bool

type Options struct {
	HostName       string   `short:"h" long:"hostname" description:"device url/ip"`
	UserName       string   `short:"u" long:"username" description:"username used for authentication"`
	PassWord       string   `short:"p" long:"password" description:"password used for authentication"`
	EnablePassword string   `short:"e" long:"enablepassword" description:"username used for authentication"`
	Port           int      `short:"P" long:"port" description:"tcp port device is listening on" default:"22"`
	TimeOut           int      `short:"t" long:"timeout" description:"timeout value in seconds for each command execution" default:"10"`
	CliArguments   string `short:"C" long:"cliarguments" description:"array containing cmdS to be run" `
	DirectoryPath       string   `short:"d" long:"filePath" description:"filePath to document containing cmdS to be run"`
	UseFile           bool     `short:"f" long:"usefile" description:"bool switch for using cmdS file"`
	SecondaryLogin bool `short:"s" long:"secondarylogin" description:"bool switch for entering login credentials again as required by some cisco switches"`
	Telnet            bool     `short:"T" long:"telnet" description:"option to use telnet instead of ssh"`
	Help           bool     `short:"H" long:"help" description:"output help"`
	Verbose bool `short:"v" long:"verbose" description:"print verbose"`
}

func main() {
	parser:= flags.NewParser(&options,flags.Default&^flags.HelpFlag)
	_, err := parser.Parse()
	if nil != err {handleError(err,true,"Error parsing CLI arguments")}
	if options.Help{parser.WriteHelp(os.Stdout);os.Exit(0)}
	//bytes := make([]byte, 32)
	//if _, err := rand.Read(bytes); err != nil {
	//	handleError(err, true, "failed to create encryption key")
	//}
	//key := hex.EncodeToString(bytes)
	verbose = options.Verbose
	secondaryLogin = options.SecondaryLogin
	switch {
	case options.UseFile:
		commandArguments := readConfigFile(options.DirectoryPath)
		switch {
		case options.Telnet:
			//
		default:
			sshDial(options.HostName,options.Port,options.UserName,options.PassWord,commandArguments,options.TimeOut,options.EnablePassword)
		}
	case options.Telnet:
		println(color.Ize(color.Red,"Telnet functions are not ready yet terminating application")); os.Exit(1)
	default:
		sshDial(options.HostName, options.Port, options.UserName, options.PassWord,stringToArray(options.CliArguments),options.TimeOut,options.EnablePassword)
	}
	fmt.Println(returnStrings)
}
