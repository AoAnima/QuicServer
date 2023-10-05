package main

import (
	"crypto/tls"
	"crypto/x509"
	"io"
	"net/http"
	"net/http/httptrace"

	"os"

	. "aoanima.ru/logger"
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
		MaxIdleConnsPerHost: 100,
		ForceAttemptHTTP2:   true,
		IdleConnTimeout:     0,
	}

	client := &http.Client{
		Transport: transport,
	}

	for i := 1; i <= 10; i++ {
		req, err := http.NewRequest("GET", "https://localhost:444", nil)
		if err != nil {
			Ошибка(" %s", err)
		}

		trace := &httptrace.ClientTrace{
			GetConn: func(hostPort string) {
				// Обрабатываем событие получения соединения
				Инфо("Получено соединение: %+v", hostPort)
			},
			GotConn: func(info httptrace.GotConnInfo) {
				// Обрабатываем событие получения установленного соединения
				Инфо ("Получено установленное соединение:  %+v", info.Conn.RemoteAddr())
			},
		}
		req = req.WithContext(httptrace.WithClientTrace(req.Context(), trace))


		resp, err := client.Do(req)
		if err != nil {
			Ошибка(" %+v", resp, err)
			return
		}

		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			Ошибка("Error:%+v\n", err)
		}
		Инфо("body :  %+v", string(body))
	}

	// ответ, err := client.Get("https://127.0.0.1:444?v=2")
	// if err != nil {
	// 	Ошибка("Error: %s\n", err)
	// }
	// Инфо(" %s", ответ)

}
