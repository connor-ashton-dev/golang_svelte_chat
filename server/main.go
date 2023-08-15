package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"golang.org/x/net/websocket"
)

func main() {
	server := NewServer()
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	http.Handle("/ws/", websocket.Handler(server.HandleWS))

	go SpamFeed(server)

	fmt.Printf("starting server on port %s ... \n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
