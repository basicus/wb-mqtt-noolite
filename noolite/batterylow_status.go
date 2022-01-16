package noolite

import "fmt"

// BatteryLowStatus Основной статус устройства
type BatteryLowStatus struct {
	// Статус разряда батареи. true - разряжена, false - в порядке
	BatteryLow bool
	// Адрес устройства
	Address [4]byte
}

func NewBatteryLowStatusStatus(batteryLow bool, address [4]byte) *BatteryLowStatus {
	return &BatteryLowStatus{BatteryLow: batteryLow, Address: address}
}

func (ds *BatteryLowStatus) String() string {
	return fmt.Sprintf("%s [0x%s] Open %t ", ds.GetDeviceModel(), ds.GetAddress(), ds.GetOn())
}

func (ds *BatteryLowStatus) GetValue() string {
	return ""

}

func (ds *BatteryLowStatus) GetValue2() string {
	return ""
}

func (ds *BatteryLowStatus) GetFwVersion() string {
	return ""
}

func (ds *BatteryLowStatus) GetDeviceModel() string {
	return ""
}

func (ds *BatteryLowStatus) GetOn() bool {
	return ds.BatteryLow
}

func (ds *BatteryLowStatus) GetAddress() string {
	return fmt.Sprintf("%02x%02x%02x%02x", ds.Address[0], ds.Address[1], ds.Address[2], ds.Address[3])
}
