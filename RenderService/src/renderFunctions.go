package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"reflect"
	"strconv"
	"time"

	. "aoanima.ru/Logger"
)

func РендерФункции() template.FuncMap {
	return template.FuncMap{
		"JS": func(s string) template.JS {
			return template.JS(s)
		},
		"Слайс": func(аргументы ...interface{}) []interface{} {
			return аргументы
		},
		"Мап": func(КлючЗначение ...interface{}) map[string]interface{} {
			if len(КлючЗначение)%2 != 0 {
				return map[string]interface{}{
					"ошибка": "не чётное количество аргументов для map",
				}
			}
			Карта := make(map[string]interface{}, len(КлючЗначение)/2)
			for i := 0; i < len(КлючЗначение); i += 2 {
				key, ok := КлючЗначение[i].(string)
				if !ok {
					return map[string]interface{}{
						"ошибка": fmt.Sprintf("Ключ %+v должен быть строкой", КлючЗначение[i]),
					}
				}
				Карта[key] = КлючЗначение[i+1]
			}
			return Карта
		},
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
