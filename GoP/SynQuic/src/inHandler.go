package main

import (
	"errors"
	"fmt"
	"net/url"
	"strings"
	"sync"

	. "aoanima.ru/ConnQuic"
	. "aoanima.ru/Logger"
)

type ОтветПолучен bool

type СтруктураКартыЗапросов struct {
	Очередь            *Очередь
	Блок               *sync.RWMutex
	ПаралельнаяОчередь map[ИмяСервиса]ОтветПолучен // [сервис]ОтветПолучен
	ФинальныеСервисы   *Очередь                    //сервисы которые выполняются после получения всех ответов от Паралельных и последовательных, Сюда добавляем такие сервисы как Рендер - который рендерит ответ, в html или json структуру заданного формата, и Сервис от которого пришёл первый зпрос
}

var БлокКартаЗапросов sync.RWMutex
var КартаЗапросов = make(map[Уид]СтруктураКартыЗапросов)

var БлокКартаОтветов sync.RWMutex

type КартаОтветов map[Уид]struct {
	Блок   *sync.RWMutex
	Ответы map[ИмяСервиса]Ответ
}

type ОчередностьСервисовИзБД struct {
	Очередь            []string
	ПаралельнаяОчередь []string
}

func СохранитьСообщение(запрос *Сообщение) {
	JSONЗапроса, err := Кодировать(запрос.Запрос)
	if err != nil {
		Ошибка("  %+v \n", err)
	}
	sql := fmt.Sprintf("INSERT INTO querys (id,id_query, query, service) VALUES (%s, %s, %s, %s)", запрос.ИдКлиента.String(), запрос.УИДСообщения, JSONЗапроса, string(запрос.Сервис))
	Инфо(" Пишем в бд >> %+v \n", sql)

}

func обработчикСообщенийHTTPсервера(сообщение Сообщение) (Сообщение, error) {
	go СохранитьСообщение(&сообщение)
	Инфо(" cообщение %+v %T \n", сообщение, сообщение)
	err := ПостроитьМаршрут(&сообщение)
	if err != nil {
		Ошибка(" не удалось построить маршрут %+v \n", err)
		return Сообщение{}, err

	} else {
		// после построения маршрутов, получим очеред запросов из карты
		очередьЗапросов, ЕстьОчередиЗапросов := КартаЗапросов[сообщение.Запрос.УИДЗапроса]
		if ЕстьОчередиЗапросов {

			картаПаралельныхОтветов := make(КартаОтветов)
			картаПаралельныхОтветов[сообщение.Запрос.УИДЗапроса] = struct {
				Блок   *sync.RWMutex
				Ответы map[ИмяСервиса]Ответ
			}{}

			ожидание := sync.WaitGroup{}
			Ошибка(" Пустая очередь запросов \n")
			if len(очередьЗапросов.ПаралельнаяОчередь) > 0 {
				ожидание.Add(1)
				go ОбработатьПаралельныеСервисы(&очередьЗапросов, сообщение, &картаПаралельныхОтветов, &ожидание)
				// Ответы от сервисов будут складываться в КартаОТветов
			}

			ответ := Сообщение{}
			if !очередьЗапросов.Очередь.Пусто() {
				// очередьЗапросовВОбработку := очередьЗапросов.Очередь
				// перезапишем сообщение, со всеми ответами от последовтаельных сервисов
				сообщение, err = ОбработатьОчередь(очередьЗапросов.Очередь, сообщение)
				if err != nil {
					Ошибка("  %+v \n", err)
				}
			}

			// Дождёмся обработки всех паралельных запросов
			// for len(очередьЗапросов.ПаралельнаяОчередь) > 0 {
			ожидание.Wait()

			if len(очередьЗапросов.ПаралельнаяОчередь) == 0 {
				// Все паралельные запросы отработали, соххраним все ответы от паралельных сервисов в Сообщение
				for имяСервиса, ДанныеОтвета := range картаПаралельныхОтветов[сообщение.Запрос.УИДЗапроса].Ответы {
					сообщение.Ответ[имяСервиса] = ДанныеОтвета[имяСервиса]
				}
				// очистим карту ответов для текущего запроса
				// БлокКартаОтветов.Lock()
				// delete(КартаОтветов, сообщение.Запрос.УИДЗапроса)
				// БлокКартаОтветов.Unlock()
				картаПаралельныхОтветов = nil // удалим карту освободим память

				ответ, err = ОбработатьОчередь(очередьЗапросов.ФинальныеСервисы, сообщение)
				if err != nil {
					Ошибка("  %+v \n", err)
				}

				БлокКартаЗапросов.RLock()
				delete(КартаЗапросов, сообщение.Запрос.УИДЗапроса)
				БлокКартаЗапросов.RUnlock()
				// отправить, err := Кодировать(ответ)
				// if err != nil {
				// 	Ошибка("  %+v \n", err)
				// }
				// потокДляОтвета.Write(отправить)
				return ответ, nil
			}
			// }
		}
	}
	return Сообщение{}, fmt.Errorf("ХЗ , не удалось построить маршрут и почему сюда вышло ? ")
}

