package server

import (
	"bufio"
	"fmt"
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
	scanner := bufio.NewScanner(conn)

	for scanner.Scan() {
		response := protocol.Execute(scanner.Text(), s.store)
		if _, err := conn.Write([]byte(response)); err != nil {
			return
		}
	}
}
