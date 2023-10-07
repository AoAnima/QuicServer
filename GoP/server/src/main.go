package main

import (
	"net/http"

	. "aoanima.ru/logger"
)

type Запрос struct {
	ИдКлиента   string
	Req         http.Request
	Запрос      interface{}
	КаналОтвета chan ОтветКлиенту
}
type Ответ struct {
	Сообщение interface{}
}

func main() {

	каналЗапросов := make(chan Запрос, 10)
	каналОтветов := make(chan Ответ, 10)
	go ListenAndServeTLS(каналЗапросов, каналОтветов)

	Инфо(" %s", "запустили сервер")

	go ИнициализацияСервисов(каналЗапросов, каналОтветов)

	ListenAndServe()

}

func ИнициализацияСервисов(каналЗапросов chan Запрос, каналОтветов <-chan Ответ) {
	go ПодключитсяКМенеджеруЗапросов(каналЗапросов, каналОтветов)
}

func ListenAndServeTLS(каналЗапросов chan<- Запрос, каналОтветов <-chan Ответ) {

	err := http.ListenAndServeTLS(":443",
		"cert/server.crt",
		"cert/server.key",
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				обработчикЗапроса(w, r, каналЗапросов, каналОтветов)
			}))

	if err != nil {
		Ошибка(" %s ", err)
	}
}
func обработчикЗапроса(w http.ResponseWriter, req *http.Request, каналЗапросов chan<- Запрос, каналОтветов <-chan Ответ) {

	Инфо(" %s \n", *req)
	// w.Write([]byte("Your response content"))
	каналОтвета := make(chan ОтветКлиенту)
	каналЗапросов <- Запрос{
		ИдКлиента:   Уид(),
		Req:         *req,
		Запрос:      req.URL.String(),
		КаналОтвета: каналОтвета,
	}
	for данныеДляОтвета := range каналОтвета {
		if данныеДляОтвета.Ответ != nil {
			Инфо(" данныеДляОтвета.Ответ %+v \n", данныеДляОтвета.Ответ)

			if f, ok := w.(http.Flusher); ok {
				i, err := w.Write([]byte(данныеДляОтвета.Ответ.(string)))
				Инфо("  %+v \n", i)
				if err != nil {
					// Handle error
				}
				f.Flush()
				break
			}
		}
	}
	// for {
	// select {
	// case данныеДляОтвета := <-каналОтвета:
	// 	if данныеДляОтвета.Ответ != nil {

	// 		Инфо(" данныеДляОтвета.Ответ %+v \n", данныеДляОтвета.Ответ)

	// 		if f, ok := w.(http.Flusher); ok {
	// 			i, err := w.Write([]byte(данныеДляОтвета.Ответ.(string)))
	// 			Инфо("  %+v \n", i)
	// 			if err != nil {
	// 				Ошибка("  %+v \n", err)
	// 			}
	// 			f.Flush()

	// 		}
	// 	}
	// }
	// }

	// for {
	// 	данныеДляОтвета := <-каналОтвета
	// 	Инфо(" данныеДляОтвета %+v \n", данныеДляОтвета)
	// 	if данныеДляОтвета.Ответ != nil {

	// 		Инфо(" данныеДляОтвета.Ответ %+v \n", данныеДляОтвета.Ответ)

	// 		if f, ok := w.(http.Flusher); ok {
	// 			i, err := w.Write([]byte(данныеДляОтвета.Ответ.(string)))
	// 			Инфо("  %+v \n", i)
	// 			if err != nil {
	// 				Ошибка("  %+v \n", err)
	// 			}
	// 			f.Flush()
	// 			break
	// 		}
	// 		// i, err := w.Write([]byte(данныеДляОтвета.Ответ.(string)))
	// 		// Инфо("  %+v \n", i)
	// 		// if err != nil {
	// 		// 	Ошибка("  %+v \n", err)
	// 		// }
	// 		// break
	// 	}
	// }
	// for {
	// 	if f, ok := w.(http.Flusher); ok {
	// 		w.Write([]byte("Ответ на запрос "))
	// 		f.Flush()
	// 	}
	// }
	// for данныеДляОтвета := range каналОтвета {
	// 	Инфо(" данныеДляОтвета %+v \n", данныеДляОтвета)
	// 	if данныеДляОтвета.Ответ != nil {

	// 		Инфо(" данныеДляОтвета.Ответ %+v \n", данныеДляОтвета.Ответ)

	// 		w.Write([]byte(данныеДляОтвета.Ответ.(string)))
	// 	}
	// }
	// var ok bool
	// for ok {
	// 	f, ok := w.(http.Flusher)
	// 	Инфо(" ok= %+v f %+v \n", ok, f)
	// }
	// ОбработчикОтветов(w, каналОтветов)
}

func ОбработчикОтветов(w http.ResponseWriter, каналОтветов <-chan Ответ) {

	Ответ := <-каналОтветов
	if Ответ.Сообщение != nil {
		w.Write([]byte(Ответ.Сообщение.(string)))
	}

}

func ListenAndServe() {
	err := http.ListenAndServe(":80", nil)
	// err := http.ListenAndServe(":80", http.HandlerFunc(

	// 	func(w http.ResponseWriter, req *http.Request) {
	// 		Инфо(" %s  %s \n", w, req)
	// 		// http.Redirect(w, req, "https://localhost:443"+req.RequestURI, http.StatusMovedPermanently)
	// 	}))

	if err != nil {
		Ошибка(" %s ", err)
	}
}

// func Рендер(каналеРендера chan interface{}) {
// 	Инфо(" %s  \n", "Рендер")
// 	каналОтправкиДанных := make(chan interface{}, 10)
// 	go СоденитьсяССервисомРендера(каналОтправкиДанных)
// 	// go СоденитьсяССервисомРендера(каналОтправкиДанных)

// 	for {
// 		if данныеДляРендера := <-каналеРендера; данныеДляРендера != nil {
// 			Инфо(" %s  \n", данныеДляРендера)
// 		}
// 	}

// }

//
