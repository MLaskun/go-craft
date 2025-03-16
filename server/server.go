package server

import (
	"fmt"
	"net"

	"github.com/MLaskun/go-craft/connection"
)

func NewServer() {
	l, err := net.Listen("tcp", "localhost:25565")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer l.Close()

	fmt.Println("Server listening on port 25565")

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error:", err)
			continue
		}
		go connection.HandleClient(conn)
	}
}
