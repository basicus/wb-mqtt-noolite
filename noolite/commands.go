package noolite

import (
	"errors"
	"strconv"
)

// NewRequestServiceMode Возвращет запрос для перевода адаптер в сервисный режим
func NewRequestServiceMode() *Request {
	rs, err := NewRequest(ModeNooliteService, CtrRequestSendCommand, 0, CmdOff, FmtMain, EmptyData, EmptyAddress)
	if err != nil {
		return nil
	} else {
		return rs
	}
}

// RequestBindChannel Запрос ручной привязки устройства mode
// channel - ячейка памяти
// mode - режим ModeNooliteTX, ModeNooliteFTX, ModeNooliteRX, ModeNooliteFRX
func RequestBindChannel(channel uint8, mode uint8) *Request {
	var ctr = CtrRequestSendCommand
	if mode == ModeNooliteRX {
		ctr = CtrRequestBindEnable
	}
	rs, err := NewRequest(mode, ctr, channel, CmdBind, FmtMain, EmptyData, EmptyAddress)
	if err != nil {
		return nil
	} else {
		return rs
	}
}

// RequestUnBindChannel Запрос ручной отвязки устройства
// channel - ячейка памяти
// mode - режим ModeNooliteTX, ModeNooliteFTX, ModeNooliteRX, ModeNooliteFRX
func RequestUnBindChannel(channel uint8, mode uint8) *Request {
	var ctr = CtrRequestSendCommand
	if mode == ModeNooliteRX {
		ctr = CtrRequestClearChannel
	}
	rs, err := NewRequest(mode, ctr, channel, CmdUnbind, FmtMain, EmptyData, EmptyAddress)
	if err != nil {
		return nil
	} else {
		return rs
	}
}

// RequestClearAllChannels Очистка всех ячеек памяти
// mode - режим ModeNooliteRX, ModeNooliteFRX
func RequestClearAllChannels(mode uint8) *Request {
	if mode == ModeNooliteRX || mode == ModeNooliteFRX {
		rs, err := NewRequest(mode, CtrRequestClearAllChannels, 0, CmdOff, FmtMain, ClearAllMemoryData, EmptyAddress)
		if err != nil {
			return nil
		} else {
			return rs
		}
	}
	return nil
}

// RequestSetSettings Установка режима включения нагрузки, управление сохранением режима работы и отключением
// channel - ячейка памяти
// mode - режим ModeNooliteFRX
func RequestSetSettings(channel uint8, mode uint8, settings DeviceSettings) *Request {
	if mode != ModeNooliteFTX {
		return nil
	}
	rs, err := NewRequest(mode, CtrRequestSendCommand, channel, CmdWriteState, FmtSettings, settings.getData(), EmptyAddress)
	if err != nil {
		return nil
	} else {
		return rs
	}

}

// RequestGetSettings Получение настроек
func RequestGetSettings(channel uint8, mode uint8, settings DeviceSettings) *Request {
	if mode != ModeNooliteFTX {
		return nil
	}
	rs, err := NewRequest(mode, CtrRequestSendCommand, channel, CmdReadState, FmtSettings, EmptyData, EmptyAddress)
	if err != nil {
		return nil
	} else {
		return rs
	}
}

// RequestSetTemperature Установка температуры srf-1-3000-t и не включать.
func RequestSetTemperature(ch uint8, temp uint8) *Request {
	rs, err := NewRequest(ModeNooliteFTX, CtrRequestSendCommand, ch, CmdSetBrightness, FmtSetAndDontOn, [4]byte{temp, 0, 0, 0}, EmptyAddress)
	if err != nil {
		return nil
	} else {
		return rs
	}
}

// RequestOn Включить устройство
func RequestOn(ch uint8, mode uint8) *Request {
	rs, err := NewRequest(mode, CtrRequestSendCommand, ch, CmdOn, FmtMain, [4]byte{0, 0, 0, 0}, EmptyAddress)
	if err != nil {
		return nil
	} else {
		return rs
	}
}

// RequestOff Выключить устройство
func RequestOff(ch uint8, mode uint8) *Request {
	rs, err := NewRequest(mode, CtrRequestSendCommand, ch, CmdOff, FmtMain, [4]byte{0, 0, 0, 0}, EmptyAddress)
	if err != nil {
		return nil
	} else {
		return rs
	}
}

// RequestSetThermostatMode Установка режима термостата
func RequestSetThermostatMode(ch uint8, mode uint8) *Request {
	var d2 uint8
	switch mode {
	case ModeManualFloorSensor:
		d2 = 127
	case ModeManualAirSensor:
		d2 = 9
	case ModeManualWirelessSensor:
		d2 = 3
	case ModeCalendarFloorSensor:
		d2 = 127
	case ModeCalendarAirSensor:
		d2 = 127
	}

	rs, err := NewRequest(ModeNooliteFTX, CtrRequestSendCommand, ch, CmdWriteState, FmtSettings, [4]byte{mode, 0, d2, 0}, EmptyAddress)
	if err != nil {
		return nil
	} else {
		return rs
	}
}

