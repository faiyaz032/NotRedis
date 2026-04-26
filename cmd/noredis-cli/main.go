package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:6370")
	if err != nil {
		fmt.Printf("Could not connect to NoRedis at localhost:6370: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close()

	reader := bufio.NewReader(os.Stdin)
	serverReader := bufio.NewReader(conn)

	for {
		fmt.Print("noredis> ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		if input == "exit" || input == "quit" {
			return
		}
		if input == "" {
			continue
		}

		fmt.Fprintf(conn, "%s", input+"\n")

		response, err := serverReader.ReadString('\n')
		if err != nil {
			fmt.Println("Connection lost.")
			return
		}
		fmt.Print(response)
	}
}
