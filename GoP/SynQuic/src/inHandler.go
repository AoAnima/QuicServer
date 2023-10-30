package main

import (
	"encoding/binary"
	"errors"
	"fmt"
	"net/url"
	"strings"
	"sync"

	connector "aoanima.ru/connector"
	. "aoanima.ru/logger"
	"github.com/quic-go/quic-go"
)

func СохранитьСообщение(запрос *connector.Сообщение) {
	JSONЗапроса, err := Кодировать(запрос.Запрос)
	if err != nil {
		Ошибка("  %+v \n", err)
	}
	sql := fmt.Sprintf("INSERT INTO querys (id,id_query, query, service) VALUES (%s, %s, %s, %s)", запрос.ИдКлиента.String(), запрос.УИДСообщения, JSONЗапроса, string(запрос.Сервис))
	Инфо(" Пишем в бд >> %+v \n", sql)

}

// Обрабатывает только запросы полученный от сервисов, в ответ ничего не отправляет
func обработчикВходящихСообщений(поток quic.Stream, сообщение []byte) {

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

		if длинаДанных < 4 {
			Ошибка(" длинаДанных нечгео декодировать %+v \n", длинаДанных)
			return
		}

		//читаем количество байт = длинаСообщения
		// var запросКлиента ЗапросКлиента
		пакетЗапроса := make([]byte, длинаДанных)
		прочитаноБайт, err = поток.Read(пакетЗапроса)

		if err != nil {
			Ошибка("Ошибка при десериализации структуры: %+v ", err)
		}

		if длинаДанных != uint32(прочитаноБайт) {
			Ошибка("Количество прочитаных байт не ранво длине данных :\n длинаДанных %+v  <> прочитаноБайт %+v ", длинаДанных, прочитаноБайт)
		}

		// Запускаем для пакета отдельную горутину, т.к. в ожном соединении будет приходить множество запросов от разных клиентов, и обработчик будт всегда один

		go ДекодироватьПакет(пакетЗапроса)

	}

}

// var КартаЗапросов map[string]Очередь = make(map[string]Очередь)
//
//	type ОчередьМаршрутов struct {
//		Очередь map[string]Очередь
//		mutex sync.Mutex
//	}
type ОтветПолучен bool

var БлокКартаЗапросов sync.Mutex
var КартаЗапросов = make(map[connector.Уид]struct {
	Очередь          *Очередь
	Паралельно       map[connector.Сервис]ОтветПолучен // [сервис]ОтветПолучен
	ФинальныеСервисы *Очередь                          //сервисы которые выполняются после получения всех ответов от Паралельных и последовательных, Сюда добавляем такие сервисы как Рендер - который рендерит ответ, в html или json структуру заданного формата, и Сервис от которого пришёл первый зпрос
})
var БлокКартаОтветов sync.Mutex
var КартаОтветов = make(map[connector.Уид]struct {
	Ответы map[connector.Сервис]connector.Ответ
})