// func ОбработатьФинальныеСервисы(очередьЗапросов *Очередь, сообщение Сообщение) Сообщение {
// 	return ОбработатьОчередь(очередьЗапросов, сообщение)
// }

func ОбработатьОчередь(очередьЗапросов *Очередь, сообщение Сообщение) (Сообщение, error) {

	for !очередьЗапросов.Пусто() {
		имяСервиса := ИмяСервиса(очередьЗапросов.Далее().(string))
		потокСессии, err := АктивныеСессии.ПолучитьПоток(имяСервиса) // получаем поток сервиса
		if err != nil {
			Ошибка("  %+v \n", err)
			return сообщение, err
		}
		сообщениеБин, err := Кодировать(сообщение)
		if err != nil {
			Ошибка("  %+v \n", err)
			return сообщение, err
		}
		// отправляем в поток сервиса
		потокСессии.Поток.Write(сообщениеБин)
		ответ := ЧитатьСообщение(потокСессии.Поток)
		АктивныеСессии.Вернуть(имяСервиса, потокСессии)

		сообщение = ответ
		сообщение.Ответ[имяСервиса] = ОтветСервиса{
			Сервис:          имяСервиса,
			УИДЗапроса:      string(сообщение.Запрос.УИДЗапроса), // потому что сервис при ответе в УИДСообщения напишет своё значение
			Данные:          nil,
			ЗапросОбработан: false,
		}
	}
	return сообщение, nil

}

func ОбработатьПаралельныеСервисы(картаЗапросов *СтруктураКартыЗапросов, сообщение Сообщение, картаПаралельныхОтветов *КартаОтветов, ожидание *sync.WaitGroup) {
	for имяПаралельногоСервиса := range картаЗапросов.ПаралельнаяОчередь {
		// ДОБАВИТЬ ВЭЙТГРУПП ИЛИ ОЖИДАЮЩИЙ КАНАЛ
		ожидание.Add(1)
		go func(имяПаралельногоСервиса ИмяСервиса) {
			потокСессии, err := АктивныеСессии.ПолучитьПоток(имяПаралельногоСервиса)
			if err != nil {
				Ошибка("  %+v \n", err)
				return
			}
			сообщениеБин, err := Кодировать(сообщение)
			if err != nil {
				Ошибка("  %+v \n", err)
			}
			потокСессии.Поток.Write(сообщениеБин)
			ответ := ЧитатьСообщение(потокСессии.Поток)
			СохранитьПаралельныйОтвет(&ответ, картаПаралельныхОтветов)
			// возвращаем поток через который общались с паралельным сервисом в очередь
			АктивныеСессии.Вернуть(имяПаралельногоСервиса, потокСессии)
			ожидание.Done()
		}(имяПаралельногоСервиса)
	}
	ожидание.Done()
}

func СохранитьПаралельныйОтвет(сообщение *Сообщение, картаПаралельныхОтветов *КартаОтветов) {
	if очередьЗапросов, ЕстьОчередиЗапросов := КартаЗапросов[сообщение.Запрос.УИДЗапроса]; ЕстьОчередиЗапросов {

		_, естьОжиданиеОтвета := очередьЗапросов.ПаралельнаяОчередь[сообщение.Сервис]
		if естьОжиданиеОтвета { // len(очередьЗапросов.ПаралельнаяОчередь) > 0 &&

			(*картаПаралельныхОтветов)[сообщение.Запрос.УИДЗапроса].Блок.RLock()
			(*картаПаралельныхОтветов)[сообщение.Запрос.УИДЗапроса].Ответы[сообщение.Сервис] = сообщение.Ответ

			(*картаПаралельныхОтветов)[сообщение.Запрос.УИДЗапроса].Блок.RUnlock()

			БлокКартаЗапросов.RLock()
			// удалим сервс из списка ожидающих ответа
			delete(КартаЗапросов[сообщение.Запрос.УИДЗапроса].ПаралельнаяОчередь, сообщение.Сервис)
			БлокКартаЗапросов.RUnlock()

			// Если паралельные запросы тоже стали пустые, то перместим все Ответы от сервисов из КратыОТветов в Изначальный Запрос. отправим в Рендер и в Сервис запросивший данные.....

		} else {
			Ошибка(" Почему то в очереди запросов нет данных ожидания ответа от сервиса : %+v ? Ответ: %+v \n", сообщение.Сервис, сообщение)
		}
	} else {
		Ошибка(" Пустая очередь запросов \n", КартаЗапросов)
	}

}

