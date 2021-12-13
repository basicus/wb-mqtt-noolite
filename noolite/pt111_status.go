package noolite

import "fmt"

// PT111Status Основной статус устройства (FMT = 7)
type PT111Status struct {
	// Type Тип устройства
	Type DeviceModel
	// Версия прошивки
	FwVersion uint8
	// Текущая влажность в %
	Humidity uint8
	// Temperature Текущая температура
	Temperature float32
	// Адрес устройства
	Address [4]byte
}

func NewPT111StatusStatus(fwVersion uint8, temperature float32, humidity uint8, address [4]byte) *PT111Status {
	return &PT111Status{Type: NewDeviceModel(DevicePT111), FwVersion: fwVersion, Temperature: temperature, Humidity: humidity, Address: address}
}

func (ds *PT111Status) String() string {
	return fmt.Sprintf("%s [0x%s] %t %.1fT %d", ds.GetDeviceModel(), ds.GetAddress(), ds.GetOn(), ds.Temperature, ds.Humidity)
}

func (ds *PT111Status) GetValue() string {
	return fmt.Sprintf("%.1f", ds.Temperature)

}

func (ds *PT111Status) GetValue2() string {
	return fmt.Sprintf("%d", ds.Humidity)
}

func (ds *PT111Status) GetFwVersion() string {
	return string(ds.FwVersion)
}

func (ds *PT111Status) GetDeviceModel() string {
	return string(ds.Type)
}

func (ds *PT111Status) GetOn() bool {
	return false
}

func (ds *PT111Status) GetAddress() string {
	return fmt.Sprintf("%02x%02x%02x%02x", ds.Address[0], ds.Address[1], ds.Address[2], ds.Address[3])
}
