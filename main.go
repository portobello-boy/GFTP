package main

import (
	"golang.org/x/crypto/ssh"
	"github.com/pkg/sftp"
	"github.com/miquella/ask"
	"time"
	"log"
)

func getPassword() (string, error) {
    err := ask.Print("Warning! I am about to ask you for a password!\n")
    if err != nil {
        return "", err
    }

    return ask.HiddenAsk("Password: ")
}

func main() {
	log.Print("Welcome to GFTP - An FTP service written in Golang")

	// Define variables needed
	host     := ""
	port     := ""
	user     := ""
	password, _ := getPassword()

	// Create SSH Client Config
	conf := &ssh.ClientConfig{
		User: user,					// Define the user in the config
		Auth: []ssh.AuthMethod{     // Set the authentication method (password, SSH key, etc.)
			ssh.Password(password),
		},
		Timeout: 2*time.Minute,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), // For now, ignore host key verification
	}

	// Connect to SSH host
	conn, err := ssh.Dial("tcp", host + port, conf)
	if err != nil {
		log.Print("Error connecting to SSH host!")
		log.Fatal(err)
	}
	defer conn.Close() // Ensure we close the connection when the program terminates

	// Create SFTP Client
	client, err := sftp.NewClient(conn)
	if err != nil {
		log.Print("Error creating SFTP client!")
		log.Fatal(err)
	}
	defer client.Close()
}