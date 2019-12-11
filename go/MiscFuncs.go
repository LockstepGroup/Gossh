package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"strings"
	"time"

	pan "github.com/zepryspet/GoPAN/utils"
)

func handleError(err error) {
	if err != nil {
		panic(err)
	}
}
func readBuffForString2(sshOut io.Reader) string {
	buf := make([]byte, 1000)
	n, err := sshOut.Read(buf) //this reads the ssh terminal
	waitingString := ""
	if err == nil {
		for _, v := range buf[:n] {
			fmt.Sprint("%c", v)
		}
		waitingString = string(buf[:n])
	}
	for err == nil {
		// this loop will not end!!
		n, err = sshOut.Read(buf)
		waitingString += string(buf[:n])
		for _, v := range buf[:n] {
			fmt.Sprint("%c", v)
		}
		if err != nil {
			fmt.Println(err)
		}
	}
	return waitingString
}
func readBuffForString1(sshOut io.Reader) string {
	buf := make([]byte, 1000)
	n, err := sshOut.Read(buf) //this reads the ssh terminal
	waitingString := ""
	if err == nil {
		for _, v := range buf[:n] {
			fmt.Printf("%c", v)
		}
		waitingString = string(buf[:n])
	}
	for err == nil {
		// this loop will not end!!
		n, err = sshOut.Read(buf)
		waitingString += string(buf[:n])
		for _, v := range buf[:n] {
			fmt.Printf("%c", v)
		}
		if err != nil {
			fmt.Println(err)
		}
	}
	return waitingString
}
func write(cmd string, sshIn io.WriteCloser) {
	_, err := sshIn.Write([]byte(cmd + "\r"))
	handleError(err)
}

//This is a wrapper function, wraps the original function and passes the password to avoid creating a public variable
func Challenge(pass string) func(string, string, []string, []bool) ([]string, error) {
	return func(user, instruction string, questions []string, echos []bool) (answers []string, err error) {
		answers = make([]string, len(questions))
		for n, q := range questions {
			fmt.Printf("%s\n", q)
			if q == "Password: " {
				answers[n] = pass
				fmt.Printf("Signing In...")
			} else if q == "Do you accept and acknowledge the statement above ? (yes/no) : " {
				answers[n] = "yes"
				fmt.Printf("yes\nauto acknowdleged banner...\n")
			} else {
				reader := bufio.NewReader(os.Stdin)
				r, _ := reader.ReadString('\n')
				answers[n] = strings.TrimSpace(r)
			}
		}
		return answers, nil
	}
}

func cmdSend(sshOut io.Reader, sshIn io.WriteCloser, cmd string, isFile bool, isConfig bool, timeout int) {
	//setting up the initial prompt
	prompt := ">"
	readBuff(prompt, sshOut, timeout)
	//disabling the CLI pager to avoid having to tab on large outputs
	if _, err := writeBuff("set cli pager off", sshIn); err != nil {
		pan.Logerror(err, true)
	}
	readBuff(prompt, sshOut, timeout)

	//verifying if the commands need to be run in config mode
	if isConfig {
		//changing the prompt to bash as due config mode
		prompt = "#"
		//Sending configuration
		if _, err := writeBuff("configure", sshIn); err != nil {
			pan.Logerror(err, true)
		}
		readBuff(prompt, sshOut, timeout)
	}

	//Sending the command (s) to the endpoints

	//checking if it's a file or a single command
	if isFile {
		file, err := os.Open(cmd)
		if err != nil {
			pan.Logerror(err, true)
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			//removing empty spaces
			newline := strings.TrimSpace(scanner.Text())
			if _, err := writeBuff(newline, sshIn); err != nil {
				pan.Logerror(err, true)
			}
			readBuff(prompt, sshOut, timeout)
		}
		if err := scanner.Err(); err != nil {
			pan.Logerror(err, true)
		}
	} else {
		if _, err := writeBuff(cmd, sshIn); err != nil {
			pan.Logerror(err, true)
		}
		readBuff(prompt, sshOut, timeout)
	}
	//Exiting from config mode
	if isConfig {
		if _, err := writeBuff("exit", sshIn); err != nil {
			pan.Logerror(err, true)
		}
	}
	//Exiting from operational mode
	if _, err := writeBuff("exit", sshIn); err != nil {
		pan.Logerror(err, true)
	}
}

