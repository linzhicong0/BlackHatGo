package main

import (
	"io"
	"log"
	"net"
)

func main(){


	port := "8888"
	listener, err := net.Listen("tcp", ":" + port)


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

	buff := make([]byte, 512)

	for {
		size, err := conn.Read(buff)

		if err == io.EOF{
			log.Println("Client disconnected")
			break
		}

		if err != nil {
			log.Println("Unexpected error")
			break
		}

		log.Printf("Received %d bytes: %s\n", size, string(buff))

		log.Println("Writing data")

		if _, err := conn.Write(buff[0:size]); err != nil {
			log.Fatalln("Unable to write data")
		}
	}

}