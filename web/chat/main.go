package main

import (
	"log"
	"net/http"
)

func main () {
	hub := newHub()

	go hub.run()

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWS(hub, w, r)
	})

	err := http.ListenAndServe(":8000", nil)
	log.Fatal("Listen and Serve: ", err)
}
