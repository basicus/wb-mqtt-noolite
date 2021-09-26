package device

import (
	"fmt"
	"strconv"
	"wb-noolite-mtrf/mqtt"
	"wb-noolite-mtrf/noolite"
)

// Device Описание устройства с которым будем работать
type Device struct {
	Type     NooliteDeviceType `json:"noolite_type"`
	Error    string
	Ch       uint8  `json:"ch"`
	Name     string `json:"name"`
	Address  string `json:"address"`
	Template string `json:"template"`
	Controls []*Control
}

// UpdateDeviceStatus Обновить органы управления
func (d *Device) UpdateDeviceStatus(ds noolite.StatusType) bool {
	var updated int
	// Сброс ошибки
	if d.Error != "" {
		d.Error = "ok"
	}
	switch v := ds.(type) {
	case *noolite.DeviceMainStatus:
		for _, control := range d.Controls {
			switch control.Name {
			case "status":
				if v.GetOn() {
					control.Value = MQTTSwitchOn
					updated++
				} else {
					control.Value = MQTTSwitchOff
					updated++
				}
			case "value":
				control.Value = v.GetValue()
				updated++
			case "address":
				control.Value = v.GetAddress()
				updated++
			case "model":
				control.Value = v.GetDeviceModel()
				updated++
			}
		}
	default:
		panic("Cant update unknown type")
	}
	return updated > 0
}

// FindControl Поиск элемента управления по control
func (d *Device) FindControl(name string) *Control {
	if name == "" {
		return nil
	}
	for _, control := range d.Controls {
		if control.Name == name {
			return control
		}
	}
	return nil
}

// GenerateMQTTPacket Подготовить топики к публикации
func (d *Device) GenerateMQTTPacket(prefix string) mqtt.Packet {
	var topics mqtt.Packet

	var deviceId = prefix + d.Type.String() + "_" + fmt.Sprintf("%d", d.Ch) + "/"
	if d.Name != "" {
		topics.Add(mqtt.Message{
			Topic:   deviceId + "meta/name",
			Retain:  true,
			Payload: d.Name,
		})
	}

	if d.Error != "" {
		if d.Error == "ok" {
			d.Error = ""
		}
		topics.Add(mqtt.Message{
			Topic:   deviceId + "meta/error",
			Retain:  false,
			Payload: d.Error,
		})
	}
	if d.Address != "" {
		topics.Add(mqtt.Message{
			Topic:   deviceId + "meta/address",
			Retain:  true,
			Payload: d.Address,
		})
	}
	if d.Error == "" {
		for _, control := range d.Controls {
			if control.Readonly || !control.sentOnce {
				if !control.sentOnce {
					control.sentOnce = true
				}
				//Main section
				controlPrefix := deviceId + "controls/" + control.Name
				topics.Add(mqtt.Message{
					Topic:   controlPrefix,
					Retain:  true,
					Payload: control.Value,
				})
				// Meta section
				if control.Type != "" {
					topics.Add(mqtt.Message{
						Topic:   controlPrefix + "/meta/type",
						Retain:  true,
						Payload: control.Type.String(),
					})
				}
				if control.Order != 0 {
					topics.Add(mqtt.Message{
						Topic:   controlPrefix + "/meta/order",
						Retain:  true,
						Payload: strconv.Itoa(control.Order),
					})
				}

				if control.Min != 0 {
					topics.Add(mqtt.Message{
						Topic:   controlPrefix + "/meta/min",
						Retain:  true,
						Payload: strconv.Itoa(control.Min),
					})
				}

				if control.Max != 0 {
					topics.Add(mqtt.Message{
						Topic:   controlPrefix + "/meta/max",
						Retain:  true,
						Payload: strconv.Itoa(control.Max),
					})
				}
				if control.Error != "" {
					topics.Add(mqtt.Message{
						Topic:   controlPrefix + "/meta/error",
						Retain:  false,
						Payload: control.Error,
					})
				}
				if control.Units != "" {
					topics.Add(mqtt.Message{
						Topic:   controlPrefix + "/meta/units",
						Retain:  true,
						Payload: control.Units,
					})
				}
				if control.Precision != "" {
					topics.Add(mqtt.Message{
						Topic:   controlPrefix + "/meta/precision",
						Retain:  true,
						Payload: control.Precision,
					})
				}

				if control.Readonly {
					topics.Add(mqtt.Message{
						Topic:   controlPrefix + "/meta/readonly",
						Retain:  true,
						Payload: "1",
					})
				} else {
					topics.Add(mqtt.Message{
						Topic:   controlPrefix + "/meta/readonly",
						Retain:  true,
						Payload: "0",
					})
				}
			}
		}
	}
	return topics
}
