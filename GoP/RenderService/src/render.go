package main

import (
	"bytes"
	"html/template"
	"log"
	"path/filepath"
	"strconv"
	"time"

	. "aoanima.ru/ConnQuic"
	. "aoanima.ru/Logger"
)

var СырыеШаблоны *template.Template

func ПостроитьДеревоШаблонов(сообщение *Сообщение) {

	// content := template.New("content")
	catalog := СырыеШаблоны.Lookup("catalog")
	ШаблонДляРендера, err := СырыеШаблоны.Clone()
	if err != nil {
		Ошибка(" Ошибка клоинровании сырых шаблонов %+v \n", err)
	}
	ШаблонДляРендера, err = ШаблонДляРендера.AddParseTree("content", catalog.Tree)
	if err != nil {
		Ошибка("  %+v \n", err)
	}

	Html := new(bytes.Buffer)
	Данные := map[string]map[string]string{
		"Кнопка": {
			"Класс": "success",
			"Тип":   "submit",
		},
	}
	if errs := ШаблонДляРендера.ExecuteTemplate(Html, "index", Данные); errs != nil {
		Ошибка("%+v\n", errs)

	}
	Инфо("  %+v \n", Html.String())

}

func Рендер(ИмяШаблона string, Данные interface{}) ([]byte, error) {
	Html := new(bytes.Buffer)

	ШаблонДляРендера, err := СырыеШаблоны.Clone()
	if err != nil {
		Ошибка(" Ошибка клоинровании сырых шаблонов %+v \n", err)
		return nil, err
	}

	if errs := ШаблонДляРендера.ExecuteTemplate(Html, ИмяШаблона, Данные); errs != nil {
		Ошибка("%+v\n", errs)
		return nil, errs
	}

	return Html.Bytes(), nil
}

func ПарсингШаблонов() {
	// "pattern": "../www/tpl/*/*.html",
	var errParseGlob error
	Инфо(" Конфиг.КаталогШаблонов  %+v \n", Конфиг.КаталогШаблонов+"*/*.html")

	filenames, err := filepath.Glob(Конфиг.КаталогШаблонов + "*/*.html")
	if err != nil {
		Ошибка(" Ошибка парсинга каталога с шаблонами HTML %+v\n", err)
	}
	Инфо(" filenames %+v \n", filenames)
	СырыеШаблоны = template.Must(template.New("raw").Funcs(РендерФункции()).ParseGlob(Конфиг.КаталогШаблонов + "*/*.html"))
	// СырыеШаблоны = template.Must(template.New("index").Funcs(РендерФункции()).ParseGlob(Конфиг.КаталогШаблонов + "*/*.html"))

	log.Printf("СырыеШаблоны %+v\n", СырыеШаблоны)
	if errParseGlob != nil {
		Ошибка("Ошибка парсинга каталога с шаблонами HTML %+v\n", errParseGlob)
	}

}

func РендерФункции() template.FuncMap {
	return template.FuncMap{
		"ВременнаяМетка": func() int64 {
			return time.Now().Unix()
		},
		"вСтроку": func(data interface{}) string {
			return data.(string)
		},
		"вЦелое": func(data float64) int {
			return int(data)
		},
		"ЧислоВСтроку": func(data float64) string {
			return strconv.Itoa(int(data))
		},
		"СтрокуВЧисло": func(data string) (int, error) {
			число, err := strconv.Atoi(data)
			if err != nil {
				Ошибка("  %+v \n", err)
				return число, err
			}
			return число, nil
		},
	}
}
