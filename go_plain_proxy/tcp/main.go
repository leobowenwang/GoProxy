package main

import (
	"fmt"
	"io"
	"log"
	"net"
)

func handleProxy(conn net.Conn, remoteHost string, remotePort int) {
	defer conn.Close()

	remoteConn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", remoteHost, remotePort))
	if err != nil {
		log.Printf("Error connecting to %s on port %d: %v", remoteHost, remotePort, err)
		return
	}
	defer remoteConn.Close()

	go func() {
		_, err := io.Copy(remoteConn, conn)
		if err != nil {
			log.Printf("Error copying from local to remote: %v", err)
		}
	}()

	_, err = io.Copy(conn, remoteConn)
	if err != nil {
		log.Printf("Error copying from remote to local: %v", err)
	}
}

func main() {
	localPort := 3333
	remoteHost := "localhost"
	remotePort := 8080

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", localPort))
	if err != nil {
		log.Fatalf("Error listening on port %d: %v", localPort, err)
	}
	defer listener.Close()

	log.Printf("Listening on port %d", localPort)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Error accepting connection: %v", err)
			continue
		}

		go handleProxy(conn, remoteHost, remotePort)
	}
}
