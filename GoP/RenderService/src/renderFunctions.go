package main

import (
	"encoding/json"
	"html/template"
	"reflect"
	"strconv"
	"time"

	. "aoanima.ru/Logger"
)

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
		"вJSON": func(data interface{}) interface{} {
			//Инфо("  reflect.TypeOf %+v",  reflect.TypeOf(data).Kind())
			//if reflect.TypeOf(data).Kind().String() == "map" {
			if data == nil {
				return nil
			}
			r, e := json.Marshal(data)
			Инфо("вJSON data %+v r %+v", data, r)
			if e != nil {
				Ошибка("  %+v", e)
			}
			return string(r)
			//}
		},
		"Пусто": func(data interface{}) interface{} {

			if data == nil {
				return "Нет данных"
			} else {
				switch data.(type) {
				case int64:
					return strconv.Itoa(int(data.(float64)))
				case map[string]interface{}:
					return data
				}

			}

			return data
		},
		"TypeOf": func(n interface{}) string {
			return reflect.TypeOf(n).String()
		},
		"Тип": func(n interface{}) string {
			return reflect.TypeOf(n).String()
		},
		"Строка": func(Строки ...string) string {
			Результат := ""
			for _, Стр := range Строки {
				Результат = Результат + Стр
			}
			return Результат
		},
	}
}
