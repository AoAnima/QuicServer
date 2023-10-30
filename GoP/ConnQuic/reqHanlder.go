package ConnQuic

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/binary"
	"os"
	"time"

	jsoniter "github.com/json-iterator/go"

	. "aoanima.ru/logger"
	"github.com/google/uuid"
)

// каналПолученияСообщений - исползуется для получения запросов от клиента, в запросе от клиента передаётся канал в который нужно отправить ответ клиенту

// func ПодключитсяКСерверуДляПолученияСообщений(каналПолученияСообщений chan []byte, адрес string, порт string, отпечатокСервиса Отпечаток) {
// 	caCert, err := os.ReadFile("cert/ca.crt")

// 	if err != nil {
// 		Ошибка(" %s ", err)
// 	}

// 	caCertPool := x509.NewCertPool()
// 	ok := caCertPool.AppendCertsFromPEM(caCert)
// 	Инфо("Корневой сертфикат создан?  %v ", ok)

// 	cert, err := tls.LoadX509KeyPair("cert/client.crt", "cert/client.key")
// 	if err != nil {
// 		Ошибка(" %s", err)
// 	}

// 	tlsConfig := &tls.Config{
// 		RootCAs:      caCertPool,
// 		Certificates: []tls.Certificate{cert},
// 	}

// 	// Подключение к TCP-серверу с TLS на localhost:8080
// 	количествоПопыток := 500
// 	задержка := 1 * time.Second
// 	var сервер *tls.Conn
// 	var errDial error
// 	for попытка := 1; попытка <= количествоПопыток; попытка++ {
// 		сервер, errDial = tls.Dial("tcp", адрес+":"+порт, tlsConfig)
// 		if errDial != nil {
// 			Ошибка("  %+v \n", err)
// 			time.Sleep(задержка)
// 		} else {
// 			break
// 		}
// 	}
// 	go ЧитатьСообщения(сервер, каналПолученияСообщений)
// 	Рукопожатие(сервер, отпечатокСервиса)
// }

func Рукопожатие(сервер *tls.Conn, ОтпечатокСервиса Отпечаток) {
	// буфер := new(bytes.Buffer)
	// Запрос{
	// 	Сервис:    []byte("КлиентСервер"),
	// 	Запрос:    "🤝",
	// 	ИдКлиента: Уид(),
	// }

	// Инфо("  %+v %+v \n", "🤝", []byte("🤝"), len([]byte("🤝")))
	// binary.Write(буфер, binary.LittleEndian, [4]byte{240, 159, 164, 157}) // [4]byte{240, 159, 164, 157} = "🤝"

	// Будет описаывать какие данные в каком виде нужно присылать в запросах для конкретного маршрута для данного сервиса
	//например сервис КлиентСервер , имеет обработчик ОтветКЛиенту : ДляЭтого метода ему нужен ИдКлиента, и ответ в виде HTML строки или json

	// КлиентСервер := Отпечаток{
	// 	Сервис: "КлиентСервер",
	// 	Маршруты: map[string]map[string]interface{}{
	// 		"ОтветКлиенту": {
	// 			"HTML": "string",
	// 			"JSON": "string",
	// 		},
	// 		"catalog": {
	// 			"HTML": "string",
	// 			"JSON": "string",
	// 		},
	// 	},
	// }
	// КлиентСервер := Отпечаток{
	// 	Сервис: "КаталогСервис",
	// 	Маршруты: map[string]map[string]interface{}{
	// 		"catalog": map[string]interface{}{
	// 			"Запрос": "string",
	// 		}
	//
	// 	},
	// }

	данныеВОтправку, err := Кодировать(ОтпечатокСервиса)
	if err != nil {
		Ошибка("  %+v \n", err)
	}
	// binary.Write(буфер, binary.LittleEndian, int32(len([]byte("КлиентСервер"))))
	// binary.Write(буфер, binary.LittleEndian, []byte("КлиентСервер"))
	сервер.Write(данныеВОтправку)

}

// func ЧитатьСообщения(сервер *tls.Conn, каналПолученияСообщений chan []byte) {

