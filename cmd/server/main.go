package main

import (
	"log"

	"github.com/faiyaz032/NotRedis/internal/server"
	"github.com/faiyaz032/NotRedis/internal/store"
)

func main() {
	st := store.NewStore()
	tcpServer := server.NewTCPServer(":6370", st)

	log.Println("Store running on :6370...")
	if err := tcpServer.Run(); err != nil {
		log.Fatal(err)
	}
}
