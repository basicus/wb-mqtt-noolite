package main

import (
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
	"os"
	"time"
	"wb-noolite-mtrf/config"
	"wb-noolite-mtrf/device"
	nl "wb-noolite-mtrf/noolite"
)

// bindTimeout Время ожидания подтверждения привязки. В режимах RX/RX-F адаптер слушает канал
// 40 секунд (см. руководство MTRF-64-USB-A §5.3), поэтому таймаут берется с запасом.
const bindTimeout = 45 * time.Second

// suggestedTemplates Сопоставление модели устройства nooLite-F (см. noolite.NewDeviceModel)
// со штатными шаблонами из templates/templates.json. Используется только как подсказка при
// генерации сниппета для devices.json после успешной привязки - не гарантирует точное совпадение.
var suggestedTemplates = map[nl.DeviceModel]string{
	nl.NewDeviceModel(nl.DeviceSRF3000T): "srf-1-3000-t",
	nl.NewDeviceModel(nl.DeviceWS1):      "ws1",
	nl.NewDeviceModel(nl.DevicePT111):    "pt111",
}

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
	pflag.Uint8VarP(&ch, "channel", "c", 0, "Set channel (1-63; not needed for clear_all)")
	pflag.Uint8VarP(&temperature, "temperature", "t", 25, "Set temperature")
	pflag.StringVarP(&command, "command", "", "", "Command: bind, unbind, clear_all, on, off, status, status_output, poweron_state, thermostat_mode, temperature")
	pflag.StringVarP(&mode, "mode", "m", "txf", "Mode: txf, tx, rxf, rx")
	pflag.StringVarP(&sensor, "sensor", "s", "", "Sensor: air, floor. Default: floor")
	pflag.StringVarP(&ponState, "state", "", "", "Power On state: on, off, last")
	pflag.StringVarP(&serviceConfig.SerialPort, "device", "d", "/dev/ttyUSB0", "Specify MTRF-64-USB-A serial port")

	pflag.Parse()

	// Show help
	if showHelp || command == "" {
		pflag.Usage()
		return
	}
	// clear_all действует на все каналы сразу, конкретный канал не нужен (руководство §6.5)
	if command != "clear_all" && ch == 0 {
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
		commandRequest = nl.RequestUnBindChannel(ch, nooliteMode)
	case "clear_all":
		commandRequest = nl.RequestClearAllChannels(nooliteMode)
		if commandRequest == nil {
			fmt.Fprintln(os.Stderr, "clear_all поддерживается только в режимах rx и rxf (см. руководство MTRF-64-USB-A §6.5)")
			os.Exit(1)
		}
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

	// Create new request - set adapter to Service Mode. Отправляется сразу при подключении,
	// чтобы не ждать 12 секунд, пока адаптер выйдет из режима обновления ПО (см. руководство
	// MTRF-64-USB-A, раздел "Внимание!").
	serviceModeRequest := nl.NewRequestServiceMode()
	if serviceModeRequest == nil {
		log.Errorf("Error when make request: %s", err)
	}

	// Define list initializations requests on connection
	var initialRequests []*nl.Request
	if serviceModeRequest != nil {
		initialRequests = append(initialRequests, serviceModeRequest)
	}

	// Init Noolite service (works with MTRF adapter)
	service, err := nl.NewNooliteService(log, &serviceConfig, initialRequests)
	if err != nil {
		log.Fatalf("Cant start noolite %s", err)
	}

	defer service.Close()

	// Для bind ждем конкретного подтверждения привязки, чтобы сразу вывести готовый
	// фрагмент для devices.json, и не виснуть в ожидании после этого.
	if command == "bind" {
		service.Send() <- commandRequest
		waitForBindResult(log, service, ch, mode)
		return
	}

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

// waitForBindResult Дожидается ответа адаптера с подтверждением привязки (CTR = Bind Success,
// см. руководство MTRF-64-USB-A §5.1/§5.3) и печатает готовый фрагмент для devices.json.
// Если подтверждение не пришло за bindTimeout - выводит сообщение об ошибке и завершает работу.
func waitForBindResult(log *logrus.Logger, service *nl.Service, ch uint8, mode string) {
	timeout := time.After(bindTimeout)
	for {
		select {
		case r := <-service.Receive():
			log.Infof("<-- %s", r.String())
			if r.Ctr == nl.CtrResponseBindSuccess && r.Ch == ch {
				printDeviceConfigSnippet(ch, mode, r)
				return
			}
		case <-timeout:
			log.Errorf("Привязка не подтверждена за %s. Убедитесь, что привязываемое устройство "+
				"переведено в режим привязки (см. инструкцию на изделие), и повторите попытку.", bindTimeout)
			return
		}
	}
}

// printDeviceConfigSnippet Печатает фрагмент для devices.json, готовый к вставке: остается
// скопировать его, заполнить name и, если модель не была определена автоматически, template.
func printDeviceConfigSnippet(ch uint8, mode string, r *nl.Response) {
	dev := device.Device{
		Name:     "TODO",
		Type:     device.NooliteDeviceType(mode),
		Ch:       ch,
		Template: "TODO",
	}

	if address := r.GetAddress(); address != "00000000" {
		dev.Address = address
		if template, ok := suggestedTemplates[nl.NewDeviceModel(r.D0)]; ok {
			dev.Template = template
			fmt.Printf("\nПривязка выполнена. Модель устройства определена автоматически: %s (проверьте template).\n", template)
		} else {
			fmt.Printf("\nПривязка выполнена. Модель устройства определить не удалось (D0=%d) - укажите template вручную.\n", r.D0)
		}
	} else {
		fmt.Println("\nПривязка выполнена.")
	}

	out, err := json.MarshalIndent(&dev, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Не удалось сформировать JSON: %s\n", err)
		return
	}

	fmt.Println("Скопируйте фрагмент ниже в devices.json и заполните name (и template, если не определен):")
	fmt.Println(string(out))
}
