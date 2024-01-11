package server

import (
	"fmt"

	"github.com/anthdm/hollywood/actor"
)

type Room struct {
	name    string
	players map[*actor.PID]*Player
}

func NewRoom(name string, maxPlayers int) actor.Producer {
	return func() actor.Receiver {
		return &Room{
			name:    name,
			players: make(map[*actor.PID]*Player),
		}
	}
}

func (r *Room) Receive(c *actor.Context) {
	switch msg := c.Message().(type) {
	case actor.Started:
		fmt.Printf("Room %s started\n", r.name)

	case RoomJoin:
		fmt.Println("Player joined room")
		r.players[msg.sessionPID] = &Player{}
	}
}
