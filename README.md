# wb-mqtt-noolite

Служба wb-mqtt-noolite позволяет интегрировать устройства Noolite (пульты, выключатели и датчики) в контроллер автоматизации [Wiren Board](http://contactless.ru/).  
Это достигается за счет соблюдения [Wiren Board MQTT Conventions](https://github.com/contactless/homeui/blob/master/conventions.md) и использования адаптера MTRF-64, с помощью которого реализуется работа с устройствами Noolite(-F).  
_Работа службы протестирована с адаптером [MTRF-64-USB-A](https://noo.by/adapter-mtrf-64-usb-a.html)._ 

**Функциональные возможности:**
* Прием управляющих команд от Wiren Board и передача их на исполнение устройствам Noolite
* Взаимодействие с адаптером NooLite MTRF-64, обработка ответов от устройств и публикация обновления статусов в MQTT.
* Опрос устройств Noolite(-F) по расписанию и отправка результатов в Wiren Board

### Сборка под Wirenboard (armv5tejl) 

```shell
# armv5tejl
export GOARCH=arm
export GOARM=5
go build -ldflags="-s -w" -o wb-mqtt-noolite-arm ../main.go
```

## Конфигурирование и запуск службы
Служба wb-mqtt-noolite конфигурируется с помощью конфигурационного файла в json формате. 
Основные параметры, которые необходимо передать: последовательный порт адаптера Noolite, реквизиты доступа к MQTT брокеру.
Также в секции device_config необходимо описать пути до отдельных JSON файлов шаблонов, описывающих модели устройств и самих устройства.
В каталоге templates находятся шаблоны. По мере развития проекта список будет пополняться.
В директории example приведен пример конфигурационного файла и списка устройств.

**Параметры командной строки:** -- config или -c с указанием пути до файла конфигурации службы

Пример конфигурационного файла:
````json
{
  "serial_port": "/dev/ttyUSB0",
  "timezone": "Europe/Moscow",
  "loglevel": "info",
  "mqtt": {
    "host": "127.0.0.1",
    "port": 1883,
    "username": "",
    "password": ""
  },
  "device_config": {
    "templates": "./templates/templates.json",
    "devices": "./example/devices.json"
  }
}
````

### Список устройств (devices.json)
Для работы службы требуется информация об устройствах, которые присутствуют в системе и описание шаблона модели устройства.   
**Пример списка устройств:**
```json
[
  {
    "name": "Теплый пол",
    "noolite_type": "txf",
    "ch": 1,
    "template": "srf-1-3000-t"
  }
]
```
_**Здесь**:_
- **name**: Наименование устройства, будет отображаться в интерфейсе
- **noolite_type**: тип Noolite (TX, TX-F, RX, RX-F)
- **ch**: Канал, 1-63
- **address**: Может указываться адрес, но не обязательно
- **template**: Шаблон устройства

###  Шаблоны устройств (templates.json)
Шаблон устройств предоставляет информацию о моделях поддерживаемых устройств, их элементов управления и информацию о командах Noolite, которые будут передаваться в устройство.
Для выполнения определенных команд на регулярной основе реализована поддержка выполнения команд по расписанию. Расписание задается в формате crontab
Таким образом, шаблон каждого из устройств описывается как наименование устройства и набор элементов управления.
Элементы управления имеют:
- **name**: имя топика
- **type**: тип см. [Wiren Board MQTT Conventions](https://github.com/contactless/homeui/blob/master/conventions.md)
- **order**: порядок сортировки
- **readonly**: определяет будет ли возможность со стороны UI возможность изменять его значение
- **get_command**: команда, которая будет выполняться по расписанию, например: _"ReadState 0"_
- **set_command**: при поступлении команды на установку данного контрола - будет отправляться соответствующая Noolite команда
- **polling**: выполнять команду по расписанию?
- **polling_cron**: расписание выполнения команды указанной в get_command
- **min**, **max**: для типа range, минимальное и максимальное принимаемое значение
- **units**: единицы измерения
- **precision**: точность
Пример:
```json
{
  "templates": [
    {
      "name": "srf-1-3000-t",
      "controls": [
        {
          "name": "on",
          "type": "switch",
          "order": 1,
          "initial_value": "0",
          "readonly": false,
          "onstart_command": "",
          "get_command": "",
          "set_command": "SetSwitch"
        },
        {
          "name": "status",
          "type": "switch",
          "order": 6,
          "initial_value": "0",
          "readonly": true
        },
        {
          "name": "setting",
          "type": "range",
          "order": 2,
          "min": 5,
          "max": 30,
          "initial_value": "23",
          "onstart_command": "",
          "get_command": "",
          "set_command": "SetTemperature"
        },
        {
          "name": "value",
          "type": "temperature",
          "order": 3,
          "readonly": true,
          "initial_value": "",
          "onstart_command": "",
          "get_command": "ReadState 0",
          "set_command": "",
          "polling": true,
          "polling_cron": "*/1 * * * *"
        },
        {
          "name": "model",
          "type": "text",
          "order": 4,
          "readonly": true
        },
        {
          "name": "address",
          "type": "text",
          "order": 5,
          "readonly": true
        }
      ]
    }
  ]
}
```
#### Поддерживаемые команды
У каждого элемента управления указываются команды, которые преобразуют MQTT сообщение в Noolite команду и обратно.\
**Поддерживаемые команды описаны ниже:**

|Команда|Описание и параметры|Пример|
|---|---|---|
|SetOn|Включить, не имеет параметров. |SetOn|
|SetOff|Выключить, не имеет параметров.|SetOff|
|SetSwitch|Включить или выключить, передается 1 или 0 соответственно.|SetSwitch|
|SetTemperature|Установить температуру (для srf-1-3000-t).|SetTemperature|
|ReadState|Запросить статус, параметр fmt (см.документацию на MTRF)|ReadState 0|

*Список команд будет расширяться.
