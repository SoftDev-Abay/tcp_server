package main

import (
	"bufio"
	"fmt"
	"net"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:8081")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close()

	msg := "Hello from client\n"
	_, err = conn.Write([]byte(msg))
	if err != nil {
		fmt.Println(err)
		return
	}

	response, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Received:", response)
}
