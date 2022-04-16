package device

import (
	"strconv"
	"wb-noolite-mtrf/mqtt"
)

// Control Описание органа управления устройства
type Control struct {
	Name          string      `json:"name"`
	Type          ControlType `json:"type"`
	Order         int         `json:"order"`
	Readonly      bool        `json:"readonly"`
	Error         string
	Value         string `json:"initial_value"`
	Min           int    `json:"min"`
	Max           int    `json:"max"`
	Units         string `json:"units"`
	Precision     string `json:"precision"`
	GetCommand    string `json:"get_command"`
	SetCommand    string `json:"set_command"`
	Polling       bool   `json:"polling"`
	PollingCron   string `json:"polling_cron"`
	sentOnce      bool
	dontUseRetain bool `json:"dont_use_retain"`
}

func (control *Control) GenerateMQTTPacket(controlPrefix string) []*mqtt.Message {
	var topics []*mqtt.Message
	if control.Error != "" {
		topics = append(topics, &mqtt.Message{
			Topic:   controlPrefix + "/meta/error",
			Retain:  false,
			Payload: control.Error,
		})
	}

	topics = append(topics, &mqtt.Message{
		Topic:   controlPrefix,
		Retain:  true,
		Payload: control.Value,
	})

	if !control.sentOnce {
		// Meta section
		if control.Type != "" {
			topics = append(topics, &mqtt.Message{
				Topic:   controlPrefix + "/meta/type",
				Retain:  true,
				Payload: control.Type.String(),
			})
		}
		if control.Order != 0 {
			topics = append(topics, &mqtt.Message{
				Topic:   controlPrefix + "/meta/order",
				Retain:  true,
				Payload: strconv.Itoa(control.Order),
			})
		}

		if control.Min != 0 {
			topics = append(topics, &mqtt.Message{
				Topic:   controlPrefix + "/meta/min",
				Retain:  true,
				Payload: strconv.Itoa(control.Min),
			})
		}

		if control.Max != 0 {
			topics = append(topics, &mqtt.Message{
				Topic:   controlPrefix + "/meta/max",
				Retain:  true,
				Payload: strconv.Itoa(control.Max),
			})
		}
		if control.Units != "" {
			topics = append(topics, &mqtt.Message{
				Topic:   controlPrefix + "/meta/units",
				Retain:  true,
				Payload: control.Units,
			})
		}
		if control.Precision != "" {
			topics = append(topics, &mqtt.Message{
				Topic:   controlPrefix + "/meta/precision",
				Retain:  true,
				Payload: control.Precision,
			})
		}

		if control.Readonly {
			topics = append(topics, &mqtt.Message{
				Topic:   controlPrefix + "/meta/readonly",
				Retain:  true,
				Payload: "1",
			})
		} else {
			topics = append(topics, &mqtt.Message{
				Topic:   controlPrefix + "/meta/readonly",
				Retain:  true,
				Payload: "0",
			})
		}
		control.sentOnce = true
	}
	return topics
}

func (control *Control) GetControlPrefix(deviceId string) string {
	return deviceId + "controls/" + control.Name
}
