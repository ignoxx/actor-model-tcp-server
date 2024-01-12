package server

import (
	"encoding/binary"
	"fmt"
	"math"
	"net"

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

	case PlayerState:
		buf := make([]byte, 1+4+4+4+4+4+4)

		buf[0] = 1
		binary.LittleEndian.PutUint32(buf[1:], math.Float32bits(msg.id))
		binary.LittleEndian.PutUint32(buf[5:], math.Float32bits(msg.x))
		binary.LittleEndian.PutUint32(buf[9:], math.Float32bits(msg.y))
		binary.LittleEndian.PutUint32(buf[13:], math.Float32bits(msg.dx))
		binary.LittleEndian.PutUint32(buf[17:], math.Float32bits(msg.dy))
		binary.LittleEndian.PutUint32(buf[21:], math.Float32bits(msg.speed))

		s.conn.Write(buf)

	case PlayerMoveRequest:
		c.Send(s.roomPID, msg)

	}

}

type PlayerMoveRequest struct {
	pid *actor.PID
	id  float64
	dx  float64
	dy  float64
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

		if msg[0] == 0 {
			fmt.Println("Joining room")

			playerID := Float32frombytes(msg[1:])

			c.Send(c.Parent(), RoomJoin{sessionPID: c.PID(), playerID: float64(playerID)})
		}

		if msg[0] == 1 {
			id := Float32frombytes(msg[1:])
			dx := Float32frombytes(msg[5:])
			dy := Float32frombytes(msg[9:])

			c.Send(c.PID(), PlayerMoveRequest{pid: c.PID(), id: float64(id), dx: float64(dx), dy: float64(dy)})
		}
	}
}

func Float32frombytes(bytes []byte) float32 {
	bits := binary.LittleEndian.Uint32(bytes)
	float := math.Float32frombits(bits)
	return float
}

func Float32bytes(float float32) []byte {
	bits := math.Float32bits(float)
	bytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(bytes, bits)
	return bytes
}

func Float64frombytes(bytes []byte) float64 {
	bits := binary.LittleEndian.Uint64(bytes)
	float := math.Float64frombits(bits)
	return float
}

func Float64bytes(float float64) []byte {
	bits := math.Float64bits(float)
	bytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(bytes, bits)
	return bytes
}
