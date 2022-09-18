package main

import (
	"io"
	"log"
	"net"
)

var target = "xxx:80"

func main(){

	port := ":80"
	listener, err := net.Listen("tcp", port)

	if err != nil {
		log.Fatalf("Unable to bind to port: %s", port[1:])
		return
	}


	log.Println("Listening on: 0.0.0.0", port)

	for {
		conn, err := listener.Accept()

		if err != nil {
			log.Fatalf("Unable to accept connection: %s", err)
		}

		log.Println("New connection")
		go handle(conn)
	}

}


func handle(src net.Conn){

	dst, err := net.Dial("tcp", target)

	if err != nil {
		log.Fatalf("Unable to connect to target: %s\n", target)
		return
	}

	go func() {
		if _, err = io.Copy(dst, src); err != nil{
			log.Fatalln(err)
		}
	}()


	if _, err = io.Copy(src, dst); err != nil{
		log.Fatalln(err)
	}

	dst.Close()


}
