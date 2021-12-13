package noolite

// Магические константы Noolite
const (
	StartRequest  uint8 = 171
	StopRequest   uint8 = 172
	StartResponse uint8 = 173
	StopResponse  uint8 = 174
)

// Режим работы адаптера
const (
	// ModeNooliteTX режим Noolite TX
	ModeNooliteTX uint8 = 0
	// ModeNooliteRX режим Noolite RX
	ModeNooliteRX uint8 = 1
	// ModeNooliteFTX режим Noolite-F TX
	ModeNooliteFTX uint8 = 2
	// ModeNooliteFRX режим Noolite-F RX
	ModeNooliteFRX uint8 = 3
	// ModeNooliteService Сервисный режим работы с Noolite-F
	ModeNooliteService uint8 = 4
	// ModeNooliteUpgrade Режим обновления ПО Noolite-F
	ModeNooliteUpgrade uint8 = 5
)

// Поле CTR в Request
const (
	// CtrRequestSendCommand Передать команду
	CtrRequestSendCommand uint8 = 0
	// CtrRequestSendBroadcastCommand Передать широковещательную команду всем устройстам на канале
	CtrRequestSendBroadcastCommand uint8 = 1
	// CtrRequestReadResponse Считать ответ из приемного буфера
	CtrRequestReadResponse uint8 = 2
	// CtrRequestBindEnable Включить привязку
	CtrRequestBindEnable uint8 = 3
	// CtrRequestBindDisable Выключить привязку
	CtrRequestBindDisable uint8 = 4
	// CtrRequestClearChannel Очистить ячейку (канал)
	CtrRequestClearChannel uint8 = 5
	// CtrRequestClearAllChannels  Очистить память (все каналы)
	CtrRequestClearAllChannels uint8 = 6
	// CtrRequestUnbindAddress Отвязать адрес от канала
	CtrRequestUnbindAddress uint8 = 7
	// CtrRequestSendCommandToAddress Передать команду по указанному адресу
	CtrRequestSendCommandToAddress uint8 = 8
	// CtrRequestSendCommandByAddress Установка настроек
	CtrRequestSendCommandByAddress uint8 = 9
)

// Поле CTR в Response
const (
	// CtrResponseSuccess Ответ: Команда выполнена
	CtrResponseSuccess uint8 = 0
	// CtrResponseNoResponse Ответ: Нет ответа от блока
	CtrResponseNoResponse uint8 = 1
	// CtrResponseError Ответ: Ошибка во время выполнения
	CtrResponseError uint8 = 2
	// CtrResponseBindSuccess Ответ: Привязка выполнена
	CtrResponseBindSuccess uint8 = 3
)

// Типы устройств в Response
const (
	// DeviceMTRF64 радиоконтроллер MTRF-64
	DeviceMTRF64 uint8 = 0
	// DeviceSLF300 реле SLF-1-300
	DeviceSLF300 uint8 = 1
	// DeviceSRF101000 реле из блока SRF-10-1000
	DeviceSRF101000 uint8 = 2
	// DeviceSRF3000Rozetka реле SRF-1-3000 (розетка)
	DeviceSRF3000Rozetka uint8 = 3
	// DeviceSRF3000Podrozetnik реле SRF-1-3000 (для подрозетника)
	DeviceSRF3000Podrozetnik uint8 = 4
	// DeviceSUF300 диммер SUF-1-300
	DeviceSUF300 uint8 = 5
	// DeviceSRF3000T терморегулятор SRF-1-3000-T
	DeviceSRF3000T uint8 = 6
	// DeviceSRF1000R блок роллет SRF-1-1000-R
	DeviceSRF1000R uint8 = 7
	// DevicePT111 датчик температуры и влажности
	DevicePT111 uint8 = 254
)

