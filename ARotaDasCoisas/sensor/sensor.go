package sensor

import (
	"math/rand"
	"time"
)

func ValueGenerator(sensors map[string]float64) {
	rand.Seed(time.Now().UnixNano())

	for {
		generatedTemperature := rand.Float64()
		generatedLuminosity := rand.Float64()

		temperature := generatedTemperature*30 + 10
		luminosity := generatedLuminosity

		sensors["temperature"] = temperature
		sensors["luminosity"] = luminosity

		time.Sleep(1 * time.Second)
	}
}
