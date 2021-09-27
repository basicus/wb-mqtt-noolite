package device

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"sync"
	"time"
	"wb-noolite-mtrf/config"
	"wb-noolite-mtrf/mqtt"
)

// PublishQueue Очередь публикации в MQTT статусов с задержкой. Необходима для того, чтобы при частых изменениях данные уходили один раз
type PublishQueue struct {
	queue  map[string]*Device
	mqtt   *mqtt.Connector
	log    *logrus.Logger
	lock   sync.Mutex
	config *config.Config
}

// NewPublishQueue Инициализация очереди для отправки.
func NewPublishQueue(mqtt *mqtt.Connector, log *logrus.Logger, config *config.Config) *PublishQueue {
	return &PublishQueue{queue: make(map[string]*Device, 10), mqtt: mqtt, config: config, log: log}
}

func (q *PublishQueue) Enqueue(device *Device) {
	q.lock.Lock()
	defer q.lock.Unlock()
	deviceId := fmt.Sprintf("%d", device.Ch) + "_" + device.Type.String() + "_" + device.Template + "_" + device.Address
	_, ok := q.queue[deviceId]
	if ok {
		return
	}
	q.queue[deviceId] = device
	go func() {
		<-time.After(q.config.Mqtt.PublishDelay)
		q.lock.Lock()
		defer q.lock.Unlock()
		mqttPacket := device.GenerateMQTTPacket(q.config.Mqtt.DevicePrefix)
		err := q.mqtt.PublishPacket(mqttPacket)
		if err != nil {
			q.log.Errorf("MQTT publish error: %s", err)
		}
		delete(q.queue, deviceId)
	}()
}
