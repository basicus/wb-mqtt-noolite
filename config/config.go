package config

import (
	"time"
)

// Config Конфигурационные параметры
type Config struct {
	//SerialPort Последовательный порт к которому подключен адаптер
	SerialPort string `json:"serial_port,omitempty" env-default:"/dev/ttyUSB0"`
	//QueueLen Глубина очереди команд принимаемых для отправки через адаптер
	QueueLen int `json:"queue_len,omitempty" env-default:"32"`
	//DevicePattern Регулярное выражение с помощью которого выполняется фильтрация входящих топиков
	DevicePattern string `json:"device_pattern,omitempty" env-default:"^/devices/mtrf_([a-z]*)_([a-zA-Z0-9]*)/controls/(\\w*)$"`
	Tz            *time.Location
	//TimeZone Часовой пояс согласно которого описываются правила crontab
	TimeZone string `json:"timezone,omitempty" env-default:"Europe/Moscow"`
	//Loglevel Уровень логирования (info, debug, trace, error)
	Loglevel     string     `json:"loglevel" env-default:"info"`
	Mqtt         MQTTConfig `json:"mqtt"`
	DeviceConfig Files      `json:"device_config"`
}

type Files struct {
	// Templates описание правил устройств
	Templates string `json:"templates,omitempty" env-default:"/etc/wb-mqtt-noolite-templates.json"`
	// Devices Устройства в системе и правила для работы с ними
	Devices string `json:"devices,omitempty" env-default:"/etc/wb-mqtt-noolite-devices.json"`
}

// MQTTConfig Конфигурация MQTT
type MQTTConfig struct {
	// Host Адрес подключения к MQTT брокеру
	Host string `json:"host,omitempty" env-default:"127.0.0.1"`
	// Port Порт подключения к MQTT брокеру
	Port int `json:"port,omitempty" env-default:"1883"`
	// Username Имя пользователя при подключении к MQTT брокеру
	Username string `json:"username,omitempty" env-default:""`
	// Password Пароль пользователя при подключении к MQTT брокеру
	Password string `json:"password,omitempty" env-default:""`
	// DevicePrefix Префикс устройства при генерации и отправке топиков в MQTT брокер
	DevicePrefix string `json:"device_prefix,omitempty" env-default:"/devices/mtrf_"`
	// Subscription Строка подписки при подключении к MQTT брокеру
	Subscription string `json:"subscription,omitempty" env-default:"/devices/#"`
	// ClientName Имя клиента при подключении к MQTT брокеру
	ClientName string `json:"client_name,omitempty" env-default:"wb-mqtt-noolite"`
	// PublishDelay время на которое откладывается публикация топиков устройства, защита от множественных отправок топиков по одному и тому же устройству
	PublishDelay time.Duration `json:"publish_delay,omitempty" env-default:"100ms"`
}
