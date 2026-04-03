GO=go
BIN_EXT=

build: build-server build-client build-sensors build-actuators
	@echo "Build completo"

build-server:
	$(GO) build -o server/server$(BIN_EXT) server/server.go

build-client:
	$(GO) build -o client/client$(BIN_EXT) client/client.go

build-sensors:
	$(GO) build -o sensor/gas/gas$(BIN_EXT) sensor/gas/gas.go
	$(GO) build -o sensor/humidity/humidity$(BIN_EXT) sensor/humidity/humidity.go
	$(GO) build -o sensor/luminosity/luminosity$(BIN_EXT) sensor/luminosity/luminosity.go
	$(GO) build -o sensor/smoke/smoke$(BIN_EXT) sensor/smoke/smoke.go
	$(GO) build -o sensor/temperature/temperature$(BIN_EXT) sensor/temperature/temperature.go

build-actuators:
	$(GO) build -o actuator/airConditioner/airConditioner$(BIN_EXT) actuator/airConditioner/airConditioner.go
	$(GO) build -o actuator/exhaustFan/exhaustFan$(BIN_EXT) actuator/exhaustFan/exhaustFan.go
	$(GO) build -o actuator/humidifier/humidifier$(BIN_EXT) actuator/humidifier/humidifier.go
	$(GO) build -o actuator/light/light$(BIN_EXT) actuator/light/light.go
	$(GO) build -o actuator/sprinkler/sprinkler$(BIN_EXT) actuator/sprinkler/sprinkler.go

run-server:
	gnome-terminal -- bash -c "go run server/server.go; exec bash"

run-client:
	gnome-terminal -- bash -c "go run client/client.go; exec bash"

run-sensors:
	gnome-terminal -- bash -c "go run sensor/gas/gas.go; exec bash"
	gnome-terminal -- bash -c "go run sensor/humidity/humidity.go; exec bash"
	gnome-terminal -- bash -c "go run sensor/luminosity/luminosity.go; exec bash"
	gnome-terminal -- bash -c "go run sensor/smoke/smoke.go; exec bash"
	gnome-terminal -- bash -c "go run sensor/temperature/temperature.go; exec bash"

run-actuators:
	gnome-terminal -- bash -c "go run actuator/airConditioner/airConditioner.go; exec bash"
	gnome-terminal -- bash -c "go run actuator/exhaustFan/exhaustFan.go; exec bash"
	gnome-terminal -- bash -c "go run actuator/humidifier/humidifier.go; exec bash"
	gnome-terminal -- bash -c "go run actuator/light/light.go; exec bash"
	gnome-terminal -- bash -c "go run actuator/sprinkler/sprinkler.go; exec bash"

run-all: run-server
	@sleep 1
	$(MAKE) run-client
	$(MAKE) run-sensors
	$(MAKE) run-actuators

build-and-run: build run-all

clean:
	rm -f server/server server/server.exe client/client client/client.exe \
		sensor/gas/gas sensor/gas/gas.exe \
		sensor/humidity/humidity sensor/humidity/humidity.exe \
		sensor/luminosity/luminosity sensor/luminosity/luminosity.exe \
		sensor/smoke/smoke sensor/smoke/smoke.exe \
		sensor/temperature/temperature sensor/temperature/temperature.exe \
		actuator/airConditioner/airConditioner actuator/airConditioner/airConditioner.exe \
		actuator/exhaustFan/exhaustFan actuator/exhaustFan/exhaustFan.exe \
		actuator/humidifier/humidifier actuator/humidifier/humidifier.exe \
		actuator/light/light actuator/light/light.exe \
		actuator/sprinkler/sprinkler actuator/sprinkler/sprinkler.exe
	@echo "Limpeza concluída"