func readBuffForString(whattoexpect string, sshOut io.Reader, buffRead chan<- string) {
	buf := make([]byte, 2000)
	n, err := sshOut.Read(buf) //this reads the ssh terminal
	waitingString := ""
	if err == nil {
		waitingString = string(buf[:n])
	}
	for (err == nil) && (!strings.Contains(waitingString, whattoexpect)) {
		n, err = sshOut.Read(buf)
		waitingString += string(buf[:n])
		//fmt.Println(waitingString) //uncommenting this might help you debug if you are coming into errors with timeouts when correct details entered
	}
	fmt.Println(waitingString)
	//pan.Wlog("output.txt", waitingString, true)
	buffRead <- waitingString
}

func readBuff(whattoexpect string, sshOut io.Reader, timeoutSeconds int) string {
	ch := make(chan string)
	go func(whattoexpect string, sshOut io.Reader) {
		buffRead := make(chan string)
		go readBuffForString(whattoexpect, sshOut, buffRead)
		select {
		case ret := <-buffRead:
			ch <- ret
		case <-time.After(time.Duration(timeoutSeconds) * time.Second):
			pan.Wlog("error.txt", "timeout waiting for command", true)
			break
		}
	}(whattoexpect, sshOut)
	return <-ch
}

func writeBuff(Config string, sshIn io.WriteCloser) (int, error) {
	returnCode, err := sshIn.Write([]byte(Config + "\r"))
	return returnCode, err
}

///////////////////////////////////////////////////
///////////Test Functions for netscan//////////////
///////////////////////////////////////////////////
func inc(ip net.IP) {
	for j := len(ip) - 1; j >= 0; j-- {
		ip[j]++
		if ip[j] > 0 {
			break
		}
	}
}
func cmdSendsonic(sshOut io.Reader, sshIn io.WriteCloser, cmd string, isFile bool, isConfig bool, timeout int) {
	//setting up the initial prompt
	prompt := ">"
	readBuff(prompt, sshOut, timeout)
	//disabling the CLI pager to avoid having to tab on large outputs
	if _, err := writeBuff("no cli pager session", sshIn); err != nil {
		pan.Logerror(err, true)
	}
	readBuff(prompt, sshOut, timeout)

	//verifying if the commands need to be run in config mode
	if isConfig {
		//changing the prompt to bash as due config mode
		prompt = "#"
		//Sending configuration
		if _, err := writeBuff("configure", sshIn); err != nil {
			pan.Logerror(err, true)
		}
		readBuff(prompt, sshOut, timeout)
	}
	//Sending the command (s) to the endpoints

	//checking if it's a file or a single command
	if isFile {
		file, err := os.Open(cmd)
		if err != nil {
			pan.Logerror(err, true)
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			//removing empty spaces
			newline := strings.TrimSpace(scanner.Text())
			if _, err := writeBuff(newline, sshIn); err != nil {
				pan.Logerror(err, true)
			}
			readBuff(prompt, sshOut, timeout)
		}
		if err := scanner.Err(); err != nil {
			pan.Logerror(err, true)
		}
	} else {
		if _, err := writeBuff(cmd, sshIn); err != nil {
			pan.Logerror(err, true)
		}
		readBuff(prompt, sshOut, timeout)
	}
	//Exiting from config mode
	if isConfig {
		if _, err := writeBuff("exit", sshIn); err != nil {
			pan.Logerror(err, true)
		}
	}
	//Exiting from operational mode
	if _, err := writeBuff("exit", sshIn); err != nil {
		pan.Logerror(err, true)
	}
}
