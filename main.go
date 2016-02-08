package main

import (
	"fmt"
	"os"

	"github.com/micahhausler/go-ftp/server"
	flag "github.com/spf13/pflag"
)

var addressP = flag.StringP("address", "a", "127.0.0.1", "Default listen address to use")
var portP = flag.IntP("port", "p", 2121, "Default port to use")
var versionP = flag.BoolP("version", "v", false, "Print version and exit")

// The binary version
const Version = "0.0.2"

func main() {
	if *versionP {
		fmt.Printf("go-ftp %s\n", Version)
		os.Exit(0)
	}

	flag.Parse()

	fmt.Println("Starting up FTP server")

	serv := server.Server{
		Port: *portP,
	}
	serv.Run()

}
