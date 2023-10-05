package main

import (
	"crypto/tls"
	"crypto/x509"
	"io"
	"net/http"

	"os"

	. "aoanima.ru/logger"
)

var (
	PORT    = 81
	TLSPort = 444
)

func main() {
	брокерСообщений := make(chan interface{}, 10)
	go ListenAndServeTLS(брокерСообщений)

	Инфо(" %s", "запустили сервер")

	ListenAndServe()

}

type Writer interface {
	Write(p []byte) (n int, err error)
}

type Ty struct{}

func ListenAndServeTLS(брокерСообщений chan interface{}) {

	caCert, err := os.ReadFile("cert/ca.crt")
	if err != nil {
		Ошибка(" %s", err)
	}
	caCertPool := x509.NewCertPool()
	ok := caCertPool.AppendCertsFromPEM(caCert)
	Инфо("Корневой сертфикат создан?  %v ", ok)
	cfg := &tls.Config{
		ClientAuth: tls.RequireAndVerifyClientCert,
		ClientCAs:  caCertPool,
	}
	srv := &http.Server{
		Addr:      ":444",
		Handler:   http.HandlerFunc(обработчикЗапроса),
		TLSConfig: cfg,
	}

	err = srv.ListenAndServeTLS("cert/server.crt", "cert/server.key")
	if err != nil {
		Ошибка(" %s ", err)
	}
}
func обработчикЗапроса(w http.ResponseWriter, req *http.Request) {
	// Инфо(" %s  %s \n", w, *req)
	// АнализЗапроса(w, req)

	Инфо("RemoteAddr %s \n", req.RemoteAddr)
	Инфо("req body %s \n", req.le)

	
	w.Write([]byte("Ответ на запрос"))

}

func ListenAndServe() {
	err := http.ListenAndServe(":81", http.HandlerFunc(

		func(w http.ResponseWriter, req *http.Request) {
			Инфо(" %s  %s \n", w, req)
			// http.Redirect(w, req, "https://localhost:443"+req.RequestURI, http.StatusMovedPermanently)
		}))

	if err != nil {
		Ошибка(" %s ", err)
	}
}
