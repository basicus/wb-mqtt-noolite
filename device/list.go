package device

import (
	"encoding/json"
	"errors"
	"github.com/go-co-op/gocron"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"regexp"
	"strconv"
	"strings"
	"wb-noolite-mtrf/config"
	"wb-noolite-mtrf/mqtt"
	"wb-noolite-mtrf/noolite"
)

var ErrDeviceNotFound = errors.New("device not found")
var ErrDeviceTemplateIsNil = errors.New("device templates is nil")

// List Список устройств
type List struct {
	devices      []*Device
	templates    *Templates
	log          *logrus.Logger
	config       *config.Config
	noolite      *noolite.Service
	mqtt         *mqtt.Connector
	topicRegex   *regexp.Regexp
	cron         *gocron.Scheduler
	publishQueue *PublishQueue
}

func (l *List) Devices() []*Device {
	return l.devices

}

// InitNoolite Инициализирует получение и обработку ответов
func (l *List) InitNoolite(service *noolite.Service) {
	l.noolite = service

	go func() { // Goroutine for receive responses
		for {
			r := <-l.noolite.Receive()
			deviceState := r.GetDeviceState()
			l.log.Debugf("<-- %s", r.String())
			if deviceState != nil {
				l.log.Tracef("<-- STATE %s", deviceState.String())
				var device *Device
				device, err := l.FindByAddress(r.GetAddress())
				if err != nil {
					device, err = l.FindByChannel(r.Ch)
				}
				if err == nil && device != nil {
					//Found device state for device. Publish it state
					if device.UpdateDeviceStatus(deviceState) {
						l.log.Tracef("Updating controls. Enqueue it for update")
						l.publishQueue.Enqueue(device)
					}
				}
			} else if r.IsDeviceNoResponse() {
				l.log.Tracef("<-- DEVICE NO RESPONSE ch %s address %s", string(r.Ch), r.GetAddress())
				var device *Device
				device, err := l.FindByAddress(r.GetAddress())
				if err != nil {
					device, err = l.FindByChannel(r.Ch)
				}
				if err == nil && device != nil {
					//Device found. Publish it state
					device.Error = "No response or not found"
					l.log.Tracef("Updating device state. Enqueue it for update")
					l.publishQueue.Enqueue(device)
				}
			}
		}
	}()
}

func (l *List) InitMQTT(connector *mqtt.Connector) {
	l.mqtt = connector
	l.publishQueue = NewPublishQueue(l.mqtt, l.log, l.config)
	go func() { // Goroutine for receive MQTT
		for {
			r := <-l.mqtt.Receive()
			// 1. Фильтрация и поиск устройства
			// 2. Найти соответствующий элемент управления
			// 3. Проверить поддерживает ли он отправку команд. Закодировать команды
			// 4. Отправить в очередь команд
			// 5. Добавить шедулер, который будет выполнять команду по расписанию
			if l.topicRegex.MatchString(r.Topic) {
				found := l.topicRegex.FindAllStringSubmatch(r.Topic, -1)
				if len(found) >= 1 && len(found[0]) == 4 {
					// 0 - all, 1 - noolite (txf, tx, rx, rxf), 2 - address or channel, 3 control
					//nooliteMode := found[0][1]
					address := found[0][2]
					controlName := found[0][3]

					var device *Device
					device, err := l.FindByAddress(address)
					if err != nil {
						ch, err := strconv.Atoi(address)
						if err != nil {
							continue
						}
						device, _ = l.FindByChannel(uint8(ch))
					}
					if device != nil {
						if control := device.FindControl(controlName); control != nil {
							if !control.Readonly && control.SetCommand != "" {

								l.log.Tracef("%+v", control)
								command := strings.Split(control.SetCommand, " ")
								command = append(command, r.Payload)
								nooliteRequest, err := noolite.RequestMQTTCommand(device.Ch, device.Type.GetMode(), command...)
								if err != nil {
									l.log.Errorf("Cant create request to Noolite device")
									continue
								} else {
									l.log.Tracef("Received command. Send it to Noolite device")
									l.noolite.Send() <- nooliteRequest
								}
							}

						}
					}
				}

			}

		}
	}()
}

// FindByChannel Найти устройство по каналу
func (l *List) FindByChannel(ch uint8) (*Device, error) {
	for _, device := range l.devices {
		if device.Ch == ch {
			return device, nil
		}
	}
	return nil, ErrDeviceNotFound
}

// FindByAddress Найти устройство по адресу
func (l List) FindByAddress(address string) (*Device, error) {
	for _, device := range l.devices {
		if device.Address == address {
			return device, nil
		}
	}
	return nil, ErrDeviceNotFound
}

// NewDeviceList Создает новый список устройств и загружает его из JSON файла
func NewDeviceList(log *logrus.Logger, config *config.Config, path string, templates *Templates) (*List, error) {
	jsonFile, err := os.Open(path)

	if err != nil {
		return nil, err
	}

	defer jsonFile.Close()

	readJson, err := ioutil.ReadAll(jsonFile)

	if err != nil {
		return nil, err
	}

	var dev []*Device
	err = json.Unmarshal(readJson, &dev)
	if err != nil {
		return nil, err
	}

	return &List{
		devices:    dev,
		log:        log,
		config:     config,
		templates:  templates,
		topicRegex: regexp.MustCompile(config.DevicePattern),
	}, nil
}

// InitDeviceTemplates Инициализация элементов управления у устройств
func (l *List) InitDeviceTemplates() error {
	if l.templates == nil {
		return ErrDeviceTemplateIsNil
	}
	for _, device := range l.devices {
		controls, err := l.templates.FindTemplateByName(device.Template)
		if err != nil {
			return err
		}
		device.Controls = controls
	}
	return nil
}

// InitDeviceScheduler Инициализация планировщика выполняющего запросы
func (l *List) InitDeviceScheduler() error {
	l.log.Infof("Initialize scheduler for devices")
	l.cron = gocron.NewScheduler(l.config.Tz)
	for _, device := range l.devices {
		for _, control := range device.Controls {
			if control.Polling && control.PollingCron != "" {
				if control.GetCommand != "" {
					l.log.Tracef("%+v", control)
					command := strings.Split(control.GetCommand, " ")
					nooliteRequest, err := noolite.RequestMQTTCommand(device.Ch, device.Type.GetMode(), command...)
					if err != nil {
						l.log.Errorf("Cant create request to Noolite device. Cant apply crontab")
						continue
					} else {
						l.log.Infof("Apply crontab for device %d ch control %s %s crontab: %s", device.Ch, device.Template, control.Name, control.PollingCron)

						_, err := l.cron.Cron(control.PollingCron).Do(func() {
							l.log.Tracef("Starting crontab for device %d ch control %s %s", device.Ch, device.Template, control.Name)
							l.noolite.Send() <- nooliteRequest
						})
						if err != nil {
							l.log.Errorf("Error when apply crontab %s", err)
						}
					}
				}
			}
		}
	}

	l.log.Infof("Scheduled %d jobs", len(l.cron.Jobs()))
	l.cron.SingletonMode()
	l.cron.StartAsync()
	l.log.Infof("Crontab jobs started")
	return nil
}
