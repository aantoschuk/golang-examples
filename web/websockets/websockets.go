package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize: 1024,
	WriteBufferSize: 1024,
	CheckOrigin: func (r *http.Request) bool {
		origin := r.Header.Get("Origin")
		return origin  ==  "http://localhost:8000"
	},
	// Wait for the WS handshake for no more then 10secs
	// If the Handshake takes longer, abort the connection attemp
	HandshakeTimeout: time.Duration(10 * time.Second),
	EnableCompression: true,
	// Error  - I can define a custom error here.
}

func main () {
	http.HandleFunc("/echo", func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Fatal("Unable to establish connection", err)
		}

		// Do not forget to close connection 
		defer conn.Close()

		// It means after 1 minute from the start connection will be corrupt
		conn.SetReadDeadline(time.Now().Add(1 * time.Minute))

		inactive := make(chan bool)
		go func () {
			<-time.After(1 * time.Minute)
			inactive <-  true
		}()

		for  {
			select {
			case <-inactive:
				log.Println("Connection closed du to inactivity")
				conn.Close()
				return
			default:
				msgType, msg, err := conn.ReadMessage()
				if err != nil {
					log.Print("Unable to read the message", err)
					return
				}
				fmt.Printf("%s sent: %s\n", conn.RemoteAddr(), string(msg))

				if err = conn.WriteMessage(msgType, msg); err != nil {
					log.Print("Unable to send the message", err)
					return
				}
			}
		}
	})

	http.HandleFunc("/", func (w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "websockets.html")
	})

	http.ListenAndServe(":8000", nil)
}

