package device

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
)

var ErrNoTemplateFound = errors.New("template not found")

// Templates Список шаблонов
type Templates struct {
	Templates []Template `json:"templates,omitempty"`
}

func (r *Templates) FindTemplateByName(name string) ([]*Control, error) {
	for _, template := range r.Templates {
		if template.Name == name {
			return template.Controls, nil
		}
	}
	return nil, ErrNoTemplateFound
}

// Template Правило. Описывает модель устройства и его органы управления.
type Template struct {
	Name     string     `json:"name,omitempty"`
	Controls []*Control `json:"controls"`
}

// NewTemplateEngine Создает новый объект со списками правил устройств и загружает их из указанного файла
func NewTemplateEngine(path string) (*Templates, error) {
	jsonFile, err := os.Open(path)

	if err != nil {
		return nil, err
	}

	defer jsonFile.Close()

	readJson, err := ioutil.ReadAll(jsonFile)

	if err != nil {
		return nil, err
	}

	var templates Templates
	err = json.Unmarshal(readJson, &templates)
	if err != nil {
		return nil, err
	}
	return &templates, nil
}
