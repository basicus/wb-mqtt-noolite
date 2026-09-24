package device

import "fmt"

// deprecatedTypeUnits Единицы измерения ("units") для замены устаревших специфичных типов величин
// на общий WbControlTypeGeneric, см. Wiren Board MQTT Conventions (Control Types).
var deprecatedTypeUnits = map[ControlType]string{
	WbControlTypeTemperature:      "deg C",
	WbControlTypeRelHumidity:      "%",
	WbControlTypePower:            "W",
	WbControlTypePowerConsumption: "kWh",
}

// legacyControlRenames Переименования идентификаторов каналов из шаблонов версии 1.x в канон WB-STD-001.
var legacyControlRenames = map[string]string{
	"value":       "temperature",
	"setting":     "setpoint_t_ambient",
	"low_battery": "alarm_low_battery",
}

// MigrateTemplates Приводит шаблоны устройств версии 1.x к формату 2.0:
//   - заменяет устаревшие типы величин (temperature/rel_humidity/power/power_consumption)
//     на WbControlTypeGeneric с соответствующим units;
//   - переименовывает устаревшие идентификаторы каналов ("value", "setting", "low_battery")
//     в канонические, согласно WB-STD-001 (temperature, setpoint_t_ambient, alarm_low_battery);
//   - проставляет тип alarm индикаторам разряда батареи;
//   - добавляет временный title (en/ru), если он не был задан, чтобы шаблон не остался без
//     названий - его текст нужно перевести/поправить вручную после конвертации.
//
// Функция мутирует переданные templates и возвращает список внесенных изменений для лога.
// Разрешено применять повторно (идемпотентна): если шаблон уже соответствует формату 2.0,
// список изменений будет пустым.
func MigrateTemplates(templates *Templates) []string {
	var changes []string
	for i := range templates.Templates {
		template := &templates.Templates[i]
		for _, control := range template.Controls {
			oldName := control.Name
			oldType := control.Type

			if units, deprecated := deprecatedTypeUnits[control.Type]; deprecated {
				control.Type = WbControlTypeGeneric
				if control.Units == "" {
					control.Units = units
				}
			}

			if oldName == "low_battery" {
				control.Type = WbControlTypeAlarm
			}

			if newName, ok := legacyControlRenames[oldName]; ok {
				control.Name = newName
			}

			if control.Name != oldName || control.Type != oldType {
				changes = append(changes, fmt.Sprintf(
					"шаблон %q: канал %q -> %q, тип %q -> %q",
					template.Name, oldName, control.Name, oldType, control.Type))
			}

			if len(control.Title) == 0 {
				control.Title = Title{"en": control.Name, "ru": control.Name}
				changes = append(changes, fmt.Sprintf(
					"шаблон %q: каналу %q добавлен временный title - переведите вручную",
					template.Name, control.Name))
			}
		}

		if len(template.Title) == 0 {
			template.Title = Title{"en": template.Name, "ru": template.Name}
			changes = append(changes, fmt.Sprintf(
				"шаблон %q: добавлен временный title - переведите вручную", template.Name))
		}
	}
	return changes
}
