package server

import (
	"fmt"
	"math"
	"time"

	"github.com/anthdm/hollywood/actor"
)

type Room struct {
	name    string
	players map[*actor.PID]*Player
}

type PlayerState struct {
	speed float32
	dx    float32
	dy    float32
	x     float32
	y     float32
	id    float32
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
		go r.roomGameLoop(c)
		go r.sendState(c)

	case RoomJoin:
		fmt.Println("Player joined room", msg.playerID)
		fmt.Println("Player joined room, msg.roomPID", msg.roomPID)
		fmt.Println("Player joined room, c.PID()", c.PID())
		r.players[msg.sessionPID] = &Player{
			id:    msg.playerID,
			pid:   msg.sessionPID,
			name:  fmt.Sprintf("player %d", len(r.players)),
			x:     100,
			y:     100,
			dx:    100,
			dy:    100,
			speed: 35,
		}

	case PlayerMoveRequest:
		fmt.Println("Player moved 3!!!!", msg.id)
		if player, ok := r.players[msg.pid]; ok {
			fmt.Println("Player moved", msg.id)
			player.dx = msg.dx
			player.dy = msg.dy
		}
	}
}

func (r *Room) sendState(c *actor.Context) {
	for {
		for _, player := range r.players {
			for _, p := range r.players {
				if p.id != player.id {
					c.Send(p.pid, PlayerState{
						speed: float32(player.speed),
						dx:    float32(player.dx),
						dy:    float32(player.dy),
						x:     float32(player.x),
						y:     float32(player.y),
						id:    float32(player.id),
					})
				}
			}
		}

		time.Sleep(time.Second / 60)
	}
}

func (r *Room) roomGameLoop(c *actor.Context) {
	// game loop with delta time
	var deltaTime float64
	for {
		start := time.Now()

		// update game state
		for _, player := range r.players {
			spd := player.speed * deltaTime
			xDiff := player.dx - player.x
			yDiff := player.dy - player.y
			angle := math.Atan2(yDiff, xDiff)

			if spd*spd >= xDiff*xDiff+yDiff*yDiff {
				player.x = player.dx
				player.y = player.dy

			} else {
				player.x += spd * math.Cos(angle)
				player.y += spd * math.Sin(angle)
			}
		}

		// sleep until next tick
		elapsed := time.Since(start)
		time.Sleep(time.Second/60 - elapsed)

		// calculate delta time
		deltaTime = time.Since(start).Seconds()

		// reset delta time if it's too high
		if deltaTime > 1 {
			deltaTime = 0
		}

		// reset delta time if it's too low
		if deltaTime < 0 {
			deltaTime = 0
		}

		// fmt.Println("deltaTime:", deltaTime)
		// fmt.Println("FPS:", 1/deltaTime)
	}
}
