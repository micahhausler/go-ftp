package server

import (
	"fmt"
	"os"
)

// A Server handles configuration
type Server struct {
	Port           int
	Callback       PasswordCallback
	PostProcessors []FileProcessor
}

// A PasswordCallback is a fuc that accepts a username and password, and returns a bool if the user was authenticated
type PasswordCallback func(username, password []byte) bool

// A FileProcessor is a func that accepts a path and returns an optional error
type FileProcessor func(string) error

// CleanupFile removes the file
func CleanupFile(path string) error {
	return os.Remove(path)
}
