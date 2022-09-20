package main

import (
	"flag"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"html/template"
	"log"
	"net/http"
)

var (
	upgrader = websocket.Upgrader{CheckOrigin: func(r *http.Request) bool {
		return true
	}}

	listenAddr string
	wsAddr     string
	jsTemplate *template.Template
)

func init() {

	fmt.Println("init")
	flag.StringVar(&listenAddr, "listen-addr", "", "Address to listen on")
	flag.StringVar(&wsAddr, "ws-addr", "", "Address for websocket connection")
	flag.Parse()

	var err error
	jsTemplate, err = template.ParseFiles("logger.js")
	if err != nil {
		panic(err)
	}
}
func main() {

	r := mux.NewRouter()

	r.HandleFunc("/ws", serverWS)
	r.HandleFunc("/k.js", serverFile)

	log.Println("starting server")
	log.Fatal(http.ListenAndServe(":8080", r))

}

func serverWS(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		http.Error(w, "Error", 500)
		return
	}

	defer conn.Close()

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			return
		}

		fmt.Printf("From %s: %s\n", conn.RemoteAddr().String(), string(msg))
	}
}

func serverFile(w http.ResponseWriter, r *http.Request) {
	log.Println("file")

	w.Header().Set("Content-Type", "application/javascript")

	jsTemplate.Execute(w, wsAddr)

}