// ИСПРАИВТЬ: ПолучитьОчередьСервисовИзБД наврено нужно переименовать,потому что получет очередь из баз и адресной строки
func ПолучитьОчередьСервисовИзБД(Запрос *Сообщение) (ОчередностьСервисовИзБД, error) {

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
					ПаралельнаяОчередь: данныеИзБД["ПаралельнаяОчередь"],
				}

	*/
	очередь := ОчередностьСервисовИзБД{
		Очередь:            маршрут,
		ПаралельнаяОчередь: make([]string, 0),
	}
	return очередь, nil
}

func ПостроитьМаршрут(Сообщение *Сообщение) error {
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

	ОчередьМаршрутов, err := ПолучитьОчередьСервисовИзБД(Сообщение) // МОК. не работает, получаем очередь из БД, предварительно прасим URL path? и на его основе получаем маршруты
	Инфо(" ОчередьМаршрутов : %+v \n", ОчередьМаршрутов)
	if err != nil {
		Ошибка("  %+v \n", err)
	}
	ОчередьСервисов := НоваяОчередь()

	Инфо(" Пока условно полагаем что маршрут строиться из url path , но по факту в path нас интересует только 1 элемент, он указывает на сервис в кооторый адресован запрос, тогда по идее, остальной маршрут должен вытсролиться из БД  %+v \n", Сообщение)

	if ОчередьМаршрутов.Очередь == nil && ОчередьМаршрутов.ПаралельнаяОчередь == nil {
		Инфо("  %+v ??\n", "404")
		// ОчередьСеревисов.Добавить("404")
	}
	группаОжидания := &sync.WaitGroup{}
	группаОжидания.Add(1)
	go func() {

		for _, маршрут := range ОчередьМаршрутов.Очередь {
			// Проверим если в активных сессиях есть Сессия с соответсвующим сервисом, то добавим в Очередь сервисов имя сервиса в окторый будет отправлен запрос,

			имяСервиса, ок := КартаОбработчиков[Маршрут(маршрут)]
			if !ок {
				Ошибка("в КартаОбработчиков нет сервиса для маршрута  %+v \n", маршрут)
				continue
			}

			if _, есть := АктивныеСессии[имяСервиса]; есть {
				Инфо(" добавляем %+v \n", имяСервиса)
				ОчередьСервисов.Добавить(имяСервиса)
			} else {
				Инфо(" Нет в маршрутизаторе имяСервиса %+v \n", имяСервиса)
				// ОчередьСеревисов.Добавить("404")
			}
		}
		группаОжидания.Done()
	}()

	ПаралельныеСервисы := make(map[ИмяСервиса]ОтветПолучен)
	группаОжидания.Add(1)
	go func() {

		for _, маршрутПаралельногоСервиса := range ОчередьМаршрутов.ПаралельнаяОчередь {
			имяПаралельногоСервиса, ок := КартаОбработчиков[Маршрут(маршрутПаралельногоСервиса)]
			if !ок {
				Ошибка("в КартаОбработчиков нет сервиса для маршрута  %+v \n", маршрутПаралельногоСервиса)
				continue
			}

			if _, есть := АктивныеСессии[ИмяСервиса(имяПаралельногоСервиса)]; есть {
				Инфо(" добавляем в ПаралельныеСервисы -  %+v \n", имяПаралельногоСервиса)
				ПаралельныеСервисы[ИмяСервиса(имяПаралельногоСервиса)] = false
			} else {
				Инфо(" Нет в маршрутизаторе имяПаралельногоСервиса %+v \n", имяПаралельногоСервиса)
				// ОчередьСеревисов.Добавить("404")
			}
		}
		группаОжидания.Done()
	}()

	ОчередьФинальныеСервисы := НоваяОчередь()

	ОчередьФинальныеСервисы.Добавить("Рендер")                 // Добавляем сервис Рендер для пожготовки ответа клиенту, в соответствии с запрошенным форматом
	ОчередьФинальныеСервисы.Добавить(string(Сообщение.Сервис)) // Последним добавим Сервис от которого пришло изначальное сообщение чтобы в него отправить ответ
	группаОжидания.Wait()

	БлокКартаЗапросов.Lock()
	КартаЗапросов[Сообщение.Запрос.УИДЗапроса] = СтруктураКартыЗапросов{
		Очередь:            ОчередьСервисов,
		ПаралельнаяОчередь: ПаралельныеСервисы,
		ФинальныеСервисы:   ОчередьФинальныеСервисы,
	}
	// КартаЗапросов[Запрос.УИДСообщения].Очередь = *ОчередьСеревисов
	БлокКартаЗапросов.Unlock()

	Инфо(" ПолучитьМаршрут анализируем запрос, обращаемся в БД за получением маршрутизации %+v \n")
	// return "маршрут запроса... ДОДЕЛАТЬ"
	return nil
}

