package noolite

import "fmt"

// ToggleStatus Событие переключения (команда Switch, CMD=4) от устройства, которое не
// сообщает целевое состояние явно - оно означает "переключить", а не "включить"/"выключить"
// (см. руководство MTRF-64-USB-A, Таблица 3, команда Switch). Итоговое значение вычисляется
// относительно текущего состояния соответствующего канала (см. device.Device.UpdateDeviceStatus).
type ToggleStatus struct {
	// Address Адрес устройства (для nooLite-F; для классического nooLite - нулевой)
	Address [4]byte
}

func NewToggleStatus(address [4]byte) *ToggleStatus {
	return &ToggleStatus{Address: address}
}

func (ds *ToggleStatus) String() string {
	return fmt.Sprintf("[0x%s] toggle", ds.GetAddress())
}

func (ds *ToggleStatus) GetValue() string {
	return ""
}

func (ds *ToggleStatus) GetValue2() string {
	return ""
}

func (ds *ToggleStatus) GetFwVersion() string {
	return ""
}

func (ds *ToggleStatus) GetDeviceModel() string {
	return ""
}

func (ds *ToggleStatus) GetOn() bool {
	return false
}

func (ds *ToggleStatus) GetAddress() string {
	return fmt.Sprintf("%02x%02x%02x%02x", ds.Address[0], ds.Address[1], ds.Address[2], ds.Address[3])
}

func (ds *ToggleStatus) GetBatteryLow() bool {
	return false
}
