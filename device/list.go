package device

import (
	"encoding/json"
	"errors"
	"github.com/go-co-op/gocron"
	"github.com/jinzhu/copier"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
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
	retainRegex  *regexp.Regexp
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
	go func() {
		<-time.After(l.config.Mqtt.PublishNewDeviceDelay)
		// Scan for Noolite TX Devices for publish to mqtt info about it
		for _, device := range l.devices {
			if device.Type.GetMode() == noolite.ModeNooliteTX && !device.receiveOnce {
				l.log.Tracef("Send TX device init to MQTT: %+v", device)
				l.publishQueue.Enqueue(device)
			}
		}
	}()
	go func() { // Goroutine for receive MQTT
		for {
			r := <-l.mqtt.Receive()
			newValue := l.topicRegex.MatchString(r.Topic)
			retainValue := l.retainRegex.MatchString(r.Topic) && r.Retain
			if newValue || retainValue {
				var found [][]string
				if newValue {
					found = l.topicRegex.FindAllStringSubmatch(r.Topic, -1)
				} else {
					found = l.retainRegex.FindAllStringSubmatch(r.Topic, -1)
				}

				if len(found) >= 1 && len(found[0]) == 4 {
					// 0 - all, 1 - noolite (txf, tx, rx, rxf), 2 - address or channel, 3 control
					//nooliteMode := found[0][1]
					address := found[0][2]
					controlName := found[0][3]

					l.log.Tracef("Searching device by channel/address %s control name %s", address, controlName)
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
						if newValue {
							l.log.Tracef("New value for device: %+v", device)
						} else {
							l.log.Tracef("Retain value for device: %+v", device)
						}
						if !device.receiveOnce {
							device.receiveOnce = true
						} // Check for receive data from MQTT for this device

						if control := device.FindControl(controlName); control != nil {
							if r.Payload != control.Value {
								l.log.Tracef("Received retained value for Noolite device %s control %s setting to it %s. Old value was %s", device.String(), control.Name, r.Payload, control.Value)
							} else {
								retainValue = false
							}
							if (newValue || retainValue) && !control.dontUseRetain {
								control.Value = r.Payload
								if !control.Readonly && control.SetCommand != "" {
									l.log.Tracef("Found control for send to Noolite device: %+v", control)
									command := strings.Split(control.SetCommand, " ")
									command = append(command, control.Value)
									nooliteRequest, err := noolite.RequestMQTTCommand(device.Ch, device.Type.GetMode(), command...)
									if err != nil {
										l.log.Errorf("Cant create request to Noolite device")
										continue
									} else {
										l.log.Tracef("Received command. Send it to Noolite device")
										l.noolite.Send() <- nooliteRequest
										mqttPacket := control.GenerateMQTTPacket(control.GetControlPrefix(device.GetDeviceId(l.config.Mqtt.DevicePrefix)))
										err := l.mqtt.PublishPacket(mqtt.NewPacket(mqttPacket...))
										if err != nil {
											continue
										}
									}
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
		devices:     dev,
		log:         log,
		config:      config,
		templates:   templates,
		topicRegex:  regexp.MustCompile(config.DevicePattern),
		retainRegex: regexp.MustCompile(config.RetainSettingsPattern),
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
		var cp []*Control
		for _, c := range controls {
			newControl := &Control{}
			copier.Copy(newControl, c)
			cp = append(cp, newControl)
		}
		device.Controls = cp
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
						d := *device
						c := *control
						_, err := l.cron.Cron(control.PollingCron).Do(func() {
							l.log.Tracef("Starting crontab for channel %d control %s %s", d.Ch, d.Template, c.Name)
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
