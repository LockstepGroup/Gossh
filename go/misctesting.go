package main
//
//import (
//	"bufio"
//	"crypto/aes"
//	"crypto/cipher"
//	"crypto/rand"
//	"encoding/hex"
//	"fmt"
//	"github.com/TwiN/go-color"
//	"github.com/jessevdk/go-flags"
//	"golang.org/x/crypto/ssh"
//	"io"
//	"os"
//	"regexp"
//	"strconv"
//	"strings"
//	"time"
//)
//
////Misc Functions Begin
////Global Variables
//
//var TesterrConnectionTimeout = regexp.MustCompile(".*dial tcp \\d+\\.\\d+\\.\\d+\\.\\d+:?\\d+?: i/o timeout.*")
//var TesterrKeyexchange = regexp.MustCompile(".*no common algorithm for key exchange.*")
//
//func TestsshHandleError(err error, fatal bool) {
//	type error interface {
//		Error() string
//	}
//	var errMatched = false
//	if nil != err {
//		ConnectionTimeout := TesterrConnectionTimeout.FindStringSubmatch(err.Error())
//		if ConnectionTimeout != nil {
//			println(color.Red + "Connection failed check Hostname, Port, Credentials, and network access." + color.Reset)
//			println(color.Yellow + "Original go err can be found bellow: \r\n\r\n" + color.Reset)
//			println(color.Red + err.Error() + color.Reset)
//			errMatched = true
//		}
//		if true != errMatched {
//			println(color.Red + "Error message was not matched if you have any information on this error \r\n " +
//				"please let the development team know on github \r\n\r\n" + color.Reset)
//			println(color.Red + err.Error())
//		}
//	}
//	if io.EOF == err {
//		println(color.Ize(color.Red, "Encountered EOF"))
//		println(color.Ize(color.Yellow, "Original go error can be found bellow:\n\n"))
//		println(color.Ize(color.Red, err.Error()))
//	}
//	if fatal {
//		os.Exit(1)
//	}
//}
//
//func TesthandleError(err error, fatal bool, message string) {
//	if nil != err {
//		type error interface {
//			Error() string
//		}
//		switch {
//		case len([]rune(message)) == 0:
//			println(color.Ize(color.Yellow, "uncategorized error raw error bellow:\n\n"))
//			println(color.Ize(color.Red, err.Error()))
//		case fatal:
//			println(color.Ize(color.Yellow, "Fatal error encountered terminating execution : "+message+"\n\n"))
//			println(color.Ize(color.Red, err.Error()))
//			os.Exit(1)
//
//		case len([]rune(message)) >= 1:
//			println(color.Ize(color.Yellow, message+"\n\n"))
//			println(color.Ize(color.Red, err.Error()))
//		}
//	}
//}
//func testreadConfigFile(filepath string) []string {
//	var ConfigArray []string
//	ConfigFile, err := os.Open(filepath)
//	TestsshHandleError(err, true)
//	defer func(ConfigFile *os.File) {
//		err := ConfigFile.Close()
//		TestsshHandleError(err, true)
//	}(ConfigFile)
//	FileScanner := bufio.NewScanner(ConfigFile)
//	for FileScanner.Scan() {
//		ConfigArray = append(ConfigArray, FileScanner.Text())
//	}
//	return ConfigArray
//}
//
//func TestwriteBuffer2Txt(fileName string, text string, append bool) {
//	if append {
//		text += "\n"
//	}
//	f, err := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
//	if nil != err {
//		TestsshHandleError(err, true)
//	}
//	if _, err := f.Write([]byte(text)); err != nil {
//		TestsshHandleError(err, true)
//	}
//	if err := f.Close(); err != nil {
//		TestsshHandleError(err, true)
//	}
//
//}
//
////Misc Functions end
//
////SSH Function begin
////Global Variables//
//
//var TestcmdEnd = regexp.MustCompile(".*@?.*([#>])\\s?$")
//var TestipRegex = regexp.MustCompile("\\d+\\.\\d+\\.\\d+\\.\\d+")
//
////Functions
//func TestencryptString(string2Encrypt string, keyString string) string {
//	key, _ := hex.DecodeString(keyString)
//	plaintext := []byte(string2Encrypt)
//	block, err := aes.NewCipher(key)
//	if nil != err {
//		TesthandleError(err, true, "Credential encryption failed terminating execution")
//	}
//	aesGCM, err := cipher.NewGCM(block)
//	if nil != err {
//		TesthandleError(err, true, "Credential encryption failed terminating execution")
//	}
//	nonce := make([]byte, aesGCM.NonceSize())
//	if _, err = io.ReadFull(rand.Reader, nonce); nil != err {
//		TesthandleError(err, true, "Credential encryption failed terminating execution")
//	}
//	ciphertext := aesGCM.Seal(nonce, nonce, plaintext, nil)
//	return fmt.Sprintf("%x", ciphertext)
//}
//
//func TestdecryptString(encryptedString string, keyString string) string {
//	key, _ := hex.DecodeString(keyString)
//	enc, _ := hex.DecodeString(encryptedString)
//	block, err := aes.NewCipher(key)
//	if err != nil {
//		TesthandleError(err, true, "Decryption failed terminating execution")
//	}
//	aesGCM, err := cipher.NewGCM(block)
//	if err != nil {
//		TesthandleError(err, true, "Decryption failed terminating execution")
//	}
//	nonceSize := aesGCM.NonceSize()
//	nonce, ciphertext := enc[:nonceSize], enc[nonceSize:]
//	plaintext, err := aesGCM.Open(nil, nonce, ciphertext, nil)
//	if err != nil {
//		TesthandleError(err, true, "Decryption failed terminating execution")
//	}
//	return fmt.Sprintf("%s", plaintext)
//}
//func TestsshWriter(cmd string, SshIN io.WriteCloser) {
//	_, err := SshIN.Write([]byte(cmd + "\r"))
//	TestsshHandleError(err, false)
//}
//var returnStrings string
//func TestsshReadBuffer(sshOut io.Reader, timeOutValue int, cmd string, hostName string) {
//	buf := make([]byte, 2048)
//	StartTime := time.Now()
//	bufferResults, err := sshOut.Read(buf);TesthandleError(err,false,"")
//	currentLine := ""
//	currentLine = string(buf[:bufferResults])
//ReaderLoop:
//	for nil == err && TestcmdEnd.FindStringSubmatch(currentLine) == nil {
//		n, err := sshOut.Read(buf)
//		CurrentTime := time.Now()
//		TimeElapsed := CurrentTime.Sub(StartTime)
//		TimeOutValue := time.Duration(timeOutValue) * time.Second
//		switch {
//		case nil != err:
//			TestsshHandleError(err, false)
//			break ReaderLoop
//		case TimeElapsed >= TimeOutValue:
//			directoryPath, err := os.UserHomeDir()
//			if err != nil {
//				TestsshHandleError(err, false)
//			}
//			println(color.Ize(color.Red, "Timeout Value exceeded please check error log for more details "+
//				"canceling read operation for cmd:= "+cmd+" | "+directoryPath+"/"+hostName+"-error.txt"))
//			TestwriteBuffer2Txt(directoryPath+"/"+hostName+"error.txt", returnStrings, true)
//			break ReaderLoop
//		default:
//			currentLine += string(buf[:n])
//		}
//	}
//	//println(color.Ize(color.Cyan, currentLine))   //uncomment while troubleshoot to see where the buffer is getting stuck
//
//	returnStrings += currentLine
//}
//func TestStringToArray(CLIAurguments string) (Commands []string) {
//	StringSplit := strings.Split(CLIAurguments, "||")
//	var Count int
//	for i := range StringSplit {
//		Count = i
//	}
//	var array = make([]string, Count+1)
//	for i, stringssplit := range StringSplit {
//		if stringssplit != " " {
//			array[i] = strings.TrimRight(stringssplit, " ")
//		}
//	}
//	fmt.Println(array)
//	return array
//}
//func TestsshDial(hostname string, port int, username string, password string, cmdS []string, timeOutValue int, enablePassword string) {
//	var portNumber = strconv.Itoa(port)
//	conn, err := ssh.Dial("tcp", hostname+":"+portNumber, &ssh.ClientConfig{
//		User:            username,
//		Timeout:         5 * time.Second,
//		Auth:            []ssh.AuthMethod{ssh.Password(password)},
//		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
//	})
//	TesthandleError(err, true, "SSH dial failed check device information and network")
//	session, err := conn.NewSession()
//	TesthandleError(err, true, "SSH dial failed check device information and network")
//	sshOut, err := session.StdoutPipe()
//	TesthandleError(err, true, "SSH out pipe creation failed terminating session")
//	sshIn, err := session.StdinPipe();TesthandleError(err, true, "SSH in pipe creation failed terminating session")
//	modes := ssh.TerminalModes{
//		ssh.ECHO:          0,     // disable echoing
//		ssh.TTY_OP_ISPEED: 14400, // input speed = 14.4kbaud
//		ssh.TTY_OP_OSPEED: 14400, // output speed = 14.4kbaud
//	}
//	// Request pseudo terminal
//	if err := session.RequestPty("xterm", 40, 80, modes);nil != err {TesthandleError(err,false,"")}
//	err = session.Shell();TesthandleError(err, true, "SSH in pipe creation failed terminating session")
//	defer func(conn *ssh.Client) {
//		err := conn.Close()
//		TesthandleError(err, false, "Connection failed to close")
//	}(conn)
//	defer func(session *ssh.Session) {
//		err := session.Close()
//		TesthandleError(err, false, "Session failed to close")
//	}(session)
//	switch { //created as a switch to allow the addition of buffer reading for click through prompts
//	case len([]rune(enablePassword)) >= 1:
//		TestsshWriter("en", sshIn)
//		TestsshWriter(enablePassword, sshIn)
//	}
//	for _, CommandString := range cmdS {
//		TestsshWriter(CommandString,sshIn)
//		TestsshReadBuffer(sshOut, timeOutValue, CommandString, hostname)
//	}
//}
//type Options struct {
//	HostName       string   `short:"h" long:"hostname" description:"device url/ip"`
//	UserName       string   `short:"u" long:"username" description:"username used for authentication"`
//	PassWord       string   `short:"p" long:"password" description:"password used for authentication"`
//	EnablePassword string   `short:"e" long:"enablepassword" description:"username used for authentication"`
//	Port           int      `short:"P" long:"port" description:"tcp port device is listening on" default:"22"`
//	TimeOut           int      `short:"t" long:"timeout" description:"timeout value in seconds for each command execution" default:"10"`
//	CliArguments   string `short:"C" long:"cliarguments" description:"array containing cmdS to be run" `
//	DirectoryPath       string   `short:"d" long:"filePath" description:"filePath to document containing cmdS to be run"`
//	UseFile           bool     `short:"f" long:"usefile" description:"bool switch for using cmdS file"`
//	Telnet            bool     `short:"T" long:"telnet" description:"option to use telnet instead of ssh" default:"false"`
//	Help           bool     `short:"H" long:"help" description:"output help"`
//}
//	var options Options
//func main() {
//	parser:= flags.NewParser(&options,flags.Default&^flags.HelpFlag)
//	_, err := parser.Parse()
//	if nil != err {os.Exit(0)}
//	if options.Help{parser.WriteHelp(os.Stdout);os.Exit(0)}
//	//bytes := make([]byte, 32)
//	//if _, err := rand.Read(bytes); err != nil {
//	//	TesthandleError(err, true, "failed to create encryption key")
//	//}
//	//key := hex.EncodeToString(bytes)
//	switch {
//	case options.Telnet:
//		println(color.Ize(color.Red,"Telnet functions are not ready yet terminating application")); os.Exit(1)
//	default:
//		TestsshDial(options.HostName, options.Port, options.UserName, options.PassWord,TestStringToArray(options.CliArguments),options.TimeOut,options.EnablePassword)
//	}
//	println(color.Ize(color.Bold + color.Purple,returnStrings))
//}