// АнализЗапроса -  анализируем url.Path с целью обнаружения сервиса в который отправляем данные
func АнализСообщения(Сообщение *connector.Сообщение) {

	/*
		Если же несколько сервисов не нуждаютс в ажанных от других сервисов и обрабатываюют только изначальный запрос, то к таким сервисам можно отправлять запрос паралельно.
	*/

	// Попытаемся получить Маршурт(Сервис) для текущего сообщения по его УИДСообщения, если в КартеЗАпросов Есть очередь для текущего Сообщения.УИДСообщения то получаем сервис и отправляем ав него данные, иначе если такого УИЛ в карте нету то строим очереди из сервисоов

	/*
		Не получается обработать паралельные запросы

		не правильная логика

		сервисы должны отправлять Сообщения, если сервис отправляет Сообщение, то поле Сервис должно содержать имя сервиса отправившего сообщение,.... нужно изменить стурктуру и ллогику,

	*/

	var (
		сервис string
		err    error
	)

	// Если в карте запросов есть УДИЗапроса то проверим, если очередь - последовательных запросов не пуста, то получим следующий сервис и отправим в него запрос
	// Тух проверяем по уид апроса, потому что поле запрос при отправке с другой сервис для получения ответа останеться не изменным
	if очередьЗапросов, ЕстьОчередиЗапросов := КартаЗапросов[Сообщение.Запрос.УИДЗапроса]; ЕстьОчередиЗапросов {

		// Вначале проверим не пришёл ли ответ от Сервиса обрабатывающего паралельный запрос, если очерель паралельнызх запросов не пуста и Такой сервис есть в Очереди Ожидающих ответов, то сохраним ответ, и продолжим ожидание сообщений, если же очередь паралельных запросв пуста или в очереди нет Сервиса, то перейдём к последовательным запросам
		Инфо(" очередьЗапросов %+v \n", очередьЗапросов)
		_, естьОжиданиеОтвета := очередьЗапросов.Паралельно[connector.Сервис(Сообщение.Сервис)]
		if естьОжиданиеОтвета { // len(очередьЗапросов.Паралельно) > 0 &&

			БлокКартаОтветов.Lock()
			КартаОтветов[Сообщение.Запрос.УИДЗапроса].Ответы[Сообщение.Сервис] = Сообщение.Ответ
			БлокКартаОтветов.Unlock()

			БлокКартаЗапросов.Lock()
			// удалим сервс из списка ожидающих ответа
			delete(КартаЗапросов[Сообщение.Запрос.УИДЗапроса].Паралельно, Сообщение.Сервис)
			БлокКартаЗапросов.Unlock()

			// Если паралельные запросы тоже стали пустые, то перместим все Ответы от сервисов из КратыОТветов в Изначальный Запрос. отправим в Рендер и в Сервис запросивший данные.....

		} else {
			// полгаем что следующее сообщение из очерди запросов не пустое....
			/**

			 */
			if неПусто := !очередьЗапросов.Очередь.Пусто(); неПусто {

				сервис, err = ПолучитьСледующийСервис(&Сообщение.Запрос.УИДЗапроса)

				Инфо("  %+v \n")
				if сервис != "" && err == nil {
					// Перед отправкой в сервис, предзаполним Поле ответ, Флаг ОтветПолучен выставим в лож, внутри сервиса, сервис должен поменять влаг на false, при паралельной обработке , нужно дождаться чтобы все ответы были получены перед отправкой клиенту...
					// if сервис != string(Запрос.Сервис) { следующий сервис == Запрос.Сервис - из которого пришло первое сообщение
					// if сервис != "Рендер" {

					Инфо(" Запрос.Ответ проверка до изменения ответов %+v \n", Сообщение.Ответ)
					Сообщение.Ответ[connector.Сервис(сервис)] = connector.ОтветСервиса{
						Сервис:          connector.Сервис(сервис),
						УИДЗапроса:      string(Сообщение.Запрос.УИДЗапроса), // потому что сервис при ответе в УИДСообщения напишет своё значение
						Данные:          nil,
						ЗапросОбработан: false,
					}
					Инфо(" Запрос.Ответ проверка после изменения  %+v \n", Сообщение.Ответ)

					ОтправитьЗапросВСервис(сервис, Сообщение)
					// } else {
					// следующий сервис == Рендер
					// Значит следующи сервис юудет равен Запрос.Сервис - из которого пришло первое сообщение, провери мвсе ли сервисы которые выполняют паралельную обработку запросов вернули ответы,

				}
			}

			// Проверим если обе очереди пусты, то объеденим ответы от всех сервисов и отправим в финальные сервисы
			if очередьЗапросов.Очередь.Пусто() && len(очередьЗапросов.Паралельно) == 0 {
				БлокКартаЗапросов.Lock()
				сервис := очередьЗапросов.ФинальныеСервисы.Далее()
				БлокКартаЗапросов.Unlock()

				if сервис != nil {
					for имяСервиса, ДанныеОтвета := range КартаОтветов[Сообщение.Запрос.УИДЗапроса].Ответы {
						Сообщение.Ответ[имяСервиса] = ДанныеОтвета[имяСервиса]
					}
					БлокКартаОтветов.Lock()
					delete(КартаОтветов, Сообщение.Запрос.УИДЗапроса)
					БлокКартаОтветов.Unlock()
					ОтправитьЗапросВСервис(сервис.(string), Сообщение)
					if очередьЗапросов.ФинальныеСервисы.Пусто() {
						БлокКартаЗапросов.Lock()
						delete(КартаЗапросов, Сообщение.Запрос.УИДЗапроса)
						БлокКартаЗапросов.Unlock()
					}
				}
			}

		}
	} else {
		// Нет УИДЗапроса в Карте Запросов , создадим очередь
		err = ПостроитьМаршрут(Сообщение)
		// тут нужно перейти в начало
	}

	// if err != nil {
	// 	if err.Error() == "нет очереди для УИДСообщения" {
	// 		Инфо("  %+v, в Маршрутизаторе нет маршрута для запроса с УИДСообщения  %+v Построим Новый маршрут \n", err, &Запрос.УИДСообщения)

	// 		err = ПостроитьМаршрут(Запрос) // всегда создаёт маршрут для запроса, если не нашёл сервис соответсвующий маршруту то создаст очередь из 404 , и Рендер

	// 		if err != nil {
	// 			Ошибка(" НЕ удалось построить маршрут, нужно что то вернуть клиенту %+v \n", err)
	// 		} else {
	// 			Инфо(" получаем следующий сервис и отправляем в него данные запроса %+v \n", "ПолучитьСледующийСервис")
	// 			сервис, err = ПолучитьСледующийСервис(&Запрос.Запрос.УИДЗапроса)
	// 		}
	// 	} else {
	// 		Ошибка(" %+v \n", err)
	// 	}
	// }

}

