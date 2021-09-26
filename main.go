package main

import (
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
	"os"
	"time"
	"wb-noolite-mtrf/config"
	"wb-noolite-mtrf/device"
	"wb-noolite-mtrf/mqtt"
	nl "wb-noolite-mtrf/noolite"
)

// Основной сервис выполняющий управление устройствами, периодический опрос
func main() {

	var showHelp bool
	var configFile string
	var err error
	var serviceConfig config.Config

	// Set logger
	log := logrus.New()
	log.SetLevel(logrus.TraceLevel)
	log.SetOutput(os.Stdout)
	log.SetFormatter(&logrus.TextFormatter{
		ForceQuote:      false,
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05.000",
	})

	// Set environment files
	pflag.BoolVarP(&showHelp, "help", "h", false, "Show help message")
	pflag.StringVarP(&configFile, "config", "c", "/etc/wb-mqtt-noolite.json", "Please specify config json file")
	pflag.Parse()

	err = cleanenv.ReadConfig(configFile, &serviceConfig)
	if err != nil {
		log.Fatalf("Cant read config %s", err)
	}
	serviceConfig.Mqtt.PublishDelay = time.Millisecond * 100 // 100 мсек, задержка перед отправкой статуса устройства в MQTT
	// Show help
	if showHelp {
		pflag.Usage()
		return
	}
	logLevel := logrus.InfoLevel
	err = logLevel.UnmarshalText([]byte(serviceConfig.Loglevel))
	if err != nil {
		log.SetLevel(logrus.InfoLevel)

	} else {
		log.SetLevel(logLevel)
	}
	log.Infof("Logging level: %s", log.GetLevel().String())

	// Initiation of Template Engine
	templateEngine, err := device.NewTemplateEngine(serviceConfig.DeviceConfig.Templates)
	if err != nil {
		log.Fatalf("Error when read templates JSON file: %s", err)
	}

	// Initiation of Devices List
	deviceList, err := device.NewDeviceList(log, &serviceConfig, serviceConfig.DeviceConfig.Devices, templateEngine)
	if err != nil {
		log.Fatalf("Error when read device JSON file: %s", err)
	}

	// Scheduler
	tzLocation, err := time.LoadLocation(serviceConfig.TimeZone)
	if err != nil {
		log.Fatalf("Error when read TimeZone of Location: %s", err)
	}
	serviceConfig.Tz = tzLocation

	// Create new request - set adapter to Service Mode
	request := nl.NewRequestServiceMode()
	if request == nil {
		log.Errorf("Error when make request: %s", err)
	}

	// Define list initializations requests on connection
	var initialRequests []*nl.Request

	// Init Noolite service (works with MTRF adapter)
	service, err := nl.NewNooliteService(log, &serviceConfig, initialRequests)
	if err != nil {
		log.Fatalf("Cant start noolite %s", err)
	}

	defer service.Close()

	mqttConnector := mqtt.NewConnector(log, &serviceConfig.Mqtt)
	err = mqttConnector.Init()
	if err != nil {
		log.Fatalf("Cant initialize MQTT connection")
	}

	// Init devices
	deviceList.InitMQTT(mqttConnector)
	deviceList.InitNoolite(service)
	err = deviceList.InitDeviceTemplates()
	if err != nil {
		log.Fatalf("Init device templates error: %s", err)
	}
	err = deviceList.InitDeviceScheduler()
	if err != nil {
		log.Fatalf("Init device scheduler error: %s", err)
	}

	wait := make(chan struct{})
	<-wait
}
