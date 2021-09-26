package noolite

import (
	"errors"
	"fmt"
)

// PacketLen Длина пакета для отправки в адаптер MTRF
const PacketLen = 17

var (
	ErrorInvalidPacket = errors.New("invalid packet")
	ErrorInvalidCrc    = errors.New("crc invalid")
)

// Response Ответ от адаптера
type Response struct {
	// St Стартовый байт
	ST uint8

	// Mode Режим работы адаптера (const Mode**)
	Mode uint8

	// Ctr Управление адаптеров
	Ctr uint8

	// Togl Количество оставшихся ответов от адаптера, значение Togl
	Togl uint8

	// Ch Адрес канала, ячейки привязки 0..63
	Ch uint8

	// Cmd Команда (const Cmd**)
	Cmd uint8

	// Fmt Формат
	Fmt uint8

	// D0 Байт данных 0
	D0 uint8

	// D1 Байт данных 1
	D1 uint8

	// D2 Байт данных 2
	D2 uint8

	// D3 Байт данных 3
	D3 uint8

	// ID0 идентификатор блока, биты 31...24
	ID0 uint8

	// ID1 идентификатор блока, биты 23...16
	ID1 uint8

	// ID2 идентификатор блока, биты 15...8
	ID2 uint8

	// ID3 идентификатор блока, биты 7...0
	ID3 uint8

	// Crc Контрольная сумма (младший байт от суммы первых 15 байт )
	Crc uint8

	// SP Стоповый байт, значение 172
	SP uint8
}

// Parse Парсинг принятого пакета
func (r *Response) Parse(recv []byte) error {
	if len(recv) != PacketLen {
		return ErrorInvalidPacket
	}
	if recv[0] != StartResponse || recv[PacketLen-1] != StopResponse {
		return ErrorInvalidPacket
	}
	r.ST = recv[0]
	r.SP = recv[16]
	r.Mode = recv[1]
	r.Ctr = recv[2]
	r.Togl = recv[3]
	r.Ch = recv[4]
	r.Cmd = recv[5]
	r.Fmt = recv[6]
	r.D0 = recv[7]
	r.D1 = recv[8]
	r.D2 = recv[9]
	r.D3 = recv[10]
	r.ID0 = recv[11]
	r.ID1 = recv[12]
	r.ID2 = recv[13]
	r.ID3 = recv[14]
	r.Crc = recv[15]
	if crc(recv[:PacketLen-2]) != recv[PacketLen-2] {
		return ErrorInvalidCrc
	}

	return nil
}

// CheckCRC Проверка корректности CRC принятого пакета
func (r *Response) CheckCRC() bool {
	var s = uint(0) + uint(r.ST) + uint(r.Mode) + uint(r.Ctr) + uint(r.Togl) + uint(r.Ch) + uint(r.Cmd) + uint(r.Fmt) +
		uint(r.D0) + uint(r.D1) + uint(r.D2) + uint(r.D3) + uint(r.ID0) + uint(r.ID1) + uint(r.ID2) + uint(r.ID3)

	return r.Crc == byte(s&0xFF)
}

func (r *Response) String() string {
	return fmt.Sprintf("Mode: %d Control: %d Command: %d Togl: %d Channel: %d Fmt: %d Data: %x Address %x CRC %d", r.Mode, r.Ctr,
		r.Cmd, r.Togl, r.Ch, r.Fmt, [4]byte{r.D0, r.D1, r.D2, r.D3}, [4]byte{r.ID0, r.ID1, r.ID2, r.ID3}, r.Crc)
}

func crc(buf []byte) byte {
	s := uint(0)
	for _, value := range buf {
		s += uint(value)
	}
	return byte(s & 0xFF)
}

// GetDeviceState Получение информации об устройстве
func (r *Response) GetDeviceState() StatusType {
	if r.CheckCRC() && r.Ctr == CtrResponseSuccess &&
		r.Cmd == CmdSendState {
		switch r.Fmt {
		case FmtMain:
			deviceModel := NewDeviceModel(r.D0)
			var isBindState bool
			var isTemporaryOn bool
			var isOn bool
			isBindState = (r.D2 >> 7 & 0b00000001) == 1
			loadStatus := r.D2 & 0b00000011
			switch loadStatus {
			case 0:
				isTemporaryOn = false
				isOn = false
			case 1:
				isTemporaryOn = false
				isOn = true
			case 2:
				isTemporaryOn = true
				isOn = false
			default:
			}
			deviceMainStatus := NewDeviceMainStatus(
				deviceModel, r.D1, isBindState, isTemporaryOn, isOn, r.D3, [4]byte{r.ID0, r.ID1, r.ID2, r.ID3})
			return deviceMainStatus
		case FmtSettings:
		default:
			return nil
		}

	}
	return nil
}

// IsDeviceNoResponse Если от устройства нет ответа
func (r *Response) IsDeviceNoResponse() bool {
	if r.CheckCRC() && r.Ctr == CtrResponseNoResponse {
		return true
	}
	return false
}

func (r *Response) GetAddress() string {
	return fmt.Sprintf("%02x%02x%02x%02x", r.ID0, r.ID1, r.ID2, r.ID3)
}
