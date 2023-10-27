package main

import (
	"crypto/tls"
	"crypto/x509"
	"net/http"

	_ "net/http/pprof"
	"os"

	. "aoanima.ru/logger"
)

var (
	ВходящийПорт  = ":81"
	ИсходящийПорт = ":82"
	// каналОтправкиОтветов     = make(chan ОтветКлиенту, 10)
	// КаналыИсходящихСообщений = map[string]chan ОтветКлиенту{}
)

// type ОтветКлиенту struct {
// 	Сервис    []byte
// 	Ответ     []byte
// 	ИдКлиента []byte
// }

// type ЗапросКлиента struct {
// 	Сервис       []byte
// 	Запрос       *ЗапросОтКлиента
// 	ИдКлиента    uuid.UUID
// 	ТокенКлиента []byte // JWT сериализованный
// }
// type ЗапросОтКлиента struct {
// 	СтрокаЗапроса string
// 	Форма         map[string][]string
// 	Файл          string
// }

func main() {
	go func() {
		http.ListenAndServe("localhost:6061", nil)
	}()
	// Вероятно нужно откуда то получить список Сервисов с которомы предстоит общаться
	//  Или !!!! ОбработчикВходящихСообщений

	go ЗапуститьСерверВходящихСообщений()

	Инфо(" %s", "запустили сервер")
	ЗапуститьСерверИсходящихСообщений()
}

// Сервер для обработки взодящих запросов, принимает только входящие сообщения, не отвечает на запросы
func ЗапуститьСерверВходящихСообщений() {

	cert, err := tls.LoadX509KeyPair("cert/server.crt", "cert/server.key")
	if err != nil {
		Ошибка(" %s", err)

	}
	caCert, err := os.ReadFile("cert/ca.crt")
	if err != nil {
		Ошибка(" %s", err)
	}
	caCertPool := x509.NewCertPool()
	ok := caCertPool.AppendCertsFromPEM(caCert)
	Инфо("Корневой сертфикат создан?  %v ", ok)

	конфиг := &tls.Config{
		ClientAuth:   tls.RequireAndVerifyClientCert,
		ClientCAs:    caCertPool,
		Certificates: []tls.Certificate{cert},
	}
	сервер, err := tls.Listen("tcp", ВходящийПорт, конфиг)
	if err != nil {
		Ошибка("  %+v \n", err)
	}

	defer сервер.Close()
	for {
		клиент, err := сервер.Accept()
		if err != nil {
			Ошибка("  %+v \n", err)
			defer клиент.Close()
		}
		// Инфо(" %+v %+v \n", "клиент подключен ко входящему серверу ", клиент)
		go обработчикВходящихСообщений(клиент)
		// go ТестВХодящихСообщенийСнизкойСкоростью(клиент)
	}
}

// Сервер через который отправляются запросы в сервисы или ответы клиенту.
func ЗапуститьСерверИсходящихСообщений() {
	cert, err := tls.LoadX509KeyPair("cert/server.crt", "cert/server.key")
	if err != nil {
		Ошибка(" %s", err)

	}
	caCert, err := os.ReadFile("cert/ca.crt")
	if err != nil {
		Ошибка(" %s", err)
	}
	caCertPool := x509.NewCertPool()
	ok := caCertPool.AppendCertsFromPEM(caCert)
	Инфо("Корневой сертфикат создан?  %v ", ok)

	конфиг := &tls.Config{
		ClientAuth:   tls.RequireAndVerifyClientCert,
		ClientCAs:    caCertPool,
		Certificates: []tls.Certificate{cert},
	}
	сервер, err := tls.Listen("tcp", ИсходящийПорт, конфиг)
	if err != nil {
		Ошибка("  %+v \n", err)
	}

	defer сервер.Close()
	for {
		клиент, err := сервер.Accept()
		if err != nil {
			Ошибка("  %+v \n", err)
		}
		defer клиент.Close()
		Инфо(" %+v \n", "клиент подключен к сиходящему серверу")
		go обработчикИсходящихСоединений(клиент)
	}
}
