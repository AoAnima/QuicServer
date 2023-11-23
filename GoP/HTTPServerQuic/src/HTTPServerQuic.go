package main

import (
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"sync"
	"time"

	_ "net/http/pprof"

	. "aoanima.ru/ConnQuic"
	. "aoanima.ru/Logger"
	jsoniter "github.com/json-iterator/go"
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

// var каталогСтатичныхФайлов string

func init() {
	Инфо(" проверяем какие аргументы переданы при запуске, если пусто то читаем конфиг, если конфига нет то устанавливаем значения по умолчанию %+v \n")

	// каталогСтатичныхФайлов = "../../HTML/static/"
	ЧитатьКонфиг()
}

type конфигурация struct {
	КаталогСтатичныхФайлов string
}

var Конфиг = &конфигурация{}

func ЧитатьКонфиг() {
	конфиг, err := os.ReadFile("config.json")
	if err != nil {
		Ошибка("  %+v \n", err)
	}

	err = jsoniter.Unmarshal(конфиг, Конфиг)
	if err != nil {
		Ошибка("  %+v \n", err)
	}
}

func main() {

	/* каналЗапросовОтКлиентов - передаём этот канал в в функци  ЗапуститьСерверТЛС , когда прийдёт сообщение из браузера, функция обработчик запишет данные в этот канал
	 */
	// каналЗапросовОтКлиентов := make(chan http.Request, 10)
	/*
	   Запускаем сервер передаём в него канал, в который запишем обработанный запрос из браузера
	*/
	go ЗапуститьСерверТЛС()
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

func обработчикЗапроса(w http.ResponseWriter, req *http.Request) {

	Инфо(" %s \n", *req)
	Инфо(" %s \n", req.URL.Path)
	if origin := req.Header.Get("Origin"); origin != "" {
		w.Header().Set("Access-Control-Allow-Origin", origin)
	}
	Инфо(" path.Ext(req.URL.Path) %+v \n", path.Ext(req.URL.Path))
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
	ОтправитьСообщениеКлиенту(ответ, w)
}

func ОтправитьСообщениеКлиенту(сообщение Сообщение, w http.ResponseWriter) {
	ответ := КодироватьСообщениеОтвет(сообщение)

	if f, ok := w.(http.Flusher); ok {

		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

		// Нужно проанализировать сообщение.ТипОтвета и установить Content-Type
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		Инфо("  %+v \n", string(ответ))
		i, err := w.Write(ответ)
		Инфо("  %+v \n", i)
		if err != nil {
			Ошибка(" %s ", err)
		}
		f.Flush()
	}
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
		Ошибка(" %s ", err)
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

		f, err := http.Dir(Конфиг.КаталогСтатичныхФайлов + каталог.Каталог).Open(req.URL.Path)

		if err == nil {
			content := io.ReadSeeker(f)
			http.ServeContent(w, req, файл, time.Now(), content)
			return
		} else {
			log.Printf("%+v\n", err)
		}
	}
	http.NotFound(w, req)
}