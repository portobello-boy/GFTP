# GFTP
A simple SFTP application written in Golang. This tool was based on the content found here: [Golang SSH Client](https://networkbit.ch/golang-ssh-client/), [Golang SFTP Client](http://networkbit.ch/golang-sftp-client/)

## Installation
Ensure you have downloaded and installed a Golang interpreter/compiler. 

The dependencies are listed in the `go.mod` file. To download the dependencies, run `go mod download`. 

## Usage
At the moment, this application is quite rudimentary. To connect to a specific SSH/SFTP host, define the information in the `main()` function as so:
```
	host     := "123.456.123.456"
	port     := ":12345"
	user     := "johnsmith"
```
The application will prompt the user for a password, which will _not_ be displayed on the screen or stored anywhere. The service also supports the following commands:
```
    pwd             - Display working directory on remote host
    lpwd            - Display working directory on local machine
    ls              - List contents of directory on remote host
    lls             - List contents of directory on local machine
    lcd <directory> - Change directories on local machine
    put <file>      - Upload file to remote host
    get <file>      - Download file to local machine
```

## To Do
Most to-do items are listed in this repository's 'Issues' section, but some are listed here:
- Add `cd` support for remote host (might require different SSH library)
- Add command line arguments for host/port/user
- Separate command functionality into independent functions
- Add interpreter for non-SFTP specific commands such as running programs (check `os` libary, particularly `os.Run(...)`)
- Add unit tests using Golang's built-in test suite
