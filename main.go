package main

import (
	"log"
	"net/http"
)

func main() {

	s := newServer()
	go s.run()

	http.HandleFunc("/ws", s.handleConnection)
	log.Printf("Server starting at localhost:8888")

	if err := http.ListenAndServe(":8888", nil); err != nil {
		log.Fatal(err)
	}

}
