package main

import (
	"encoding/json"
	"fmt"
	"maps"
	"net"
	"os"
	"sync"
	"time"
)

type Response struct {
	Status string `json:"status"`
	Data   Sensor `json:"data"`
	Error  string `json:"error"`
}

type Sensor struct {
	ID    string `json:"id"`
	Type  string `json:"type"`
	Value int    `json:"value"`
}

type Request struct {
	ID     string `json:"id"`
	Action string `json:"action"`
}

type Actuator struct {
	ID   string `json:"id"`
	On   bool   `json:"on"`
	Type string `json:"type"`
}

type ActuatorConn struct {
	Actuator Actuator `json:"actuator"`
	Conn     net.Conn `json:"conn"`
}

var (
	sensors    = make(map[string]Sensor)
	actuators  = make(map[string]ActuatorConn)
	muSensor   sync.Mutex
	muActuator sync.Mutex
)

func receiveRequest(conn net.Conn, request Request) error {
	decoder := json.NewDecoder(conn)

	if err := decoder.Decode(&request); err != nil {
		fmt.Println("\nErro na requisição do cliente: ", err)
		return err
	}
	return nil
}

func sendResponse(conn net.Conn, response Response) error {
	encoder := json.NewEncoder(conn)

	if err := encoder.Encode(response); err != nil {
		fmt.Println("\nErro ao enviar resposta: ", err)
		return err
	}
	return nil
}

func checkList(conn net.Conn) bool {
	muSensor.Lock()
	copySensors := maps.Clone(sensors)
	muSensor.Unlock()

	if len(copySensors) == 0 {
		response := Response{
			Status: "error",
			Error:  "Lista de sensores vazia",
		}

		_ = sendResponse(conn, response)
		return false
	}
	return true
}

func listActuators(conn net.Conn) {

}

func verifyActuators(conn net.Conn) {

}

func selectActuator() {

}

func sensorClientRequest(conn net.Conn, request Request) {
	if !checkList(conn) {
		return
	}

	switch request.Action {
	case "listSensor":
		muSensor.Lock()
		copySensors := maps.Clone(sensors)
		muSensor.Unlock()

		for _, sensor := range copySensors {
			response := Response{
				Status: "success",
				Data:   sensor,
			}

			if sendResponse(conn, response) != nil {
				return
			}
		}
		response := Response{
			Status: "end",
		}
		_ = sendResponse(conn, response)

	case "verifySensor", "selectSensor":
		start := time.Now()

		for {
			if time.Since(start) >= 10*time.Second {
				response := Response{
					Status: "end",
				}
				_ = sendResponse(conn, response)
				return
			}

			muSensor.Lock()
			copySensors := maps.Clone(sensors)
			muSensor.Unlock()

			if request.Action == "verifySensor" {
				for _, sensor := range copySensors {
					response := Response{
						Status: "success",
						Data:   sensor,
					}
					if sendResponse(conn, response) != nil {
						return
					}
				}

				response := Response{
					Status: "endOfRound",
				}

				if sendResponse(conn, response) != nil {
					return
				}
			} else if request.Action == "selectSensor" {
				sensor, ok := copySensors[request.ID]
				if !ok {
					response := Response{
						Status: "error",
						Error:  "Sensor não encontrado",
					}
					_ = sendResponse(conn, response)
					return
				}

				response := Response{
					Status: "success",
					Data:   sensor,
				}

				if sendResponse(conn, response) != nil {
					return
				}
			}

			time.Sleep(1 * time.Second)
		}
	}
}

func sendActuatorCommand(id, command string) {
	muActuator.Lock()
	actuator, ok := actuators[id]
	muActuator.Unlock()

	if !ok {
		fmt.Println("\nAtuador não encontrado")
		return
	}

	request := Request{
		ID:     id,
		Action: command,
	}

	if err := json.NewEncoder(actuator.Conn).Encode(request); err != nil {
		fmt.Println("\nErro ao enviar comando para atuador: ", err)
		return
	}
}

func actuatorControl() {
	for {
		muSensor.Lock()
		copySensors := maps.Clone(sensors)
		muSensor.Unlock()

		for _, s := range copySensors {
			if len(s.Luminosities) > 0 {
				luminosity := s.Luminosities[len(s.Luminosities)-1]

				if luminosity >= 50 {
					sendActuatorCommand(s.ID, "off")
				} else {
					sendActuatorCommand(s.ID, "on")
				}
			}

		}
		time.Sleep(1 * time.Second)
	}
}

func handleActuator(conn net.Conn) {
	var actuator Actuator
	if err := json.NewDecoder(conn).Decode(&actuator); err != nil {
		fmt.Println("\nErro ao registrar atuador no servidor: ", err)
		conn.Close()
		return
	}

	muActuator.Lock()
	actuators[actuator.ID] = ActuatorConn{
		Actuator: actuator,
		Conn:     conn,
	}
	muActuator.Unlock()

	fmt.Printf("\nAtuador registrado: %s (%s)", actuator.Type, actuator.ID)
}

func handleClient(conn net.Conn) {
	defer conn.Close()

	for {
		var request Request

		if receiveRequest(conn, request) != nil {
			return
		}

		switch request.Action {
		case "listSensors", "verifySensors", "selectSensor":
			sensorClientRequest(conn, request)
		case "listActuators":
			listActuators(conn)
		}
	}
}

func listenSensor() {
	bufferSensors := make([]byte, 1024)

	connSensor, err := net.ListenPacket("udp", "127.0.0.1:7000")
	if err != nil {
		fmt.Println("\nErro ao iniciar servidor UDP:", err)
		return
	}
	defer connSensor.Close()

	for {
		n, _, err := connSensor.ReadFrom(bufferSensors)
		if err != nil {
			fmt.Println("\nErro ao se comunicar com sensor: ", err)
			continue
		}

		var received Sensor
		err = json.Unmarshal(bufferSensors[:n], &received)
		if err != nil {
			fmt.Println("\nErro ao descompactar sensor: ", err)
			continue
		}

		muSensor.Lock()
		sensors[received.ID] = received
		muSensor.Unlock()
	}
}

func listenActuator() {
	listenerActuator, err := net.Listen("tcp", ":9000")
	if err != nil {
		panic(err)
	}
	defer listenerActuator.Close()

	for {
		connActuator, err := listenerActuator.Accept()
		if err != nil {
			fmt.Println("\nErro na conexão com atuador: ", err)
			continue
		}

		fmt.Println("\nAtuador conectado.")

		go handleActuator(connActuator)
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
			fmt.Println("\nErro na conexão com o cliente: ", err)
			continue
		}

		fmt.Println("\nCliente conectado.")
		go handleClient(connClient)
	}
}

func saveFile() {
	for {
		muSensor.Lock()
		copySensors := maps.Clone(sensors)
		muSensor.Unlock()

		file, err := os.Create("../dataBase.json")
		if err != nil {
			fmt.Println("\nErro ao criar arquivo JSON.")
			return
		}

		encoder := json.NewEncoder(file)
		encoder.SetIndent("", "  ")
		encoder.Encode(copySensors)

		file.Close()

		time.Sleep(5 * time.Second)
	}
}

func main() {
	fmt.Println("\nServidor inicializado.")

	go listenSensor()
	go listenActuator()
	go actuatorControl()
	go saveFile()
	listenClient()
}

//Parei na parte da mudança das três funções pra uma, aparentemente terminei, tem que testar.
//Depois é pra mudar no cliente essa lógica. Fazer a lógica do cliente com atuador.
//Consertar o controle do atuador pelo servidor.
//Verificar se falta algum consertar algo na mudança da struct que eu fiz.
