package main

import (
	"crypto/tls"
	"crypto/x509"
	"net/http"
	"os"

	. "aoanima.ru/logger"
)

func main() {

	go ListenAndServeTLS()

	Инфо(" %s", "запустили рендер сервер")

	ListenAndServe()

}

type Writer interface {
	Write(p []byte) (n int, err error)
}

type Ty struct{}

func ListenAndServeTLS() {
	caCert, err := os.ReadFile("cert/client.crt")
	if err != nil {
		Ошибка(" %s", err)
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)
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
func ListenAndServe() {
	err := http.ListenAndServe(":81", http.HandlerFunc(

		func(w http.ResponseWriter, req *http.Request) {
			Инфо(" %s  %s \n", w, req)
			http.Redirect(w, req, "https://localhost:80"+req.RequestURI, http.StatusMovedPermanently)
		}))

	if err != nil {
		Ошибка(" %s ", err)
	}
}

func обработчикЗапроса(w http.ResponseWriter, req *http.Request) {
	// Инфо(" %s  %s \n", w, *req)
	Инфо("  %s \n", *req)
}
