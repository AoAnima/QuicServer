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
	"path/filepath"
	"strconv"
	"time"

	. "aoanima.ru/Logger"
	"github.com/dgryski/go-metro"
	"github.com/google/uuid"
	jsoniter "github.com/json-iterator/go"
)

type ТипФайла struct {
	ТипКонтента string
	Каталог     string
}

var ТипыСтатическихФайлов = map[string]ТипФайла{
	".css": {
		ТипКонтента: "text/css",
		Каталог:     "./css/",
	},
	".js": {
		ТипКонтента: "text/javascript",
		Каталог:     "./js/",
	},
	".jpeg": {
		ТипКонтента: "image/jpeg",
		Каталог:     "./images/",
	},
	".jpg": {
		ТипКонтента: "image/jpeg",
		Каталог:     "./images/",
	},
	".png": {
		ТипКонтента: "image/png",
		Каталог:     "./images/",
	},
	".svg": {
		ТипКонтента: "image/svg+xml",
		Каталог:     "./images/",
	},
	".gif": {
		ТипКонтента: "image/gif",
		Каталог:     "./images/",
	},
	".ico": {
		ТипКонтента: "image/x-icon",
		Каталог:     "./images/",
	},
	// ".zip":    {
	// 	ТипКонтента:"application/zip",
	// 	Каталог:     "./images/",
	// },
	// ".pdf":   "application/pdf",
	// ".doc":   "application/msword",
	// ".xls":   "application/vnd.ms-excel",
	// ".ppt":   "application/vnd.ms-powerpoint",
	// ".mp3":   "audio/mpeg",
	// ".mp4":   "video/mp4",
	// ".wav":   "audio/wav",
	// ".ogg":   "audio/ogg",
	// ".webm":  "video/webm",
	".ttf": {
		ТипКонтента: "font/ttf",
		Каталог:     "./fonts/",
	},
	".woff": {
		ТипКонтента: "font/woff",
		Каталог:     "./fonts/",
	},
	".woff2": {
		ТипКонтента: "font/woff2",
		Каталог:     "./fonts/",
	},
	".eot": {
		ТипКонтента: "font/eot",
		Каталог:     "./fonts/",
	},
	".otf": {
		ТипКонтента: "font/otf",
		Каталог:     "./fonts/",
	},
	".ttc": {
		ТипКонтента: "font/ttc",
		Каталог:     "./fonts/",
	},
}

type ТипОтвета int

const (
	AjaxHTML ТипОтвета = iota
	AjaxJSON
	HTML
	Error
)

type ТипЗапроса int

const (
	GET ТипЗапроса = iota
	POST
	AJAX
	AJAXPost
)

//	type Отпечаток struct {
//		Сервис   string
//		Маршруты map[string]*СтруктураМаршрута
//	}
//
//	type СтруктураМаршрута struct {
//		Запрос map[string]interface{} // описывает  данные которые нужны для обработки маршрута
//		Ответ  map[string]interface{} // описывает формат в котором вернёт данные
//	}
type ИмяСервиса string
type Маршрут string // часть url.Path после домена, например example.ru/catalog/item?list=1 catalog это маршрут который сответсвует обработчику какого тое сервиса

type ДанныеAjaxHTML struct {
	Цель          string
	HTML          string
	СпособВставки string
}

const (
	Удалить         = "delete"
	ВставитьВКонец  = "append"
	ВставитьВНачало = "prepend"
	ВставитьПеред   = "before"
	ВставитьПосле   = "after"
	Заменить        = "replaceWith"
	Обновить        = "replaceWith"
)

type Сообщение struct {
	Сервис      ИмяСервиса // Имя Сервиса который шлёт Сообщение, каждый сервис пишет своё имя в не зависимости что это ответ или запрос
	Регистрация bool
	Пинг        bool
	Понг        bool
	// Маршруты     map[Маршрут]*СтруктураМаршрута
	Маршруты     []Маршрут
	Запрос       Запрос
	Ответ        Ответ
	ОтветКлиенту ОтветКлиенту
	ИдКлиента    uuid.UUID
	УИДСообщения Уид          // ХЗ по логике каждый сервис должен вставлять сюбда своё УИД
	ТокенКлиента ТокенКлинета // JWT сериализованный
	JWT          string
}

