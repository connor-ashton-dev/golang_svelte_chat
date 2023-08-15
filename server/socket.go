package main

import (
	"fmt"
	"io"
	"sync"
	"time"

	"golang.org/x/net/websocket"
)

// ------------------  SERVER  ------------------
type Server struct {
	channels map[string]*Channel
}

func NewServer() *Server {
	return &Server{
		channels: map[string]*Channel{
			"default": newChannel(),
			"feed":    newChannel(),
		},
	}
}

// ------------------  CHANNEL  ------------------

type Channel struct {
	mut   sync.Mutex
	conns map[*websocket.Conn]*Socket
}

func newChannel() *Channel {
	return &Channel{
		conns: map[*websocket.Conn]*Socket{},
		mut:   sync.Mutex{},
	}
}

// ------------------  SOCKET  ------------------

type Socket struct {
	id        string
	connected bool
}

func (s *Socket) newSocket(id string) *Socket {
	return &Socket{
		id:        id,
		connected: true,
	}

}

// ------------------  SOCKET HANDLERS  ------------------

// HandleWS handles websocket connections
func (s *Server) HandleWS(ws *websocket.Conn) {
	path := ws.Request().URL.Path //format = /ws/default
	fmt.Println("new incomming connection from client: ", ws.RemoteAddr())
	channel := path[len("/ws/"):]
	if channel == "" {
		channel = "global"
	}
	_, ok := s.channels[channel]
	if !ok {
		fmt.Printf("channel %s does not exist. creating...\n", channel)
		s.channels[channel] = newChannel()
	}
	s.channels[channel].mut.Lock()
	s.channels[channel].conns[ws] = &Socket{
		id:        "foo",
		connected: true,
	}
	s.channels[channel].mut.Unlock()
	s.readLoop(ws, channel)
}

// readLoop reads messages from the client and broadcasts them to all other clients
func (s *Server) readLoop(ws *websocket.Conn, channel string) {
	buf := make([]byte, 1024)
	for {
		n, err := ws.Read(buf)
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Println("error reading from client: ", err)
			continue
		}
		msg := buf[:n]

		s.broadcast(msg, channel)
	}
}

// broadcast sends a message to all clients in a channel
func (s *Server) broadcast(b []byte, channel string) {
	for ws := range s.channels[channel].conns {
		go func(ws *websocket.Conn) {
			if _, err := ws.Write(b); err != nil {
				if _, ok := s.channels[channel].conns[ws]; ok {
					fmt.Println("in channel but not avaliable. deleting...")
					s.channels[channel].mut.Lock()
					delete(s.channels[channel].conns, ws)
					s.channels[channel].mut.Unlock()
				} else {
					fmt.Println("error writing to client: ", err)
				}
			}
		}(ws)
	}
}

// ------------------  FEED  ------------------

// SpamFeed sends a message to the feed channel every 2 seconds
func SpamFeed(s *Server) {
	greetings := []string{
		"hello",
		"hi",
		"how are you",
		"good",
		"bad",
		"ok",
		"bye",
		"goodbye",
	}
	for {
		for _, greeting := range greetings {
			s.broadcast([]byte(greeting), "feed")
			time.Sleep(time.Second * 2)
		}
	}
}
