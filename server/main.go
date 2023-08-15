package main

import (
	"fmt"
	"log"
	"net/http"

	"golang.org/x/net/websocket"
)

func main() {
	server := NewServer()
	http.Handle("/ws/", websocket.Handler(server.HandleWS))
	go SpamFeed(server)
	fmt.Println("starting server on port 3000")
	log.Fatal(http.ListenAndServe(":3000", nil))
}
