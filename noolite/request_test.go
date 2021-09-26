package noolite

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

//  TestNewRequestFTXBind Привязка Noolite-F TX
func TestNewRequestFTXBind(t *testing.T) {
	request := RequestBindChannel(5, ModeNooliteFTX)
	assert.NotNil(t, request)
	assert.True(t, request.Crc == 193)
}

// TestNewRequestFTXRemoteBind Удаленная привязка в режиме Noolite-F TX
func TestNewRequestFTXRemoteBind(t *testing.T) {
	// Формирование команды сервисного режима. Этап 1
	requestService, err := NewRequest(ModeNooliteFTX, CtrRequestSendCommand, 5, CmdService, 0, [4]byte{1, 0, 0, 0}, EmptyAddress)
	assert.NoError(t, err)
	assert.True(t, requestService.Crc == 54)

	// Отправка команды Bind. Этап 2
	requestBind, err := NewRequest(ModeNooliteFTX, CtrRequestSendCommand, 10, CmdBind, 0, [4]byte{0, 0, 0, 0}, EmptyAddress)
	assert.NoError(t, err)
	assert.True(t, requestBind.Crc == 198)
}

// TestNewRequestRxBind Привязка в режиме Noolite RX
func TestNewRequestRxBind(t *testing.T) {
	request, err := NewRequest(ModeNooliteRX, CtrRequestBindEnable, 5, CmdOff, 0, EmptyData, EmptyAddress)
	assert.NoError(t, err)
	assert.True(t, request.Crc == 180)
}

//  TestNewRequestFUnbind Ручная отвязка в режиме Noolite-F TX
func TestNewRequestFUnbind(t *testing.T) {
	request, err := NewRequest(ModeNooliteTX, CtrRequestSendCommand, 5, CmdUnbind, 0, EmptyData, EmptyAddress)
	assert.NoError(t, err)
	assert.True(t, request.Crc == 185)
}

//  TestNewRequestFTxUnbind Ручная отвязка в режиме Noolite-F TX
func TestNewRequestFTxUnbind(t *testing.T) {
	request, err := NewRequest(ModeNooliteFTX, CtrRequestSendCommand, 5, CmdUnbind, 0, EmptyData, EmptyAddress)
	assert.NoError(t, err)
	assert.True(t, request.Crc == 187)
}

//  TestNewRequestFTxRemoteUnbind Удаленная отвязка в режиме Noolite-F TX
func TestNewRequestFTxRemoteUnbind(t *testing.T) {
	// Формирование команды сервисного режима. Этап 1
	requestService, err := NewRequest(ModeNooliteFTX, CtrRequestSendCommand, 5, CmdService, 0, [4]byte{1, 0, 0, 0}, EmptyAddress)
	assert.NoError(t, err)
	assert.True(t, requestService.Crc == 54)

	// Отправка команды Bind. Этап 2
	request, err := NewRequest(ModeNooliteFTX, CtrRequestSendCommand, 5, CmdUnbind, 0, [4]byte{0, 0, 0, 0}, EmptyAddress)
	assert.NoError(t, err)
	assert.True(t, request.Crc == 187)
}

//  TestNewRequestFTClearChannel Отвязка (очистка канала) в режимах Noolite-F RX
func TestNewRequestFTClearChannel(t *testing.T) {
	request, err := NewRequest(ModeNooliteFRX, CtrRequestClearChannel, 5, CmdOff, 0, EmptyData, EmptyAddress)
	assert.NoError(t, err)
	assert.True(t, request.Crc == 184)
}

//  TestNewRequestRXClearChannel Очистка памяти  режимах Noolite RX
func TestNewRequestRXClearChannel(t *testing.T) {
	request, err := NewRequest(ModeNooliteRX, CtrRequestClearAllChannels, 0, CmdOff, 0, ClearAllMemoryData, EmptyAddress)
	assert.NoError(t, err)
	assert.True(t, request.Crc == 176)
}

//  TestNewRequestFRXClearChannel Очистка памяти  режимах Noolite-F RX
func TestNewRequestFRXClearChannel(t *testing.T) {
	request, err := NewRequest(ModeNooliteFRX, CtrRequestClearAllChannels, 0, CmdOff, 0, ClearAllMemoryData, EmptyAddress)
	assert.NoError(t, err)
	assert.True(t, request.Crc == 178)
}

// TestNewRequestSendFOn Передача команды включения
func TestNewRequestSendFOn(t *testing.T) {
	request, err := NewRequest(ModeNooliteFTX, CtrRequestSendCommand, 10, CmdOn, 0, EmptyData, EmptyAddress)
	assert.NoError(t, err)
	assert.True(t, request.Crc == 185)
}

// TestNewRequestSendFOff Передача команды выключения
func TestNewRequestBroadcastSendFOff(t *testing.T) {
	request, err := NewRequest(ModeNooliteFTX, CtrRequestSendBroadcastCommand, 10, CmdOff, 0, EmptyData, EmptyAddress)
	assert.NoError(t, err)
	assert.True(t, request.Crc == 184)
}
