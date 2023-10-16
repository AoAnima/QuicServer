package main

import (
	"net/http"

	_ "net/http/pprof"

	. "aoanima.ru/logger"
	"github.com/google/uuid"
)

type Запрос struct {
	Сервис      []byte
	ИдКлиента   uuid.UUID
	Req         http.Request
	Запрос      ЗапросОтКлиента
	КаналОтвета chan ОтветКлиенту
}

type ЗапросОтКлиента struct {
	Строка string
	Форма  map[string][]string
	Файл   string
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

	// парсим пост put и прочие подобные запросы и запихуем данные в url с учётом данных get строки
	Инфо(" req.Method %+v \n", req.Method)
	запросОтКлиента := ЗапросОтКлиента{
		Строка: req.URL.String(),
	}
	if req.Method == http.MethodPost {
		contentType := req.Header.Get("Content-Type")
		if contentType == "multipart/form-data" {
			Инфо("нужно реализовать декодирование, я так понимаю тут передаются файлы через форму %+v \n", "multipart/form-data")
		}

		req.ParseForm()
		запросОтКлиента.Форма = req.Form
	}

	каналЗапросов <- Запрос{
		ИдКлиента:   ИД,
		Req:         *req,
		Запрос:      запросОтКлиента,
		КаналОтвета: каналОтвета,
		Сервис:      []byte("КлиентСервер"),
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
