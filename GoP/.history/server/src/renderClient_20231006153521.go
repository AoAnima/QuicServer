package main

import (
	"bufio"
	"bytes"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"sync"

	"os"

	. "aoanima.ru/logger"
	"github.com/dgrr/fastws"
	"github.com/google/uuid"
)

// https://github.com/jcbsmpsn/golang-https-example/blob/master/https_client.go
func СоденитьсяССервисомРендера(каналОтправкиДанных chan interface{}) {
	Инфо(" %s", "СоденитьсяССервисомРендера")

	caCert, err := os.ReadFile("cert/ca.crt")

	if err != nil {
		Ошибка(" %s ", err)
	}

	caCertPool := x509.NewCertPool()
	ok := caCertPool.AppendCertsFromPEM(caCert)
	Инфо("Корневой сертфикат создан?  %v ", ok)

	cert, err := tls.LoadX509KeyPair("cert/client.crt", "cert/client.key")
	if err != nil {
		Ошибка(" %s", err)
	}

	tlsConfig := &tls.Config{
		RootCAs:      caCertPool,
		Certificates: []tls.Certificate{cert},
	}

	transport := &http.Transport{
		TLSClientConfig:     tlsConfig,
		MaxIdleConnsPerHost: 10,
		ForceAttemptHTTP2:   true,
		IdleConnTimeout:     0,
	}

	client := &http.Client{
		Transport: transport,
	}
	size := 1024

	// Создаем срез указанного размера
	slice := make([]byte, size)

	// Заполняем срез случайными данными
	_, err = rand.Read(slice)
	if err != nil {
		Ошибка(" %s", err)
	}
	// base64.StdEncoding.EncodeToString(slice)
	req, err := http.NewRequest("POST", "https://localhost:444", bytes.NewBuffer(slice))
	if err != nil {
		Ошибка(" %s", err)
	}
	req.Header.Set("Content-Type", "text/html; charset=utf-8")

	for i := 1; i <= 3; i++ {
		// n := i
		go func() {
			// if n == 2 {
			// 	req.Body = io.NopCloser(bytes.NewBuffer([]byte("Новое тело запроса")))
			// }

			resp, err := client.Do(req)
			if err != nil {
				Ошибка(" %+v", resp, err)
				return
			}
			// defer resp.Body.Close()
			// go func() {
			for {
				reader := bufio.NewReader(resp.Body)
				line, err := reader.ReadString('\n')
				if err != nil {
					// Ошибка("Error:%+v\n", err)
					if err == io.EOF {
						Ошибка("Error:%+v\n", err)

						// break
					}
					return
				}

				Инфо("lines  %+v", string(line))
			}
			// }()
			// defer resp.Body.Close()
			// body, err := io.ReadAll(resp.Body)
			// if err != nil {
			// 	Ошибка("Error:%+v\n", err)
			// }
			// Инфо("body :  %+v", string(body))
		}()
	}

}

type Клиент struct {
	ид             string
	ResponseWriter http.ResponseWriter
}

var клиенты = make(map[string]Клиент)
var мьютекс = sync.Mutex{}

func Уид() string {
	id := uuid.New()
	return id.String()
}

func ПодключитсяКМенеджеруЗапросов(каналЗапросов chan Запрос, каналОтветов <-chan Ответ) {
	// устанавливаем websocket соединение с сервисом
	// вести счётчик количество отправленных запросов, и количество полученных ответов,
	// настроить метрики и понять сколько одновременных запросов может обрабатывать одно содеинение
	// если количество запросов выходит за пределы вохмоэностей одного соединения - то можно создать новое соединение
	// при уменьшении нагрузки, лишние соединения можно закрывать
	caCert, err := os.ReadFile("cert/ca.crt")

	if err != nil {
		Ошибка(" %s ", err)
	}
	strconv.Itoa()
	caCertPool := x509.NewCertPool()
	ok := caCertPool.AppendCertsFromPEM(caCert)
	Инфо("Корневой сертфикат создан?  %v ", ok)

	cert, err := tls.LoadX509KeyPair("cert/client.crt", "cert/client.key")
	if err != nil {
		Ошибка(" %s", err)
	}

	tlsConfig := &tls.Config{
		RootCAs:      caCertPool,
		Certificates: []tls.Certificate{cert},
	}

	conn, err := fastws.DialTLS("wss://localhost:444/echo", tlsConfig)
	if err != nil {
		Ошибка("  %+v \n", err)
	}
	Инфо("  %+v \n", conn)
	// conn.WriteString("Hello")

	go ОтправитьСообщение(conn, каналЗапросов)

	go ПрочитатьСообщение(conn, каналЗапросов)
	// var msg []byte
	// // for i := 0; i < 5; i++ {
	// _, msg, err = conn.ReadMessage(msg)
	// if err != nil {
	// 	Ошибка("  %+v \n", err)
	// 	// break
	// }
	// Инфо("Client: %s \n", msg)
	// conn.Write([]byte(" и тебе привет "))
	// // 	// time.Sleep(time.Second)
	// // }
	// _, msg, err = conn.ReadMessage(msg)
	// if err != nil {
	// 	Ошибка("  %+v \n", err)
	// 	// break
	// }
	// Инфо("Client: %s \n", msg)
}

// читаем сообщение с ответом, от сервера менеджера сообщений для отправки клиенту

func ПрочитатьСообщение(conn *fastws.Conn, каналЗапросов chan Запрос) {
	Инфо(" ОтправитьСообщение \n")
	for {
		var msg []byte
		_, message, err := conn.ReadMessage(msg)
		if err != nil {
			Ошибка("%+v", err)

		}
		json
		Инфо("message: %s \n", message)
	}
}

// Отправляем сообщение на сервер менеджера сообщений
func ОтправитьСообщение(conn *fastws.Conn, каналЗапросов chan Запрос) {
	Инфо(" ОтправитьСообщение \n")
	message := <-каналЗапросов
	if message.Запрос == nil {
		мьютекс.Lock()
		// генерируем новый ид для запроса клиента
		уид := Уид()
		клиенты[уид] = Клиент{
			ResponseWriter: message.ResponseWriter,
			ид:             уид,
		}

		запрос, err := json.Marshal(
			map[string]string{
				"запрос":message.Запрос.(string),
				"ид": уид,
			})

		_, err = conn.WriteMessage(fastws.ModeBinary, запрос)
		if err != nil {
			Ошибка("Write error: %+v", err)

		}
	}

}
