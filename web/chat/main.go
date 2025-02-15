package main

import (
	"log"
	"net/http"
)

func isPreflight(r *http.Request) bool {
	return r.Method == "OPTIONS" &&
		r.Header.Get("Origin") != "" &&
		r.Header.Get("Access-Control-Request-Method") != ""
}

func main() {
	hub := newHub()

	go hub.run()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		EnableCors(&w, r)
		w.Write([]byte("Server Works"))
	})

	http.HandleFunc("/chat", func(w http.ResponseWriter, r *http.Request) {
		serveWS(hub, w, r)
	})

	err := http.ListenAndServe(":8000", nil)
	log.Fatal("Listen and Serve: ", err)
}

// TODO: Improve porject structure
// FEATURE: add registration
// FEATURE: Each message should be associated with the user who sent it
