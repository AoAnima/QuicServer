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

func генераторИд() string {
	timestamp := time.Now().UnixNano() / int64(time.Millisecond)
	return strconv.FormatInt(timestamp, 36)
}

func Сумма(числа ...interface{}) float64 {
	var сумма float64
	for _, число := range числа {
		switch v := число.(type) {
		case int:
			сумма += float64(v)
		case int8:
			сумма += float64(v)
		case int16:
			сумма += float64(v)
		case int32:
			сумма += float64(v)
		case int64:
			сумма += float64(v)
		case uint:
			сумма += float64(v)
		case uint8:
			сумма += float64(v)
		case uint16:
			сумма += float64(v)
		case uint32:
			сумма += float64(v)
		case uint64:
			сумма += float64(v)
		case float32:
			сумма += float64(v)
		case float64:
			сумма += v
		default:
			return 0
		}
	}
	return сумма
}

func РендерФункции() template.FuncMap {
	return template.FuncMap{
		"jsStr": func(s string) template.JSStr {
			return template.JSStr(s)
		},
		"JS": func(s string) template.JS {
			return template.JS(s)
		},
		"Сумма": Сумма,
		"ИД": func() string {
			timestamp := time.Now().UnixNano() / int64(time.Millisecond)
			return strconv.FormatInt(timestamp, 36)
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
				// Инфо(" %+v %#T \n", КлючЗначение[i+1], КлючЗначение[i+1])

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
