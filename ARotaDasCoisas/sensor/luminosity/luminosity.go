package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"math/rand"
	"net"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"
)

type Sensor struct {
	ID    string `json:"id"`
	Type  string `json:"type"`
	Value int    `json:"value"`
}

func step(value int) int {
	r := rand.Float64()
	if r > 0.5 {
		value += 1
	} else {
		value -= 1
	}

	if value > 100 {
		value = 100
	}
	if value < 0 {
		value = 0
	}
	return value
}

func main() {
	clearTerminal()
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("\nDigite o ID do sensor luminosidade: ")
	id, _ := reader.ReadString('\n')
	id = strings.TrimSpace(id)

	lumi := rand.Intn(101)

	clearTerminal()
	fmt.Printf("\nSensor de luminosidade %s inicializado.\n", id)

	for {
		conn, err := net.Dial("udp", "127.0.0.1:7000")
		if err != nil {
			fmt.Println("\nErro ao conectar o sensor de luminosidade: ", id, err)
			continue
		}

		for {
			lumi = step(lumi)

			data := Sensor{
				ID:    id,
				Type:  "Luminosidade",
				Value: lumi,
			}

			values, _ := json.Marshal(data)

			_, err := conn.Write(values)
			if err != nil {
				fmt.Println("\nErro no envio do sensor de luminosidade: ", id, err)
				conn.Close()
				break
			}

			time.Sleep(1 * time.Millisecond)
		}
	}
}

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
