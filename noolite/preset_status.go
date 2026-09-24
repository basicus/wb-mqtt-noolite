package noolite

import "fmt"

// PresetStatus Событие вызова сценария (команда Load_Preset, CMD=7), например по нажатию
// запрограммированной кнопки пульта (см. руководство MTRF-64-USB-A, Таблица 3).
//
// Руководство не описывает формат поля данных для этой команды применительно к пультам -
// Preset содержит "как есть" байт D0 принятого пакета. Конкретный смысл значения (номер
// сценария/пресета) зависит от того, как было настроено конкретное устройство, и должен
// быть проверен опытным путем перед использованием в автоматизациях (wb-rules и т.п.).
type PresetStatus struct {
	// Preset Байт данных D0 принятого пакета
	Preset uint8
	// Address Адрес устройства (для nooLite-F; для классического nooLite - нулевой)
	Address [4]byte
}

func NewPresetStatus(preset uint8, address [4]byte) *PresetStatus {
	return &PresetStatus{Preset: preset, Address: address}
}

func (ds *PresetStatus) String() string {
	return fmt.Sprintf("[0x%s] load_preset %d", ds.GetAddress(), ds.Preset)
}

func (ds *PresetStatus) GetValue() string {
	return fmt.Sprintf("%d", ds.Preset)
}

func (ds *PresetStatus) GetValue2() string {
	return ""
}

func (ds *PresetStatus) GetFwVersion() string {
	return ""
}

func (ds *PresetStatus) GetDeviceModel() string {
	return ""
}

func (ds *PresetStatus) GetOn() bool {
	return false
}

func (ds *PresetStatus) GetAddress() string {
	return fmt.Sprintf("%02x%02x%02x%02x", ds.Address[0], ds.Address[1], ds.Address[2], ds.Address[3])
}

func (ds *PresetStatus) GetBatteryLow() bool {
	return false
}
