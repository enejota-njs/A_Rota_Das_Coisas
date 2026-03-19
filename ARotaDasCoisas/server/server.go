package main

import (
	"ARotaDasCoisas/sensor"
	"fmt"
	"net"
)

func main() {
	listener, err := net.Listen("tcp", ":5555")
	if err != nil {
		panic(err)
	}
	defer listener.Close()

	fmt.Println("Servidor inicializado.")

	sensors := make(map[string]float64)
	go sensor.ValueGenerator(sensors)

	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}

		fmt.Println("Cliente conectado.")
		go HandleConnection(conn, sensors)
	}
}
