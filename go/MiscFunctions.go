package main

import (
	"bufio"
	"fmt"
	"github.com/TwiN/go-color"
	"os"
	"strings"
)

func handleError(err error, fatal bool, message string) {
	if nil != err {
		switch {
		case len([]rune(message)) == 0:
			fmt.Println(color.Ize(color.Bold + color.Red,"uncategorized error raw error bellow:"))
			fmt.Println(err)
		case fatal:
			fmt.Println(color.Ize(color.Bold + color.Red,"Fatal error encountered terminating execution : "+message))
			fmt.Println(err)
			os.Exit(1)

		case len([]rune(message)) >= 1:
			fmt.Println(color.Ize(color.Bold + color.Red,message))
			fmt.Println(err)
		}
	}
}

func readConfigFile(filepath string) []string {
	var ConfigArray []string
	ConfigFile, err := os.Open(filepath)
	handleError(err, true,"Error encountered opening config file")
	defer func(ConfigFile *os.File) {
		err := ConfigFile.Close()
		handleError(err, true,"Error encountered closing config file")
	}(ConfigFile)
	FileScanner := bufio.NewScanner(ConfigFile)
	for FileScanner.Scan() {
		ConfigArray = append(ConfigArray, FileScanner.Text())
	}
	return ConfigArray
}

func writeBuffer2Txt(fileName string, text string, append bool) {
	if append {
		text += "\n"
	}
	f, err := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if nil != err {
		handleError(err, true,"Error encountered opening log file")
	}
	if _, err := f.Write([]byte(text)); err != nil {
		handleError(err, true, "Error encountered writing to log file")
	}
	if err := f.Close(); err != nil {
		handleError(err, true,"Error encountered closing log file")
	}

}

func stringToArray(CLIArguments string) (Commands []string) {
	StringSplit := strings.Split(CLIArguments, "||")
	var Count int
	for i := range StringSplit {
		Count = i
	}
	var array = make([]string, Count+1)
	for i, stringsSplit := range StringSplit {
		if stringsSplit != " " {
			array[i] = strings.TrimRight(stringsSplit, " ")
		}
	}
	if verbose {println(color.Ize(color.Yellow,"Commands array:\n" + strings.Join(array,"\n")))}
	return array
}