func ПолучитьСледующийСервис(УИДСообщения *Уид) (string, error) {
	// БлокКартаЗапросов.RLock()
	// defer БлокКартаЗапросов.RUnlock()

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

// АнализЗапроса -  анализируем url.Path с целью обнаружения сервиса в который отправляем данные
// func ПолучитьДанные(Сообщение Сообщение) *Сообщение {

// 	/*
// 		Если же несколько сервисов не нуждаютс в ажанных от других сервисов и обрабатываюют только изначальный запрос, то к таким сервисам можно отправлять запрос ПаралельнаяОчередь.
// 	*/

// 	// Попытаемся получить Маршурт(Сервис) для текущего сообщения по его УИДСообщения, если в КартеЗАпросов Есть очередь для текущего Сообщения.УИДСообщения то получаем сервис и отправляем ав него данные, иначе если такого УИЛ в карте нету то строим очереди из сервисоов

// 	var (
// 		сервис string
// 		err    error
// 	)

// 	// Если в карте запросов есть УДИЗапроса то проверим, если очередь - последовательных запросов не пуста, то получим следующий сервис и отправим в него запрос
// 	// Тух проверяем по уид апроса, потому что поле запрос при отправке с другой сервис для получения ответа останеться не изменным
// 	очередьЗапросов, ЕстьОчередиЗапросов := КартаЗапросов[Сообщение.Запрос.УИДЗапроса]
// 	if !ЕстьОчередиЗапросов {
// 		err = ПостроитьМаршрут(&Сообщение)
// 		if err != nil {
// 			Ошибка(" %+v \n", err)
// 		} else {
// 			// после построения маршрутов, получим очеред запросов из карты
// 			очередьЗапросов, ЕстьОчередиЗапросов = КартаЗапросов[Сообщение.Запрос.УИДЗапроса]
// 		}
// 	}

// 	if ЕстьОчередиЗапросов {

// 		// Вначале проверим не пришёл ли ответ от Сервиса обрабатывающего паралельный запрос, если очерель паралельнызх запросов не пуста и Такой сервис есть в Очереди Ожидающих ответов, то сохраним ответ, и продолжим ожидание сообщений, если же очередь паралельных запросв пуста или в очереди нет Сервиса, то перейдём к последовательным запросам
// 		Инфо(" очередьЗапросов %+v \n", очередьЗапросов)
// 		_, естьОжиданиеОтвета := очередьЗапросов.ПаралельнаяОчередь[ИмяСервиса(Сообщение.Сервис)]
// 		if естьОжиданиеОтвета { // len(очередьЗапросов.ПаралельнаяОчередь) > 0 &&

// 			БлокКартаОтветов.RLock()
// 			КартаОтветов[Сообщение.Запрос.УИДЗапроса].Ответы[Сообщение.Сервис] = Сообщение.Ответ
// 			БлокКартаОтветов.RUnlock()

// 			БлокКартаЗапросов.RLock()
// 			// удалим сервс из списка ожидающих ответа
// 			delete(КартаЗапросов[Сообщение.Запрос.УИДЗапроса].ПаралельнаяОчередь, Сообщение.Сервис)
// 			БлокКартаЗапросов.RUnlock()

// 			// Если паралельные запросы тоже стали пустые, то перместим все Ответы от сервисов из КратыОТветов в Изначальный Запрос. отправим в Рендер и в Сервис запросивший данные.....

// 		} else {
// 			// полгаем что следующее сообщение из очерди запросов не пустое....
// 			/**

// 			 */
// 			неПусто := !очередьЗапросов.Очередь.Пусто()
// 			if неПусто {

// 				сервис, err = ПолучитьСледующийСервис(&Сообщение.Запрос.УИДЗапроса)

// 				Инфо("  %+v \n")
// 				if сервис != "" && err == nil {
// 					// Перед отправкой в сервис, предзаполним Поле ответ, Флаг ОтветПолучен выставим в лож, внутри сервиса, сервис должен поменять влаг на false, при ПаралельнаяОчередьй обработке , нужно дождаться чтобы все ответы были получены перед отправкой клиенту...
// 					// if сервис != string(Запрос.Сервис) { следующий сервис == Запрос.Сервис - из которого пришло первое сообщение
// 					// if сервис != "Рендер" {

// 					Инфо(" Запрос.Ответ проверка до изменения ответов %+v \n", Сообщение.Ответ)
// 					Сообщение.Ответ[ИмяСервиса(сервис)] = ОтветСервиса{
// 						Сервис:          ИмяСервиса(сервис),
// 						УИДЗапроса:      string(Сообщение.Запрос.УИДЗапроса), // потому что сервис при ответе в УИДСообщения напишет своё значение
// 						Данные:          nil,
// 						ЗапросОбработан: false,
// 					}
// 					Инфо(" Запрос.Ответ проверка после изменения  %+v \n", Сообщение.Ответ)

// 					// return &Сообщение // Возврааем ответ , для отправки
// 					// } else {
// 					// следующий сервис == Рендер
// 					// Значит следующи сервис юудет равен Запрос.Сервис - из которого пришло первое сообщение, провери мвсе ли сервисы которые выполняют паралельную обработку запросов вернули ответы,

// 				}
// 			}

// 			// Проверим если обе очереди пусты, то объеденим ответы от всех сервисов и отправим в финальные сервисы
// 			if !неПусто && len(очередьЗапросов.ПаралельнаяОчередь) == 0 {

// 				БлокКартаЗапросов.RLock()
// 				сервис := очередьЗапросов.ФинальныеСервисы.Далее()
// 				БлокКартаЗапросов.RUnlock()

// 				if сервис != nil {
// 					for имяСервиса, ДанныеОтвета := range КартаОтветов[Сообщение.Запрос.УИДЗапроса].Ответы {
// 						Сообщение.Ответ[имяСервиса] = ДанныеОтвета[имяСервиса]
// 					}
// 					БлокКартаОтветов.Lock()
// 					delete(КартаОтветов, Сообщение.Запрос.УИДЗапроса)
// 					БлокКартаОтветов.Unlock()

// 					// ОтправитьЗапросВСервис(сервис.(string), Сообщение)

// 					if очередьЗапросов.ФинальныеСервисы.Пусто() {
// 						БлокКартаЗапросов.RLock()
// 						delete(КартаЗапросов, Сообщение.Запрос.УИДЗапроса)
// 						БлокКартаЗапросов.RUnlock()
// 					}
// 					return &Сообщение
// 				}
// 			}

// 		}
// 	} else {
// 		// Нет УИДЗапроса в Карте Запросов , создадим очередь
// 		Ошибка("Нет очереди запросов %+v но и получается маршрут не смогли построить.... ??? или почему не вышли из функции \n", Сообщение.Запрос.УИДЗапроса)
// 		// тут нужно перейти в начало
// 		return nil
// 	}

// 	Ошибка(" Тут не должно было случится выхода... что то не обработкано ????? %+v \n", err)
// 	return nil
// }

// По УидЗапроса проверяет есть ли в КартеЗАпросов маршрут , и если есть возвращает следующий сервис

// ЧПросто отправляет сообщение в указаный сервис
// func ОтправитьЗапросВСервис(сервис string, Сообщение *Сообщение) {

// 	if Сервис, есть := АктивныеСессии[сервис]; есть {
// 		Инфо(" Сервис %+v \n", Сервис)
// 		// Сервис.КаналСообщения читается в функции ОтправитьСообщениеВСервис()
// 		Сервис.КаналСообщения <- Сообщение
// 	} else {

// 		Ошибка(" Нет маршрута для  %+v есть %+v\n", сервис, есть)
// 		Ошибка("Маршрутизатор  %+v \n", Маршрутизатор)
// 	}

// }
