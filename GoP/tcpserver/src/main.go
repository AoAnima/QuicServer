package main

import (
	"net/http"

	_ "net/http/pprof"

	. "aoanima.ru/logger"
	"github.com/google/uuid"
)

type Запрос struct {
	ИдКлиента   uuid.UUID
	Req         http.Request
	Запрос      string
	КаналОтвета chan ОтветКлиенту
}

type Ответ struct {
	Сообщение interface{}
}
type ОтветКлиенту struct {
	ИдКлиента uuid.UUID
	Ответ     string
}

func Уид() uuid.UUID {
	// id := uuid.New()
	return uuid.New()
}

func main() {
	go func() {
		http.ListenAndServe("localhost:6060", nil)
	}()
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

	каналОтвета := make(chan ОтветКлиенту)
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

	каналЗапросов <- Запрос{
		ИдКлиента:   ИД,
		Req:         *req,
		Запрос:      req.URL.String(),
		КаналОтвета: каналОтвета,
	}
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
	if err != nil {
		Ошибка(" %s ", err)
	}
}
