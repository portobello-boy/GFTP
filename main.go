package main

import (
	"golang.org/x/crypto/ssh"
	"github.com/pkg/sftp"
	"github.com/miquella/ask"

	"fmt"
	"os"
	"os/exec"
	"io"
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

func interpret(input string, client *sftp.Client) {
	splits := strings.Split(input, " ")

	switch splits[0] {
		case "pwd":
			wd, _ := client.Getwd()
			println(wd)
			break
		case "lpwd":
			wd, _ := os.Getwd()
			println(wd)
			break
		case "ls":
			wd, _ := client.Getwd()
			fileInfos, _ := client.ReadDir(wd)
			for _, f := range fileInfos {
				println(f.Name())
			}
			break
		case "lls":
			out, _ := exec.Command("ls").Output()
			println(string(out[:]))
			break
		// case "cd" has issues with SSH session - only one command can be run per session
		case "lcd":
			err := os.Chdir(splits[1])
			if err != nil {
				log.Fatal(err)
			}
		case "put":
			// Create source and destination file pointers
			dest, err := client.Create(splits[1])
			if err != nil {
				log.Fatal(err)
			}
			defer dest.Close()

			src, err := os.Open(splits[1])
			if err != nil {
				log.Fatal(err)
			}
			defer src.Close()

			// Copy into destination the data from source
			bytes, err := io.Copy(dest, src)
			if err != nil {
				log.Fatal(err)
			}
			
			log.Printf("%d bytes uploaded\n", bytes)
			break
		case "get":
			// Create source and destination file pointers
			dest, err := os.Create(splits[1])
			if err != nil {
				log.Fatal(err)
			}
			defer dest.Close()

			src, err := client.Open(splits[1])
			if err != nil {
				log.Fatal(err)
			}
			defer src.Close()

			// Copy into destination the data from source
			bytes, err := io.Copy(dest, src)
			if err != nil {
				log.Fatal(err)
			}

			// Flush in-memory copy
			err = dest.Sync()
			if err != nil {
				log.Fatal(err)
			}
			log.Printf("%d bytes downloaded\n", bytes)
			break
		
	}
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
		Timeout: 30*time.Second,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), // For now, ignore host key verification
	}

	// Connect to SSH host
	log.Print("Attempting to connect to: ", host + port, " as user ", user)
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

		input := scanner.Text()
		if strings.TrimSpace(input) == "" {
			continue
		}

		if input == "exit" {
			break
		}

		interpret(input, client)

		fmt.Print("GFTP > ")
	}
}