package main

import (
	"golang.org/x/crypto/ssh"
	"github.com/pkg/sftp"
	"github.com/miquella/ask"

	"fmt"
	"os"
	"bufio"
	"time"
	"strings"
	"log"
)

func getPassword() (string, error) {
    err := ask.Print("Warning! I am about to ask you for a password!\n")
    if err != nil {
        return "", err
    }

    return ask.HiddenAsk("Password: ")
}

func interpret(command string, client *sftp.Client) {
	switch command {
		case "exit":
			break
		case "pwd":
			wd, _ := client.Getwd()
			println(wd)
		case "lpwd":
			wd, _ := os.Getwd()
			println(wd)
	}
}

func main() {
	log.Print("Welcome to GFTP - An FTP service written in Golang")

	// Define variables needed
	host     := "159.230.226.130"
	port     := ":238"
	user     := "danielmillson"
	password, _ := getPassword()

	// Create SSH Client Config
	conf := &ssh.ClientConfig{
		User: user,					// Define the user in the config
		Auth: []ssh.AuthMethod{     // Set the authentication method (password, SSH key, etc.)
			ssh.Password(password),
		},
		Timeout: 30*time.Second,
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

	// Initialiez an interactive shell
	fmt.Print("GFTP > ")
	scanner := bufio.NewScanner(os.Stdin)

	for {
		if !scanner.Scan() {
			continue
		}

		command := scanner.Text()
		if strings.TrimSpace(command) == "" {
			continue
		}

		// println(command)
		interpret(command, client)

		if command == "exit" {
			break
		}

		fmt.Print("GFTP > ")
	}
}