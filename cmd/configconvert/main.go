package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"wb-noolite-mtrf/device"
)

// configconvert Конвертирует templates.json версии 1.x (до релиза 2.0) в новый формат:
// заменяет устаревшие типы величин и идентификаторы каналов, добавляет заготовки title (en/ru).
// devices.json конвертации не требует - формат этого файла в 2.0 не менялся.
func main() {
	var inPath, outPath string
	flag.StringVar(&inPath, "in", "", "Путь к старому templates.json (версия 1.x)")
	flag.StringVar(&outPath, "out", "", "Путь для сохранения результата; по умолчанию перезаписывает -in")
	flag.Parse()

	if inPath == "" {
		fmt.Fprintln(os.Stderr, "Использование: configconvert -in <старый templates.json> [-out <новый templates.json>]")
		fmt.Fprintln(os.Stderr, "devices.json конвертировать не нужно - его формат в версии 2.0 не изменился.")
		os.Exit(1)
	}
	if outPath == "" {
		outPath = inPath
	}

	raw, err := os.ReadFile(inPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Не удалось прочитать %s: %s\n", inPath, err)
		os.Exit(1)
	}

	var templates device.Templates
	if err := json.Unmarshal(raw, &templates); err != nil {
		fmt.Fprintf(os.Stderr, "Не удалось разобрать %s: %s\n", inPath, err)
		os.Exit(1)
	}

	changes := device.MigrateTemplates(&templates)
	if len(changes) == 0 {
		fmt.Println("Изменений не потребовалось: шаблон уже соответствует формату 2.0.")
	} else {
		fmt.Println("Внесены изменения:")
		for _, c := range changes {
			fmt.Println(" - " + c)
		}
		fmt.Println()
		fmt.Println("Проверьте и переведите добавленные временные title (en/ru) вручную,")
		fmt.Println("а также пути к шаблонам в конфигурации wb-mqtt.noolite.json при необходимости.")
	}

	out, err := json.MarshalIndent(&templates, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Не удалось сформировать результат: %s\n", err)
		os.Exit(1)
	}

	if err := os.WriteFile(outPath, append(out, '\n'), 0644); err != nil {
		fmt.Fprintf(os.Stderr, "Не удалось записать %s: %s\n", outPath, err)
		os.Exit(1)
	}

	fmt.Printf("\nШаблоны сохранены в %s\n", outPath)
}
