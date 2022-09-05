package _tcp

import (
	"bufio"
	"fmt"
	"net"
)

func process(conn net.Conn) {
	defer conn.Close()

	for {
		reader := bufio.NewReader(conn)

		var buf [128]byte
		n, err := reader.Read(buf[:])
		if err != nil {
			fmt.Println("read data from client error[]", err)
			break
		}
		rcvStr := string(buf[:n])
		fmt.Println("receive data from client: ", rcvStr)
		conn.Write([]byte("reply same: " + rcvStr))
	}
}

func BuildServer() {
	listener, err := net.Listen("tcp","127.0.0.1:1234")
	if err != nil {
		return
	}
	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		go process(conn)
	}
}