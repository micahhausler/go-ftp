package server

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"path"
	"strconv"
	"strings"
	"time"
)

const (
	DataUrl    = "https://data.ambition.io"
	storageDir = "uploads"
)

type ConnectionConfig struct {
	DataConnectionAddr string
	Filename           string
}

// Run runs the server
func (s *Server) Run() {

	port := fmt.Sprintf(":%s", strconv.Itoa(s.Port))
	listener, err := net.Listen("tcp", port)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	for {
		c, err := listener.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}
		fmt.Printf("Connection from %v established.\n", c.RemoteAddr())
		conn := Connection{
			Conn: c,
		}

		go s.handleConnection(conn)
	}
}

// A type for handling Credentials
type CredentialSet struct {
	Username string
	Password string
	Valid    bool
}

// A type for handling uploaded files
type UploadedFile struct {
	Filename string
	Dirctory string
}

// A Connection handles an individual connection from a client
type Connection struct {
	Conn               net.Conn
	Credential         CredentialSet
	DataConnectionAddr string
	Filename           string
}

// Send a message over the connection
func (c Connection) SendMsg(msg string) {
	fmt.Printf("Sending: %s", msg)
	io.WriteString(c.Conn, msg)
}

// Recieve a message over the connection
func (c Connection) GetMsg() string {
	bufc := bufio.NewReader(c.Conn)
	for {
		line, err := bufc.ReadString('\n')
		if err != nil {
			c.Conn.Close()
			break
		}
		fmt.Printf("Received: %s", line)
		return strings.TrimRight(line, "\r")
	}
	return ""
}

// A AuthProcessor is a func that accepts a CredentialSet and returns that CredentialSet
type AuthProcessor func(*CredentialSet) CredentialSet

func (s Server) HandleConn(c net.Conn) {
	defer c.Close()
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Recovered in HandleConn() for remote %s: %s", c.RemoteAddr(), r)
		}
	}()
	s.handleConn(c)
}

func (s Server) handleConn(c net.Conn) {
	set := CredentialSet{Valid: false}
	c.Credential = set
	c.SendMsg(FtpServerReady)
	for {
		response := handleAuth(&c, &set, s.AuthHandler)
		c.SendMsg(response)
		fmt.Printf("Credential \"%s\" valid: %v\n", c.Credential.Username, c.Credential.Valid)
		if c.Credential.Valid == true {
			break
		}
	}
	for {
		cmd := c.GetMsg()
		response, err := c.handleCmd(cmd)
		if err != nil {
			break
		}
		c.SendMsg(response)
		time.Sleep(100 * time.Millisecond)
	}

}

// Handles input post-authentication
func (c *Connection) handleCmd(input string) (string, error) {

	input = strings.TrimSpace(input)
	cmd, args, err := parseCommand(input)
	if err != nil {
		fmt.Printf("%s from %v: %s\n", SyntaxErr, c.Conn.RemoteAddr(), input)
		return SyntaxErr, err
	}

	ignoredCommands := []string{
		"CDUP", // cd to parent dir
		"RMD",  // remove directory
		"RNFR", // rename file from
		"RNTO", // rename file to
		"SITE", // execute arbitrary command
		"SIZE", // Size of a file
		"STAT", // Get status of FTP server
	}
	notImplemented := []string{
		"EPSV",
		"EPRT",
	}
	switch {
	case stringInList(cmd, ignoredCommands):
		return CmdNotImplmntd, nil
	case stringInList(cmd, notImplemented):
		return CmdNotImplmntd, nil
	case cmd == "NOOP":
		return CmdOk, nil
	case cmd == "SYST":
		return SysType, nil
	case cmd == "STOR":
		c.Filename = stripDirectory(args)

		c.readPort()
		// Don't upload for now
		//go uploadData(user, getFileName(user.username, ch.Filename))
		//TODO loop through post processing
		return TxfrCompleteOk, nil
	case cmd == "FEAT":
		return FeatResponse, nil
	case cmd == "PWD":
		return PwdResponse, nil
	case cmd == "TYPE" && args == "I":
		return TypeSetOk, nil
	case cmd == "PORT":
		c.DataConnectionAddr = parsePortArgs(args)
		return PortOk, nil
	case cmd == "PASV":
		// todo set up PASV mode
		//return EnteringPasvMode, nil
		return CmdNotImplmntd, nil
	case cmd == "QUIT":
		return GoodbyeMsg, nil
	}
	return "", nil
}

// HandleConnection handles a connection
func HandleConnection(c net.Conn) {

	defer c.Close()

	sendMsg(c, FtpServerReady)
	user := AuthUser{}

	for {
		message := getMsg(c)
		response := handleLogin(message, &user)
		sendMsg(c, response)
		if user.valid == true {
			break
		}
	}

	connConfig := ConnectionConfig{}

	for {
		cmd := getMsg(c)
		response, err := handleCommand(cmd, &connConfig, &user, c)
		if err != nil {
			break
		}
		sendMsg(c, response)
		time.Sleep(100 * time.Millisecond)
	}
}

