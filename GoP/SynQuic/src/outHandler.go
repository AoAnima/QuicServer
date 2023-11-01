package main

import (
	"encoding/binary"
	"net"

	. "aoanima.ru/ConnQuic"
	. "aoanima.ru/logger"
	jsoniter "github.com/json-iterator/go"
)

func обработчикИсходящихСоединений(клиент net.Conn) { //, данныеДляОтвета chan []byte
	go РукопожатиеИсходящегоКанала(клиент)
}

func (о *ОтпечатокСервиса) ЧитатьКаналСообщений() {
	for данныеИзКанала := range о.КаналСообщения {
		go о.ОтправитьСообщениеВСервис(данныеИзКанала)
	}

}

func (о *ОтпечатокСервиса) ОтправитьСообщениеВСервис(данныеИзКанала interface{}) {
	// берём соединение из буферизированного канала - своего рода пулл соедиениений

	// Проверим соответствуют ли полученные данные нужному формату

	клиент := <-о.Клиент.пулл

	Инфо(" взяли соединение из пулла  \n")

	данныеДляОтправки, err := Кодировать(данныеИзКанала)

	Инфо("данныеДляОтправки  %+v \n", данныеДляОтправки)
	if err != nil {
		Ошибка("  %+v \n", err)
	}
	// Инфо("клиент  %+v \n", клиент)

	отправленноБайт, err := клиент.Write(данныеДляОтправки)

	Инфо(" отправленноБайт клиенту  %+v \n", отправленноБайт)
	if err != nil {
		Ошибка("  %+v  отправленноБайт %+v \n", err, отправленноБайт)
	} else {

		о.Клиент.пулл <- клиент
		Инфо(" возвратили соединение в пулл  \n")
	}

}

// исключительно для рукопожатия и сохранения в пул Сервисов/ когда сервис присылает запрос на рукопожатие, он присылает маршруты которые он обрабатывает !!!
// напрмиер сервис каталогов обрабатывает запросы  начинающиеся на /catalog /product ...
func РукопожатиеИсходящегоКанала(клиент net.Conn) {
	длинаСообщения := make([]byte, 4)
	// рукопожатие := [4]byte{}
	var прочитаноБайт int
	var err error
	for {
		// получаем длину сообщения рукопожатия
		прочитаноБайт, err = клиент.Read(длинаСообщения)
		if err != nil {
			Ошибка("  %+v \n", err)
			break
		}
		Инфо("  %+v \n", прочитаноБайт)
		// читаем всё остальное сообщение
		// создадим буфер куда поместим сообщение
		// сообщениеРукопожатия := make([]byte, binary.LittleEndian.Uint32(длинаСообщения))

		// copy(рукопожатие[0:], длинаСообщения[:4])

		// if длинаСообщенияФикс == [4]byte{240, 159, 164, 157} { //"🤝"
		// 	Инфо(" %+v \n", string(длинаСообщения))
		сообщениеОтСервиса := make([]byte, binary.LittleEndian.Uint32(длинаСообщения))
		_, err = клиент.Read(сообщениеОтСервиса)
		if err != nil {
			Ошибка("  %+v \n", err)
		}

		ОтпечатокСервиса := ОтпечатокСервиса{}

		err := jsoniter.Unmarshal(сообщениеОтСервиса, &ОтпечатокСервиса)
		if err != nil {
			Ошибка("  %+v \n", err)
		}

		РегистрацияСервиса(&ОтпечатокСервиса, клиент)

	}
}