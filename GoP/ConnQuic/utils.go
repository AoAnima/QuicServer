package ConnQuic

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/binary"
	"encoding/pem"
	"fmt"
	"math/big"
	"net/url"
	"os"
	"strconv"
	"time"

	. "aoanima.ru/logger"
	"github.com/dgryski/go-metro"
	"github.com/google/uuid"
	jsoniter "github.com/json-iterator/go"
)

func ДекодироватьПакет(пакет []byte) (Сообщение, error) {
	Инфо(" ДекодироватьПакет пакет %+s \n", пакет)

	// var запросОтКлиента = ЗапросКлиента{
	// 	Сервис:    []byte{},
	// 	Запрос:    &ЗапросОтКлиента{},
	// 	ИдКлиента: uuid.UUID{},
	// }
	var Сообщение Сообщение

	// TODO тут лишний парсинг, нужно получить только URL patch чтобы определить сервис, которому принадлежит запрос, потому nxj дальше весь запрос опять сериализуйется

	err := jsoniter.Unmarshal(пакет, &Сообщение)
	if err != nil {
		Ошибка("  %+v \n", err)
		return Сообщение, err
	}
	Инфо(" Сообщение входящее %+s \n", Сообщение)

	return Сообщение, err
}
func Кодировать(данныеДляКодирования interface{}) ([]byte, error) {

	b, err := jsoniter.Marshal(&данныеДляКодирования)
	if err != nil {
		Ошибка("  %+v \n", err)
		return nil, err
	}
	данные := make([]byte, len(b)+4)
	binary.LittleEndian.PutUint32(данные, uint32(len(b)))
	copy(данные[4:], b)
	return данные, nil

}

type ТипОтвета int

const (
	AjaxHTML ТипОтвета = iota
	AjaxJSON
	HTML
)

type ТипЗапроса int

const (
	GET ТипЗапроса = iota
	POST
	AJAX
	AJAXPost
)

type Отпечаток struct {
	Сервис   string
	Маршруты map[string]*СтруктураМаршрута
}
type СтруктураМаршрута struct {
	Запрос map[string]interface{} // описывает  данные которые нужны для обработки маршрута
	Ответ  map[string]interface{} // описывает формат в котором вернёт данные
}
type ИмяСервиса string
type Маршрут string // часть url.Path после домена, например example.ru/catalog/item?list=1 catalog это маршрут который сответсвует обработчику какого тое сервиса
type ОтветСервиса struct {
	Сервис          ИмяСервиса // Имя сервиса который отправляет ответ
	УИДЗапроса      string     // Копируется из запроса
	Данные          []byte     // Ответ в бинарном формате
	ЗапросОбработан bool       // Признак того что запросы был получен и обработан соответсвуюбщим сервисом, в не зависимоти есть ли данные в ответе или нет, если данных нет, знаичт они не нужны... Выставляем в true в сеорвисе перед отправкой ответа
}

type Ответ map[ИмяСервиса]ОтветСервиса

type Сообщение struct {
	Сервис      ИмяСервиса // Имя Сервиса который шлёт Сообщение, каждый сервис пишет своё имя в не зависимости что это ответ или запрос
	Регистрация bool
	// Маршруты     map[Маршрут]*СтруктураМаршрута
	Маршруты     []Маршрут
	Запрос       Запрос
	Ответ        Ответ
	ИдКлиента    uuid.UUID
	УИДСообщения Уид    // ХЗ по логике каждый сервис должен вставлять сюбда своё УИД
	ТокенКлиента []byte // JWT сериализованный
}
type Уид string

type Запрос struct {
	ТипОтвета      ТипОтвета
	ТипЗапроса     ТипЗапроса
	СтрокаЗапроса  *url.URL // url Path Query
	МаршрутЗапроса string   // url Path Query
	Форма          map[string][]string
	Файл           string
	УИДЗапроса     Уид
}

var (
	ПортДляОтправкиСообщений  = "81"
	ПортДляПолученияСообщений = "82"
)

func УИДЗапроса(ИдКлиента *uuid.UUID, UrlPath []byte) Уид {
	return Уид(fmt.Sprintf("%+s.%+s.%+s", strconv.FormatInt(time.Now().Unix(), 10), ИдКлиента, metro.Hash64(UrlPath, 0)))
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
