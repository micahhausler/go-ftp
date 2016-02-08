package server

import (
	"errors"
	"fmt"
	"path"
	"strconv"
	"strings"
)

func stringInList(search string, data []string) bool {
	for _, d := range data {
		if search == d {
			return true
		}
	}
	return false
}

func parsePortArgs(arg string) string {
	// Parses a PORT argument and returns an IP addr and port
	// Ex: "10,0,0,1,192,127" and returns "10.0.0.1:49279"
	parts := strings.Split(arg, ",")
	ip := strings.Join(parts[:4], ".")
	p1, _ := strconv.Atoi(parts[4])
	p2, _ := strconv.Atoi(parts[5])
	port := p1*256 + p2

	return fmt.Sprintf("%s:%d", ip, port)
}

func stripDirectory(remoteName string) string {
	_, filename := path.Split(remoteName)
	return filename
}

// FtpCommands returns a map of FTP commands with the keys being all known FTP
// commands and the values stating if the command is implemented
// https://en.wikipedia.org/wiki/List_of_FTP_commands
func FtpCommands() map[string]bool {
	cmds := map[string]bool{
		"ABOR": true,
		"ACCT": false,
		"ADAT": false,
		"ALLO": false,
		"APPE": false,
		"AUTH": true,
		"CCC":  false,
		"CDUP": true,
		"CONF": false,
		"CWD":  true,
		"DELE": true,
		"ENC":  false,
		"EPRT": false,
		"EPSV": false,
		"FEAT": true,
		"HELP": false,
		"LANG": false,
		"LIST": false,
		"LPRT": false,
		"LPSV": false,
		"MDTM": false,
		"MIC":  false,
		"MKD":  true,
		"MLSD": false,
		"MLST": false,
		"MODE": false,
		"NLST": true,
		"NOOP": true,
		"OPTS": true,
		"PASS": true,
		"PASV": true,
		"PBSZ": false,
		"PORT": true,
		"PROT": false,
		"PWD":  true,
		"QUIT": true,
		"REIN": true,
		"REST": true,
		"RETR": false,
		"RMD":  true,
		"RNFR": false,
		"RNTO": false,
		"SITE": false,
		"SIZE": true,
		"SMNT": false,
		"STAT": true,
		"STOR": true,
		"STOU": false,
		"STRU": false,
		"SYST": true,
		"TYPE": true,
		"USER": true,
		"XCUP": false,
		"XMKD": false,
		"XMKD": false,
		"XPWD": false,
		"XRCP": false,
		"XRMD": false,
		"XRSQ": false,
		"XSEM": false,
		"XSEN": false,
	}
	return cmds
}

func parseCommand(input string) (string, string, error) {
	// Split out command and arguments
	var command, args string
	var err error

	// commands are all 3 or 4 characters
	if len(input) < 3 {
		return command, args, errors.New(SyntaxErr)
	}

	response := strings.SplitAfterN(input, " ", 2)

	switch {
	case len(response) == 2:
		command = strings.TrimSpace(response[0])
		args = strings.TrimSpace(response[1])
	case len(response) == 1:
		command = response[0]
	}
	return command, args, err
}
