package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:8000")
	if err != nil {
		fmt.Println("Erro ao conectar no servidor: ", err)
		return
	}
	defer conn.Close()

	fmt.Println("Conectado ao servidor.")

	input := bufio.NewReader(os.Stdin)
	reader := bufio.NewReader(conn)

	for {
		fmt.Println("|----------- MENU -----------|")
		fmt.Println("| [ 1 ] - Listar sensores    |")
		fmt.Println("| [ 2 ] - Verificar sensores |")
		fmt.Println("| [ 3 ] - Selecionar sensor  |")
		fmt.Println("|----------------------------|")

		fmt.Print("\nSelecione uma opção: ")

		option, _ := input.ReadString('\n')
		option = strings.TrimSpace(option)

		if option == "1" || option == "2" || option == "3" {
			fmt.Fprintf(conn, "%s\n", option)
		} else {
			fmt.Println("Opção inválida.")
			continue
		}

		for {
			result, err := reader.ReadString('\n')
			if err != nil {
				fmt.Println("Conexão encerrada: ", err)
				return
			}

			if strings.TrimSpace(result) == "ok" {
				break
			}

			fmt.Print(result)
		}
	}
}
