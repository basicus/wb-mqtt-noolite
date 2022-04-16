package device

import (
	"fmt"
	"wb-noolite-mtrf/mqtt"
	"wb-noolite-mtrf/noolite"
)

// Device Описание устройства с которым будем работать
type Device struct {
	Type        NooliteDeviceType `json:"noolite_type"`
	Error       string
	Ch          uint8  `json:"ch"`
	Name        string `json:"name"`
	Address     string `json:"address"`
	Template    string `json:"template"`
	Controls    []*Control
	sentOnce    bool
	receiveOnce bool
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
			case ControlSetting:
			case ControlStatus:
				if v.GetOn() {
					control.Value = MQTTSwitchOn
					updated++
				} else {
					control.Value = MQTTSwitchOff
					updated++
				}
			case ControlValue:
				control.Value = v.GetValue()
				updated++
			case ControlAddress:
				control.Value = v.GetAddress()
				updated++
			case ControlModel:
				control.Value = v.GetDeviceModel()
				updated++
			}
		}
	case *noolite.PT111Status:
		for _, control := range d.Controls {
			switch control.Name {
			case ControlHumidity:
				control.Value = v.GetValue2()
				updated++
			case ControlTemperature:
				control.Value = v.GetValue()
				updated++
			case ControlAddress:
				control.Value = v.GetAddress()
				updated++
			case ControlModel:
				control.Value = v.GetDeviceModel()
				updated++
			case ControlBatteryLow:
				if v.GetBatteryLow() {
					control.Value = MQTTSwitchOn
					updated++
				} else {
					control.Value = MQTTSwitchOff
					updated++
				}
			}
		}
	case *noolite.BinarySensorStatus:
		for _, control := range d.Controls {
			switch control.Name {
			case ControlStatus:
				if v.GetOn() {
					control.Value = MQTTSwitchOn
					updated++
				} else {
					control.Value = MQTTSwitchOff
					updated++
				}
			case ControlAddress:
				control.Value = v.GetAddress()
				updated++
			}
		}
	case *noolite.BatteryLowStatus:
		for _, control := range d.Controls {
			switch control.Name {
			case ControlBatteryLow:
				if v.GetOn() {
					control.Value = MQTTSwitchOn
					updated++
				} else {
					control.Value = MQTTSwitchOff
					updated++
				}
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
func (d *Device) GenerateMQTTPacket(prefix string) *mqtt.Packet {
	topics := mqtt.NewPacket()

	var deviceId = d.GetDeviceId(prefix)
	if !d.sentOnce {
		if d.Name != "" {
			topics.Add(&mqtt.Message{
				Topic:   deviceId + "meta/name",
				Retain:  true,
				Payload: d.Name,
			})
		}
		if d.Address != "" {
			topics.Add(&mqtt.Message{
				Topic:   deviceId + "meta/address",
				Retain:  true,
				Payload: d.Address,
			})
		}
		d.sentOnce = true
	}

	if d.Error != "" {
		if d.Error == "ok" {
			d.Error = ""
		}
		topics.Add(&mqtt.Message{
			Topic:   deviceId + "meta/error",
			Retain:  false,
			Payload: d.Error,
		})
	}

	if d.Error == "" {
		for _, control := range d.Controls {
			controlPrefix := control.GetControlPrefix(deviceId)
			topics.Add(control.GenerateMQTTPacket(controlPrefix)...)
		}
	}
	return topics
}

func (d *Device) GetDeviceId(prefix string) string {
	return prefix + d.Type.String() + "_" + fmt.Sprintf("%d", d.Ch) + "/"

}

func (d *Device) String() string {
	return fmt.Sprintf("ch:%d %s [%s] %s %s", d.Ch, d.Type.String(), d.Template, d.Name, d.Address)
}
