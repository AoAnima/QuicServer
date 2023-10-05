package main

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"net/http"
	"os"

	. "aoanima.ru/logger"
)

// https://github.com/jcbsmpsn/golang-https-example/blob/master/https_client.go
func СоденитьсяССервисомРендера(каналОтправкиДанных chan interface{}) {
	Инфо(" %s", "СоденитьсяССервисомРендера")

	caCert, err := os.ReadFile("cert/render_server.crt")

	if err != nil {
		Ошибка(" %s ", err)
	}

	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)
	cert, err := tls.LoadX509KeyPair("cert/client.crt", "cert/client.key")
	err =
	tlsConfig := &tls.Config{
		RootCAs:      caCertPool,
		Certificates: []tls.Certificate{cert},
	}

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig:     tlsConfig,
			MaxIdleConnsPerHost: 10,
		},
	}

	req, err := http.NewRequest("GET", "https://localhost:444", nil)
	if err != nil {
		Ошибка(" %s", err)
	}
	resp, err := client.Do(req)
	if err != nil {
		Ошибка(" %s", resp, err)
		return
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		Ошибка("Error: %s\n", err)
	}
	Инфо(" %s", string(body))

	req, err = http.NewRequest("GET", "https://127.0.0.1:444?v=1", nil)
	if err != nil {
		Ошибка(" %s", err)
	}
	resp, err = client.Do(req)
	if err != nil {
		Ошибка(" %s", err)
	}
	defer resp.Body.Close()

	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		Ошибка("Error: %s\n", err)
	}
	Инфо(" %s", string(body))

}
