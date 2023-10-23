package connector

import (
	_ "net/http/pprof"

	. "aoanima.ru/logger"
	"github.com/google/uuid"
)

type ТипОтвета int

const (
	AjaxHTML ТипОтвета = iota
	AjaxJSON
	HTTP
	HTTPJson
)

type ТипЗапроса int

const (
	GET ТипЗапроса = iota
	POST
	AJAX
	AJAXPost
)

type Ответ struct {
	Сервис          []byte // Имя сервиса который отправляет ответ
	ИдКлиента       []byte // Копируется из запроса
	УИДЗапроса      string // Копируется из запроса
	Ответ           []byte // Ответ в бинарном формате
	СледующийСервис []byte // елси нужно обработать ответ в ещё в каком то сервисе, то сообщим об этом synqTCP передав имя сервиса в который нуно отправить Ответ от текущего сервса в поле СледующийСервис
}

type Отпечаток struct {
	Сервис   string
	Маршруты map[string]*СтруктураМаршрута
}
type СтруктураМаршрута struct {
	Запрос map[string]interface{} // описывает  данные которые нужны для обработки маршрута
	Ответ  map[string]interface{} // описывает формат в котором вернёт данные
}

type Сообщение struct {
	Сервис    []byte
	Запрос    *Запрос
	Ответ     *Ответ
	ИдКлиента uuid.UUID

	ТокенКлиента []byte // JWT сериализованный
}
type Запрос struct {
	ТипОтвета     ТипОтвета
	ТипЗапроса    ТипЗапроса
	СтрокаЗапроса string // url Path Query
	Форма         map[string][]string
	Файл          string
}

var (
	ПортДляОтправкиСообщений  = "81"
	ПортДляПолученияСообщений = "82"
)

// ПортИсходящихСообщений, ПортВходящихСообщений указывается те порты которые были исопльзованы в synqTCP сервер
// ПортДляОтправкиСообщений - соответсвует ВходящийПорт(synqTCP) - в этот порт серввис отправлят сообщения в synqTCP
// ПортДляПолученияСообщений - соответсвует ИсходящийПорт(synqTCP) -  из этого порта сервысы получают соощения из synqTCP
func ИнициализацияСервиса(
	адрес string,
	ПортИсходящихСообщений string,
	ПортВходящихСообщений string,
	отпечатокСервиса Отпечаток) (chan []byte,
	chan []byte) {
	Инфо(" ИнициализацияСервисов %+v \n", отпечатокСервиса)
	каналПолученияСообщений := make(chan []byte, 10)
	каналОтправкиСообщений := make(chan []byte, 10)
	//
	go ПодключитсяКСерверуДляПолученияСообщений(каналПолученияСообщений, адрес, ПортВходящихСообщений, отпечатокСервиса)
	go ПодключитьсяКСерверуДляОтправкиСообщений(каналОтправкиСообщений, адрес, ПортИсходящихСообщений)
	return каналПолученияСообщений, каналОтправкиСообщений

}

// Сервис := Отпечаток{
// 	Сервис: "Каталог",
// 	Маршруты: map[string]*СтруктураМаршрута{

// 	 "/": {
// 		Запрос: {
// 			"ТипЗпроса": "int", // в заивисмости от типа запроса например ajax или обычный request будет возвращён ответ...
// 			"Строка": "string", // url Path Query
// 			"Форма": "map[string][]string",
// 			"Файл":   "string",
// 		},
// 		Ответ:	{
// 			"HTML": "string",
// 			"JSON": "string",
// 			},
// 		} ,
// 		"catalog": {
// 		Запрос: {
// 			"ТипЗпроса": "int", // в заивисмости от типа запроса например ajax или обычный request будет возвращён ответ...
// 			"Строка": "string", // url Path Query
// 			"Форма": "map[string][]string",
// 			"Файл":   "string",
// 		},
// 		Ответ:	{
// 			"HTML": "string",
// 			"JSON": "string",
// 			},
// 		} ,
// 	},
// }
