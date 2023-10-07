package main

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	
	"sync"

	"os"

	. "aoanima.ru/logger"
	"github.com/dgrr/fastws"
	"github.com/google/uuid"
)

// https://github.com/jcbsmpsn/golang-https-example/blob/master/https_client.go
// func СоденитьсяССервисомРендера(каналОтправкиДанных chan interface{}) {
// 	Инфо(" %s", "СоденитьсяССервисомРендера")

// 	caCert, err := os.ReadFile("cert/ca.crt")

// 	if err != nil {
// 		Ошибка(" %s ", err)
// 	}

// 	caCertPool := x509.NewCertPool()
// 	ok := caCertPool.AppendCertsFromPEM(caCert)
// 	Инфо("Корневой сертфикат создан?  %v ", ok)

// 	cert, err := tls.LoadX509KeyPair("cert/client.crt", "cert/client.key")
// 	if err != nil {
// 		Ошибка(" %s", err)
// 	}

// 	tlsConfig := &tls.Config{
// 		RootCAs:      caCertPool,
// 		Certificates: []tls.Certificate{cert},
// 	}

// 	transport := &http.Transport{
// 		TLSClientConfig:     tlsConfig,
// 		MaxIdleConnsPerHost: 10,
// 		ForceAttemptHTTP2:   true,
// 		IdleConnTimeout:     0,
// 	}

// 	client := &http.Client{
// 		Transport: transport,
// 	}
// 	size := 1024

// 	// Создаем срез указанного размера
// 	slice := make([]byte, size)

// 	// Заполняем срез случайными данными
// 	_, err = rand.Read(slice)
// 	if err != nil {
// 		Ошибка(" %s", err)
// 	}
// 	// base64.StdEncoding.EncodeToString(slice)
// 	req, err := http.NewRequest("POST", "https://localhost:444", bytes.NewBuffer(slice))
// 	if err != nil {
// 		Ошибка(" %s", err)
// 	}
// 	req.Header.Set("Content-Type", "text/html; charset=utf-8")

// 	for i := 1; i <= 3; i++ {
// 		// n := i
// 		go func() {

// 			resp, err := client.Do(req)
// 			if err != nil {
// 				Ошибка(" %+v", resp, err)
// 				return
// 			}
// 			// defer resp.Body.Close()
// 			// go func() {
// 			for {
// 				reader := bufio.NewReader(resp.Body)
// 				line, err := reader.ReadString('\n')
// 				if err != nil {
// 					// Ошибка("Error:%+v\n", err)
// 					if err == io.EOF {
// 						Ошибка("Error:%+v\n", err)

// 						// break
// 					}
// 					return
// 				}

// 				Инфо("lines  %+v", string(line))
// 			}
// 			// }()
// 			// defer resp.Body.Close()
// 			// body, err := io.ReadAll(resp.Body)
// 			// if err != nil {
// 			// 	Ошибка("Error:%+v\n", err)
// 			// }
// 			// Инфо("body :  %+v", string(body))
// 		}()
// 	}

// }



var клиенты = make(map[string]Запрос)
var мьютекс = sync.Mutex{}

func Уид() string {
	id := uuid.New()
	return id.String()
}

func ПодключитсяКМенеджеруЗапросов(каналЗапросов chan Запрос) {
	// устанавливаем websocket соединение с сервисом
	// вести счётчик количество отправленных запросов, и количество полученных ответов,
	// настроить метрики и понять сколько одновременных запросов может обрабатывать одно содеинение
	// если количество запросов выходит за пределы вохмоэностей одного соединения - то можно создать новое соединение
	// при уменьшении нагрузки, лишние соединения можно закрывать
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

	conn, err := fastws.DialTLS("wss://localhost:444", tlsConfig)
	if err != nil {
		Ошибка("  %+v \n", err)
	}
	Инфо("  %+v \n", conn)
	// conn.WriteString("Hello")

	// каналЗапросов - исползуется для получения запросов от клиента, в запросе от клиента передаётся канал в который нужно отправить 
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
type ОтветКлиенту struct {
	ИдКлиента string
	Ответ     interface{}
}

// читаем сообщение от сервера менеджера сообщений, и отправляем ответ клиенту
func ПрочитатьСообщение(conn *fastws.Conn, каналЗапросов chan Запрос) {
	Инфо(" ПрочитатьСообщение \n")
	var сообщ []byte
	var ответ ОтветКлиенту
	for {

		frame, err := conn.NextFrame()
		// _, сообщение, err := conn.ReadMessage(сообщ[:0])
		Инфо("сообщение: %s  %s \n", сообщ, frame)
		// if сообщение != nil {

		if err != nil {
			Ошибка("  %+v \n", err)

			if err == fastws.EOF {
				Ошибка(" соединение закрыто походу  %+v \n", err)
				return
			}

		}
		defer fastws.ReleaseFrame(frame)
		if frame != nil {
			Инфо("сообщение: %+v \n", string(frame.Payload()))
			err = json.Unmarshal(frame.Payload(), &ответ)
			if err != nil {
				Ошибка("  %+v \n", err)
			}

			Инфо("сообщение: %+v \n", ответ)
			клиенты[ответ.ИдКлиента].КаналОтвета <- ответ

		}

	}
}

type ЗапросКлиента struct {
	Запрос    interface{}
	ИдКлиента string
}

// Отправляем сообщение на сервер менеджера сообщений
func ОтправитьСообщение(conn *fastws.Conn, каналЗапросов chan Запрос) {
	Инфо(" ОтправитьСообщение \n")
	for {
		запросКлиента := <-каналЗапросов
		Инфо(" ОтправитьСообщение %+v \n", запросКлиента)
		if запросКлиента.Запрос != nil {
			мьютекс.Lock()
			// генерируем новый ид для запроса клиента
			// уид := Уид()
			клиенты[запросКлиента.ИдКлиента] = запросКлиента
			мьютекс.Unlock()
			Инфо(" ОтправитьСообщение Уид %+v \n", Уид)
			// новыйЗапрос := map[string]string{
			// 	"Запрос": сообщение.Запрос.(string),
			// 	"ид":     уид,
			// }
			новыйЗапрос := ЗапросКлиента{
				Запрос:    запросКлиента.Запрос,
				ИдКлиента: запросКлиента.ИдКлиента,
			}

			Инфо(" ОтправитьСообщение новыйЗапрос  %+v \n", новыйЗапрос)

			запросБайт, err := json.Marshal(новыйЗапрос)

			if err != nil {
				Ошибка("v error: %+v", err)
			}

			Инфо(" ОтправитьСообщение в synqest  запрос %+v \n", string(запросБайт))
			_, _, errRead := conn.ReadMessage(nil)
			if errRead != nil {
				Ошибка("errRead : %+v", errRead)
			}
			i, err := conn.WriteMessage(fastws.ModeText, запросБайт)

			if err != nil {
				Ошибка("Write error: %+v %+v", i, err)

			}
		}
	}

}
