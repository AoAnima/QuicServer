package main

import (
	"net/http"

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

func ListenAndServeTLS(каналеРендера chan interface{}) {

	err := http.ListenAndServeTLS(":443", "cert/cert.pem", "cert/key.pem", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        обработчикЗапроса(w, r, каналеРендера)
    }))

	if err != nil {
		Ошибка(" %s ", err)
	}
}
func обработчикЗапроса(w http.ResponseWriter, req *http.Request, каналеРендера chan interface{}) {
	// Инфо(" %s  %s \n", w, *req)
	// АнализЗапроса(w, req)
	Инфо(" %s \n", *req)
	каналеРендера <- *req
}

func ListenAndServe() {
	err := http.ListenAndServe(":80", http.HandlerFunc(

		func(w http.ResponseWriter, req *http.Request) {
			Инфо(" %s  %s \n", w, req)
			// http.Redirect(w, req, "https://localhost:443"+req.RequestURI, http.StatusMovedPermanently)
		}))

	if err != nil {
		Ошибка(" %s ", err)
	}
}



func Рендер(каналеРендера chan interface{}) {
	Инфо(" %s  \n", "Рендер")
	каналОтправкиДанных := make(chan interface{}, 10)
	go СоденитьсяССервисомРендера(каналОтправкиДанных)

	for {
		if данныеДляРендера := <-каналеРендера; данныеДляРендера != nil {
			Инфо(" %s  \n", данныеДляРендера)
		}
	}

}


// openssl x509 \
//     -req \
//     -days 365 \
//     -sha256 \
//     -in server.csr \
//     -CA ca.crt \
//     -CAkey ca.key \
//     -CAcreateserial \
//     -out server.crt \
//     -extfile <(echo subjectAltName = DNS:localhost,IP:127.0.0.1)