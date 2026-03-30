package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

func clearTerminal() {
	var cmd *exec.Cmd

	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/c", "cls")
	} else {
		cmd = exec.Command("clear")
	}

	cmd.Stdout = os.Stdout
	cmd.Run()
}

type Request struct {
	ID     string `json:"id"`
	Action string `json:"action"`
}
type ResponseSensor struct {
	Status string `json:"status"`
	Data   Sensor `json:"data"`
	Error  string `json:"error"`
}

type Sensor struct {
	ID          string `json:"id"`
	Temperature *int   `json:"temperature"`
	Luminosity  *int   `json:"luminosity"`
	Humidity    *int   `json:"humidity"`
}

func main() {
	clearTerminal()
	conn, err := net.Dial("tcp", "localhost:8000")
	if err != nil {
		fmt.Println("Erro ao conectar no servidor: ", err)
		return
	}
	defer conn.Close()

	input := bufio.NewReader(os.Stdin)

	for {
		clearTerminal()
		fmt.Println("\n|------------------------------|")
		fmt.Println("|             MENU             |")
		fmt.Println("|------------------------------|")
		fmt.Println("|           SENSORES           |")
		fmt.Println("|                              |")
		fmt.Println("| [ 1 ] - Listar sensores      |")
		fmt.Println("| [ 2 ] - Verificar sensores   |")
		fmt.Println("| [ 3 ] - Selecionar sensor    |")
		fmt.Println("|                              |")
		fmt.Println("|------------------------------|")
		fmt.Println("|           ATUADORES          |")
		fmt.Println("|                              |")
		fmt.Println("| [ 4 ] - Listar atuadores     |")
		fmt.Println("| [ 5 ] - Verificar atuadores  |")
		fmt.Println("| [ 6 ] - Selecionar atuador   |")
		fmt.Println("|                              |")
		fmt.Println("|------------------------------|")
		fmt.Println("| [ 7 ] - Fechar               |")
		fmt.Println("|------------------------------|")

		fmt.Print("\nSelecione uma opção: ")

		option, _ := input.ReadString('\n')
		option = strings.TrimSpace(option)

		var request Request

		switch option {
		case "1":
			request = Request{
				Action: "list",
			}

			if err := json.NewEncoder(conn).Encode(request); err != nil {
				fmt.Println("\nErro ao enviar requisição para o servidor: ", err)
				continue
			}

			decoder := json.NewDecoder(conn)

			for {
				var responseSensor ResponseSensor

				if err := decoder.Decode(&responseSensor); err != nil {
					fmt.Println("Erro ao receber resposta do servidor: ", err)
					break
				}

				if responseSensor.Status == "end" {
					break
				}

				if responseSensor.Status != "success" {
					fmt.Println("Erro: ", responseSensor.Error)
					break
				}

				sensorResult := responseSensor.Data

				if sensorResult.Temperature != nil {
					fmt.Printf("\nTemperatura (%s)", sensorResult.ID)
				}
				if sensorResult.Humidity != nil {
					fmt.Printf("\nUmidade (%s)", sensorResult.ID)
				}
				if sensorResult.Luminosity != nil {
					fmt.Printf("\nLuminosidade (%s)", sensorResult.ID)
				}
			}

			fmt.Print("\n\nPressione ENTER para voltar ao menu.")
			fmt.Scanln()

		case "2":
			request = Request{
				Action: "verify",
			}

			if err := json.NewEncoder(conn).Encode(request); err != nil {
				fmt.Println("\nErro ao enviar requisição para o servidor: ", err)
				continue
			}

			decoder := json.NewDecoder(conn)
			latest := make(map[string]Sensor)

			for {
				var responseSensor ResponseSensor

				if err := decoder.Decode(&responseSensor); err != nil {
					fmt.Println("\nErro ao receber resposta do servidor: ", err)
					break
				}

				if responseSensor.Status == "end" {
					break
				}

				if responseSensor.Status == "endOfRound" {
					clearTerminal()
					fmt.Println("\nSensores: ")

					for _, sensor := range latest {
						if sensor.Temperature != nil {
							fmt.Printf("\nTemperatura (%s) = %d ", sensor.ID, *sensor.Temperature)
						}
						if sensor.Humidity != nil {
							fmt.Printf("\nUmidade (%s) = %d ", sensor.ID, *sensor.Humidity)
						}
						if sensor.Luminosity != nil {
							fmt.Printf("\nLuminosidade (%s) = %d ", sensor.ID, *sensor.Luminosity)
						}
					}

					continue
				}

				if responseSensor.Status == "error" {
					fmt.Println("\nFalha: ", responseSensor.Error)
					break
				}

				if responseSensor.Status == "success" {
					sensorResult := responseSensor.Data
					latest[sensorResult.ID] = sensorResult
				}
			}

		case "3":
			fmt.Print("\nDigite o ID do sensor: ")
			id, _ := input.ReadString('\n')
			id = strings.TrimSpace(id)

			request = Request{
				ID:     id,
				Action: "select",
			}

			if err := json.NewEncoder(conn).Encode(request); err != nil {
				fmt.Println("\nErro ao enviar requisição para o servidor: ", err)
				return
			}

			decoder := json.NewDecoder(conn)

			for {
				var responseSensor ResponseSensor

				if err := decoder.Decode(&responseSensor); err != nil {
					fmt.Println("\nErro ao receber resposta do servidor: ", err)
					return
				}

				if responseSensor.Status == "end" {
					break
				}

				if responseSensor.Status == "error" {
					fmt.Println("\nFalha: ", responseSensor.Error)
					return
				}

				if responseSensor.Status == "success" {
					clearTerminal()

					sensor := responseSensor.Data

					fmt.Println("\nSensor: ")

					if sensor.Temperature != nil {
						fmt.Printf("\nTemperatura (%s) = %d ", sensor.ID, *sensor.Temperature)
					}
					if sensor.Humidity != nil {
						fmt.Printf("\nUmidade (%s) = %d ", sensor.ID, *sensor.Humidity)
					}
					if sensor.Luminosity != nil {
						fmt.Printf("\nLuminosidade (%s) = %d ", sensor.ID, *sensor.Luminosity)
					}
				}
			}

		case "7":
			conn.Close()
			return

		default:
			fmt.Println("\nOpção inválida.")

			fmt.Print("\nPressione ENTER para voltar ao menu.")
			fmt.Scanln()
			continue
		}
	}
}