func handleCommand(input string, ch *ConnectionConfig, user *AuthUser, c net.Conn) (string, error) {
	// Handles input after authentication

	input = strings.TrimSpace(input)
	cmd, args, err := parseCommand(input)
	if err != nil {
		fmt.Printf("%s from %v: %s\n", SyntaxErr, c.RemoteAddr(), input)
		return SyntaxErr, err
	}

	ignoredCommands := []string{
		"CDUP", // cd to parent dir
		"RMD",  // remove directory
		"RNFR", // rename file from
		"RNTO", // rename file to
		"SITE", // execute arbitrary command
		"SIZE", // Size of a file
		"STAT", // Get status of FTP server
	}
	notImplemented := []string{
		"EPSV",
		"EPRT",
	}

	switch {
	case stringInList(cmd, ignoredCommands):
		return CmdNotImplmntd, nil
	case stringInList(cmd, notImplemented):
		return CmdNotImplmntd, nil
	case cmd == "NOOP":
		return CmdOk, nil
	case cmd == "SYST":
		return SysType, nil
	case cmd == "STOR":
		ch.Filename = stripDirectory(args)
		readPortData(ch, user.username, c)
		// Don't upload for now
		//go uploadData(user, getFileName(user.username, ch.Filename))
		return TxfrCompleteOk, nil
	case cmd == "FEAT":
		return FeatResponse, nil
	case cmd == "PWD":
		return PwdResponse, nil
	case cmd == "TYPE" && args == "I":
		return TypeSetOk, nil
	case cmd == "PORT":
		ch.DataConnectionAddr = parsePortArgs(args)
		return PortOk, nil
	case cmd == "PASV":
		// todo set up PASV mode
		//return EnteringPasvMode, nil
		return CmdNotImplmntd, nil
	case cmd == "QUIT":
		return GoodbyeMsg, nil
	}
	return "", nil
}

func uploadData(user *AuthUser, filePath string) {
	// Upload to data.ambition

	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		fmt.Printf("Error reading file %s: %s\n", filePath, err)
		return
	}
	_, filename := path.Split(filePath)

	uri := fmt.Sprint(DataUrl, "/api/upload/ftp-upload/")

	resp, err := http.PostForm(uri,
		url.Values{
			"username":       {user.username},
			"password":       {user.password},
			"local_filename": {filename},
			"data":           {string(content)}})

	if resp.StatusCode == http.StatusCreated && err == nil {
		fmt.Printf("File %s uploaded to data!\n", filename)
		err = os.Remove(filePath)
		if err != nil {
			fmt.Printf("Error removing file %s: %s\n", filePath, err)
			return
		}
	} else {
		fmt.Printf("Could not upload '%s' to data! %s\n", filePath, err)
	}

}

func getFileName(username, filename string) string {
	return path.Join(storageDir, username, filename)
}

func (c *Connection) readPort() {
	fmt.Printf("connecting to %s\n", c.DataConnectionAddr)
	var err error
	reverseConn, err := net.Dial("tcp", c.DataConnectionAddr)
	reverseConn.SetReadDeadline(time.Now().Add(time.Minute))
	defer reverseConn.Close()
	if err != nil {
		fmt.Printf("Connection to %s errored out: %s\n", c.DataConnectionAddr, err)
		return
	}
	c.SendMsg(DataCnxAlreadyOpenStartXfr)

	err = os.MkdirAll(path.Join(storageDir, c.Credential.Username), 0777)
	if err != nil {
		fmt.Printf("error creating dir: %s\n", err)
		return
	}

	outputName := getFileName(c.Credential.Username, c.Filename)
	file, err := os.Create(outputName)
	defer file.Close()
	if err != nil {
		fmt.Printf("error creating file '%s': %s\n", outputName, err)
		return
	}

	reader := bufio.NewReader(reverseConn)
	buf := make([]byte, 1024) // big buffer
	for {
		n, err := reader.Read(buf)
		if err != nil && err != io.EOF {
			fmt.Println("read error:", err)
			break
		}
		if n == 0 {
			break
		}
		if _, err := file.Write(buf[:n]); err != nil {
			fmt.Println("read error:", err)
			break
		}
	}
}

func readPortData(ch *ConnectionConfig, username string, out net.Conn) {
	// Read data from the client, write out to file
	fmt.Printf("connecting to %s\n", ch.DataConnectionAddr)

	var err error

	c, err := net.Dial("tcp", ch.DataConnectionAddr)
	// set timeout of one minute
	c.SetReadDeadline(time.Now().Add(time.Minute))
	defer c.Close()
	if err != nil {
		fmt.Printf("connection to %s errored out: %s\n", ch.DataConnectionAddr, err)
		return
	}
	sendMsg(out, DataCnxAlreadyOpenStartXfr)

	err = os.MkdirAll(path.Join(storageDir, username), 0777)
	if err != nil {
		fmt.Printf("error creating dir: %s\n", err)
		return
	}

	outputName := getFileName(username, ch.Filename)
	file, err := os.Create(outputName)
	defer file.Close()
	if err != nil {
		fmt.Printf("error creating file '%s': %s\n", outputName, err)
		return
	}

	reader := bufio.NewReader(c)
	buf := make([]byte, 1024) // big buffer
	for {
		n, err := reader.Read(buf)
		if err != nil && err != io.EOF {
			fmt.Println("read error:", err)
			break
		}
		if n == 0 {
			break
		}
		if _, err := file.Write(buf[:n]); err != nil {
			fmt.Println("read error:", err)
			break
		}
	}
}

// Read a message from the connection
func getMsg(conn net.Conn) (string, error) {
	// Split the response into CMD and ARGS
	bufc := bufio.NewReader(conn)
	for {
		line, err := bufc.ReadString('\n')
		if err != nil {
			return nil, err
		}
		log.Printf("Received: %s\n", line)
		return strings.TrimRight(line, "\r"), nil
	}
}

func sendMsg(c net.Conn, message string) error {
	log.Printf("Sending: %s\n", message)
	_, err := io.WriteString(c, message)
	return err
}
