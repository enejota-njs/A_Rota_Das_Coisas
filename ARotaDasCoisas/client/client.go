package main

import (
	"bufio"
	"fmt"
	"net"
)

func main() {
	conn, err := net.Dial("tcp", ":5555")

	if err != nil {
		panic(err)
	}
	defer conn.Close()

	fmt.Println("Conectado.")

	reader := bufio.NewReader(conn)

	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			return
		}

		fmt.Print(message)
	}
}
