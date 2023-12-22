package main

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"net/http"

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
		Handler:   http.HandlerFunc(fastws.NetUpgrade(обработчикВебСокет)),
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

	if upgrade, ok := req.Header["Upgrade"]; ok {
		switch upgrade[0] {
		case "WebSocket":
			Инфо(" Есть заголовок upgrade = %+v \n", upgrade)
			fastws.NetUpgrade(обработчикВебСокет)
		default:
			Инфо("  заголовок upgrade = %+v , есть или отличается о WebSocket \n", upgrade)
			// обработчикЗапроса(w, req)
		}

	}
}

func обработчикВебСокет(conn *fastws.Conn) {
	Инфо(" обработчикВебСокет %+v \n", conn)
	// канал - канал для обмена сообщением между фукнциями ПрочитатьСообщение и ОТправитьСообщение
	// функция прочитать сообщние, читает сообщение из ws соединенияя, обрабатывает, и результат отправляет в канал, функция ОТправитьСообщение читает данные из канала, и отправляет сообщение в ws соединение
	канал := make(chan ОтветКлиенту, 10)

	go ПрочитатьСообщение(conn, канал)
	ОтправитьСообщение(conn, канал)
	Инфо(" %+v \n", "выход из обработчика ")
}

type ЗапросКлиента struct {
	Запрос    interface{}
	ИдКлиента string
}

func ПрочитатьСообщение(conn *fastws.Conn, канал chan ОтветКлиенту) {
	Инфо(" ПрочитатьСообщение \n %+v", conn)
	var сообщ []byte
	for {
		h := conn.ReadTimeout.Hours()
		frame, err := conn.NextFrame()
		// _, сообщение, err := conn.ReadMessage(сообщ[:0])
		Инфо("frame: %s  %s %+v\n", сообщ, frame, h)
		if err != nil {
			Ошибка("  %+v \n", err)

			if err == fastws.EOF {
				Ошибка(" соединение закрыто походу  %+v \n", err)
				return
			}

		}
		defer fastws.ReleaseFrame(frame)

		if frame != nil {
			Инфо("Received: %s", frame.Status().String())
			Инфо("Received: %s", string(frame.Payload()))
			// if сообщение != nil {

			var запрос ЗапросКлиента
			err = json.Unmarshal(frame.Payload(), &запрос)
			if err != nil {
				Ошибка("  %+v \n", err)
			}

			Инфо("сообщение: %s \n", запрос)
			канал <- ОтветКлиенту{
				ИдКлиента: запрос.ИдКлиента,
				Ответ:     запрос.Запрос,
			}
			Инфо("данные отправлены в канал в функцию ОТправитьСОобщение \n")
		}
		// Выведите полученное сообщение.

		// }

		// time.Sleep(time.Second)

	}
}

type ОтветКлиенту struct {
	Ответ     interface{}
	ИдКлиента string
}

func ОтправитьСообщение(conn *fastws.Conn, канал chan ОтветКлиенту) {
	Инфо(" ОтправитьСообщение читаем данные из канала, с результатом обработки зароса\n")
	for ответ := range канал {
		if ответ.Ответ != nil {
			сообщение, err := json.Marshal(ответ)
			if err != nil {
				Ошибка("  %+v \n", err)
			}

			Инфо(" отправляем данные в w соединение %+v \n", сообщение)
			_, _, errRead := conn.ReadMessage(nil)
			if errRead != nil {
				Ошибка("errRead : %+v", errRead)
			}
			_, err = conn.WriteMessage(fastws.ModeBinary, сообщение)
			if err != nil {
				Ошибка("Write error: %+v", err)

			}
		}

	}
}

// func обработчикЗапроса(w http.ResponseWriter, req *http.Request) {
// 	// Инфо(" %s  %s \n", w, *req)
// 	// АнализЗапроса(w, req)

// 	w.Header().Set("Access-Control-Allow-Origin", "*")
// 	w.Header().Set("Cache-Control", "no-cache")
// 	w.Header().Set("Connection", "keep-alive")
// 	w.Header().Set("Content-Type", "text/event-stream")

// 	// прочитать тело запроса
// 	body, err := io.ReadAll(req.Body)
// 	if err != nil {
// 		http.Error(w, "Ошибка чтения тела запроса", http.StatusBadRequest)
// 		return
// 	}

// 	Инфо("RemoteAddr %+v   body %+v  \n", req.RemoteAddr, body)
// 	Инфо("req ContentLength %+v \n", req.ContentLength)

// 	// defer req.Body.Close()
// 	nx := 0

// 	w.Write([]byte("первый Ответ на запрос"))

// 	ctx := req.Context()

// 	for {
// 		select {
// 		case <-ctx.Done():
// 			Инфо("Запрос был отменен, останавливаем отправку данных %+v", "1")
// 			// Запрос был отменен, останавливаем отправку данных

// 			return
// 		default:
// 			nx++
// 			if f, ok := w.(http.Flusher); ok {
// 				Инфо(" %+v  отправляем %+v \n", req.RemoteAddr, nx)
// 				w.Write([]byte("Ответ на запрос " + strconv.Itoa(nx) + "\n"))
// 				f.Flush()
// 				time.Sleep(2 * time.Second)
// 				if nx == 5 {
// 					break
// 				}
// 				// req.Body.Close()
// 			} else {
// 				// Соединение не установлено, обработка ошибки
// 				Ошибка("Соединение не установлено %+v", http.StatusInternalServerError)
// 			}
// 		}

// 	}

// 	w.Write([]byte("Ответ на запрос в конце"))

// }
