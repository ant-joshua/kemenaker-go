package main

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

func main() {
	const SSH_ADDRESS = "103.193.176.131:22"
	const SSH_USERNAME = "hacktiv8"
	const SSH_PASSWORD = "Cobacoba123!"

	sshConfig := &ssh.ClientConfig{
		User:            SSH_USERNAME,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Auth: []ssh.AuthMethod{
			ssh.Password(SSH_PASSWORD),
		},
	}

	client, err := ssh.Dial("tcp", SSH_ADDRESS, sshConfig)

	if client != nil {
		defer client.Close()
	}

	if err != nil {
		fmt.Println("Failed to dial: ", err.Error())
	}

	session, err := client.NewSession()
	if session != nil {
		defer session.Close()
	}

	if err != nil {
		fmt.Println("Failed to create session: ", err.Error())
	}

	// session.Stdin = os.Stdin
	session.Stdout = os.Stdout
	session.Stderr = os.Stderr

	// err = session.Run("ls -l ~/")

	// if err != nil {
	// 	fmt.Println("Failed to run: " + err.Error())
	// }

	commands := []string{
		"ls -l ~/",
		"cd ~/ini_project",
		"ls -l",
		"cat joshua.txt",
	}

	stdin, err := session.StdinPipe()
	if err != nil {
		fmt.Println(err)
	}

	// err = session.Shell()
	err = session.Start("/bin/bash")

	if err != nil {
		fmt.Println("Failed to start: " + err.Error())
	}

	for _, cmd := range commands {
		if _, err := fmt.Fprintln(stdin, cmd); err != nil {
			fmt.Println(err)
		}
	}

	sftpConnect(client)

	err = session.Wait() // wait for session to finish executing , only triggered by exit command

	if err != nil {
		fmt.Println("Failed to wait: " + err.Error())
	}

}

func sftpConnect(client *ssh.Client) {
	stfpClient, err := sftp.NewClient(client)

	if err != nil {
		fmt.Println("Failed to create sftp client: ", err.Error())
	}

	fDestination, err := stfpClient.Create("/home/hacktiv8/ini_project/joshua.txt")

	if err != nil {
		fmt.Println("Failed to create file: ", err.Error())
	}

	fSource, err := os.Open("/Users/joshua/joshua.txt")

	if err != nil {
		fmt.Println("Failed to open file: ", err.Error())
	}

	_, err = io.Copy(fDestination, fSource)

	if err != nil {
		fmt.Println("Failed to copy file: ", err.Error())
	}

	log.Println("Success")
}
