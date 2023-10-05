package main

import (
	"crypto/tls"
	"crypto/x509"
	"io"
	"net/http"
	"strconv"
	"time"

	"os"

	. "aoanima.ru/logger"
)

var (
	PORT    = 81
	TLSPort = 444
)

func main() {

	go ListenAndServeTLS()

	Инфо(" %s", "запустили сервер")

	ListenAndServe()

}

func ListenAndServeTLS() {

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

func обработчикЗапроса(w http.ResponseWriter, req *http.Request) {
	// Инфо(" %s  %s \n", w, *req)
	// АнализЗапроса(w, req)

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Content-Type", "text/event-stream")

	Инфо("RemoteAddr %+v   req %+v  \n", req.RemoteAddr, req.Body)
	Инфо("req ContentLength %+v \n", req.ContentLength)
	_, err := io.ReadAll(req.Body)
	if err != nil {
		http.Error(w, "Ошибка чтения тела запроса", http.StatusBadRequest)
		return
	}
	// defer req.Body.Close()
	nx := 0

	w.Write([]byte("первый Ответ на запрос"))

	ctx := req.Context()

	for {
		select {
		case <-ctx.Done():
			Инфо("Запрос был отменен, останавливаем отправку данных %+v", "1")
			// Запрос был отменен, останавливаем отправку данных

			return
		default:
			nx++
			if f, ok := w.(http.Flusher); ok {
				Инфо(" %+v  отправляем %+v \n", req.RemoteAddr, nx)
				w.Write([]byte("Ответ на запрос " + strconv.Itoa(nx) + "\n"))
				f.Flush()
				time.Sleep(2 * time.Second)
				// req.Body.Close()
			} else {
				// Соединение не установлено, обработка ошибки
				Ошибка("Соединение не установлено %+v", http.StatusInternalServerError)
			}
		}

	}

	w.Write([]byte("Ответ на запрос"))

}
