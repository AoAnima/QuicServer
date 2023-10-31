package ConnQuic

import (
	"encoding/binary"
	"fmt"
	"net/url"
	"time"

	. "aoanima.ru/logger"
	"github.com/dgryski/go-metro"
	"github.com/google/uuid"
	jsoniter "github.com/json-iterator/go"
)

func ДекодироватьПакет(пакет []byte) Сообщение {
	Инфо(" ДекодироватьПакет пакет %+s \n", пакет)

	// var запросОтКлиента = ЗапросКлиента{
	// 	Сервис:    []byte{},
	// 	Запрос:    &ЗапросОтКлиента{},
	// 	ИдКлиента: uuid.UUID{},
	// }
	var Сообщение Сообщение

	// TODO тут лишний парсинг, нужно получить только URL patch чтобы определить сервис, которому принадлежит запрос, потому nxj дальше весь запрос опять сериализуйется

	err := jsoniter.Unmarshal(пакет, &Сообщение)
	if err != nil {
		Ошибка("  %+v \n", err)
	}
	Инфо(" Сообщение входящее %+s \n", Сообщение)

	return Сообщение
}
func Кодировать(данныеДляКодирования interface{}) ([]byte, error) {

	b, err := jsoniter.Marshal(&данныеДляКодирования)
	if err != nil {
		Ошибка("  %+v \n", err)
		return nil, err
	}
	данные := make([]byte, len(b)+4)
	binary.LittleEndian.PutUint32(данные, uint32(len(b)))
	copy(данные[4:], b)
	return данные, nil

}

type ТипОтвета int

const (
	AjaxHTML ТипОтвета = iota
	AjaxJSON
	HTML
)

type ТипЗапроса int

const (
	GET ТипЗапроса = iota
	POST
	AJAX
	AJAXPost
)

type Отпечаток struct {
	Сервис   string
	Маршруты map[string]*СтруктураМаршрута
}
type СтруктураМаршрута struct {
	Запрос map[string]interface{} // описывает  данные которые нужны для обработки маршрута
	Ответ  map[string]interface{} // описывает формат в котором вернёт данные
}
type Сервис string
type ОтветСервиса struct {
	Сервис          Сервис // Имя сервиса который отправляет ответ
	УИДЗапроса      string // Копируется из запроса
	Данные          []byte // Ответ в бинарном формате
	ЗапросОбработан bool   // Признак того что запросы был получен и обработан соответсвуюбщим сервисом, в не зависимоти есть ли данные в ответе или нет, если данных нет, знаичт они не нужны... Выставляем в true в сеорвисе перед отправкой ответа
}

type Ответ map[Сервис]ОтветСервиса

type Сообщение struct {
	Сервис       Сервис // Имя Сервиса который шлёт Сообщение, каждый сервис пишет своё имя в не зависимости что это ответ или запрос
	Регистрация  bool
	Маршруты     map[string]*СтруктураМаршрута
	Запрос       Запрос
	Ответ        Ответ
	ИдКлиента    uuid.UUID
	УИДСообщения Уид    // ХЗ по логике каждый сервис должен вставлять сюбда своё УИД
	ТокенКлиента []byte // JWT сериализованный
}
type Уид string

type Запрос struct {
	ТипОтвета      ТипОтвета
	ТипЗапроса     ТипЗапроса
	СтрокаЗапроса  *url.URL // url Path Query
	МаршрутЗапроса string   // url Path Query
	Форма          map[string][]string
	Файл           string
	УИДЗапроса     Уид
}

var (
	ПортДляОтправкиСообщений  = "81"
	ПортДляПолученияСообщений = "82"
)

func УИДЗапроса(ИдКлиента *uuid.UUID, UrlPath []byte) Уид {
	return Уид(fmt.Sprintf("%+s.%+s.%+s", time.Now().Unix(), ИдКлиента, metro.Hash64(UrlPath, 0)))
}
