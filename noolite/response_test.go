package noolite

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

//  TestSetModeResponse Разбор ответа на команду установки
func TestSetModeResponse(t *testing.T) {
	buf := [17]byte{173, 4, 0, 0, 0, 0, 0, 8, 0, 0, 0, 0, 1, 145, 71, 146, 174}
	response := Response{}
	err := response.Parse(buf[:])
	if err != nil {
		return
	}
	assert.NoError(t, err)
	assert.True(t, response.Crc == 146)
}

func TestResponse_GetDeviceState(t *testing.T) {

	//Receive response Mode: 2 Control: 3 Command: 0 Togl: 0 Channel: 1 Fmt: 0 Data: 06000000 Address  CRC 27
	buf := [17]byte{173, 2, 1, 0, 1, 128, 0, 6, 0, 0, 0, 0, 1, 139, 214, 153, 174}
	response := Response{}

	err := response.Parse(buf[:])
	deviceState := response.GetDeviceState()
	if err != nil {
		return
	}
	assert.NoError(t, err)
	assert.Nil(t, deviceState)

}

func TestResponse_LoadOnState(t *testing.T) {
	//Receive response Mode: ,Control: ,Command: 1,Togl: ,Channel: ,Fmt: ,Data: 0600011b Address 00018bd,CRC 182
	// Load is on,
	buf := [17]byte{173, 2, 0, 0, 1, 130, 0, 6, 0, 1, 27, 0, 1, 139, 214, 182, 174}
	response := Response{}

	err := response.Parse(buf[:])
	deviceState := response.GetDeviceState()
	if err != nil {
		return
	}
	assert.NoError(t, err)
	assert.NotNil(t, deviceState)
	assert.True(t, deviceState.GetOn())
	value := deviceState.GetValue()
	assert.True(t, value == "27")
}

func TestResponse_LoadOffState(t *testing.T) {
	//Receive response Mode: ,Control: ,Command: 1,Togl: ,Channel: ,Fmt: ,Data: 0600011b Address 00018bd,CRC 182
	// Load is on,
	buf := [17]byte{173, 2, 0, 0, 1, 130, 0, 6, 0, 0, 26, 0, 1, 139, 214, 180, 174}
	response := Response{}

	err := response.Parse(buf[:])
	deviceState := response.GetDeviceState()
	if err != nil {
		return
	}
	assert.NoError(t, err)
	assert.NotNil(t, deviceState)
	assert.False(t, deviceState.GetOn())
	value := deviceState.GetValue()

	assert.True(t, value == "26")
}

func TestResponse_PT111GetState(t *testing.T) {
	buf := [17]byte{173, 1, 0, 6, 4, 21, 7, 4, 33, 61, 255, 0, 0, 149, 241, 187, 174}
	response := Response{}

	err := response.Parse(buf[:])
	if err != nil {
		return
	}

	deviceState := response.GetDeviceState()
	if err != nil {
		return
	}
	assert.NoError(t, err)
	assert.NotNil(t, deviceState)
	assert.Equal(t, deviceState.GetValue(), "26")
	assert.Equal(t, deviceState.GetValue2(), "61")
}

func TestResponse_PT111GetState2(t *testing.T) {
	buf := [17]byte{173, 1, 0, 6, 4, 21, 7, 9, 33, 45, 255, 0, 0, 149, 241, 176, 174}

	response := Response{}

	err := response.Parse(buf[:])
	if err != nil {
		return
	}

	deviceState := response.GetDeviceState()
	if err != nil {
		return
	}
	assert.NoError(t, err)
	assert.NotNil(t, deviceState)
	assert.Equal(t, deviceState.GetValue(), "26.5")
	assert.Equal(t, deviceState.GetValue2(), "45")

}
