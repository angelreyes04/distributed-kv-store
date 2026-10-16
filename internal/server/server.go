package main

import (
	"fmt"
	"log"
	"net"
	"os"
)

const (
	CONN_TYPE   = "tcp"
	CONN_HOST   = "localhost"
	CONN_PORT   = ":8080"
	BUFFER_SIZE = 8
)

func main() {
	listener, err := net.Listen(CONN_TYPE, CONN_HOST+CONN_PORT)
	if err != nil {
		log.Printf("Error establishing TCP connection: %s", err.Error())
		os.Exit(1)
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Error accepting TCP connection: %s", err.Error())
			os.Exit(1)
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	buf := make([]byte, BUFFER_SIZE)
	for {
		n, err := conn.Read(buf)
		if err != nil {
			log.Printf("Error while reading from connection: %s", err.Error())
			os.Exit(1)
		}
		fmt.Printf("Message: %s Bytes read: %d", string(buf), n)
	}
}
