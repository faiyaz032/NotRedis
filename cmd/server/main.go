package main

import (
	"log"

	"github.com/faiyaz032/NotRedis/internal/server"
	"github.com/faiyaz032/NotRedis/internal/store"
	"github.com/faiyaz032/NotRedis/internal/wal"
)

func main() {
	w, err := wal.Open("notredis.wal")
	if err != nil {
		log.Fatal("failed to open WAL:", err)
	}
	defer w.Close()

	st := store.NewStore()

	log.Println("Replaying WAL...")
	if err := w.Replay(func(entry wal.Entry) {
		switch entry.Command {
		case "SET":
			st.Apply(entry.Key, entry.Value)
		}
	}); err != nil {
		log.Fatal("WAL replay failed:", err)
	}

	tcpServer := server.NewTCPServer(":6370", st, w)
	log.Println("Store running on :6370...")
	if err := tcpServer.Run(); err != nil {
		log.Fatal(err)
	}
}
