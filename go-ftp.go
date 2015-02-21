package main

import (
	"fmt"
	"github.com/micahhausler/go-ftp/server"
	"net"
	"os"
)

func main() {
	fmt.Println("Starting up FTP server")
	port := ":2121"
	ln, err := net.Listen("tcp", port)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	for {
		c, err := ln.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}

		fmt.Printf("Connection from %v established.\n", c.RemoteAddr())
		go server.HandleConnection(c)
	}

}