type ОтветКлиенту struct {
	ТипОтвета ТипОтвета
	HTML      []byte                    // Ответ в бинарном формате
	AjaxHTML  map[string]ДанныеAjaxHTML // Типа map[id селектор для вставки в HTML]<html для вставки>
	JSON      interface{}
}

// КартаДанныхШаблона, скорей всего буду в БД хранить какому шаблоны, из каокго сервиса нужно получить данные, и пперед отправкой в сервис рендер, или перед рендером буду их складывать в эту стрктуру
type ИмяШаблона string

type КартаДанныхШаблона struct {
	ИмяБазовогоШаблона ИмяШаблона
	Данные             map[ИмяСервиса]interface{}
}

type Запрос struct {
	ТипОтвета          ТипОтвета
	ТипЗапроса         ТипЗапроса
	ИмяБазовогоШаблона ИмяШаблона //
	СпособВставки      string
	СтрокаЗапроса      *url.URL // url Path Query
	МаршрутЗапроса     string   // url Path Query
	КартаМаршрута      []string // url Path Query
	Форма              map[string][]string
	Файл               string
	УИДЗапроса         Уид
	Шаблонизатор       map[ИмяШаблона]КартаДанныхШаблона
}

type Ответ map[ИмяСервиса]ОтветСервиса

type ОтветСервиса struct {
	Сервис          ИмяСервиса                // Имя сервиса который отправляет ответ
	УИДЗапроса      string                    // Копируется из запроса
	HTML            string                    // Ответ в бинарном формате
	AjaxHTML        map[string]ДанныеAjaxHTML // Типа map[id селектор для вставки в HTML]<html для вставки>
	Данные          interface{}               // данные в виде структуры какойто
	ТипОтвета       ТипОтвета
	ЗапросОбработан bool // Признак того что запросы был получен и обработан соответсвуюбщим сервисом, в не зависимоти есть ли данные в ответе или нет, если данных нет, знаичт они не нужны... Выставляем в true в сеорвисе перед отправкой ответа
}

type ТокенКлинета struct {
	ИдКлиента uuid.UUID `json:"UID"`
	Роль      []string  `json:"role"`
	Токен     string    `json:"token"`
	Права     []string  `json:"access"`
	Истекает  int64     `json:"expires"`
	Создан    int64     `json:"created"`
}

type Уид string
type СистемноеСообщение struct {
	Сервис       ИмяСервиса
	Запрос       interface{}
	Ответ        interface{}
	УИДСообщения Уид
}

var (
	ПортДляОтправкиСообщений  = "81"
	ПортДляПолученияСообщений = "82"
)
var ДирректорияЗапуска string

func ОпределитьДирректориюЗапуска() {
	if ДирректорияЗапуска != "" {
		Инфо("Директория из которой запущен текущий файл уже определена: %+v \n", ДирректорияЗапуска)
		return
	}
	exePath, err := os.Executable()
	if err != nil {
		Ошибка("Ошибка при получении пути к исполняемому файлу:", err)

	}
	ДирректорияЗапуска = filepath.Dir(exePath)
	Дир := ДирректорияЗапуска[len(ДирректорияЗапуска)-3:]
	if Дир == "src" {
		ДирректорияЗапуска = ДирректорияЗапуска[:len(ДирректорияЗапуска)-3] + "bin"
	}
	Инфо("Директория, из которой запущен текущий файл: %+v \n", ДирректорияЗапуска)

}
func УИДЗапроса(ИдКлиента *uuid.UUID, UrlPath []byte) Уид {
	времяГенерации := time.Now().Unix()
	return Уид(fmt.Sprintf("%+s.%+s.%d", strconv.FormatInt(времяГенерации, 10), ИдКлиента, metro.Hash64(UrlPath, uint64(времяГенерации))))
}

