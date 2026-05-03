package server

import (
	"fmt"
	"io"
	"net"

	"github.com/faiyaz032/NotRedis/internal/protocol"
	"github.com/faiyaz032/NotRedis/internal/store"
)

type TCPServer struct {
	addr  string
	store *store.Store
}

func NewTCPServer(addr string, st *store.Store) *TCPServer {
	return &TCPServer{
		addr:  addr,
		store: st,
	}
}

func (s *TCPServer) Run() error {
	listener, err := net.Listen("tcp", s.addr)
	if err != nil {
		return fmt.Errorf("listen on %s: %w", s.addr, err)
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			return fmt.Errorf("accept connection: %w", err)
		}
		go s.handleConnection(conn)
	}
}

func (s *TCPServer) handleConnection(conn net.Conn) {
	defer conn.Close()

	reader := protocol.NewRESPReader(conn)

	for {
		value, err := reader.ReadValue()
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Println("Error reading from client: ", err.Error())
			break
		}

		if value.Type != protocol.Array {
			fmt.Println("Invalid request: expected array")
			continue
		}

		response := protocol.Execute(value.Array, s.store)
		if _, err := conn.Write(response.Marshal()); err != nil {
			fmt.Println("Error writing to client: ", err.Error())
			break
		}
	}
}
