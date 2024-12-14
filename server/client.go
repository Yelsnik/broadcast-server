package server

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

type Client struct {
	Conn net.Conn
	Name string
}

func ConnectToServer(address string) {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		log.Fatalf("Error connecting to server: %v", err.Error())
	}
	defer conn.Close()

	fmt.Println("Connected to the server. Type your messages below:")

	// Handle incoming messages from the server in a separate goroutine
	go func() {
		scanner := bufio.NewScanner(conn)
		for scanner.Scan() {
			fmt.Println(scanner.Text())
		}
		if err := scanner.Err(); err != nil {
			log.Printf("Error reading from server: %v", err)
		}
		os.Exit(0)
	}()

	// Read input from the user and send it to the server
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		message := scanner.Text()
		_, err := fmt.Fprintln(conn, message)
		if err != nil {
			log.Printf("Error sending message to server: %v", err)
			return
		}
	}

	if err := scanner.Err(); err != nil {
		log.Printf("Error reading from input: %v", err)
	}
}
