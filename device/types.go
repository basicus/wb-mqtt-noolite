package device

import (
	"encoding/json"
	"errors"
	"wb-noolite-mtrf/noolite"
)

// ControlType MQTT Wirebboard типы органов управления
type ControlType string

// Wirenboard control types See full list https://github.com/wirenboard/conventions
const (
	WbControlTypeSwitch  ControlType = "switch" // Possible values 0, 1
	WbControlTypeAlarm   ControlType = "alarm"  // Possible values 0, 1
	WbControlTypeRange   ControlType = "range"  // Possible values 0..255, min 0
	WbControlTypeRGB     ControlType = "rgb"    // Possible values R;G;B where R,G,B 0...255
	WbControlTypeText    ControlType = "text"   // R/O Text, anything
	WbControlTypeGeneric ControlType = "value"  // Generic value, физическая величина уточняется метаключом units

	// Deprecated: специфичные типы величин объявлены устаревшими в Wiren Board MQTT Conventions.
	// Новые шаблоны должны использовать WbControlTypeGeneric + метаключ units ("deg C", "%", "W", "kWh").
	// Константы оставлены только для чтения старых конфигураций конвертером (см. cmd/configconvert).
	WbControlTypeTemperature      ControlType = "temperature"
	WbControlTypeRelHumidity      ControlType = "rel_humidity"
	WbControlTypePower            ControlType = "power"
	WbControlTypePowerConsumption ControlType = "power_consumption"
)

// Title Локализованное название (device/control), ключи - коды языков ("en", "ru")
type Title map[string]string

// UnmarshalJSON Десериализация JSON
func (ct *ControlType) UnmarshalJSON(b []byte) error {

	type CT ControlType
	var res *CT = (*CT)(ct)
	err := json.Unmarshal(b, &res)
	if err != nil {
		return err
	}
	switch *ct {
	case WbControlTypeSwitch, WbControlTypeAlarm, WbControlTypeRange, WbControlTypeRGB,
		WbControlTypeText, WbControlTypeGeneric, WbControlTypeTemperature, WbControlTypeRelHumidity,
		WbControlTypePower, WbControlTypePowerConsumption:
		return nil
	}
	return errors.New("invalid Wirenboard control type")
}

func (ct *ControlType) String() string {
	return string(*ct)
}

// NooliteDeviceType Режим работы устройства Noolite Noolite-F или Noolite. Режим приема коаманд (RX), режим отправки команд (TX)
type NooliteDeviceType string

// MQTT Noolite device types
const (
	TypeTx  NooliteDeviceType = "tx"
	TypeRx  NooliteDeviceType = "rx"
	TypeTxF NooliteDeviceType = "txf"
	TypeRxF NooliteDeviceType = "rxf"
)

func (nd *NooliteDeviceType) String() string {
	return string(*nd)
}

// UnmarshalJSON Десериализация JSON
func (nd *NooliteDeviceType) UnmarshalJSON(b []byte) error {

	type ND NooliteDeviceType
	var res *ND = (*ND)(nd)
	err := json.Unmarshal(b, &res)
	if err != nil {
		return err
	}
	switch *nd {
	case TypeTx, TypeRx, TypeTxF, TypeRxF:
		return nil
	}
	return errors.New("invalid noolite device type")
}

// GetMode Получаем MTRF режим
func (nd *NooliteDeviceType) GetMode() uint8 {
	switch *nd {
	case TypeTx:
		return noolite.ModeNooliteTX
	case TypeTxF:
		return noolite.ModeNooliteFTX
	case TypeRx:
		return noolite.ModeNooliteRX
	case TypeRxF:
		return noolite.ModeNooliteFRX
	}
	return noolite.ModeNooliteFTX

}
