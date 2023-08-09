package main

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/andiblas/sockets-exercise/pkg/ratelimiter"
)

func main() {
	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()

	go func() {
		// listen for interrupt signal
		interrupt := make(chan os.Signal, 1)
		signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)
		<-interrupt
		cancelFunc()
	}()

	// Listen for incoming connections on a specific port
	server, err := net.Listen("tcp", "localhost:8080")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer server.Close()

	fmt.Println("Server listening on", server.Addr())

	wg := sync.WaitGroup{}
	rateLim := ratelimiter.NewLeakyBucket(2, 5*time.Second)
	rateLim.Start()

	for {
		if errors.Is(ctx.Err(), context.Canceled) {
			break
		}

		// Accept incoming connection
		conn, err := server.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}

		if !rateLim.Allow(conn) {
			fmt.Println("Rate limit exceeded for", conn.RemoteAddr().String())
			message := []byte("RATE_LIMIT")
			_, err := conn.Write(message)
			if err != nil {
				fmt.Println("Error writing:", err)
			}
			err = conn.Close()
			if err != nil {
				return
			}
			continue
		}

		fmt.Println("Accepted connection from", conn.RemoteAddr())

		wg.Add(1)
		go func() {
			defer wg.Done()
			handleConnection(conn)
		}()
	}

	// wait until all connections finished
	wg.Wait()

}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	err := conn.SetReadDeadline(time.Now().Add(200 * time.Millisecond))
	if err != nil {
		fmt.Println("could not set read deadline to connection")
		return
	}
	connReader := bufio.NewReader(conn)
	receivedMessage, err := connReader.ReadString(byte('\n'))
	if err != nil && !errors.Is(err, io.EOF) {
		fmt.Println("error reading:", err)
		return
	}

	fmt.Printf("Received: %s\n", receivedMessage)
}
