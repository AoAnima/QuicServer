package main

import (
	"fmt"
	"net/http"
	"time"

	_ "net/http/pprof"

	. "aoanima.ru/logger"
	"github.com/dgryski/go-metro"
	"github.com/google/uuid"
)

type Запрос struct {
	УИДЗапроса   string // УИД запроса, будет складывать из временной метки, УИД клиента и хэш фукнцией от запроса
	Сервис       []byte
	ИдКлиента    uuid.UUID
	ТокенКлиента []byte // JWT сериализованный
	Req          http.Request
	УрлПуть      []byte
	Запрос       ЗапросОтКлиента
	КаналОтвета  chan ОтветКлиенту
}

type ЗапросОтКлиента struct {
	СтрокаЗапроса string
	Форма         map[string][]string
	Файл          string
}

type Ответ struct {
	Сообщение interface{}
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

	каналЗапросов := make(chan Запрос, 10)

	go ListenAndServeTLS(каналЗапросов)

	Инфо(" %s", "запустили сервер")

	go ИнициализацияСервисов(каналЗапросов)

	ListenAndServe()

}

func ИнициализацияСервисов(каналЗапросов chan Запрос) {
	go ПодключитсяКМенеджеруЗапросов(каналЗапросов)

}

func ListenAndServeTLS(каналЗапросов chan<- Запрос) {

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
	var ИД uuid.UUID

	if cookieUuid, err := req.Cookie("uuid"); err != nil {
		Ошибка(" %+v \n", err)
		ИД = Уид()
	} else {
		ИД, err = uuid.Parse(cookieUuid.Value)
		if err != nil {
			Ошибка("  %+v \n", err)
		}
	}

	// парсим пост put и прочие подобные запросы и запихуем данные в url с учётом данных get строки
	Инфо(" req.Method %+v \n", req.Method)
	Инфо(" req %+v \n", req)
	запросОтКлиента := ЗапросОтКлиента{
		СтрокаЗапроса: req.URL.String(),
	}
	типДанных := req.Header.Get("Content-Type")
	if req.Method == http.MethodPost {
		if типДанных == "multipart/form-data" {
			Инфо("нужно реализовать декодирование, я так понимаю тут передаются файлы через форму %+v \n", "multipart/form-data")
		}
		req.ParseForm()
		запросОтКлиента.Форма = req.Form
	}

	if req.Method == "AJAX" || req.Method == "AJAXPost" {
		if типДанных == "application/json" {

		}
	}
	//каналЗапросов читается в функции ОтправитьЗапросВОбработку, которая отправляет данные в synqTCP поэтому если нужно обраьботкть запро сперед отправкой, то его можно либо обрабатывать тут, перед отправкой в каналЗапросов, лобо внутри фкнции ОтправитьЗапросВОбработку перед записью данный в соединение с synqTCp
	хэшЗапроса := metro.Hash64([]byte(запросОтКлиента.СтрокаЗапроса), 0)
	timestamp := time.Now().Unix()
	УИДЗапроса := fmt.Sprintf("%+s.%+s.%+s", timestamp, ИД, хэшЗапроса)
	каналЗапросов <- Запрос{
		УИДЗапроса:  УИДЗапроса,
		ИдКлиента:   ИД,
		Req:         *req,
		Запрос:      запросОтКлиента,
		УрлПуть:     []byte(req.URL.Path),
		КаналОтвета: каналОтвета,
		Сервис:      []byte("КлиентСервер"),
	}
	//WARNING: БЛОКИРОВКА? читаем данные для ответа в отдельном потоке чтобы не
	// go func() {
	for данныеДляОтвета := range каналОтвета {
		if данныеДляОтвета.Ответ != "" {
			Инфо(" данныеДляОтвета.Ответ %+v \n", данныеДляОтвета.Ответ)

			if f, ok := w.(http.Flusher); ok {
				i, err := w.Write([]byte(данныеДляОтвета.Ответ))
				Инфо("  %+v \n", i)
				if err != nil {
					// Handle error
				}
				f.Flush()
				break
			}
		}
	}
	// }()

}

func ОбработчикОтветов(w http.ResponseWriter, каналОтветов <-chan Ответ) {

	Ответ := <-каналОтветов
	if Ответ.Сообщение != nil {
		w.Write([]byte(Ответ.Сообщение.(string)))
	}

}

func ListenAndServe() {
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
