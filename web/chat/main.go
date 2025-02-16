package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"example/chat/api"
	"example/chat/websocket"

	_ "github.com/lib/pq"
)

var db *sql.DB

func connect() {
	fmt.Println("Connecting to db")
	var err error
	conStr := "postgres://alexander:123456@localhost/chat?sslmode=disable"
	db, err = sql.Open("postgres", conStr)
	if err != nil {
		panic(err)
	}

	defer db.Close()

	if err = db.Ping(); err != nil {
		panic(err)
	}

	fmt.Println("db connected")
}

func main() {
	connect()
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

// FEATURE: add registration
// FEATURE: Each message should be associated with the user who sent it