// По УидЗапроса проверяет есть ли в КартеЗАпросов маршрут , и если есть возвращает следующий сервис
func ПолучитьСледующийСервис(УИДСообщения *connector.Уид) (string, error) {
	БлокКартаЗапросов.Lock()
	defer БлокКартаЗапросов.Unlock()

	if очередьСервисов, ok := КартаЗапросов[*УИДСообщения]; ok {
		// Получим следующий элемент из очереди
		Очередь := очередьСервисов.Очередь
		СледующийСервис := Очередь.Далее()
		// если он не пустой
		if СледующийСервис != nil {
			// проверим есть ли в очередьСервисов ещё элемнеты, еслли нету, то очистим КартуЗАпросов , удалив данные для текущего сообщения, полагаем что это последний сервис, и дальше результат будет отправлен клиенту
			if Очередь.Пусто() {
				Инфо(" Больше нет последовательных сервисов %+v \n", Очередь)
				// удалим и КартыЗАпросов УИДСообщения - если удалим то как же очередь паралельных запросов? вдруг они ещё н еобработались
				// delete(КартаЗапросов, *УИДСообщения)

			}
			return СледующийСервис.(string), nil
		} else {
			// по идее сюда мы никогда не попадём
			Ошибка("Сюда не должны были попасть СледующийСервис == nil  %+v \n", СледующийСервис)
			// delete(КартаЗапросов, *УИДСообщения)
			return "", errors.New("cюда не должны были попасть? нет больше сервисов в очереди")
		}
	}
	return "", errors.New("нет очереди для УИДСообщения")

}

type ОчередностьСервисовИзБД struct {
	Очередь    []string
	Паралельно []string
}

func ПолучитьОчередьСервисовИзБД(Запрос *connector.Сообщение) (ОчередностьСервисовИзБД, error) {

	Инфо(" ПолучитьОчередьСервисов %+v \n", "Обращаемся к БД, получаем очередь обработки запроса, для указаного запроса. Некоторый запросы могут выполнятся не зависимо от других сервисов")

	параметрыЗапроса, err := url.Parse(Запрос.Запрос.МаршрутЗапроса)
	Инфо(" параметрыЗапроса %+v \n", параметрыЗапроса)
	if err != nil {
		Ошибка("Ошибка при парсинге СтрокаЗапроса запроса:", err)
		return ОчередностьСервисовИзБД{}, err
	}
	параметрыЗапроса.Path = strings.Trim(параметрыЗапроса.Path, "/")
	маршрут := strings.Split(параметрыЗапроса.Path, "/")
	// TODO: Тут нужно сделать обращение к БД для получения очереди сервисов которые должны обработать запрос еред возвратом клиенту
	if len(маршрут) == 0 {
		Инфо(" Пустой маршрут, добавляем в маршруты обработку по умолчанию.... \n")
		маршрут = append(маршрут, "/")
	}

	/* TODO: Обращение к БД за маршрутами
	данныеИзБД := ЗапросВБД()
	очередь := ОчередностьСервисов{
					Очередь: данныеИзБД["очередь"],
					Паралельно: данныеИзБД["паралельно"],
				}

	*/
	очередь := ОчередностьСервисовИзБД{
		Очередь:    маршрут,
		Паралельно: make([]string, 0),
	}
	return очередь, nil
}

