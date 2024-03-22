package main

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"

	. "aoanima.ru/Logger"
)

func handleReverseProxy(w http.ResponseWriter, r *http.Request) {
	// Получаем URL из параметра запроса "query"
	queryURL, err := url.QueryUnescape(r.URL.Query().Get("query"))
	Инфо(" queryURL %+v \n", queryURL)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Создаем новый HTTP-запрос для целевого URL
	targetURL, err := url.Parse(queryURL)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Создаем обратный прокси-запрос
	proxy := http.DefaultTransport
	proxyReq, err := http.NewRequest(r.Method, targetURL.String(), r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Копируем заголовки из исходного запроса
	for header, values := range r.Header {
		for _, value := range values {
			proxyReq.Header.Add(header, value)
		}
	}
	Инфо(" proxyReq %+v \n", proxyReq)

	// Отправляем прокси-запрос и получаем ответ
	resp, err := proxy.RoundTrip(proxyReq)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	// Копируем заголовки ответа и передаем тело ответа клиенту
	// for header, values := range resp.Header {
	// 	for _, value := range values {
	// 		w.Header().Add(header, value)
	// 	}
	// }
	Инфо(" resp %+v \n", resp)
	b, e := json.Marshal(resp.Body)
	if e != nil {
		Ошибка(" ОписаниеОшибки %+v \n", e)
	}
	Инфо(" b %+v \n", b)

	w.WriteHeader(resp.StatusCode)
	w.Write(b)
}

func main() {
	http.HandleFunc("/", handleReverseProxy)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// var customTransport = http.DefaultTransport

// func main() {
// 	// Настройка прокси-сервера

// 	err := http.ListenAndServeTLS(":30443",
// 		"cert/server.crt",
// 		"cert/server.key",
// 		http.HandlerFunc(
// 			func(w http.ResponseWriter, r *http.Request) {
// 				handleRequest(w, r)
// 			}))

// 	// err := http.ListenAndServe(":30443",
// 	// 	http.HandlerFunc(
// 	// 		func(w http.ResponseWriter, r *http.Request) {
// 	// 			обработчикЗапроса(w, r)
// 	// 		}))

// 	if err != nil {
// 		Ошибка(" %s ", err)
// 	}
// 	// server := http.Server{
// 	// 	Addr:    ":8080",
// 	// 	Handler: http.HandlerFunc(handleRequest),
// 	// }

// 	// // Start the server and log any errors
// 	// log.Println("Starting proxy server on :8080")
// 	// err := server.ListenAndServe()
// 	// if err != nil {
// 	// 	log.Fatal("Error starting proxy server: ", err)
// 	// }

// }

// func handleRequest(w http.ResponseWriter, r *http.Request) {
// 	// Create a new HTTP request with the same method, URL, and body as the original request
// 	// targetURL := r.URL
// 	// запрос := r.URL.Query().Get("запрос")

// 	proxyReq, err := http.NewRequest(r.Method, "https://google.com", r.Body)
// 	if err != nil {
// 		http.Error(w, "Error creating proxy request", http.StatusInternalServerError)
// 		return
// 	}

// 	// Copy the headers from the original request to the proxy request
// 	for name, values := range r.Header {
// 		Инфо(" name %+v  values %+v \n", name, values)

// 		for _, value := range values {
// 			Инфо(" name1 %+v  value %+v \n", name, value)
// 			proxyReq.Header.Add(name, value)
// 		}
// 	}
// 	Инфо(" proxyReq %+v \n", proxyReq)

// 	// Send the proxy request using the custom transport
// 	resp, err := customTransport.RoundTrip(proxyReq)
// 	if err != nil {
// 		http.Error(w, "Error sending proxy request", http.StatusInternalServerError)
// 		return
// 	}
// 	defer resp.Body.Close()
// 	Инфо(" resp %+v \n", resp)

// 	// Copy the headers from the proxy response to the original response
// 	for name, values := range resp.Header {
// 		Инфо(" name %+v  values %+v \n", name, values)
// 		for _, value := range values {
// 			Инфо(" name %+v  values %+v \n", name, value)
// 			w.Header().Add(name, value)
// 		}
// 	}

// 	// Set the status code of the original response to the status code of the proxy response
// 	w.WriteHeader(resp.StatusCode)

// 	// Copy the body of the proxy response to the original response
// 	io.Copy(w, resp.Body)
// }

// func обработчикЗапроса(w http.ResponseWriter, r *http.Request) {

// 	запрос := r.URL.Query().Get("запрос")
// 	Инфо("запрос %+v \n", запрос)

// 	proxyURL, err := url.Parse("htМне нужна помощь в написаноо прокси сервера, с функцией обратного прокси. Суть такая: Я отправляю запрос на свой адрес например myproxy.ru/query=https://google.com/?q=quest , Мой прокси должен получить ответ от сайта заданого в параметре query , и вернуть его мне в бразуер, при этом адресная строка в браузуру не должна изменятся.tps://google.ru/")

// 	if err != nil {
// 		log.Fatalf("Failed to parse proxy URL: %v", err)
// 	}
// 	Инфо(" proxyURL %+v \n", proxyURL)

// 	proxy := httputil.NewSingleHostReverseProxy(proxyURL)
// 	// Инфо(" proxy %+v \n", proxy)
// 	// Обработчик HTTP-запросов
// 	// http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
// 	// Обработка WebSocket-запросов
// 	if strings.ToLower(r.Header.Get("Upgrade")) == "websocket" {
// 		Инфо(" r.Header %+v \n", r.Header)

// 		handleWebSocket(w, r, proxyURL)
// 		return
// 	}

// 	// Обработка HTTP-запросов
// 	r.URL.Scheme = proxyURL.Scheme
// 	r.URL.Host = proxyURL.Host
// 	proxy.ServeHTTP(w, r)
// 	// })

// 	// fmt.Println("Proxy server is running on :8080")
// 	// log.Fatal(http.ListenAndServe(":8080", nil))
// }

// func handleWebSocket(w http.ResponseWriter, r *http.Request, proxyURL *url.URL) {
// 	// Установка соединения с целевым сервером
// 	targetURL := *proxyURL
// 	targetURL.Scheme = "ws"
// 	targetConn, _, err := websocket.DefaultDialer.Dial(targetURL.String(), r.Header)
// 	if err != nil {
// 		http.Error(w, "Failed to connect to target server", http.StatusBadGateway)
// 		return
// 	}
// 	defer targetConn.Close()

// 	// Установка соединения с клиентом
// 	upgrader := websocket.Upgrader{}
// 	clientConn, err := upgrader.Upgrade(w, r, nil)
// 	if err != nil {
// 		http.Error(w, "Failed to upgrade to WebSocket", http.StatusBadRequest)
// 		return
// 	}
// 	defer clientConn.Close()

// 	// Передача данных между клиентом и целевым сервером
// 	go func() {
// 		for {
// 			_, message, err := clientConn.ReadMessage()
// 			if err != nil {
// 				log.Printf("Error reading from client: %v", err)
// 				return
// 			}
// 			if err := targetConn.WriteMessage(websocket.TextMessage, message); err != nil {
// 				log.Printf("Error writing to target: %v", err)
// 				return
// 			}
// 		}
// 	}()

// 	for {
// 		messageType, message, err := targetConn.ReadMessage()
// 		if err != nil {
// 			log.Printf("Error reading from target: %v", err)
// 			return
// 		}
// 		if err := clientConn.WriteMessage(messageType, message); err != nil {
// 			log.Printf("Error writing to client: %v", err)
// 			return
// 		}
// 	}
// }
