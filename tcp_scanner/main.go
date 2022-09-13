package main

import (
	"fmt"
	"net"
)

//var host = "scanme.nmap.org"
var host = "127.0.0.1"

func main() {

	portCh := make(chan int, 100)
	resultCh := make(chan int)

	// start workers
	fmt.Println("starting workers...")
	for i := 0; i < cap(portCh); i++ {
		go worker(portCh, resultCh)
	}

	go func() {
		fmt.Println("putting ports into channel...")
		for i := 0; i < 1024; i++ {
			portCh <- i
		}
	}()


	results := make([]int, 0)

	fmt.Println("reading results...")
	for i := 0; i < 1024; i++ {
		result := <- resultCh

		if result != 0{
			results = append(results, result)
		}
	}

	fmt.Println("finished")
	close(portCh)
	close(resultCh)

	for _, p := range results {
		fmt.Printf("Connection successful: %d\n", p)
	}

}

func worker(ports, results chan int) {

	for p := range ports {
		address := fmt.Sprintf("%s:%d", host, p)
		conn, err := net.Dial("tcp", address)
		if err == nil {
			conn.Close()
			results <- p
		} else {
			results <- 0
		}
	}

}
