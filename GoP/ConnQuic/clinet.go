package ConnQuic

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/binary"
	"log"
	"os"

	. "aoanima.ru/logger"
	quic "github.com/quic-go/quic-go"
)

//	type Сессия struct {
//		Блок   sync.RWMutex
//		Сессия quic.Connection
//		Потоки []*quic.Stream
//	}
//
// type КартаСессий struct {
// Блок   sync.RWMutex
//
//		СессииСервисов []Сессия        // кладём соовтетсвие сессий и потоков
//		ОчередьПотоков *ОчередьПотоков // все потоки всех сессий кладём в одну очередь
//	}
/*
Создаём новый Серверв сервер := Сервер{
	Адрес: "localhost:4242",
	Имя:   "ИмяСервера",
}
Создаём новый Сессии := КартаСессий{}
клиент := Клиент{
	серврер: Сессии
}
Вызываем мтед Соединиться
Клиент.Соединиться(Адрес string, обработчикСообщений func(поток quic.Stream, сообщение Сообщение))

дальше если нужно ещё одно соединение, то повторям,

*/

// где string это адрес или имя сервиса.. лучше наверное адрес
type Сервер string
type СхемаСервера struct {
	Имя   Сервер
	Адрес string
	КартаСессий
}

type Клиент map[Сервер]СхемаСервера

func (к Клиент) Соединиться(сервер Сервер, обработчикСообщений func(поток quic.Stream, сообщение Сообщение)) {
	конфигТлс, err := клиентскийТлсКонфиг("root.crt")
	if err != nil {
		Ошибка("  %+v \n", err)
	}
	// Адрес = "localhost:4242"
	if сервер == "" {
		сервер = "SynQuic"

		к[сервер] = СхемаСервера{
			Имя:         сервер,
			Адрес:       "localhost:4242",
			КартаСессий: КартаСессий{},
		}
	}

	сессия, err := quic.DialAddr(context.Background(), к[сервер].Адрес, конфигТлс, &quic.Config{})
	if err != nil {
		log.Fatal(err)
	}

	// к.Сессии = append(к.Сессии, &сессия)
	// к.Блок.RUnlock()
	for {
		поток, err := сессия.AcceptStream(context.Background())

		c := Сессия{
			Соединение: сессия,
			Потоки:     []quic.Stream{поток},
		}

		схемаСервера := к[сервер]
		схемаСервера.Lock()
		схемаСервера.СессииСервисов = append(схемаСервера.СессииСервисов, c)
		схемаСервера.RUnlock()
		схемаСервера.ОчередьПотоков.Добавить(поток)

		if err != nil {
			Ошибка("  %+v \n", err)
		}
		go к.ЧитатьСообщения(поток, обработчикСообщений)

	}

}

func (к Клиент) ЧитатьСообщения(поток quic.Stream, обработчикСообщений func(поток quic.Stream, сообщение Сообщение)) {

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
		сообщениеБинарное := make([]byte, длинаДанных)
		прочитаноБайт, err = поток.Read(сообщениеБинарное)
		if err != nil {
			Ошибка("Ошибка при десериализации структуры: %+v ", err)
		}

		if длинаДанных != uint32(прочитаноБайт) {
			Ошибка("Количество прочитаных байт не ранво длине данных :\n длинаДанных %+v  <> прочитаноБайт %+v ", длинаДанных, прочитаноБайт)
		} else {

			сообщение := ДекодироватьПакет(сообщениеБинарное)

			go обработчикСообщений(поток, сообщение)

			break
		}

		// каналПолученияСообщений <- пакетОтвета

	}

}

func клиентскийТлсКонфиг(caCertFile string) (*tls.Config, error) {
	caCert, err := os.ReadFile(caCertFile)
	if err != nil {
		return nil, err
	}

	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	return &tls.Config{
		RootCAs:    caCertPool,
		NextProtos: []string{"http/1.1", "h2", "h3", "quic", "websocket"},
	}, nil
}
