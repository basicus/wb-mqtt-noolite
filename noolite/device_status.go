package noolite

import "fmt"

// StatusType Статус устройства
type StatusType interface {
	// String Текстовое представление информации об устройстве
	String() string
	// GetValue Получение текущего значения Мощности, температуры текущего датчика
	GetValue() string
	// GetValue2 string Второе значение Влажность, например
	GetValue2() string
	// GetFwVersion Получение версии прошивки
	GetFwVersion() string
	// GetDeviceModel Получить модель устройства
	GetDeviceModel() string
	// GetOn Получения статуса включено или нет
	GetOn() bool
	// GetAddress Получение адреса устройства
	GetAddress() string
	// GetBatteryLow Получение статуса батареи
	GetBatteryLow() bool
}

// DeviceMainStatus Основной статус устройства (FMT = 0)
type DeviceMainStatus struct {
	// Type Тип устройства
	Type DeviceModel
	// Версия прошивки
	FwVersion uint8
	// Сервисный режим: привязка включена
	IsBindState bool
	// Временное включение
	IsTemporaryOn bool
	// Нагрузка включена
	IsOn bool
	// Текущий уровень 0... 25
	Value uint8
	// Адрес устройства
	Address [4]byte
}

func NewDeviceMainStatus(deviceType DeviceModel, fwVersion uint8, isBindState bool, isTemporaryOn bool, isOn bool, value uint8, address [4]byte) *DeviceMainStatus {
	return &DeviceMainStatus{Type: deviceType, FwVersion: fwVersion, IsBindState: isBindState, IsTemporaryOn: isTemporaryOn, IsOn: isOn, Value: value, Address: address}
}

func (ds *DeviceMainStatus) String() string {
	return fmt.Sprintf("%s [0x%s] %t %d", ds.GetDeviceModel(), ds.GetAddress(), ds.GetOn(), ds.Value)
}

func (ds *DeviceMainStatus) GetValue() string {
	return fmt.Sprintf("%d", ds.Value)

}

func (ds *DeviceMainStatus) GetValue2() string {
	return ""
}

func (ds *DeviceMainStatus) GetFwVersion() string {
	return string(ds.FwVersion)
}

func (ds *DeviceMainStatus) GetDeviceModel() string {
	return string(ds.Type)
}

func (ds *DeviceMainStatus) GetOn() bool {
	return ds.IsTemporaryOn || ds.IsOn
}

func (ds *DeviceMainStatus) GetAddress() string {
	return fmt.Sprintf("%02x%02x%02x%02x", ds.Address[0], ds.Address[1], ds.Address[2], ds.Address[3])
}

func (ds *DeviceMainStatus) GetBatteryLow() bool {
	return false
}