// Команды
const (
	// CmdOff Выключить нагрузку
	CmdOff uint8 = 0
	// CmdBrightDown Запускает плавное понижение яркости
	CmdBrightDown uint8 = 1
	// CmdOn Включить нагрузку
	CmdOn uint8 = 2
	// CmdBrightUp Запускает плавное повышение яркости
	CmdBrightUp uint8 = 3
	// CmdSwitch Включает или выключает нагрузку
	CmdSwitch uint8 = 4
	// CmdBrightBack Запускает плавное изменение яркости в обратном направлении
	CmdBrightBack uint8 = 5
	// CmdSetBrightness Установить заданную в расширении команды яркость (количество данных зависит от устройства).
	CmdSetBrightness uint8 = 6
	// CmdLoadPreset Вызвать записанный сценарий
	CmdLoadPreset uint8 = 7
	// CmdSavePreset Записать сценарий в память
	CmdSavePreset uint8 = 8
	// CmdUnbind Запускает процедуру стирания адреса управляющего устройства из памяти исполнительного
	CmdUnbind uint8 = 9
	// CmdStopReg Прекращает действие команд
	//  Bright_Down, Bright_Up, Bright_Back.
	CmdStopReg uint8 = 10
	// CmdBrightStepDown Понизить яркость на шаг. При отсутствии поля
	//данных увеличивает отсечку на 64 мкс, при
	//наличии поля данных на величину в микросекундах (0 соответствует 256 мкс).
	CmdBrightStepDown uint8 = 11
	// CmdBrightStepUp Повысить яркость на шаг. При отсутствии поля данных увеличивает отсечку на 64 мкс, при
	//наличии поля данных на величину в микросекундах (0 соответствует 256 мкс).
	CmdBrightStepUp uint8 = 12
	// CmdBrightReg Запускает плавное изменение яркости с направлением и скоростью, заданными в расширении
	CmdBrightReg uint8 = 13
	// CmdBind Сообщает исполнительному устройству, что управляющее хочет активировать режим привязки.
	//При привязке также передаётся тип устройства в данных.
	CmdBind uint8 = 15
	// CmdRollColour Запускает плавное изменение цвета в RGB-контроллере по радуге
	CmdRollColour uint8 = 16
	// CmdSwitchColour Переключение между стандартными цветами в RGB-контроллере
	CmdSwitchColour uint8 = 17
	// CmdSwitchMode Переключение между режимами RGB-контроллера
	CmdSwitchMode uint8 = 18
	// CmdSpeedModeBack Запускает изменение скорости работы режимов RGB-контроллера в обратном направлении.
	CmdSpeedModeBack uint8 = 19
	// CmdBatteryLow У устройства, которое передало данную команду, разрядился элемент питания
	CmdBatteryLow uint8 = 20
	// CmdSensTempHumi Передает данные о температуре, влажности и состоянии элементов.
	CmdSensTempHumi uint8 = 21
	// CmdTemporaryOn Включить свет на заданное время. Время в 5-и секундных тактах передается в расширении
	CmdTemporaryOn uint8 = 25
	// CmdModes Установка режимов работы исполнительного устройства
	CmdModes uint8 = 26
	// CmdReadState Получение состояния исполнительного устройства
	CmdReadState uint8 = 128
	// CmdWriteState Установка состояния исполнительного устройства.
	CmdWriteState uint8 = 129
	// CmdSendState Ответ от исполнительного устройства
	CmdSendState uint8 = 130
	// CmdService Включение сервисного режима на заранее привязанном устройстве
	CmdService uint8 = 131
	// CmdClearMemory Очистка памяти устройства nooLite.
	//Для выполнения команды используется ключ 170-85-170-85 (записывается в поле данных D0…D3).
	CmdClearMemory uint8 = 132
)

// Термостат. Константы управления режимом и датчиком
const (
	// ModeManualFloorSensor Режим работы ручной по датчику пола
	ModeManualFloorSensor uint8 = 1
	// ModeManualAirSensor Режим работы ручной по датчику воздуха
	ModeManualAirSensor uint8 = 9
	// ModeManualWirelessSensor Режим работы ручной по беспроводному датчику
	ModeManualWirelessSensor uint8 = 3
	// ModeCalendarFloorSensor Режим работы по календарю и датчику пола
	ModeCalendarFloorSensor uint8 = 0
	// ModeCalendarAirSensor Режим работы по календарю и датчику воздуха
	ModeCalendarAirSensor uint8 = 8
)

const (
	// FmtMain Основная информация о блоке (ro)
	FmtMain uint8 = 0
	// FmtSetAndDontOn Установить и не включать
	FmtSetAndDontOn = 1
	// FmtSetAndOn Установить и включить
	FmtSetAndOn = 2
	// FmtSettings Настройки блока (rw)
	FmtSettings uint8 = 16
	// FmtExternalLoadInfo Получение данных с выхода для нагрузки устройства
	FmtExternalLoadInfo uint8 = 19
	// FmtSensTempHumi Данные в формате температура и влажность
	FmtSensTempHumi uint8 = 7
)
