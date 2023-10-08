package main

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/binary"
	"fmt"
	"io"
	"sync"
	"unsafe"

	"os"

	. "aoanima.ru/logger"
)

var клиенты = make(map[string]Запрос)
var мьютекс = sync.Mutex{}

// каналЗапросов - исползуется для получения запросов от клиента, в запросе от клиента передаётся канал в который нужно отправить ответ клиенту
func ПодключитсяКМенеджеруЗапросов(каналЗапросов chan Запрос) {
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
	сервер, err := tls.Dial("tcp", "localhost:81", tlsConfig)
	if err != nil {
		Ошибка(" %s", err)
		return
	}
	// каналЗапросов - исползуется для получения запросов от клиента, в запросе от клиента передаётся канал в который нужно отправить ответ клиенту
	go ОтправитьЗапросВОбработку(сервер, каналЗапросов)
	go ОтправитьОтветКлиенту(сервер, каналЗапросов)

}

type ЗапросВОбработку struct {
	ИдКлиента []byte
	Запрос    []byte
}

func ОтправитьЗапросВОбработку(сервер *tls.Conn, каналЗапросов chan Запрос) {
	for ЗапросОтКлиента := range каналЗапросов {
		// Отправка сообщений серверу

		мьютекс.Lock()
		клиенты[ЗапросОтКлиента.ИдКлиента] = ЗапросОтКлиента
		мьютекс.Unlock()

		// буфер := new(bytes.Buffer)
		ЗапросВОбработку := ЗапросВОбработку{
			ИдКлиента: []byte(ЗапросОтКлиента.ИдКлиента),
			Запрос:    []byte(ЗапросОтКлиента.Запрос),
		}
		Инфо(" ЗапросВОбработку %+v \n", ЗапросВОбработку)

		БинарныйЗапрос, err := ЗапросВОбработку.Кодировать()

		if err != nil {
			Ошибка("  %+v \n", err)
		}
		Инфо(" БинарныйЗапрос %+v \n", БинарныйЗапрос)
		// err = binary.Write(буфер, binary.LittleEndian, БинарныйЗапрос)
		// if err != nil {
		// 	Ошибка("  %+v \n", err)
		// }

		сервер.Write(БинарныйЗапрос)

	}
}

func ОтправитьОтветКлиенту(сервер *tls.Conn, каналЗапросов chan Запрос) {

	for {
		var ОтветКлиенту ОтветКлиенту
		длина := make([]byte, 128)
		n, err := io.ReadFull(сервер, длина)
		Инфо("  %+v \n", n)
		if err != nil {
			Ошибка("  %+v \n", err)
		}
		lenData := binary.LittleEndian.Uint64(длина)

		буфер := make([]byte, lenData)
		i, err := io.ReadFull(сервер, буфер)
		Инфо("  %+v \n", i)
		if err != nil {
			Ошибка("  %+v \n", err)
		}
		err = binary.Read(bytes.NewReader(буфер), binary.LittleEndian, &ОтветКлиенту)
		if err != nil {
			Ошибка("Ошибка при десериализации структуры: %+v ", err)
		}

		клиенты[ОтветКлиенту.ИдКлиента].КаналОтвета <- ОтветКлиенту

	}
}

func (з ЗапросВОбработку) Кодировать() ([]byte, error) {
	// Вычисляем размер буфера, который необходим для сериализации структуры
	// размер := binary.Size(з)
	размер := unsafe.Sizeof(з)
	Инфо(" размер  %+v %+v %+v %+v \n", з, размер, len(з.Запрос), len(з.ИдКлиента))
	// Создаем буфер нужного размера для сериализации
	// буфер := make([]byte, размер)
	буфер := new(bytes.Buffer)

	// Сериализуем поля структуры в буфер
	// бинарныеДанные := bytes.NewBuffer(буфер)
	// err := binary.Write(бинарныеДанные, binary.LittleEndian, з)
	// if err != nil {
	// 	Ошибка("  %+v \n", err)
	// 	return nil, err
	// }
	binary.Write(буфер, binary.LittleEndian, int32(len(з.Запрос)))
	binary.Write(буфер, binary.LittleEndian, з.Запрос)

	binary.Write(буфер, binary.LittleEndian, int32(len(з.ИдКлиента)))
	binary.Write(буфер, binary.LittleEndian, з.ИдКлиента)
	// err = binary.Write(бинарныеДанные, binary.LittleEndian, з.Запрос)
	// if err != nil {
	// 	Ошибка("  %+v \n", err)
	// 	return nil, err
	// }

	Инфо("бинарныеДанные  %+v ;Bytes %+v \n", буфер, буфер.Bytes())

	// Возвращаем сериализованные бинарные данные и ошибку (если есть)
	return буфер.Bytes(), nil
}

func ДеКодироватьОтветКлиенту(бинарныеДанные []byte) (*ОтветКлиенту, error) {
	буфер := bytes.NewReader(бинарныеДанные)
	var длинаИдКлиента int32
	if err := binary.Read(буфер, binary.LittleEndian, &длинаИдКлиента); err != nil {
		Ошибка("  %+v \n", err)
	}
	идКлиентаBytes := make([]byte, длинаИдКлиента)
	if err := binary.Read(буфер, binary.LittleEndian, &идКлиентаBytes); err != nil {
		return nil, fmt.Errorf("ошибка чтения ИдКлиента: %v", err)
	}
	идКлиента := string(идКлиентаBytes)

	var значениеBytes []byte
	if err := binary.Read(буфер, binary.LittleEndian, &значениеBytes); err != nil {
		return nil, fmt.Errorf("ошибка чтения значения типа string: %v", err)
	}
	ответ := string(значениеBytes)
	ответКлиенту := &ОтветКлиенту{
		ИдКлиента: идКлиента,
		Ответ:     ответ,
	}

	return ответКлиенту, nil
}
