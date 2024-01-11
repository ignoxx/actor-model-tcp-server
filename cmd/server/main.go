package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/anthdm/hollywood/actor"
	"github.com/ignoxx/game-server-actors/actors/server"
)

func main() {
	engine, err := actor.NewEngine(nil)
	if err != nil {
		panic(err)
	}

	serverPID := engine.Spawn(server.NewServer(":8080"), "server")

	sigchan := make(chan os.Signal)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	<-sigchan

	engine.Poison(serverPID).Wait()
}
