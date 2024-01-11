package server

import (
	"fmt"
	"net"
	"strings"

	"github.com/anthdm/hollywood/actor"
)

type Session struct {
	conn    net.Conn
	roomPID *actor.PID
}

func NewSession(conn net.Conn) actor.Producer {
	return func() actor.Receiver {
		return &Session{
			conn:    conn,
			roomPID: nil,
		}
	}
}

func (s *Session) Receive(c *actor.Context) {
	switch msg := c.Message().(type) {
	case actor.Started:
		go s.readLoop(c)
	case RoomJoin:
		s.roomPID = msg.roomPID
		c.Send(msg.roomPID, msg)
	}
}

func (s *Session) readLoop(c *actor.Context) {
	buf := make([]byte, 1024)

	for {
		n, err := s.conn.Read(buf)
		if err != nil {
			break
		}

		msg := make([]byte, n)
		copy(msg, buf[:n])

		fmt.Println("Received message:", string(msg))
        
		if strings.TrimSpace(string(msg)) == "join" {
			fmt.Println("Joining room")
			c.Send(c.Parent(), RoomJoin{sessionPID: c.PID()})
		}
	}
}
