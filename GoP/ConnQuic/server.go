package ConnQuic

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/binary"
	"encoding/pem"
	"math/big"
	"os"
	"time"

	. "aoanima.ru/logger"
	quic "github.com/quic-go/quic-go"
)

// var Адрес = "localhost:4242"
// Запускаем сервер который слушает на адресе,
// принимает соединиеие, и отправляет его в обработчик Сессии
// обработчикСообщений - функция которая релаизует логику обработки входящего сообщения, сообщщение уэе прочитано, , логику обработки сообщений пишем в самом сервисе... И отправляем ответ в поток, который передаётся в фукнцию

func ЗапуститьСервер(Адрес string, обработчикСообщений func(поток quic.Stream, сообщение []byte)) {
	кофигТлс, err := серверныйТлсКонфиг()
	if err != nil {
		Ошибка("  %+v \n", err)
	}
	listener, err := quic.ListenAddr(Адрес, кофигТлс, nil)
	if err != nil {
		Ошибка(" %+v ", err)
	}

	for {
		сессия, err := listener.Accept(context.Background())
		if err != nil {
			Ошибка(" %+v ", err)
		}

		go ОбработчикСессии(сессия, обработчикСообщений)
	}
}

func ОбработчикСессии(сессия quic.Connection, обработчикСообщений func(поток quic.Stream, сообщение []byte)) {
	for {
		поток, err := сессия.AcceptStream(context.Background())
		if err != nil {
			Ошибка(" %+v ", err)
		}
		go ЧитатьСообщения(поток, обработчикСообщений)
	}

}

func ЧитатьСообщения(поток quic.Stream, обработчикСообщений func(поток quic.Stream, сообщение []byte)) {

	длинаСообщения := make([]byte, 4)
	var прочитаноБайт int
	var err error
	for {
		прочитаноБайт, err = поток.Read(длинаСообщения)
		Инфо(" длинаСообщения %+v , прочитаноБайт %+v \n", длинаСообщения, прочитаноБайт)

		if err != nil {
			Ошибка(" прочитаноБайт %+v  err %+v \n", прочитаноБайт, err)
			break
		}

		// получаем число байткоторое нужно прочитать
		длинаДанных := binary.LittleEndian.Uint32(длинаСообщения)

		Инфо(" длинаДанных  %+v \n", длинаДанных)
		Инфо(" длинаСообщения %+v ,  \n прочитаноБайт %+v ,  \n длинаДанных %+v \n", длинаСообщения,
			прочитаноБайт, длинаДанных)

		//читаем количество байт = длинаСообщения
		// var запросКлиента ЗапросКлиента
		сообщение := make([]byte, длинаДанных)
		прочитаноБайт, err = поток.Read(сообщение)
		if err != nil {
			Ошибка("Ошибка при десериализации структуры: %+v ", err)
		}
		if длинаДанных != uint32(прочитаноБайт) {
			Ошибка("Количество прочитаных байт не ранво длине данных :\n длинаДанных %+v  <> прочитаноБайт %+v ", длинаДанных, прочитаноБайт)
		}

		// каналПолученияСообщений <- пакетОтвета
		обработчикСообщений(поток, сообщение)

	}

}
func серверныйТлсКонфиг() (*tls.Config, error) {
	СоздатьКорневойСертификат()
	caCert, err := os.ReadFile("root.crt")
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

func генерироватьТлсКонфиг() *tls.Config {
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		panic(err)
	}
	template := x509.Certificate{SerialNumber: big.NewInt(1)}
	certDER, err := x509.CreateCertificate(rand.Reader, &template, &template, &key.PublicKey, key)
	if err != nil {
		panic(err)
	}
	keyPEM := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(key)})
	certPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: certDER})

	tlsCert, err := tls.X509KeyPair(certPEM, keyPEM)
	if err != nil {
		panic(err)
	}
	return &tls.Config{Certificates: []tls.Certificate{tlsCert}, InsecureSkipVerify: true}
}

func СоздатьКорневойСертификат() {
	_, err := os.Stat("root.crt")
	if !os.IsNotExist(err) { // если файл существует выходим
		return
	}
	// Генерация приватного ключа
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		panic(err)
	}

	// Создание шаблона для корневого сертификата
	rootTemplate := x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			Organization: []string{"Alсazar AO"},
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(10, 0, 0), // Срок действия 10 лет
		KeyUsage:              x509.KeyUsageCertSign | x509.KeyUsageDigitalSignature,
		BasicConstraintsValid: true,
		IsCA:                  true,
	}

	// Создание самоподписного корневого сертификата
	certBytes, err := x509.CreateCertificate(rand.Reader, &rootTemplate, &rootTemplate, &privateKey.PublicKey, privateKey)
	if err != nil {
		panic(err)
	}

	// Сохранение сертификата в файл
	certFile, err := os.Create("root.crt")
	if err != nil {
		panic(err)
	}
	pem.Encode(certFile, &pem.Block{Type: "CERTIFICATE", Bytes: certBytes})
	certFile.Close()

	// Сохранение приватного ключа в файл
	keyFile, err := os.Create("root.key")
	if err != nil {
		panic(err)
	}
	pem.Encode(keyFile, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(privateKey)})
	keyFile.Close()

}
