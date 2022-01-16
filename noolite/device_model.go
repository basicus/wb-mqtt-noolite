package noolite

// DeviceModel Модель устройства
type DeviceModel string

// NewDeviceModel Новая модель устройства
func NewDeviceModel(modelId uint8) DeviceModel {
	var model DeviceModel
	switch modelId {
	case DeviceMTRF64:
		model = "mtrf64" // радиоконтроллер MTRF-64
	case DeviceSLF300:
		model = "slf-1-300" // реле  SLF-1-300
	case DeviceSRF101000:
		model = "srf-10-1000" // реле из блока SRF-10-1000
	case DeviceSRF3000Rozetka:
		model = "srf-1-3000r" // SRF-1-3000(розетка)
	case DeviceSRF3000Podrozetnik:
		model = "srf-1-3000p" // SRF-1-3000 (для подрозетника)
	case DeviceSUF300:
		model = "suf-1-300" // SUF-1-300 диммер
	case DeviceSRF3000T:
		model = "srf-1-3000t" // терморегулятор SRF-1-3000-T
	case DeviceSRF1000R:
		model = "srf-1-1000r" // блок роллет SRF-1-1000-R
	case DevicePT111:
		model = "pt111" // датчик температуры и влажности
	case DeviceDS1:
		model = "ds1" // датчик открытия окон/дверей
	case DeviceWS1:
		model = "ws1" // датчик протечки
	case DeviceBatteryLow:
		model = "battery_low"
	default:
		model = "unknown"
	}
	return model
}
