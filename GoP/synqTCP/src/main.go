package main

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/binary"
	"io"
	"net"

	"os"

	. "aoanima.ru/logger"
)

var (
	PORT    = 81
	TLSPort = 444
)

type ОтветКлиенту struct {
	Ответ     interface{}
	ИдКлиента string
}

type ЗапросКлиента struct {
	Запрос    interface{}
	ИдКлиента string
}

func main() {

	ЗапуститьСервер()

	Инфо(" %s", "запустили сервер")

}

func ServerWSS() {}

func ЗапуститьСервер() {
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
	сервер, err := tls.Listen("tcp", ":81", конфиг)
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
		Инфо(" %+v \n", "клиент подключен")
		go обработчикСоединения(клиент)
	}
}

func обработчикСоединения(клиент net.Conn) {
	defer клиент.Close()
	Инфо(" обработчикСоединения \n")
	длина := make([]byte, 8)
	_, err := io.ReadFull(клиент, длина)
	if err != nil {
		Ошибка("  %+v \n", err)
	}
	lenData := binary.LittleEndian.Uint64(длина)

	var ОтветКлиенту ОтветКлиенту
	буфер := make([]byte, lenData)
	err = binary.Read(bytes.NewReader(буфер), binary.LittleEndian, &ОтветКлиенту)
	if err != nil {
		Ошибка("Ошибка при десериализации структуры: %+v ", err)
	}
	Инфо("  %+v \n", ОтветКлиенту)
}
