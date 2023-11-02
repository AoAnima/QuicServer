package ConnQuic

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"log"
	"os"

	. "aoanima.ru/logger"
	quic "github.com/quic-go/quic-go"
)

// Подключение к Сервре quic
/*

 */
// Адрес = "localhost:4242"
func ПокдлючитьсяКСерверу(Адрес string, обработчикСообщений func(поток *quic.Stream, сообщение Сообщение)) (quic.Connection, error) {
	конфигТлс, err := КлиентскийТлсКонфиг("root.crt")
	if err != nil {
		Ошибка("  %+v \n", err)
	}
	// Адрес = "localhost:4242"
	сессия, err := quic.DialAddr(context.Background(), Адрес, конфигТлс, &quic.Config{})
	if err != nil {
		log.Fatal(err)
	}

	for {
		поток, err := сессия.AcceptStream(context.Background())
		if err != nil {
			Ошибка("  %+v \n", err)
		}
		go ЧитатьСообщения(&сессия, поток, обработчикСообщений)

	}
}


func КлиентскийТлсКонфиг(caCertFile string) (*tls.Config, error) {
	caCert, err := os.ReadFile(caCertFile)
	if err != nil {
		Ошибка("  %+v \n", err)
	}

	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	return &tls.Config{
		RootCAs:    caCertPool,
		NextProtos: []string{"http/1.1", "h2", "h3", "quic", "websocket"},
	}, nil
}



