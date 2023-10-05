package main

import (
	"bufio"
	"bytes"
	"crypto/rand"
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
		MaxIdleConnsPerHost: 1,
		ForceAttemptHTTP2:   true,
		IdleConnTimeout:     0,
	}

	client := &http.Client{
		Transport: transport,
	}

	for i := 1; i <= 20; i++ {
		go func() {
			size := 1024

			// Создаем срез указанного размера
			slice := make([]byte, size)

			// Заполняем срез случайными данными
			_, err := rand.Read(slice)
			if err != nil {
				Ошибка(" %s", err)
			}
			// base64.StdEncoding.EncodeToString(slice)
			req, err := http.NewRequest("POST", "https://localhost:444", bytes.NewBuffer(slice))
			if err != nil {
				Ошибка(" %s", err)
			}
			req.Header.Set("Content-Type", "text/html; charset=utf-8")

			
		

			resp, err := client.Do(req)
			if err != nil {
				Ошибка(" %+v", resp, err)
				return
			}

			defer resp.Body.Close()

			reader := bufio.NewReader(resp.Body)
			for {
				line, err := reader.ReadString('\n')
				if err != nil {
					Ошибка("Error:%+v\n", err)
					if err == io.EOF {
						break
					}
					return
				}

				Инфо("lines  %+v", string(line))
			}
			// body, err := io.ReadAll(resp.Body)
			// if err != nil {
			// 	Ошибка("Error:%+v\n", err)
			// }
			// Инфо("body :  %+v", string(body))
		}()
	}

	// ответ, err := client.Get("https://127.0.0.1:444?v=2")
	// if err != nil {
	// 	Ошибка("Error: %s\n", err)
	// }
	// Инфо(" %s", ответ)

}
