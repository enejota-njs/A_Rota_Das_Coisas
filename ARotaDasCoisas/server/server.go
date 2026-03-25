package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"maps"
	"net"
	"strings"
	"sync"
	"time"
)

type Sensor struct {
	ID          string `json:"id"`
	Temperature *int   `json:"temperature"`
	Luminosity  *int   `json:"luminosity"`
	Humidity    *int   `json:"humidity"`
}

var (
	sensors = make(map[string]Sensor)
	mu      sync.Mutex
)

func listenSensor() {
	bufferSensors := make([]byte, 1024)

	connSensor, err := net.ListenPacket("udp", "127.0.0.1:7000")
	if err != nil {
		fmt.Println("Erro ao iniciar servidor UDP:", err)
		return
	}
	defer connSensor.Close()

	for {
		n, _, err := connSensor.ReadFrom(bufferSensors)
		if err != nil {
			fmt.Println("Erro no ReadFrom:", err)
			continue
		}

		var received Sensor
		err = json.Unmarshal(bufferSensors[:n], &received)
		if err != nil {
			fmt.Println("Erro no Unmarshal:", err)
			continue
		}

		mu.Lock()
		current := sensors[received.ID]

		if received.Temperature != nil {
			current.Temperature = received.Temperature
		}
		if received.Luminosity != nil {
			current.Luminosity = received.Luminosity
		}
		if received.Humidity != nil {
			current.Humidity = received.Humidity
		}

		current.ID = received.ID
		sensors[received.ID] = current
		mu.Unlock()

		fmt.Println("Recebido:", received.ID)
	}
}

func handleClient(conn net.Conn) {
	defer conn.Close()

	reader := bufio.NewReader(conn)

	for {
		option, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Cliente desconectou.")
			return
		}

		option = strings.TrimSpace(option)

		switch option {
		case "1":
			listSensors(conn)
		case "2":
			verifySensors(conn)
		case "3":
			selectSensor(conn)
		}
	}

}

func listenClient() {
	listenerClient, err := net.Listen("tcp", ":8000")
	if err != nil {
		panic(err)
	}
	defer listenerClient.Close()

	for {
		connClient, err := listenerClient.Accept()
		if err != nil {
			fmt.Println("Erro no Accept:", err)
			continue
		}

		fmt.Println("Cliente conectado.")
		go handleClient(connClient)
	}
}

func selectSensor(conn net.Conn) {
	mu.Lock()
	copySensors := maps.Clone(sensors)
	mu.Unlock()

	//Parei na função de seleção de um sensor específico.

}

func verifySensors(conn net.Conn) {
	for {
		mu.Lock()
		copySensors := maps.Clone(sensors)
		mu.Unlock()

		for id, s := range copySensors {
			values := fmt.Sprintf("%s -> ", id)

			if s.Temperature != nil {
				values += fmt.Sprintf("%d ", *s.Temperature)
			}
			if s.Luminosity != nil {
				values += fmt.Sprintf("%d", *s.Luminosity)
			}
			if s.Humidity != nil {
				values += fmt.Sprintf("%d", *s.Humidity)
			}

			values += "\n"
			conn.Write([]byte(values))
		}

		time.Sleep(1 * time.Second)
	}
}

func listSensors(conn net.Conn) {
	mu.Lock()
	copySensors := maps.Clone(sensors)
	mu.Unlock()

	for id, s := range copySensors {
		if s.Temperature != nil || s.Luminosity != nil || s.Humidity != nil {
			fmt.Fprintf(conn, "%s\n", id)
		}
	}
}

func main() {
	fmt.Println("Servidor inicializado.")

	go listenSensor()
	go listenClient()

	select {}
}
