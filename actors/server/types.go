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
	roomPID    *actor.PID
	sessionPID *actor.PID
}

type Player struct {
	x float32
	y float32
}
