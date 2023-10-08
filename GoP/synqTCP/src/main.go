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
	Ответ     []byte
	ИдКлиента []byte
}

type ЗапросКлиента struct {
	Запрос    []byte
	ИдКлиента []byte
}

func main() {

	ЗапуститьСервер()

	Инфо(" %s", "запустили сервер")

}

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
	Инфо(" длина %+v \n", длина)

	lenData := binary.LittleEndian.Uint64(длина)

	Инфо(" lenData  %+v \n", lenData)

	var запросКлиента ЗапросКлиента
	буфер := make([]byte, lenData)

	err = binary.Read(bytes.NewReader(буфер), binary.LittleEndian, &запросКлиента)
	if err != nil {
		Ошибка("Ошибка при десериализации структуры: %+v ", err)
	}
	Инфо("  %+v \n", запросКлиента)
}
