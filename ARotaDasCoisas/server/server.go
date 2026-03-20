package main

import (
	"fmt"
	"net"
)

func handler(conn net.Conn) {
	defer conn.Close()


}

func connectionSensor(listenerSensor net.Listen) {
	for {
		connSensor, err := listenerSensor.Accept()
		if err != nil {
			continue
		}

		go handle
	}
}

func connectionClient(listener *net.Listen) {
	for {
		connClient, err := listenerClient.Accept()
		if err != nil {
			continue
		}

		fmt.Println("Cliente conectado.")
		go HandleConnection(connClient, sensors)
	}
}

func main() {
	fmt.Println("Servidor inicializado.")

	connSensor, err := net.ListenPacket("udp", "localhost:5050")
	if err != nil {
		panic(err)
	}
	defer connSensor.Close()

	for {
		sensors := connSensor.ReadFrom()
	}

	listenerClient, err := net.Listen("tcp", ":5555")
	if err != nil {
		panic(err)
	}
	defer listenerClient.Close()

	go connectionClient()

	
}
