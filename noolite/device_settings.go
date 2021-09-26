package noolite

// DeviceSettings Настройки устройства. Включать нагрузку при подаче питания. Восстанавливать предыдущее состояние. Отключение приема команд Noolite-F
type DeviceSettings struct {
	PowerOnEnable                 bool
	RecoverLastState              bool
	DisableReceiveNooliteCommands bool
}

// NewDeviceSettings Создать новые настройки
func NewDeviceSettings(powerOnEnable bool, recoverLastState bool, disableReceiveNooliteCommands bool) *DeviceSettings {
	return &DeviceSettings{PowerOnEnable: powerOnEnable, RecoverLastState: recoverLastState, DisableReceiveNooliteCommands: disableReceiveNooliteCommands}
}
func (s *DeviceSettings) getData() [4]byte {
	var d = [4]uint8{0, 0, 0, 0}

	if s.RecoverLastState {
		d[0] += 0b00000001
		d[2] += 0b00000001
	} else {
		if s.PowerOnEnable {
			d[0] += 0b00100000
			d[2] += 0b00100000
		} else {
			d[0] += 0b00000000
			d[2] += 0b00100000
		}
	}
	if s.DisableReceiveNooliteCommands {
		d[0] += 0b00000100
		d[2] += 0b00000100
	} else {
		d[0] += 0b00000000
		d[2] += 0b00000100
	}
	return d
}

func (s *DeviceSettings) setData(data [4]byte) error {
	s.RecoverLastState = data[0]&1 == 1
	s.DisableReceiveNooliteCommands = data[0]&4 == 4
	s.PowerOnEnable = data[0]&32 == 32
	return nil
}