// NewRequestSetPowerOnState Установка режима работы термостата после подачи питания
func NewRequestSetPowerOnState(ch uint8, mode uint8) *Request {
	const powerOnBit = 5
	const saveLastStateBit = 0
	var settings uint8 = 0
	var mask uint8 = 0
	switch mode {
	case PowerOnModeOn:
		settings |= 1 << powerOnBit
		mask |= 1 << powerOnBit
	case PowerOnModeLast:
		settings |= 1 << saveLastStateBit
		mask |= 1 << saveLastStateBit
	default:
		// PowerOnModeOff
		mask |= 1 << powerOnBit
	}

	rs, err := NewRequest(ModeNooliteFTX, CtrRequestSendCommand, ch, CmdWriteState, FmtSettings, [4]byte{settings, 0, mask, 0}, EmptyAddress)
	if err != nil {
		return nil
	} else {
		return rs
	}
}

// RequestReadState Чтение состояния устройства
func RequestReadState(ch uint8, fmt uint8) *Request {
	rs, err := NewRequest(ModeNooliteFTX, CtrRequestSendCommand, ch, CmdReadState, fmt, [4]byte{0, 0, 0, 0}, EmptyAddress)
	if err != nil {
		return nil
	} else {
		return rs
	}
}

// RequestReadStatOutputLoad Получение данных с выхода для нагрузки устройства
// Пакеты данных:
//
//Байт:	                ST	MODE	CTR	RES	CH	CMD	FMT	D0	D1	D2	D3	ID0	ID1	ID2	ID3	CRC	SP
//Передача (17 байт):	171  2	    9	0	0	128	19	0	0	0	0	id0	id1	id2	id3	crc	172
//Прием (17 байт):	    173  2      ctr	0	0	130	19	d0	d1	0	0	id0	id1	id2	id3	crc	174
//параметры пакетов:
//
//ctr = 0 – код ответа:
//
//ctr	код ответа
//0	команда выполнена
//1	нет ответа от блока
//2	ошибка во время выполнения
//3	привязка выполнена
//d0 = 0..255 – текущая мощность (яркость) на нагрузке устройства;
//d1 = 0..255 – мощность (яркость) на которую будет включена нагрузка устройства;
//id0, id1, id2, id3 – ID (адрес) устройства в системе nooLite-F, 0-ой, 1-ый, 2-ой и 3-ий байты адреса соответственно, например: 1, 171, 205, 139;
//
//crc – контрольная сумма, младший байт от суммы первых 15 байт (ST..ID3).
// Example response: Mode: 2 Control: 0 Command: 130 Togl: 0 Channel: 1 Fmt: 255 Data: 00000000 Address 00018bd6 CRC 147
func RequestReadStatOutputLoad(ch uint8) *Request {
	rs, err := NewRequest(ModeNooliteFTX, CtrRequestSendCommand, ch, CmdReadState, 19, [4]byte{0, 0, 0, 0}, EmptyAddress)
	if err != nil {
		return nil
	} else {
		return rs
	}
}

// RequestMQTTCommand Функция вызывающая и формирующая необходимый запрос
func RequestMQTTCommand(ch uint8, mode uint8, params ...string) (*Request, error) {
	if len(params) == 0 {
		return nil, errors.New("params is empty")
	}
	var commandName string
	if params[0] != "" {
		commandName = params[0]
		params = params[1:]
	}
	switch commandName {
	case MQTTReadState: // Command ReadState {fmt-int}
		var fmt uint8
		if len(params) == 1 {
			fmtInt, err := strconv.Atoi(params[0])
			if err != nil {
				return nil, errors.New("invalid command params")
			}
			fmt = uint8(fmtInt)
		}
		return RequestReadState(ch, fmt), nil
	case MQTTSetOn: // Command SetOn
		return RequestOn(ch, mode), nil
	case MQTTSetOff: // Command SefOff
		return RequestOff(ch, mode), nil
	case MQTTSetSwitch:
		if len(params) == 1 {
			status, err := strconv.Atoi(params[0])
			if err != nil {
				return nil, errors.New("invalid command params")
			}
			switch status {
			case 0:
				return RequestOff(ch, mode), nil
			case 1:
				return RequestOn(ch, mode), nil
			default:
				return nil, errors.New("invalid switch state")
			}
		}
		return nil, errors.New("invalid command params")
	case MQTTSetTemperature:
		if len(params) == 1 {
			temperature, err := strconv.Atoi(params[0])
			if err != nil {
				return nil, errors.New("invalid command params")
			}
			return RequestSetTemperature(ch, uint8(temperature)), nil
		}
		return nil, errors.New("invalid command params")
	default:
		return nil, errors.New("unknown command")
	}
}
