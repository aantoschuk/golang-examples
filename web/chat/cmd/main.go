package main

import (
	"log"
	"net/http"

	"example/chat/api"
	"example/chat/websocket"
)

func main() {
	hub := websocket.NewHub()

	go hub.Run()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		api.EnableCors(&w, r)
		w.Write([]byte("Server Works"))
	})

	http.HandleFunc("/chat", func(w http.ResponseWriter, r *http.Request) {
		websocket.ServeWS(hub, w, r)
	})

	err := http.ListenAndServe(":8000", nil)
	log.Fatal("Listen and Serve: ", err)
}

// TODO: Improve porject structure
// FEATURE: add registration
// FEATURE: Each message should be associated with the user who sent it
