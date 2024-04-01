package main

import (
	"io"
	"net/http"
	"path"
	"sync"
	"time"

	_ "net/http/pprof"

	. "aoanima.ru/ConnQuic"
	. "aoanima.ru/Logger"
	quic "github.com/quic-go/quic-go"
)

var БлокКартаSynQuic = sync.RWMutex{}
var КартаSynQuic = make(map[ИмяСервера]HTTPКлиент)
var Сервис = "КлиентСервер"

var ВремяПинга = СтатистикаПинга{
	КартаПинга:        make(map[time.Time]time.Duration),
	ПоследнееЗначение: 5,
}

type СтатистикаПинга struct {
	КартаПинга        map[time.Time]time.Duration
	ПоследнееЗначение time.Duration
}

// var МакимальноеКоличествоПотоковНаСессию = 100
type HTTPКлиент struct {
	Блок           *sync.RWMutex
	Сессии         map[НомерСессии]*СхемаСервераHTTP
	НеПолныеСессии map[НомерСессии]int // Количество открытых потоков в сессии
}

type СхемаСервераHTTP struct {
	Имя            ИмяСервера
	Адрес          string
	Блок           *sync.RWMutex
	Соединение     quic.Connection
	СистемныйПоток quic.Stream
	ОчередьПотоков *ОчередьПотоков
	НомерСессии    НомерСессии
}
type Конфигурация struct {
	КаталогСтатичныхФайлов string
	КаталогШаблонов        string
}

// var каталогСтатичныхФайлов string
var Конфиг = &Конфигурация{}

func init() {
	Инфо(" проверяем какие аргументы переданы при запуске, если пусто то читаем конфиг, если конфига нет то устанавливаем значения по умолчанию %+v \n")

	// каталогСтатичныхФайлов = "../../HTML/static/"
	ЧитатьКонфиг(Конфиг)
}

func main() {

	/* каналЗапросовОтКлиентов - передаём этот канал в в функци  ЗапуститьСерверТЛС , когда прийдёт сообщение из браузера, функция обработчик запишет данные в этот канал
	 */
	// каналЗапросовОтКлиентов := make(chan http.Request, 10)
	/*
	   Запускаем сервер передаём в него канал, в который запишем обработанный запрос из браузера
	*/
	go ЗапуститьСерверТЛС() // принимаем запрос от клиента и отправляем в обработчик, ОтправитьЗапросВОбработку
	go СоединитсяСSynQuic()
	Инфо(" %s", "запустили сервер")
	/* Инициализирум сервисы коннектора передадим в них канал, из которого Коннектор будет читать сообщение, и отправлять его в synqTCP  */

	ЗапуститьWebСервер()

}

func ЗапуститьСерверТЛС() {

	err := http.ListenAndServeTLS(":443",
		"cert/server.crt",
		"cert/server.key",
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				обработчикЗапроса(w, r)
			}))

	if err != nil {
		Ошибка(" %s ", err)
	}

}

func обработчикЗапроса(w http.ResponseWriter, req *http.Request) {

	// Инфо(" %s \n", *req)
	// Инфо(" %s \n", req.URL.Path)
	if origin := req.Header.Get("Origin"); origin != "" {
		w.Header().Set("Access-Control-Allow-Origin", origin)
	}
	// Инфо(" path.Ext(req.URL.Path) %+v \n", path.Ext(req.URL.Path))
	if каталог, статичныйФайл := ТипыСтатическихФайлов[path.Ext(req.URL.Path)]; статичныйФайл {
		ОбработчикСтатичныхФайлов(w, req, каталог)
		return
	}

	// /*  Тут мы читаем из канала  каналОтвета кторый храниться в карте клиенты , данные пишутся в канал  в функции ОтправитьОтветКлиенту */
	// Отправляем сырой запрос в функцию ОтправитьЗапросВОбработку
	ответ, err := ОтправитьЗапросВОбработку(req)
	if err != nil {
		Ошибка("  %+v \n", err)
	}
	Инфо("  Полуичли ответ от всех сервисов, возвращаем ответ клиенту %+v \n")
	ОтправитьСообщениеКлиенту(ответ, w)
}

func ОтправитьСообщениеКлиенту(сообщение Сообщение, w http.ResponseWriter) {
	Инфо(" ОтправитьСообщениеКлиенту %+v  %+v \n", сообщение, w)
	ответ := КодироватьСообщениеОтвет(сообщение)
	УстановитьЗаголовкиОтвета(&сообщение, w)

	if f, ok := w.(http.Flusher); ok {
		Инфо("   %+v  %+v \n", f, ok)
		i, err := w.Write(ответ)
		Инфо("  %+v \n", i)
		if err != nil {
			Ошибка("%v %s ", w, err.Error())
		}
		f.Flush()
	}
}

func УстановитьЗаголовкиОтвета(сообщение *Сообщение, w http.ResponseWriter) {

	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	// Нужно проанализировать сообщение.ТипОтвета и установить Content-Type

	w.Header().Set("Access-Control-Allow-Credentials", "true")
	switch сообщение.Запрос.ТипОтвета {
	case HTML:
		w.Header().Set("Content-Type", "text/html;charset=utf-8")
	case Error:
		w.Header().Set("Content-Type", "text/html;charset=utf-8")
	case AjaxJSON:
		w.Header().Set("Content-Type", "application/json")
	case AjaxHTML:
		w.Header().Set("Content-Type", "application/json")

	}
	// UID идентификатор пользователя
	// UAT уникальный токен авторизации JWT
	w.Header().Set("X-UID", сообщение.ИдКлиента.String())
	w.Header().Set("X-UAT", сообщение.JWT)

}

func ЗапуститьWebСервер() {
	err := http.ListenAndServe(":80", http.HandlerFunc(
		func(w http.ResponseWriter, req *http.Request) {
			// 	Инфо(" %s  %s \n", w, req)
			http.Redirect(w, req, "https://localhost:443"+req.RequestURI, http.StatusMovedPermanently)
		},
	))

	// err := http.ListenAndServe(":6060", nil)
	if err != nil {
		Ошибка(" %+s ", err.Error())
	}
	go func() {
		http.ListenAndServe("localhost:6060", nil)
	}()
}

func ОбработчикСтатичныхФайлов(w http.ResponseWriter, req *http.Request, каталог ТипФайла) {

	файл := req.URL.Path

	if len(файл) != 0 {

		// if static_file == "js/tpl.js"{
		// 	//log.Printf("static_file %+v\n", static_file)
		// 	jsFile := renderJS("tplsJs", nil)
		// 	fileBytes := bytes.NewReader(jsFile)
		// 	content := io.ReadSeeker(fileBytes)

		// 	http.ServeContent(w, req, static_file, time.Now(), content)
		// 	return
		// }

		f, err := http.Dir(ДирректорияЗапуска + "/" + Конфиг.КаталогСтатичныхФайлов + каталог.Каталог).Open(req.URL.Path)

		if err == nil {
			content := io.ReadSeeker(f)
			http.ServeContent(w, req, файл, time.Now(), content)
			return
		} else {
			Ошибка(" %+v\n", err)
		}
	}
	http.NotFound(w, req)
}
