package main

import (
	"io"
	"log"
	"net"
)

func main() {

	port := "8888"
	listener, err := net.Listen("tcp", ":"+port)

	if err != nil {
		log.Fatalf("Unable to bind port: %s", port)
	}

	log.Println("Listening on 0.0.0.0:", port)

	for {

		conn, err := listener.Accept()

		if err != nil {
			log.Fatalln("Unable to accept connection")
		}

		go echo(conn)

	}

}

// Read data from conn and write it back
func echo(conn net.Conn) {

	defer conn.Close()

	for {

		if _, err := io.Copy(conn, conn); err != nil {
			log.Fatalln("Unable to read/write data")

		}

	}

}
