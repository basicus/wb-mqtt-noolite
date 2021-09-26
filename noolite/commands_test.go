package noolite

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

//  TestNewDeviceSettings Тест на установку восстановления предыдущего состояния. Не должна влиять настройка включать или не включать при восстановлении
func TestNewDeviceSettings(t *testing.T) {
	settingsRecover1 := NewDeviceSettings(true, true, false)
	settingsRecover2 := NewDeviceSettings(false, true, false)
	assert.NotNil(t, settingsRecover1)
	assert.True(t, settingsRecover1.getData() == settingsRecover2.getData())
}

// TestDeviceSettingsSetData Тест проверки обратного преобразования из битовых данных в структуру настроек устройства.
func TestDeviceSettingsSetData(t *testing.T) {
	s1 := NewDeviceSettings(false, true, true)
	s2 := NewDeviceSettings(true, false, false)

	err := s2.setData(s1.getData())

	assert.NoError(t, err)
	assert.EqualValues(t, s1, s2)

}
