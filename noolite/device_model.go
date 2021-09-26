package noolite

// DeviceModel Модель устройства
type DeviceModel string

// NewDeviceModel Новая модель устройства
func NewDeviceModel(modelId uint8) DeviceModel {
	var model DeviceModel
	switch modelId {
	case 0:
		model = "mtrf64" // радиоконтроллер MTRF-64
	case 1:
		model = "slf-1-300" // реле  SLF-1-300
	case 2:
		model = "srf-10-1000" // реле из блока SRF-10-1000
	case 3:
		model = "srf-1-3000r" // SRF-1-3000(розетка)
	case 4:
		model = "srf-1-3000p" // SRF-1-3000 (для подрозетника)
	case 5:
		model = "suf-1-300" // SUF-1-300 диммер
	case 6:
		model = "srf-1-3000t" // терморегулятор SRF-1-3000-T
	case 7:
		model = "srf-1-1000r" // блок роллет SRF-1-1000-R

	}
	return model
}
