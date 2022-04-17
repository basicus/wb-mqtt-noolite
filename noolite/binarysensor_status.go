package noolite

import "fmt"

// BinarySensorStatus Статус устройства
type BinarySensorStatus struct {

	// Статус
	// Датчика протечки WS-1 . OFF (false) - протечки нет. ON (true) - протечка есть
	// Датчика двери DS-1. OFF (false) - закрыта. ON (true) - открыта
	// Выключатель PU212-2 OFF (false) - отключено. ON (true) - включено
	Alarm bool
	// Адрес устройства
	Address [4]byte
}

func NewBinarySensorStatus(alarm bool, address [4]byte) *BinarySensorStatus {
	return &BinarySensorStatus{Alarm: alarm, Address: address}
}

func (ds *BinarySensorStatus) String() string {
	return fmt.Sprintf("%s [0x%s] switch (alarm) status  %t ", ds.GetDeviceModel(), ds.GetAddress(), ds.GetOn())
}

func (ds *BinarySensorStatus) GetValue() string {
	return ""
}

func (ds *BinarySensorStatus) GetValue2() string {
	return ""
}

func (ds *BinarySensorStatus) GetFwVersion() string {
	return ""
}

func (ds *BinarySensorStatus) GetDeviceModel() string {
	return ""
}

func (ds *BinarySensorStatus) GetOn() bool {
	return ds.Alarm
}

func (ds *BinarySensorStatus) GetAddress() string {
	return fmt.Sprintf("%02x%02x%02x%02x", ds.Address[0], ds.Address[1], ds.Address[2], ds.Address[3])
}
func (ds *BinarySensorStatus) GetBatteryLow() bool {
	return false
}
