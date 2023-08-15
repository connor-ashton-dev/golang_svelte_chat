package main

import (
	"fmt"
	"io"
	"sync"
	"time"

	"golang.org/x/net/websocket"
)

type Server struct {
	channels map[string]*Connection
}

type Socket struct {
	id        string
	connected bool
}

func (s *Socket) NewSocket(id string) *Socket {
	return &Socket{
		id:        id,
		connected: true,
	}

}

type Connection struct {
	mut   sync.Mutex
	conns map[*websocket.Conn]*Socket
}

func NewConnection() *Connection {
	return &Connection{
		conns: map[*websocket.Conn]*Socket{},
		mut:   sync.Mutex{},
	}
}

func NewServer() *Server {
	return &Server{
		channels: map[string]*Connection{
			"default": NewConnection(),
			"feed":    NewConnection(),
		},
	}
}

func (s *Server) handleFeed(ws *websocket.Conn) {
	fmt.Println("new incomming commectino to feed")
	s.channels["feed"].mut.Lock()
	s.channels["feed"].conns[ws].connected = true
	s.channels["feed"].mut.Unlock()
	s.readLoop(ws, "feed")
}

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
		s.channels[channel] = NewConnection()
	}
	s.channels[channel].mut.Lock()
	s.channels[channel].conns[ws] = &Socket{
		id:        "foo",
		connected: true,
	}
	s.channels[channel].mut.Unlock()
	s.readLoop(ws, channel)
}

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

func (s *Server) broadcast(b []byte, channel string) {

	for ws := range s.channels[channel].conns {
		go func(ws *websocket.Conn) {
			if _, err := ws.Write(b); err != nil {
				if _, ok := s.channels[channel].conns[ws]; ok {
					fmt.Printf("in channel but not avaliable. deleting...")
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
