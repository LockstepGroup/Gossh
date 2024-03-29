package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type StringArray struct {
	commands []string
}

func handleError(err error) {
	if err != nil {
		fmt.Println("Connection Failed check HostName and Credentials\\r\\n")
		fmt.Println(err)
	}
}
func StringToArray1(CLIAurguments string) (CLIReturn []string) {
	StringSplit := strings.Split(CLIAurguments, ",")
	var StringReturn []string
	for _, stringssplit := range StringSplit {
		StringReturn = append(StringReturn, stringssplit)
	}
	return StringReturn
}
func StringToArray(CLIAurguments string) (Commands []string) {
	StringSplit := strings.Split(CLIAurguments, "||")
	var Count int
	for i, _ := range StringSplit {
		Count = i
	}
	var array = make([]string, Count+1)
	for i, stringssplit := range StringSplit {
		if stringssplit != " " {
			array[i] = strings.TrimRight(stringssplit, " ")
		}
	}
	fmt.Println(array)
	return array
}
func readLines(path string) []string {
	file, err := os.Open(path)
	if err != nil {
		return nil
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines
}
