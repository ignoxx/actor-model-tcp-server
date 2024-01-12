package server

import (
	"net"

	"github.com/anthdm/hollywood/actor"
)

type ConnAdd struct {
	pid  *actor.PID
	conn net.Conn
}

type ConnRemove struct {
	pid *actor.PID
}

type RoomAdd struct {
	pid        *actor.PID
	name       string
	maxPlayers int
}

type RoomRemove struct {
	pid *actor.PID
}

type RoomJoin struct {
	playerID   float64
	roomPID    *actor.PID
	sessionPID *actor.PID
}

type Player struct {
	id    float64
	pid   *actor.PID
	name  string
	x     float64
	y     float64
	dx    float64
	dy    float64
	speed float64
}
