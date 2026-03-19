package main

import (
	"fmt"
	"net"
	"time"
)

func HandleConnection(conn net.Conn, sensors map[string]float64) {
	defer conn.Close()

	for {
		message := fmt.Sprintf("Temperatura: %.2f | Luminosidade: %.2f\n",
			sensors["temperature"],
			sensors["luminosity"],
		)

		_, err := conn.Write([]byte(message))
		if err != nil {
			return
		}

		time.Sleep(1 * time.Second)
	}
}
