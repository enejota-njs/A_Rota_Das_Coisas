package main

import (
	"fmt"
	"math/rand"
	"net"
	"time"
)

func handleConnection(conn net.Conn, sensors map[string]float64) {
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

func valueGenerator(sensors map[string]float64) {
	for i := 0; i < 15; i++ {
		generatedTemperature := rand.Float64()
		generatedLuminosity := rand.Float64()

		temperature := generatedTemperature*30 + 10
		luminosity := generatedLuminosity

		sensors["temperature"] = temperature
		sensors["luminosity"] = luminosity

		time.Sleep(1 * time.Second)
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())

	listener, err := net.Listen("tcp", "26.200.141.87:5555")
	if err != nil {
		panic(err)
	}
	defer listener.Close()

	fmt.Println("Servidor inicializado.")

	sensors := make(map[string]float64)
	go valueGenerator(sensors)

	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}

		fmt.Println("Cliente conectado.")
		go handleConnection(conn, sensors)
	}

}
