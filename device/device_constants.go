package device

// DriverName Имя драйвера, публикуемое в /devices/<id>/meta согласно Wiren Board MQTT Conventions
const DriverName = "wb-mqtt-noolite"

const (
	MQTTSwitchOn  = "1"
	MQTTSwitchOff = "0"
)

// Флаги ошибок meta/error согласно Wiren Board MQTT Conventions (Error Handling):
// r - ошибка чтения/связи, w - ошибка записи, p - пропущен период опроса.
const (
	ErrorRead  = "r"
	ErrorWrite = "w"
	ErrorPoll  = "p"
)

const (
	ControlStatus      = "on"
	ControlSetting     = "setpoint_t_ambient"
	ControlAddress     = "address"
	ControlModel       = "model"
	ControlTemperature = "temperature"
	ControlHumidity    = "humidity"
	ControlBatteryLow  = "alarm_low_battery"
	ControlPreset      = "preset"
)