// 	длинаСообщения := make([]byte, 4)
// 	var прочитаноБайт int
// 	var err error
// 	for {
// 		прочитаноБайт, err = сервер.Read(длинаСообщения)
// 		Инфо(" длинаСообщения %+v , прочитаноБайт %+v \n", длинаСообщения, прочитаноБайт)

// 		if err != nil {
// 			Ошибка(" прочитаноБайт %+v  err %+v \n", прочитаноБайт, err)
// 			break
// 		}

// 		// получаем число байткоторое нужно прочитать
// 		длинаДанных := binary.LittleEndian.Uint32(длинаСообщения)

// 		Инфо(" длинаДанных  %+v \n", длинаДанных)
// 		Инфо(" длинаСообщения %+v ,  \n прочитаноБайт %+v ,  \n длинаДанных %+v \n", длинаСообщения,
// 			прочитаноБайт, длинаДанных)

// 		//читаем количество байт = длинаСообщения
// 		// var запросКлиента ЗапросКлиента
// 		пакетОтвета := make([]byte, длинаДанных)
// 		прочитаноБайт, err = сервер.Read(пакетОтвета)
// 		if err != nil {
// 			Ошибка("Ошибка при десериализации структуры: %+v ", err)
// 		}
// 		if длинаДанных != uint32(прочитаноБайт) {
// 			Ошибка("Количество прочитаных байт не ранво длине данных :\n длинаДанных %+v  <> прочитаноБайт %+v ", длинаДанных, прочитаноБайт)
// 		}

// 		каналПолученияСообщений <- пакетОтвета

// 	}

// }

func ПодключитьсяКСерверуДляОтправкиСообщений(каналОтправкиОтветов chan []byte, адрес string, порт string) {
	caCert, err := os.ReadFile("cert/ca.crt")

	if err != nil {
		Ошибка(" %s ", err)
	}

	caCertPool := x509.NewCertPool()
	ok := caCertPool.AppendCertsFromPEM(caCert)
	Инфо("Корневой сертфикат создан?  %v ", ok)

	cert, err := tls.LoadX509KeyPair("cert/client.crt", "cert/client.key")
	if err != nil {
		Ошибка(" %s", err)
	}

	tlsConfig := &tls.Config{
		RootCAs:      caCertPool,
		Certificates: []tls.Certificate{cert},
	}

	// Подключение к TCP-серверу с TLS на localhost:8080
	количествоПопыток := 500
	задержка := 1 * time.Second
	var сервер *tls.Conn
	var errDial error
	for попытка := 1; попытка <= количествоПопыток; попытка++ {
		сервер, errDial = tls.Dial("tcp", адрес+":"+порт, tlsConfig)
		if errDial != nil {
			Ошибка("  %+v \n", err)
			time.Sleep(задержка)
		} else {
			break
		}
	}

	ОтправитьОтветНаЗапрос(сервер, каналОтправкиОтветов)
}

type ЗапросВОбработку struct {
	Сервис    []byte
	ИдКлиента uuid.UUID
	Запрос    Запрос
}

func ОтправитьОтветНаЗапрос(сервер *tls.Conn, каналОтправкиОтветов chan []byte) {
	for ОтветКлиенту := range каналОтправкиОтветов {
		// Отправка сообщений серверу
		Инфо(" ОтветКлиенту %+v \n", ОтветКлиенту)

		// БинарныйОтветКлиенту, err := Кодировать(ОтветКлиенту)

		// if err != nil {
		// 	Ошибка("  %+v \n", err)
		// }
		// Инфо(" БинарныйЗапрос %+s \n", БинарныйОтветКлиенту)

		// int, err := сервер.Write(БинарныйОтветКлиенту)
		int, err := сервер.Write(ОтветКлиенту)
		if err != nil {
			Ошибка("  %+v %+v \n", int, err)
		}
		Инфо(" отправленно  %+v \n", int)
	}
}

// func (з ЗапросВОбработку) Кодировать(T any) ([]byte, error) {
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
