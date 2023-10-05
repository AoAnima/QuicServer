package main

import (
	http "github.com/valyala/fasthttp"
	"github.com/dgrr/fastws"
	. "aoanima.ru/logger"
)

func main() {
	брокерСообщений := make(chan interface{}, 10)
	go ListenAndServeTLS(брокерСообщений)

	Инфо(" %s", "запустили сервер")

	go Рендер(брокерСообщений)

	ListenAndServe()

}

type Writer interface {
	Write(p []byte) (n int, err error)
}

type Ty struct{}

func ListenAndServeTLS(брокерСообщений chan interface{}) {

	err := http.ListenAndServeTLS(":443", "cert/server.crt", "cert/server.key", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		обработчикЗапроса(w, r, брокерСообщений)
	}))

	if err != nil {
		Ошибка(" %s ", err)
	}
}
func обработчикЗапроса(w http.ResponseWriter, req *http.Request, брокерСообщений chan interface{}) {
	// Инфо(" %s  %s \n", w, *req)
	// АнализЗапроса(w, req)
	Инфо(" %s \n", *req)
	брокерСообщений <- *req
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

func Рендер(каналеРендера chan interface{}) {
	Инфо(" %s  \n", "Рендер")
	каналОтправкиДанных := make(chan interface{}, 10)
	go СоденитьсяССервисомРендера(каналОтправкиДанных)
	go СоденитьсяССервисомРендера(каналОтправкиДанных)

	for {
		if данныеДляРендера := <-каналеРендера; данныеДляРендера != nil {
			Инфо(" %s  \n", данныеДляРендера)
		}
	}

}

//
