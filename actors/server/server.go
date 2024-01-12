package server

import (
	"fmt"
	"net"

	"github.com/anthdm/hollywood/actor"
)

type Server struct {
	listenAddr string
	listener   net.Listener
	rooms      map[*actor.PID]*Room
	sessions   map[*actor.PID]net.Conn
}

func NewServer(listenAddr string) actor.Producer {
	return func() actor.Receiver {
		return &Server{
			listenAddr: listenAddr,
			rooms:      make(map[*actor.PID]*Room),
			sessions:   make(map[*actor.PID]net.Conn),
		}
	}
}

func (s *Server) Receive(c *actor.Context) {
	switch msg := c.Message().(type) {
	case actor.Initialized:
		l, err := net.Listen("tcp", s.listenAddr)
		if err != nil {
			panic(err)
		}

		s.listener = l

	case actor.Started:
		// add a room
		c.Send(c.PID(), RoomAdd{name: "room1", maxPlayers: 2})
		go s.acceptLoop(c)

	case ConnAdd:
		fmt.Println("Connection added with PID:", msg.pid)
		s.sessions[msg.pid] = msg.conn

	case ConnRemove:
		if conn, ok := s.sessions[msg.pid]; ok {
			fmt.Println("Connection removed with PID:", msg.pid)
			conn.Close()
			delete(s.sessions, msg.pid)
		}

	case RoomAdd:
		fmt.Println("Room added with name:", msg.name)
		roomPID := c.SpawnChild(NewRoom(msg.name, msg.maxPlayers), "room", actor.WithID(msg.name))
		s.rooms[roomPID] = nil

	case RoomRemove:
		if room, ok := s.rooms[msg.pid]; ok {
			fmt.Println("Room removed with name:", room.name)
			delete(s.rooms, msg.pid)
		}
	case RoomJoin:
		// get first room
		fmt.Println("RoomJoin")
		for pid := range s.rooms {
			fmt.Println("RoomJoin", pid)
			msg.roomPID = pid
			c.Send(c.Sender(), msg)
			break
		}

    case PlayerMoveRequest:
        fmt.Println("PlayerMoveRequest")

	}

}

func (s *Server) acceptLoop(c *actor.Context) {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			panic(err)
		}

		fmt.Println("New connection: ", conn.RemoteAddr())

		pid := c.SpawnChild(NewSession(conn), "session", actor.WithID(conn.RemoteAddr().String()))

		c.Send(c.PID(), ConnAdd{pid: pid, conn: conn})
	}
}
