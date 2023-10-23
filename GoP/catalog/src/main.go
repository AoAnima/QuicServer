package main

import (
	"net/http"

	_ "net/http/pprof"

	connector "aoanima.ru/connector/src"
	. "aoanima.ru/logger"
)

// type ОтветКлиенту struct {
// 	Сервис    []byte
// 	Ответ     []byte
// 	ИдКлиента []byte
// }

// type ЗапросКлиента struct {
// 	Запрос    ЗапросОтКлиента
// 	ИдКлиента uuid.UUID
// 	ТокенКлиента []byte // JWT сериализованный
// }
// type ЗапросОтКлиента struct {
// 	Строка string
// 	Форма  map[string][]string
// 	Файл   string
// }

func main() {
	go func() {
		http.ListenAndServe("localhost:6061", nil)
	}()
	// Вероятно нужно откуда то получить список Сервисов с которомы предстоит общаться
	//  Или !!!! ОбработчикВходящихСообщений

	Сервис := connector.Отпечаток{
		Сервис: "Каталог",
		Маршруты: map[string]*connector.СтруктураМаршрута{
			"/": {
				Запрос: map[string]interface{}{
					"ТипЗпроса":     "int",    // в заивисмости от типа запроса например ajax или обычный request будет возвращён ответ...
					"СтрокаЗапроса": "string", // url Path Query
					"Форма":         "map[string][]string",
					"Файл":          "string",
				},
				Ответ: map[string]interface{}{
					"HTML": "string",
					"JSON": "string",
				},
			},
			"catalog": {
				Запрос: map[string]interface{}{
					"ТипЗпроса":     "int",    // в заивисмости от типа запроса например ajax или обычный request будет возвращён ответ...
					"СтрокаЗапроса": "string", // url Path Query
					"Форма":         "map[string][]string",
					"Файл":          "string",
				},
				Ответ: map[string]interface{}{
					"HTML": "string",
					"JSON": "string",
				},
			},
			"product": {
				Запрос: map[string]interface{}{
					"ТипЗпроса":      "int",    // в заивисмости от типа запроса например ajax или обычный request будет возвращён ответ...
					"СтрокаЗапросаы": "string", // url Path Query
					"Форма":          "map[string][]string",
					"Файл":           "string",
				},
				Ответ: map[string]interface{}{
					"HTML": "string",
					"JSON": "string",
				},
			},
		},
	}

	каналПолученияСообщений, каналОтправкиСообщений := connector.ИнициализацияСервиса("localhost", connector.ПортДляОтправкиСообщений, connector.ПортДляПолученияСообщений, Сервис)

	// читаем входящий запрос от клиента, обрабатывает,, и отправляем ответ обратно
	for ВходящееСообщение := range каналПолученияСообщений {
		go ОбработатьЗапрос(ВходящееСообщение, каналОтправкиСообщений)
	}

}

func ОбработатьЗапрос(ВходящееСообщение []byte, каналОтправкиСообщений chan<- []byte) {

	Инфо(" ОбработатьЗапрос %+v \n", ВходящееСообщение)

	Ответ, err := Обработать(ДекодироватьПакет(ВходящееСообщение))
	if err != nil {
		Ошибка("  %+v \n", err)
	}
	// TODO: Реализум логику обработки запроса от клиента, и генерацию ответа

	// ОтветКлиенту := ВходящееСообщение
	каналОтправкиСообщений <- Ответ
}
