package mqtt

import (
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/sirupsen/logrus"
	"log"
	"wb-noolite-mtrf/config"
)

type Connector struct {
	config   *config.MQTTConfig
	log      *logrus.Logger
	client   mqtt.Client
	rcvQueue chan *Message
}

func NewConnector(log *logrus.Logger, config *config.MQTTConfig) *Connector {
	return &Connector{config: config, log: log}
}

func (c *Connector) Init() error {
	// Connect to MQTT
	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tcp://%s:%d", c.config.Host, c.config.Port))
	opts.SetClientID(c.config.ClientName)
	if c.config.Username != "" {
		opts.SetUsername(c.config.Username)
		opts.SetPassword(c.config.Password)
	}

	opts.SetDefaultPublishHandler(c.messagePublishHandler)
	opts.OnConnect = c.onConnectHandler
	opts.OnConnectionLost = c.connectionLostHandler
	c.client = mqtt.NewClient(opts)
	if token := c.client.Connect(); token.Wait() && token.Error() != nil {
		log.Fatalln("[MQTT] Error: ", token.Error())
		return token.Error()
	}
	c.rcvQueue = make(chan *Message, 1)
	return nil
}

// Обработчик при получении соединения
func (c *Connector) messagePublishHandler(client mqtt.Client, msg mqtt.Message) {
	c.log.Tracef("[MQTT] Received message:  %s %s", msg.Topic(), msg.Payload())
	c.rcvQueue <- &Message{msg.Topic(), msg.Retained(), string(msg.Payload())}
}

// Обработчик при подключении
func (c *Connector) onConnectHandler(client mqtt.Client) {
	c.log.Infof("[MQTT] Connected")
	client.Subscribe(c.config.Subscription, 0, c.messagePublishHandler)
}

// Обработчик потери соединения
func (c *Connector) connectionLostHandler(client mqtt.Client, err error) {
	c.log.Errorf("[MQTT] Connect lost: %v", err)
}

func (c *Connector) PublishPacket(packet Packet) error {
	var wasError error
	for _, message := range packet.messages {
		token := c.client.Publish(message.Topic, 0, message.Retain, message.Payload)
		topic := message
		go func() {
			<-token.Done()
			if token.Error() != nil {
				wasError = token.Error()
				c.log.Errorf("[MQTT] Publish topic error %s", token.Error())
				return
			}
			c.log.Tracef("[MQTT] Successfully published topic: %s %s", topic.Topic, topic.Payload)
		}()
		if wasError != nil {
			break
		}
	}
	return wasError
}

func (c *Connector) Receive() <-chan *Message {
	return c.rcvQueue
}
