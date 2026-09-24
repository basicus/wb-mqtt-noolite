package device

import (
	"encoding/json"
	"wb-noolite-mtrf/mqtt"
)

// Control Описание органа управления устройства
type Control struct {
	Name        string      `json:"name"`
	Type        ControlType `json:"type"`
	Title       Title       `json:"title,omitempty"`
	Order       int         `json:"order"`
	Readonly    bool        `json:"readonly"`
	Error       string      `json:"-"`
	Value       string      `json:"initial_value"`
	Min         int         `json:"min"`
	Max         int         `json:"max"`
	Units       string      `json:"units"`
	Precision   string      `json:"precision"`
	GetCommand  string      `json:"get_command"`
	SetCommand  string      `json:"set_command"`
	Polling     bool        `json:"polling"`
	PollingCron string      `json:"polling_cron"`
	// DontUseRetain Игнорировать ретейн-топик при восстановлении значения после переподключения
	DontUseRetain bool `json:"dont_use_retain"`
	sentOnce      bool
	hadError      bool
}

// controlMeta JSON тела топика /controls/<control>/meta, см. Wiren Board MQTT Conventions
type controlMeta struct {
	Type      string `json:"type,omitempty"`
	Title     Title  `json:"title,omitempty"`
	Order     int    `json:"order,omitempty"`
	Min       int    `json:"min,omitempty"`
	Max       int    `json:"max,omitempty"`
	Units     string `json:"units,omitempty"`
	Precision string `json:"precision,omitempty"`
	Readonly  bool   `json:"readonly"`
}

func (control *Control) GenerateMQTTPacket(controlPrefix string) []*mqtt.Message {
	var topics []*mqtt.Message

	// meta/error ретейнится, чтобы новые подписчики видели актуальное состояние ошибки;
	// пустой payload снимает ранее опубликованный флаг (см. ErrorRead/ErrorWrite/ErrorPoll)
	if control.Error != "" {
		topics = append(topics, &mqtt.Message{
			Topic:   controlPrefix + "/meta/error",
			Retain:  true,
			Payload: control.Error,
		})
		control.hadError = true
	} else if control.hadError {
		topics = append(topics, &mqtt.Message{
			Topic:   controlPrefix + "/meta/error",
			Retain:  true,
			Payload: "",
		})
		control.hadError = false
	}

	topics = append(topics, &mqtt.Message{
		Topic:   controlPrefix,
		Retain:  true,
		Payload: control.Value,
	})

	if !control.sentOnce {
		meta := controlMeta{
			Type:      control.Type.String(),
			Title:     control.Title,
			Order:     control.Order,
			Min:       control.Min,
			Max:       control.Max,
			Units:     control.Units,
			Precision: control.Precision,
			Readonly:  control.Readonly,
		}
		payload, err := json.Marshal(meta)
		if err == nil {
			topics = append(topics, &mqtt.Message{
				Topic:   controlPrefix + "/meta",
				Retain:  true,
				Payload: string(payload),
			})
		}
		control.sentOnce = true
	}
	return topics
}

func (control *Control) GetControlPrefix(deviceId string) string {
	return deviceId + "controls/" + control.Name
}
