package noolite

import "fmt"

// Request Пакет запроса к адаптеру MTRF
type Request struct {
	// St Стартовый байт
	ST uint8
	// Mode Режим работы адаптера (const Mode**)
	Mode uint8
	// Ctr Управление адаптеров
	Ctr uint8
	// Res зарезервирован, не используется
	Res uint8
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

var ClearAllMemoryData = [4]byte{170, 85, 170, 85}
var EmptyData = [4]byte{0, 0, 0, 0}
var EmptyAddress = [4]byte{0, 0, 0, 0}

// NewRequest Сформировать запрос для отправки в MTRF
func NewRequest(mode uint8, control uint8, channel uint8, command uint8, fmt uint8, data [4]byte, address [4]byte) (*Request, error) {
	r := &Request{
		ST:   StartRequest,
		Mode: mode,
		Ctr:  control,
		Res:  0,
		Ch:   channel,
		Cmd:  command,
		Fmt:  fmt,
		SP:   StopRequest,
	}
	r.D0 = data[0]
	r.D1 = data[1]
	r.D2 = data[2]
	r.D3 = data[3]
	r.ID0 = address[0]
	r.ID1 = address[1]
	r.ID2 = address[2]
	r.ID3 = address[3]

	if err := r.verify(); err != nil {
		return nil, err
	}
	r.CalcCrc() // Recalculate CRC
	return r, nil
}

// verify Проверка параметров запроса
func (r *Request) verify() error {
	return nil
}

// CalcCrc Рассчитать и обновить контрольную сумму
func (r *Request) CalcCrc() {
	var s = uint(0) + uint(r.ST) + uint(r.Mode) + uint(r.Ctr) + uint(r.Res) + uint(r.Ch) + uint(r.Cmd) + uint(r.Fmt) +
		uint(r.D0) + uint(r.D1) + uint(r.D2) + uint(r.D3) + uint(r.ID0) + uint(r.ID1) + uint(r.ID2) + uint(r.ID3)

	r.Crc = byte(s & 0xFF)
}

// BuildBytes Формирует пакет для отправки в последовательный порт
func (r *Request) BuildBytes() []byte {
	return []byte{
		r.ST, r.Mode, r.Ctr, r.Res, r.Ch, r.Cmd, r.Fmt, r.D0, r.D1, r.D2, r.D3, r.ID0, r.ID1, r.ID2, r.ID3, r.Crc, r.SP}
}

func (r *Request) String() string {
	return fmt.Sprintf("Mode: %d Control: %d Command: %d Channel: %d Fmt: %d Data: %x Address: %x CRC: %d", r.Mode, r.Ctr,
		r.Cmd, r.Ch, r.Fmt, [4]byte{r.D0, r.D1, r.D2, r.D3}, [4]byte{r.ID0, r.ID1, r.ID2, r.ID3}, r.Crc)
}
