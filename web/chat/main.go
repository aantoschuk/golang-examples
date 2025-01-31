package main

import (
	"fmt"
	"log"
	"net/http"
)

// TODO: learn more about CORS and preflight in golang.
func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}

func main () {
	hub := newHub()

	go hub.run()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		enableCors(&w)
		fmt.Println("worked")
		w.Write([]byte("Server Works"))

	})
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWS(hub, w, r)
	})

	err := http.ListenAndServe(":8000", nil)
	log.Fatal("Listen and Serve: ", err)
}