func ПостроитьМаршрут(Сообщение *connector.Сообщение) error {
	Инфо(" ПостроитьМаршрут %+v \n", Сообщение)
	// TODO: Доделать ПостроитьМаршрут, нет обработчика для БД
	// параметрыЗапроса, err := url.Parse(Запрос.Запрос.МаршрутЗапроса)
	// Инфо(" параметрыЗапроса %+v \n", параметрыЗапроса)
	// if err != nil {
	// 	Ошибка("Ошибка при парсинге СтрокаЗапроса запроса:", err)
	// 	return err
	// }
	// параметрыЗапроса.Path = strings.Trim(параметрыЗапроса.Path, "/")

	// маршрут := strings.Split(параметрыЗапроса.Path, "/")
	// Инфо("  %+v \n")
	// // TODO: Тут нужно сделать обращение к БД для получения очереди сервисов которые должны обработать запрос еред возвратом клиенту
	// if len(маршрут) == 0 {
	// 	Инфо(" Пустой маршрут, добавляем в маршруты обработку по умолчанию.... \n")
	// 	маршрут = append(маршрут, "/")
	// } else {
	// }
	ОчередьМаршрутов, err := ПолучитьОчередьСервисовИзБД(Сообщение) // МОК. не работает, получаем очередь из БД, предварительно прасим URL path? и на его основе получаем маршруты
	Инфо(" ОчередьМаршрутов : %+v \n", ОчередьМаршрутов)
	if err != nil {
		Ошибка("  %+v \n", err)
	}
	ОчередьСеревисов := НоваяОчередь()

	Инфо(" Пока условно полагаем что маршрут строиться из url path , но по факту в path нас интересует только 1 элемент, он указывает на сервис в кооторый адресован запрос, тогда по идее, остальной маршрут должен вытсролиться из БД  %+v \n")

	if ОчередьМаршрутов.Очередь == nil && ОчередьМаршрутов.Паралельно == nil {
		ОчередьСеревисов.Добавить("404")
	}

	for _, сервис := range ОчередьМаршрутов.Очередь {
		// проверим есть ли обработчик для маршрута в карте маршуртизатора, если сесть добавим в очередь иначе 404
		if _, есть := Маршрутизатор[сервис]; есть {
			Инфо(" добавляем %+v \n", сервис)
			ОчередьСеревисов.Добавить(сервис)
		} else {
			Инфо(" Нет в маршрутизаторе сервиса %+v \n", сервис)
			// ОчередьСеревисов.Добавить("404")
		}
	}

	ПаралельныеСервисы := make(map[connector.Сервис]ОтветПолучен)
	for _, сервис := range ОчередьМаршрутов.Паралельно {
		ПаралельныеСервисы[connector.Сервис(сервис)] = false
	}
	ОчередьФинальныеСервисы := НоваяОчередь()

	ОчередьФинальныеСервисы.Добавить("Рендер")                 // Добавляем сервис Рендер для пожготовки ответа клиенту, в соответствии с запрошенным форматом
	ОчередьФинальныеСервисы.Добавить(string(Сообщение.Сервис)) // Последним добавим Сервис от которого пришло изначальное сообщение чтобы в него отправить ответ

	БлокКартаЗапросов.Lock()
	КартаЗапросов[Сообщение.Запрос.УИДЗапроса] = struct {
		Очередь          *Очередь
		Паралельно       map[connector.Сервис]ОтветПолучен
		ФинальныеСервисы *Очередь
	}{
		Очередь:          ОчередьСеревисов,
		Паралельно:       ПаралельныеСервисы,
		ФинальныеСервисы: ОчередьФинальныеСервисы,
	}
	// КартаЗапросов[Запрос.УИДСообщения].Очередь = *ОчередьСеревисов
	БлокКартаЗапросов.Unlock()

	Инфо(" ПолучитьМаршрут анализируем запрос, обращаемся в БД за получением маршрутизации %+v \n")
	// return "маршрут запроса... ДОДЕЛАТЬ"
	return nil
}

// ЧПросто отправляет сообщение в указаный сервис
func ОтправитьЗапросВСервис(сервис string, Сообщение *connector.Сообщение) {

	if Сервис, есть := Маршрутизатор[сервис]; есть {
		Инфо(" Сервис %+v \n", Сервис)
		// Сервис.КаналСообщения читается в функции ОтправитьСообщениеВСервис()
		Сервис.КаналСообщения <- Сообщение
	} else {

		Ошибка(" Нет маршрута для  %+v есть %+v\n", сервис, есть)
		Ошибка("Маршрутизатор  %+v \n", Маршрутизатор)
	}

}
