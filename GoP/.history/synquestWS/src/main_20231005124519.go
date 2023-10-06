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
	"github.com/dgrr/fastws"
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

func ServerWSS() {}

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
		Handler:   http.HandlerFunc(обработчикСоединений),
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

func обработчикСоединений(w http.ResponseWriter, req *http.Request) {
	Инфо(" обработчикСоединений req %+v \n", req.Header["Upgrade"])
	// проверить что заголвок Upgrade есть в запросе
	// fastws.Upgrade(обработчикВебСокет)
	fastws.NetUpgrade(обработчикВебСокет)
	if upgrade, ok := req.Header["Upgrade"]; ok {
		switch upgrade[0] {
		case "WebSocket":
			Инфо(" Есть заголовок upgrade = %+v \n", upgrade)
			fastws.Upgrade(обработчикВебСокет)
		default:
			Инфо("  заголовок upgrade = %+v , есть или отличается о WebSocket \n", upgrade)
			обработчикЗапроса(w, req)
		}

	}
}

func обработчикВебСокет(conn *fastws.Conn) {
	conn.WriteString("Hello")
	var msg []byte
	var err error
	for {
		_, msg, err = conn.ReadMessage(msg[:0])
		if err != nil {
			if err != fastws.EOF {
				Ошибка(" %+v \n", err)
			}
			break
		}
		time.Sleep(time.Second)

		_, err = conn.Write(msg)
		if err != nil {
			Ошибка(" %+v \n", err)
			break
		}
	}
}

func обработчикЗапроса(w http.ResponseWriter, req *http.Request) {
	// Инфо(" %s  %s \n", w, *req)
	// АнализЗапроса(w, req)

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Content-Type", "text/event-stream")

	// прочитать тело запроса
	body, err := io.ReadAll(req.Body)
	if err != nil {
		http.Error(w, "Ошибка чтения тела запроса", http.StatusBadRequest)
		return
	}

	Инфо("RemoteAddr %+v   body %+v  \n", req.RemoteAddr, body)
	Инфо("req ContentLength %+v \n", req.ContentLength)

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
				if nx == 5 {
					break
				}
				// req.Body.Close()
			} else {
				// Соединение не установлено, обработка ошибки
				Ошибка("Соединение не установлено %+v", http.StatusInternalServerError)
			}
		}

	}

	w.Write([]byte("Ответ на запрос в конце"))

}
