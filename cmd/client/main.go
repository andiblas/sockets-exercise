package main

import (
	"fmt"
	"net"
)

func main() {
	// Establish a TCP connection to a server
	conn, err := net.Dial("tcp", "127.0.0.1:8080")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer conn.Close()

	fmt.Println("Connected to", conn.RemoteAddr())

	// Send a message
	response := []byte("Hello from the client!")
	_, err = conn.Write(response)
	if err != nil {
		fmt.Println("Error writing:", err)
	}
}
