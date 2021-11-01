package main

import (
	"fmt"
	pan "github.com/zepryspet/GoPAN/utils"
	"io"
	"log"
	"os"
	"regexp"
	"time"
)
var prompt = regexp.MustCompile("(%{2})")
var prompt2 = regexp.MustCompile(".*@?.*(#|>).*$")

func write(cmd string, sshIn io.WriteCloser) {
	_, err := sshIn.Write([]byte(cmd + "\r%%\r"))
	handleError(err)
}
func writeblank(sshIn io.WriteCloser) {
	_, err := sshIn.Write([]byte("\r"))
	handleError(err)
}
func readBuffForString(sshOut io.Reader, buffRead chan string){
	buf := make([]byte, 1000)
	waitingString := ""
	for {
		n, err := sshOut.Read(buf) //this reads the ssh terminal
		if err == io.EOF {
			fmt.Println(err)
			break
		}
		if err != nil {
			fmt.Println(err)
			break
		}
		current := string(buf[:n])
		// add current line to result string
		returnstring += current
		waitingString += current
		m := prompt.FindStringSubmatch(current)
		if m != nil{
			break
		}

	}
	buffRead <- waitingString
}
func readBuffForStringblank(sshOut io.Reader, buffRead chan string){
	buf := make([]byte, 1000)
	waitingString := ""
	for {
		n, err := sshOut.Read(buf) //this reads the ssh terminal
		if err == io.EOF {
			fmt.Println(err)
			break
		}
		if err != nil {
			fmt.Println(err)
			break
		}
		// for every line
		current := string(buf[:n])
		// add current line to result string
		waitingString += current
		m := prompt2.FindStringSubmatch(current)
		if m != nil{
			break
		}

	}
	buffRead <- waitingString
}
func readBuff(sshOut io.Reader, timeoutSeconds int,HostName string) string {
	dirname, err := os.UserHomeDir()
	if err != nil {
		log.Fatal( err )
	}
	ch := make(chan string)
	go func(sshOut io.Reader) {
		buffRead := make(chan string)
		go readBuffForString(sshOut, buffRead)
		select {
		case ret := <-buffRead:
			ch <- ret
		case <-time.After(time.Duration(timeoutSeconds) * time.Second):
			pan.Wlog(dirname+"/"+HostName+"-error.txt", "timeout waiting for command |" +HostName+"\r\n--------" +returnstring+"\r\n--------", true)
			os.Exit(3)
		}
	}(sshOut)
	return <-ch
}
func readBuffblank(sshOut io.Reader, timeoutSeconds int,HostName string) string {
	dirname, err := os.UserHomeDir()
	if err != nil {
		log.Fatal( err )
	}
	ch := make(chan string)
	go func(sshOut io.Reader) {
		buffRead := make(chan string)
		go readBuffForStringblank(sshOut, buffRead)
		select {
		case ret := <-buffRead:
			ch <- ret
		case <-time.After(time.Duration(timeoutSeconds) * time.Second):
			pan.Wlog(dirname+"/error.txt", "timeout waiting for command |" +HostName+"\r\n--------" +returnstring+"\r\n--------", true)
			os.Exit(3)
		}
	}(sshOut)
	return <-ch
}