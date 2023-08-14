package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func main() {
	// Connect to the chat server
	conn, err := net.Dial("tcp", "127.0.0.1:12345")
	if err != nil {
		fmt.Println("Error connecting to the server:", err)
		os.Exit(1)
	}
	defer conn.Close()

	// Start a goroutine to receive messages from the server
	go func() {
		scanner := bufio.NewScanner(conn)
		for scanner.Scan() {
			received := scanner.Text()
			fmt.Println(received)
		}
		if scanner.Err() != nil {
			fmt.Println("Error reading from server:", scanner.Err())
		}
	}()

	// Read messages from the user and send them to the server
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		message := scanner.Text()
		_, err := conn.Write([]byte(message + "\n"))
		if err != nil {
			fmt.Println("Error sending message:", err)
			break
		}
	}

	if scanner.Err() != nil {
		fmt.Println("Error reading from user input:", scanner.Err())
	}
}
