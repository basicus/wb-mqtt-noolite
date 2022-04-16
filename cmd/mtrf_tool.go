package main

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
	"os"
	"wb-noolite-mtrf/config"
	nl "wb-noolite-mtrf/noolite"
)

// Утилита для привязки или отвязки устройств. Настройки датчика. Установки температуры термостата.
func main() {

	var showHelp bool
	var err error
	var ch uint8
	var command string
	var mode string
	var nooliteMode uint8
	var temperature uint8
	var sensor string
	var ponState string

	serviceConfig := config.Config{}

	// Set environment files
	pflag.BoolVarP(&showHelp, "help", "", false, "Show help message")
	pflag.Uint8VarP(&ch, "channel", "c", 0, "Set channel")
	pflag.Uint8VarP(&temperature, "temperature", "t", 25, "Set temperature")
	pflag.StringVarP(&command, "command", "", "", "Command: bind, unbind, on, off, status, poweron_state, thermostat_mode, temperature")
	pflag.StringVarP(&mode, "mode", "m", "txf", "Mode: txf, tx, rxf,rx")
	pflag.StringVarP(&sensor, "sensor", "s", "", "Sensor: air, floor. Default: floor")
	pflag.StringVarP(&ponState, "state", "", "", "Power On state: on, off, last")
	pflag.StringVarP(&serviceConfig.SerialPort, "device", "d", "/dev/ttyUSB0", "Specify MTRF-64-USB-A serial port")

	pflag.Parse()

	// Show help
	if showHelp || ch == 0 || command == "" {
		pflag.Usage()
		return
	}

	// Check for Mode
	switch mode {
	case "txf":
		nooliteMode = nl.ModeNooliteFTX
	case "tx":
		nooliteMode = nl.ModeNooliteTX
	case "rx":
		nooliteMode = nl.ModeNooliteRX
	case "rxf":
		nooliteMode = nl.ModeNooliteFRX
	default:
		pflag.Usage()
		return
	}
	// Check for command
	var commandRequest *nl.Request
	switch command {
	case "bind":
		commandRequest = nl.RequestBindChannel(ch, nooliteMode)
	case "unbind":

	case "temperature":
		commandRequest = nl.RequestSetTemperature(ch, temperature)
	case "on":
		commandRequest = nl.RequestOn(ch, nooliteMode)
	case "off":
		commandRequest = nl.RequestOff(ch, nooliteMode)
	case "status":
		commandRequest = nl.RequestReadState(ch, nl.FmtMain)
	case "status_output":
		commandRequest = nl.RequestReadStatOutputLoad(ch)
	case "thermostat_mode":
		switch sensor {
		case "air":
			commandRequest = nl.RequestSetThermostatMode(ch, nl.ModeManualAirSensor)
		default:
			commandRequest = nl.RequestSetThermostatMode(ch, nl.ModeManualFloorSensor)
		}
	case "poweron_state":
		switch ponState {
		case "on":
			commandRequest = nl.NewRequestSetPowerOnState(ch, nl.PowerOnModeOn)
		case "last":
			commandRequest = nl.NewRequestSetPowerOnState(ch, nl.PowerOnModeLast)
		default:
			commandRequest = nl.NewRequestSetPowerOnState(ch, nl.PowerOnModeOff)
		}
	default:
		pflag.Usage()
		return
	}
	// Set logger
	log := logrus.New()
	log.SetLevel(logrus.TraceLevel)
	log.SetOutput(os.Stdout)
	log.SetFormatter(&logrus.TextFormatter{
		ForceQuote:      false,
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05.000",
	})

	// Create new request - set adapter to Service Mode
	request := nl.NewRequestServiceMode()
	if request == nil {
		log.Fatalf("Error when make request: %s", err)
	}

	// Define list initializations requests on connection
	var initialRequests []*nl.Request

	// Init Noolite service (works with MTRF adapter)
	service, err := nl.NewNooliteService(log, &serviceConfig, initialRequests)
	if err != nil {
		log.Fatalf("Cant start noolite %s", err)
	}

	defer service.Close()
	// Goroutine for receive responses
	go func() {
		for {
			r := <-service.Receive()
			deviceState := r.GetDeviceState()
			log.Infof("<-- %s", r.String())
			if deviceState != nil {
				log.Infof("<-- STATE %s", deviceState.String())
				log.Infof("<--  %s", deviceState.String())
			}

		}
	}()

	service.Send() <- commandRequest

	wait := make(chan struct{})
	<-wait
}
