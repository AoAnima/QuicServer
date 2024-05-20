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
	"mime/multipart"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	. "aoanima.ru/Logger"
	. "aoanima.ru/QErrors"
	"github.com/dgryski/go-metro"
	"github.com/google/uuid"
	jsoniter "github.com/json-iterator/go"
)

const (
	истина = 0 == 0 // Untyped bool.
	ложь   = 0 != 0 // Untyped bool.
	да     = 0 == 0 // Untyped bool.
	нет    = 0 != 0 // Untyped bool.
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
	".map": {
		ТипКонтента: "application/json",
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
	УдалитьОбъект   = "delete"
	ВставитьВКонец  = "append"
	ВставитьВНачало = "prepend"
	ВставитьПеред   = "before"
	ВставитьПосле   = "after"
	ЗаменитьОбъект  = "replaceWith"
	ОбновитьОбъект  = "replaceWith"
)

type Сообщение struct {
	Сервис      ИмяСервиса // Имя Сервиса который шлёт Сообщение, каждый сервис пишет своё имя в не зависимости что это ответ или запрос
	Регистрация bool
	Пинг        bool
	Понг        bool
	// Маршруты     map[Маршрут]*СтруктураМаршрута
	Маршруты      []Маршрут
	Запрос        Запрос
	Ответ         Ответ
	ОтветКлиенту  ОтветКлиенту
	ИдКлиента     uuid.UUID
	УИДСообщения  Уид          // ХЗ по логике каждый сервис должен вставлять сюбда своё УИД
	ТокенКлиента  ТокенКлиента // JWT сериализованный
	JWT           string
	ДанныеКлиента ДанныеКлиента
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
	ИмяБазовогоШаблона   ИмяШаблона
	ИмяВложенногоШаблона ИмяШаблона
	Данные               map[ИмяСервиса]interface{}
}

//	type СоставнаяФормаForm struct {
//		Value map[string][]string
//		File  map[string][]*multipart.FileHeader
//	}
/*
Пример обработчки дейсвтий для сервисов
if сообщение.Запрос.Действие != "" {
		Действие = сообщение.Запрос.Действие
	} else {

	маршрутЗапроса, err := url.Parse(сообщение.Запрос.МаршрутЗапроса)
		Инфо(" маршрутЗапроса %+v \n", маршрутЗапроса)
		if err != nil {
			Ошибка("Ошибка при парсинге СтрокаЗапроса запроса:", err)
		}

		маршрутЗапроса.Path = strings.Trim(маршрутЗапроса.Path, "/")
		дейсвтия := strings.Split(маршрутЗапроса.Path, "/")

		if len(дейсвтия) == 0 {
			Инфо(" Пустой маршрут, добавляем в маршруты обработку по умолчанию: авторизация \n")
			// Читаем заголовки парсим и проверяем JWT
			Действие = "авторизация" //првоерим и валидируем токен, получим права доступа

		} else {
			Действие = дейсвтия[0]
		}
	}
*/
type Запрос struct {
	ТипОтвета  ТипОтвета
	ТипЗапроса ТипЗапроса
	// ИмяБазовогоШаблона ИмяШаблона //
	СпособВставки  string
	СтрокаЗапроса  *url.URL // url Path Query
	МаршрутЗапроса string   // url Path Query
	Действие       string   // Поле по которому сервис определяет что необходимо сделать, если передана форма, то необходимо лдостать поле дейсвтвие и значение поместить в данное поле, если форма не передана то функция ПостроитьМаршрут может заполнить это поле исходя из данных БД, по идее можно перед отправкой в каждый сервис менять это значение, лбо во всех сервисах которые должны обработать текущий запрос, должен быть одинаковый обработтчик ?
	КартаМаршрута  []string // url Path Query
	Форма          map[string][]string
	СоставнаяФорма *multipart.Form
	Файл           string
	УИДЗапроса     Уид
}

type Ответ map[ИмяСервиса]ОтветСервиса

type ОтветСервиса struct {
	Сервис     ИмяСервиса // Имя сервиса который отправляет ответ
	УИДЗапроса string     // Копируется из запроса
	// HTML            string                    // Ответ в бинарном формате
	// AjaxHTML        map[string]ДанныеAjaxHTML // Типа map[id селектор для вставки в HTML]<html для вставки>
	Данные interface{} // данные в виде структуры какойто
	// Шаблонизатор   map[ИмяШаблона]КартаДанныхШаблона
	ИмяШаблона      ИмяШаблона // Заполняется при получении очереди обработчиков ?! Можно задавать имя шаблона в виде пути директорий имяОсновногоКонтентБлока/ВложеныйШаблон - где имяОсновногоКонтентБлока - это имя блока который будет вставлен в блок Контент() , а ВложеныйШаблон это блок окторый будет вствален в имяБлока_контент() - наверное так
	ТипОтвета       ТипОтвета
	ЗапросОбработан bool // Признак того что запросы был получен и обработан соответсвуюбщим сервисом, в не зависимоти есть ли данные в ответе или нет, если данных нет, знаичт они не нужны... Выставляем в true в сеорвисе перед отправкой ответа
	СтатусОтвета    СтатусСервиса
}

// ТокенКлинета, запоняется веб сервером который принимает запросы от клиента
type ТокенКлиента struct {
	ИдКлиента uuid.UUID `json:"UID"`
	Роль      []string  `json:"role"`   // не передаём клиенту записываем только в момент Авторизации и обработки запроса на сервере
	Токен     string    `json:"token"`  // пока не используется навреное нафиг не нужно
	Права     []string  `json:"access"` // не передаём клиенту записываем только в момент Авторизации и обработки запроса на сервере
	Истекает  time.Time `json:"exp"`
	Создан    time.Time `json:"iat"`
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

type ДанныеКлиента struct {
	Аутентифицирован            bool           `json:"аутентифицирован,omitempty"`
	КоличествоОшибокАвторизации bool           `json:"количество_неудачных_попыток_входа,omitempty"`
	Имя                         string         `json:"имя,omitempty"`
	Фамилия                     string         `json:"фамилия,omitempty"`
	Отчество                    string         `json:"отчество,omitempty"`
	ИдКлиента                   uuid.UUID      `json:"ид_клиента"`
	ПраваДоступа                []ПраваДоступа `json:"права_доступа,omitempty"` // При получении прав пользователя их нужно доавбялть по мере уменьшения кода, код с наименьшим числом имеет наивысший приоритет. 1. Админ.... 10.Гость . О.Создатель имеет абсолютные права
	// Роль                        []Роль    `json:"роль,omitempty"`
	// Права                       []Права   `json:"права,omitempty"`
	Статус   string    `json:"статус,omitempty"`
	Аватар   string    `json:"аватар,omitempty"`
	Email    string    `json:"email,omitempty"`
	Логин    string    `json:"логин,omitempty"`
	Пароль   string    `json:"пароль,omitempty"`
	JWT      string    `json:"jwt,omitempty"`
	Телефон  string    `json:"телефон,omitempty"`
	Адрес    Адрес     `json:"адрес,omitempty"`
	Создан   time.Time `json:"создан,omitempty"`
	Обновлен time.Time `json:"обновлен,omitempty"`
	ОСебе    string    `json:"о_себе,omitempty"`
	СоцСети  []string  `json:"социальные_ссылки,omitempty"`
	Профиль  map[string]interface{}
	// Секрет                      Секрет `json:"СекретКлиента,omitempty"`
}

type Адрес struct {
	Страна        string `json:"страна,omitempty"`
	Город         string `json:"город,omitempty"`
	Район         string `json:"район,omitempty"`
	ТипУлицы      string `json:"тип_улицы,omitempty"`
	НазваниеУлицы string `json:"название_улицы,omitempty"`
	НомерДома     string `json:"номер_дома,omitempty"`
	Корпус        string `json:"корпус,omitempty"`
	НомерКвартиры string `json:"номер_квартиры,omitempty"`
}
type Секрет struct {
	ИдКлиента string    `json:"ид_клиента,omitempty"`
	Секрет    string    `json:"секрет,omitempty"`
	Обновлен  time.Time `json:"обновлен,omitempty"`
}

type КонфигурацияОбработчика struct {
	UID          string         `json:"uid,omitempty"`
	Маршрут      string         `json:"маршрут,omitempty"`
	Действие     string         `json:"действие,omitempty"`
	Обработчик   string         `json:"обработчик,omitempty"`
	ПраваДоступа []ПраваДоступа `json:"доступ"`
	Описание     string         `json:"описание,omitempty"`
	Шаблонизатор []Шаблон       `json:"шаблонизатор,omitempty"`
	Ассинхронно  bool           `json:"ассинхронно,omitempty"`
	Тип          string         `json:"dgraph.type,omitempty"`
}
type ОбработчикМаршрута struct {
	UID                             string         `json:"uid,omitempty"`
	Маршрут                         string         `json:"маршрут,omitempty"`
	Комманда                        string         `json:"комманда,omitempty"`
	ОчередьОбработчиков             []Обработчик   `json:"очередь_обработчиков,omitempty"`
	АссинхроннаяОчередьОбработчиков []Обработчик   `json:"ассинхронная_очередь_обоработчиков,omitempty"`
	ПраваДоступа                    []ПраваДоступа `json:"доступ"`
	// Роль                            Роль           `json:"роль,omitempty"`
	Описание   string `json:"описание,omitempty"`
	Шаблон     Шаблон `json:"шаблон,omitempty"`
	ИмяШаблона string `json:"имя_шаблона,omitempty"`
	Тип        string `json:"dgraph.type,omitempty"`
}

type Роль struct {
	Тип     string `json:"dgraph.type,omitempty"`
	Код     int    `json:"код.роли,omitempty"`
	ИмяРоли string `json:"имя.роли,omitempty"`
}
type Права struct {
	Тип     string `json:"dgraph.type,omitempty"`
	Код     int    `json:"код.прав,omitempty"`
	ИмяПрав string `json:"имя.прав,omitempty"`
}
type Обработчик struct {
	Тип            string `json:"dgraph.type,omitempty"`
	UID            string `json:"uid,omitempty"`
	Очередь        *int   `json:"очередь,omitempty"`
	ИмяСервиса     string `json:"сервис,omitempty"`
	ИмяОбработчика string `json:"имя_обработчика,omitempty"`
}

type Шаблон struct {
	UID        string `json:"uid,omitempty"`
	Тип        string `json:"dgraph.type,omitempty"`
	Код        *int   `json:"код,omitempty"`         // статус ответа сервиса QErrors
	ИмяШаблона string `json:"имя_шаблона,omitempty"` // полный путь к шаблону после каталога с именем роли , тоесть естли в этом поле написано рабочийСтол/настройки , то фактический путь будет контент/имяРоли/рабочийСтол/настройки
	// ПраваДоступа []ПраваДоступа `json:"доступ,omitempty"`
}
type ПраваДоступа struct {
	UID   string   `json:"uid,omitempty"`
	Тип   string   `json:"dgraph.type,omitempty"`
	Логин []string `json:"пользователи,omitempty"`
	Роль  Роль     `json:"роль"`
	Права []Права  `json:"права"`
}

func ОпределитьДирректориюЗапуска() {
	if ДирректорияЗапуска != "" {
		Инфо("Директория из которой запущен текущий файл уже определена: %+v \n", ДирректорияЗапуска)
		return
	}
	exePath, err := os.Executable()
	Инфо("  %+v  %+v \n", exePath, err)
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
	Инфо(" ДирректорияЗапуска %+v \n", ДирректорияЗапуска)
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
	Инфо("len(b) %+v \n", len(b))
	// конецСообщения := []byte("\\<")
	// b = append(b, конецСообщения...)
	данные := make([]byte, len(b)+4)

	binary.LittleEndian.PutUint32(данные, uint32(len(b)))
	// Инфо(" %+v %+v данные %+v \n", uint32(len(b)), len(b), данные)

	copy(данные[4:], b)
	// log.Print(данные, string(данные))
	// Инфо(" Кодировать %+s %+s %+v \n", данные, string(данные), len(b))
	return данные, nil
}

func Json(данныеДляКодирования interface{}) ([]byte, error) {

	данные, err := jsoniter.Marshal(&данныеДляКодирования)
	if err != nil {
		Ошибка("  %+v  %+v \n", данные, err.Error())
		return nil, err
	}

	return данные, nil
}

func ИзJson(пакет []byte, объект interface{}) error {
	Инфо(" ДекодироватьПакет ИзJson %#T \n", объект)
	switch v := объект.(type) {
	case *string:
		*v = string(пакет)
		объект = v
		return nil
	default:
		err := jsoniter.Unmarshal([]byte(пакет), &объект)
		if err != nil {
			Ошибка(" err  %+s; \n\n Тип пакета > %#T < \n\n В какой объект помещаем: > %+s < ; \n", err.Error(), пакет, объект)
			return err
		}
		// объект = v
		// Инфо("объект %+v \n", объект)

	}
	// switch (*объект).(type) {
	// 	case string:
	// 		*объект = string(пакет)
	// 		return nil
	// 	default:
	// 		Инфо("  %+v \n", *объект)
	// 	}

	// var запросОтКлиента = ЗапросКлиента{
	// 	Сервис:    []byte{},
	// 	Запрос:    &ЗапросОтКлиента{},
	// 	ИдКлиента: uuid.UUID{},
	// }

	// TODO тут лишний парсинг, нужно получить только URL patch чтобы определить сервис, которому принадлежит запрос, потому nxj дальше весь запрос опять сериализуйется

	return nil
}

func НайтиВJson(объект []byte, путь string) interface{} {

	срезПути := strings.Split(путь, ".")
	срезИнтерфейсов := make([]interface{}, len(срезПути))

	for i, v := range срезПути {
		срезИнтерфейсов[i] = v
	}
	значение := jsoniter.Get(объект, срезИнтерфейсов...)

	return значение.GetInterface()
}
