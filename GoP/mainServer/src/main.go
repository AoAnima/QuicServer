package main

import (
	"net/http"

	_ "net/http/pprof"

	connector "aoanima.ru/connector"
	. "aoanima.ru/logger"

	"github.com/google/uuid"
)

// структура для сохранения в карту клиенты
type Запрос struct {
	Req         http.Request
	КаналОтвета chan ОтветКлиенту
	ИдКлиента   uuid.UUID
	УИДЗапроса  connector.УидЗапроса
}

type ОтветКлиенту struct {
	УИДЗапроса string
	ИдКлиента  uuid.UUID
	Ответ      string
}

func Уид() uuid.UUID {
	// id := uuid.New()
	return uuid.New()
}

func main() {
	/* каналЗапросовОтКлиентов - передаём этот канал в в функци  ЗапуститьСерверТЛС , когда прийдёт сообщение из браузера, функция обработчик запишет данные в этот канал
	 */
	каналЗапросовОтКлиентов := make(chan Запрос, 10)
	/*
	   Запускаем сервер передаём в него канал, в который запишем обработанный запрос из браузера
	*/
	go ЗапуститьСерверТЛС(каналЗапросовОтКлиентов)

	Инфо(" %s", "запустили сервер")
	/* Инициализирум сервисы коннектора передадим в них канал, из которого Коннектор будет читать сообщение, и отправлять его в synqTCP  */
	go ИнициализацияСервисов(каналЗапросовОтКлиентов)

	ЗапуститьСервер()

}

func ИнициализацияСервисов(каналЗапросовОтКлиентов chan Запрос) {
	Сервис := connector.Отпечаток{
		Сервис: "КлиентСервер",
		Маршруты: map[string]*connector.СтруктураМаршрута{
			"ОтветКлиенту": {
				Запрос: map[string]interface{}{
					"ТипЗпроса":     "int",    // в заивисмости от типа запроса например ajax или обычный request будет возвращён ответ...
					"СтрокаЗапроса": "string", // url Path Query
					"Форма":         "map[string][]string",
					"Файл":          "string",
				},
				Ответ: map[string]interface{}{
					"HTML": "string",
					"JSON": "string",
				},
			},
		},
	}
	каналПолученияСообщений, каналОтправкиСообщений := connector.ИнициализацияСервиса("localhost",
		connector.ПортДляОтправкиСообщений,
		connector.ПортДляПолученияСообщений,
		Сервис)

	// for ЗапросОтКлиента := range каналЗапросовОтКлиентов {

	go ОтправитьЗапросВОбработку(каналОтправкиСообщений, каналЗапросовОтКлиентов)
	go ОтправитьОтветКлиенту(каналПолученияСообщений)
	// ОбработатьСообщение(ЗапросОтКлиента, каналОтправкиСообщений)
	// каналОтправкиСообщений <- Кодировать(ЗапросОтКлиента)
	// }

}

// обработчик сообщений от synqTCP
func ОбработатьСообщение(ВходящееСообщение Запрос, каналОтправкиСообщений chan<- []byte) {

	Инфо(" ОбработатьСообщение %+v \n", ВходящееСообщение)

	Сообщение := Кодировать(ВходящееСообщение)

	// TODO: Реализум логику обработки запроса от клиента, и генерацию ответа

	каналОтправкиСообщений <- Сообщение
}

func ЗапуститьСерверТЛС(каналЗапросов chan<- Запрос) {

	err := http.ListenAndServeTLS(":443",
		"cert/server.crt",
		"cert/server.key",
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				обработчикЗапроса(w, r, каналЗапросов)
			}))

	if err != nil {
		Ошибка(" %s ", err)
	}
}
func обработчикЗапроса(w http.ResponseWriter, req *http.Request, каналЗапросов chan<- Запрос) {

	Инфо(" %s \n", *req)

	каналОтвета := make(chan ОтветКлиенту, 10)

	// отправляем сообщение в функцию  ОтправитьЗапросВОбработку
	каналЗапросов <- Запрос{
		Req:         *req,
		КаналОтвета: каналОтвета,
	}

	// /*  Тут мы читаем из канала  каналОтвета кторый храниться в карте клиенты , данные пишутся в канал  в функции ОтправитьОтветКлиенту */
	for данныеДляОтвета := range каналОтвета {
		if данныеДляОтвета.Ответ != "" {
			Инфо(" данныеДляОтвета.Ответ %+v \n", данныеДляОтвета.Ответ)

			if f, ok := w.(http.Flusher); ok {
				i, err := w.Write([]byte(данныеДляОтвета.Ответ))
				Инфо("  %+v \n", i)
				if err != nil {
					Ошибка(" %s ", err)
				}
				f.Flush()
				break
			}
		}
	}

}

// func ОбработчикОтветов(w http.ResponseWriter, каналОтветов <-chan Ответ) {

// 	Ответ := <-каналОтветов
// 	if Ответ.Сообщение != nil {
// 		w.Write([]byte(Ответ.Сообщение.(string)))
// 	}

// }

func ЗапуститьСервер() {
	err := http.ListenAndServe(":80", http.HandlerFunc(
		func(w http.ResponseWriter, req *http.Request) {
			// 	Инфо(" %s  %s \n", w, req)
			http.Redirect(w, req, "https://localhost:443"+req.RequestURI, http.StatusMovedPermanently)
		},
	))
	// err := http.ListenAndServe(":6060", nil)
	if err != nil {
		Ошибка(" %s ", err)
	}
	go func() {
		http.ListenAndServe("localhost:6060", nil)
	}()
}