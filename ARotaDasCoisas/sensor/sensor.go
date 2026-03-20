package sensor

import (
	"encoding/json"
	"math/rand"
	"net"
	"time"
)

type Sensor struct {
	Temperature float64 `json:"temperature"`
	Luminosity  float64 `json:"luminosity"`
	Humidity    float64 `json:"humidity"`
}

func main() {
	conn, err := net.Dial("udp", "localhost:5050")
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	for {
		rand.Seed(time.Now().UnixNano())

		data := Sensor{
			Temperature: rand.Float64(),
			Luminosity:  rand.Float64(),
			Humidity:    rand.Float64(),
		}

		values, _ := json.Marshal(data)

		conn.Write(values)
		time.Sleep(1 * time.Second)
	}
}
