package connectorQUIC

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"encoding/binary"
	"encoding/pem"
	"fmt"
	"math/big"
	_ "net/http/pprof"
	"net/url"
	"time"

	. "aoanima.ru/logger"
	"github.com/dgryski/go-metro"
	"github.com/google/uuid"
	quic "github.com/quic-go/quic-go"
)

// type Глюк struct {
// 	Текст     string
// 	КодОшибки int
// }

// func (e *Глюк) Error() string {
// 	return e.Текст
// }
// func (e *Глюк) Код() int {
// 	return e.КодОшибки
// }

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

	go ПодключитсяКСерверуДляПолученияСообщений(каналПолученияСообщений, адрес, ПортВходящихСообщений, отпечатокСервиса)
	go ПодключитьсяКСерверуДляОтправкиСообщений(каналОтправкиСообщений, адрес, ПортИсходящихСообщений)

	return каналПолученияСообщений, каналОтправкиСообщений
}

type ОчередьПотоков struct {
	потоки chan quic.Stream
}

func НоваяОчередьПотоков(размер int) *ОчередьПотоков {
	return &ОчередьПотоков{
		потоки: make(chan quic.Stream, размер),
	}
}
func (о *ОчередьПотоков) Взять(поток quic.Stream) (quic.Stream, error) {
	select {
	case поток := <-о.потоки:
		return поток, nil
	default:
		return nil, fmt.Errorf("Нет свободных потоков")
	}

}

func (о *ОчередьПотоков) Вернуть(поток quic.Stream) {
	select {
	case о.потоки <- поток:
	default:
		// Если канал полон, просто закрываем поток
		поток.Close()
	}
}

func УстановитьКвикСоедиенение() {

}

func Сервер() {
	listener, err := quic.ListenAddr("localhost:4242", генерироватьТлсКонфиг(), nil)
	if err != nil {
		Ошибка(" %+v ", err)
	}

	for {
		сессия, err := listener.Accept(context.Background())
		if err != nil {
			Ошибка(" %+v ", err)
		}

		go ОбработчикСессии(сессия)
	}
}

func ОбработчикСессии(сессия quic.Connection) {
	for {
		поток, err := сессия.AcceptStream(context.Background())
		if err != nil {
			Ошибка(" %+v ", err)
		}
		go ЧитатьСообщения(поток, каналПолученияСообщений)
	}

}

func ЧитатьСообщения(поток quic.Stream, каналПолученияСообщений chan []byte) {

	длинаСообщения := make([]byte, 4)
	var прочитаноБайт int
	var err error
	for {
		прочитаноБайт, err = поток.Read(длинаСообщения)
		Инфо(" длинаСообщения %+v , прочитаноБайт %+v \n", длинаСообщения, прочитаноБайт)

		if err != nil {
			Ошибка(" прочитаноБайт %+v  err %+v \n", прочитаноБайт, err)
			break
		}

		// получаем число байткоторое нужно прочитать
		длинаДанных := binary.LittleEndian.Uint32(длинаСообщения)

		Инфо(" длинаДанных  %+v \n", длинаДанных)
		Инфо(" длинаСообщения %+v ,  \n прочитаноБайт %+v ,  \n длинаДанных %+v \n", длинаСообщения,
			прочитаноБайт, длинаДанных)

		//читаем количество байт = длинаСообщения
		// var запросКлиента ЗапросКлиента
		пакетОтвета := make([]byte, длинаДанных)
		прочитаноБайт, err = сервер.Read(пакетОтвета)
		if err != nil {
			Ошибка("Ошибка при десериализации структуры: %+v ", err)
		}
		if длинаДанных != uint32(прочитаноБайт) {
			Ошибка("Количество прочитаных байт не ранво длине данных :\n длинаДанных %+v  <> прочитаноБайт %+v ", длинаДанных, прочитаноБайт)
		}

		каналПолученияСообщений <- пакетОтвета

	}

}

func Клиент() {

}
func генерироватьТлсКонфиг() *tls.Config {
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		panic(err)
	}
	template := x509.Certificate{SerialNumber: big.NewInt(1)}
	certDER, err := x509.CreateCertificate(rand.Reader, &template, &template, &key.PublicKey, key)
	if err != nil {
		panic(err)
	}
	keyPEM := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(key)})
	certPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: certDER})

	tlsCert, err := tls.X509KeyPair(certPEM, keyPEM)
	if err != nil {
		panic(err)
	}
	return &tls.Config{Certificates: []tls.Certificate{tlsCert}, InsecureSkipVerify: true}
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