func УИДСистемногоЗапроса(Сервис string) Уид {
	времяГенерации := time.Now().Unix()
	return Уид(fmt.Sprintf("%+s.%+s.%d", strconv.FormatInt(времяГенерации, 10), Сервис, metro.Hash64([]byte(Сервис), uint64(времяГенерации))))
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

	_, err := os.Stat(ДирректорияЗапуска + "/cert/root.crt")
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
	certFile, err := os.Create(ДирректорияЗапуска + "/cert/root.crt")
	if err != nil {
		panic(err)
	}
	pem.Encode(certFile, &pem.Block{Type: "CERTIFICATE", Bytes: certBytes})
	certFile.Close()

	// Сохранение приватного ключа в файл
	keyFile, err := os.Create(ДирректорияЗапуска + "/cert/root.key")
	if err != nil {
		panic(err)
	}
	pem.Encode(keyFile, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(privateKey)})
	keyFile.Close()

}

func СерверныйТлсКонфиг() (*tls.Config, error) {

	// СоздатьКорневойСертификат()
	// dir, err := os.Getwd()
	// if err != nil {
	// 	Ошибка("Ошибка при получении текущей директории:", err)

	// }

	// Инфо(" dir %+v \n", dir)

	caCert, err := os.ReadFile(ДирректорияЗапуска + "/cert/ca.crt")
	if err != nil {
		return nil, err
	}

	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	cert, err := tls.LoadX509KeyPair(ДирректорияЗапуска+"/cert/server.crt", ДирректорияЗапуска+"/cert/server.key")
	if err != nil {
		Ошибка("  %+v \n", err)
	}

	return &tls.Config{
		// InsecureSkipVerify: true,
		Certificates: []tls.Certificate{cert},
		RootCAs:      caCertPool,
		NextProtos:   []string{"h3", "quic", "websocket"},
		// NextProtos: []string{"http/1.1", "h2", "h3", "quic", "websocket"},
	}, nil
}

// type Конфигурация interface{}

func ЧитатьКонфиг(Конфиг interface{}) {
	ОпределитьДирректориюЗапуска()
	конфиг, err := os.ReadFile(ДирректорияЗапуска + "/config.json")
	if err != nil {
		Ошибка("  %+v \n", err)
	}

	err = jsoniter.Unmarshal(конфиг, Конфиг)
	if err != nil {
		Ошибка("  %+v \n", err)
	}
}

func ДекодироватьПакет(пакет []byte) (Сообщение, error) {
	// Инфо(" ДекодироватьПакет пакет %+s \n", пакет)

	// var запросОтКлиента = ЗапросКлиента{
	// 	Сервис:    []byte{},
	// 	Запрос:    &ЗапросОтКлиента{},
	// 	ИдКлиента: uuid.UUID{},
	// }
	var сообщение Сообщение

	// TODO тут лишний парсинг, нужно получить только URL patch чтобы определить сервис, которому принадлежит запрос, потому nxj дальше весь запрос опять сериализуйется

	err := jsoniter.Unmarshal(пакет, &сообщение)
	if err != nil {
		Ошибка(" err  %+v пакет >%+s< ; \n", err.Error(), пакет)
		return сообщение, err
	}
	// Инфо(" Сообщение входящее %+s \n", Сообщение)

	return сообщение, err
}
func Кодировать(данныеДляКодирования interface{}) ([]byte, error) {

	b, err := jsoniter.Marshal(&данныеДляКодирования)
	if err != nil {
		Ошибка("  %+v \n", err)
		return nil, err
	}
	// конецСообщения := []byte("\\<")
	// b = append(b, конецСообщения...)
	данные := make([]byte, len(b)+4)

	binary.LittleEndian.PutUint32(данные, uint32(len(b)))

	copy(данные[4:], b)
	// log.Print(данные, string(данные))
	// Инфо(" Кодировать %+s %+s \n", данные, string(данные))
	return данные, nil

}
