package main

import (
	"github.com/TwiN/go-color"
	"golang.org/x/crypto/ssh"
	"io"
	"os"
	"regexp"
	"strconv"
	"time"
)
var cmdEnd = regexp.MustCompile(".*@?.*([#>])\\s?$")
var cmdUserName = regexp.MustCompile(".*(User\\s+Name|username|user\\sname).*$?")

func sshWriter(cmd string, SshIN io.WriteCloser) {
	_, err := SshIN.Write([]byte(cmd + "\r"))
	handleError(err,false,"Error writing to ssh buffer")
}

func sshWriteNewLine(sshIn io.WriteCloser){
	_, err := sshIn.Write([]byte(" \r\n"))
	handleError(err,false,"Error writing new line")
}

func sshReadLogin (sshIn io.WriteCloser,sshOut io.Reader, userName string, password string){
	buf := make([]byte, 2048)
	bufferResults, err := sshOut.Read(buf);handleError(err,false,"Error encountered reading ssh buffer")
	currentLine := string(buf[:bufferResults])
	ReaderLoop:
		for nil == err && cmdEnd.FindStringSubmatch(currentLine) == nil {
			switch {
			case cmdUserName.FindStringSubmatch(currentLine) != nil:
				sshWriter(userName,sshIn);sshWriter(password,sshIn)
				n, err := sshOut.Read(buf);if nil != err {handleError(err,true,"error reading login promt");break ReaderLoop}
				currentLine += string(buf[:n])
				break ReaderLoop
			default:
				sshWriteNewLine(sshIn);sshWriteNewLine(sshIn)
				n, err := sshOut.Read(buf);if nil != err {handleError(err,true,"error reading login promt");break ReaderLoop}
				currentLine += string(buf[:n])
			}
		}
	returnStrings += currentLine
}

func sshReadBuffer(sshOut io.Reader, timeOutValue int, cmd string, hostName string) {
	buf := make([]byte, 2048)
	StartTime := time.Now()
	bufferResults, err := sshOut.Read(buf);handleError(err,false,"Error encountered reading ssh buffer")
	currentLine := string(buf[:bufferResults])
	if verbose{println(color.Ize(color.Bold + color.Green,currentLine))}
	TimeOutValue := time.Duration(timeOutValue) * time.Second
ReaderLoop:
	for nil == err && cmdEnd.FindStringSubmatch(currentLine) == nil {
		n, err := sshOut.Read(buf)
		if verbose{println(color.Ize(color.Bold + color.Green,string(buf[:n])))}
		switch {
		case nil != err:
			handleError(err, false,"Error encountered reading ssh buffer")
			break ReaderLoop
		case time.Since(StartTime) > TimeOutValue:
			directoryPath, err := os.UserHomeDir()
			if err != nil {
				handleError(err, false,"Error encountered while trying to determine home directory")
				break ReaderLoop
			}
			println(color.Ize(color.Red, "Timeout Value exceeded please check error log for more details "+
				"canceling read operation for cmd:= "+cmd+" | "+directoryPath+"/"+hostName+"-error.txt"))
			writeBuffer2Txt(directoryPath+"/"+hostName+"error.txt", returnStrings, true)
			break ReaderLoop
		default:
			currentLine += string(buf[:n])
		}
	}
	//println(color.Ize(color.Cyan, currentLine))   //uncomment while troubleshoot to see where the buffer is getting stuck
	returnStrings += currentLine
}

func sshDial(hostname string, port int, username string, password string, cmdS []string, timeOutValue int, enablePassword string) {
	var portNumber = strconv.Itoa(port)
	conn, err := ssh.Dial("tcp", hostname+":"+portNumber, &ssh.ClientConfig{
		User:            username,
		Timeout:         5 * time.Second,
		Auth:            []ssh.AuthMethod{ssh.Password(password)},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	})
	handleError(err, true, "SSH dial failed check device information and network")
	session, err := conn.NewSession()
	handleError(err, true, "SSH dial failed check device information and network")
	sshOut, err := session.StdoutPipe()
	handleError(err, true, "SSH out pipe creation failed terminating session")
	sshIn, err := session.StdinPipe();handleError(err, true, "SSH in pipe creation failed terminating session")
	modes := ssh.TerminalModes{
		ssh.ECHO:          0,     // disable echoing
		ssh.TTY_OP_ISPEED: 14400, // input speed = 14.4kbaud
		ssh.TTY_OP_OSPEED: 14400, // output speed = 14.4kbaud
	}
	// Request pseudo terminal
	if err := session.RequestPty("xterm", 40, 80, modes);nil != err {handleError(err,false,"")}
	err = session.Shell();handleError(err, true, "SSH in pipe creation failed terminating session")
	defer func(conn *ssh.Client) {
		err := conn.Close()
		handleError(err, false, "Connection failed to close")
	}(conn)
	defer func(session *ssh.Session) {
		err := session.Close()
		handleError(err, false, "Session failed to close")
	}(session)
	switch { //created as a switch to allow the addition of buffer reading for click through prompts
	case len([]rune(enablePassword)) >= 1:
		sshWriter("en", sshIn)
		sshWriter(enablePassword, sshIn)
	case secondaryLogin:
		sshReadLogin(sshIn,sshOut,username,password)
	}
	for _, CommandString := range cmdS {
		switch {
		case verbose:
			println(color.Ize(color.Bold + color.Green,"writing command : "+CommandString))
		}
		sshWriter(CommandString,sshIn)
		sshReadBuffer(sshOut, timeOutValue, CommandString, hostname)
	}
}