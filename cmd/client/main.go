package main

import (
	"fmt"
	"net"
	"strings"
)

func main() {
	for i := 0; i < 10; i++ {
		// Establish a TCP connection to a server
		conn, err := net.Dial("tcp", "127.0.0.1:8080")
		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		fmt.Println("Connected to", conn.RemoteAddr(), "from", conn.LocalAddr())

		builder := strings.Builder{}
		for i := 0; i < 3; i++ {
			builder.WriteString(fmt.Sprintf("Hello from the client #%d!", i))
		}
		builder.WriteString("\n")

		// Send a message
		message := []byte(builder.String())
		_, err = conn.Write(message)
		if err != nil {
			fmt.Println("Error writing:", err)
		}
	}
}
