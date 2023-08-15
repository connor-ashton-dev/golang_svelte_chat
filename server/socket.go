package main

import (
	"fmt"
	"io"
	"strings"
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
		stringMsg := string(msg)
		parts := strings.Split(stringMsg, "::")

		if len(parts) < 3 {
			fmt.Println("error: invalid message format")
			continue
		}

		if parts[0] == "message" { // MESSAGE
			s.cleanInactiveSockets(channel)
			s.broadcast([]byte(msg), channel)
			s.sendUsers(ws, channel)
		}

		if parts[0] == "join" { // JOIN
			username := parts[1] + "::" + parts[2]
			s.cleanInactiveSockets(channel)
			s.addUser(ws, channel, username)
			s.sendUsers(ws, channel)
		}
	}
}

func (s *Server) addUser(ws *websocket.Conn, channel, username string) {
	// add user to channel
	s.channels[channel].mut.Lock()
	s.channels[channel].conns[ws] = &Socket{
		id: "foo",
	}
	s.channels[channel].mut.Unlock()

	// add user to users
	s.channels[channel].mut.Lock()
	s.channels[channel].users[ws] = username
	s.channels[channel].mut.Unlock()

	fmt.Println(s.channels[channel].users)
}

func (s *Server) sendUsers(ws *websocket.Conn, channel string) {
	//loop thru users
	usersString := ""
	for _, user := range s.channels[channel].users {
		username := strings.Split(user, "::")[0]
		usersString += "::" + username

	}
	s.broadcast([]byte("users"+usersString), channel)
}

// loop thru all sockets in channel and clean inactive ones
func (s *Server) cleanInactiveSockets(channel string) {
	for ws := range s.channels[channel].conns {
		if _, err := ws.Write([]byte("")); err != nil {
			if _, ok := s.channels[channel].conns[ws]; ok {
				fmt.Println("in channel but not avaliable. deleting...")
				s.channels[channel].mut.Lock()
				delete(s.channels[channel].conns, ws)
				delete(s.channels[channel].users, ws)
				s.channels[channel].mut.Unlock()
			} else {
				fmt.Println("error at cleanInactiveSockets: ", err)
			}
		}
	}
}

// broadcast sends a message to all clients in a channel
func (s *Server) broadcast(b []byte, channel string) {
	for ws := range s.channels[channel].conns {
		go func(ws *websocket.Conn) {
			if _, err := ws.Write(b); err != nil {
				// If there's an error while writing, delete the socket from the map
				s.channels[channel].mut.Lock()
				delete(s.channels[channel].conns, ws)
				delete(s.channels[channel].users, ws)
				s.channels[channel].mut.Unlock()
			}
		}(ws)
	}
}

// ------------------  CHANNEL  ------------------

type Channel struct {
	mut   sync.Mutex
	conns map[*websocket.Conn]*Socket
	users map[*websocket.Conn]string
}

func newChannel() *Channel {
	return &Channel{
		conns: map[*websocket.Conn]*Socket{},
		mut:   sync.Mutex{},
		users: map[*websocket.Conn]string{},
	}
}

// ------------------  SOCKET  ------------------

type Socket struct {
	id string
}

func (s *Socket) newSocket(id string) *Socket {
	return &Socket{
		id: id,
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
		id: "foo",
	}
	s.channels[channel].mut.Unlock()
	s.readLoop(ws, channel)
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
