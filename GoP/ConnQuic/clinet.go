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

type ПулПотоков struct {
}

func Клиент() (quic.Connection, error) {
	конфигТлс, err := клиентскийТлсКонфиг("root.crt")
	if err != nil {
		Ошибка("  %+v \n", err)
	}
	сессия, err := quic.DialAddr(context.Background(), "localhost:4242", конфигТлс, &quic.Config{})
	if err != nil {
		log.Fatal(err)
	}

	// реализцем механизм : создание пул потоков для сессии, ротацию потоков для отправки сообщений

}
func клиентскийТлсКонфиг(caCertFile string) (*tls.Config, error) {
	caCert, err := os.ReadFile(caCertFile)
	if err != nil {
		return nil, err
	}

	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	return &tls.Config{
		RootCAs:    caCertPool,
		NextProtos: []string{"http/1.1", "h2", "h3", "quic", "websocket"},
	}, nil
}
