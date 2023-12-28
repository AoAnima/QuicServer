package main

import (
	_ "net/http/pprof"
	"sync"

	. "aoanima.ru/ConnQuic"
	. "aoanima.ru/Logger"
	quic "github.com/quic-go/quic-go"
)

var клиент = make(Клиент)
var Сервис ИмяСервиса = "Рендер"

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
	ПарсингШаблонов()
	go наблюдатьЗаИзменениямиШаблонов()
	go наблюдатьЗаИзменениямиСтатичныхФайлов()
}

/*

 Для того чтобы быстрее верстать не обновляя руками страницу, сделаем websocket сервер с сервисом рендера, который будет отправлять сигнал в браузер на обновление при имзенении шаблонов или статичных файлов

*/

func main() {
	Инфо("  %+v \n", " Запуск Render сервиса ")

	// ПостроитьДеревоШаблонов(&Сообщение{})
	go WSСервер()
	сервер := &СхемаСервера{
		Имя:   "SynQuic",
		Адрес: "localhost:4242",
		ДанныеСессии: ДанныеСессии{
			Блок:   &sync.RWMutex{},
			Потоки: []quic.Stream{},
		},
	}
	сообщениеРегистрации := Сообщение{
		Сервис:      Сервис,
		Регистрация: true,
		// Маршруты:    []Маршрут{"/", "index"},
		Маршруты: []Маршрут{},
	}

	клиент.Соединиться(сервер,
		сообщениеРегистрации,
		ОбработчикОтветаРегистрации,
		ОбработчикЗапросовСервера)

}

func ОбработчикЗапросовСервера(поток quic.Stream, сообщение Сообщение) {
	Инфо("  ОбработчикЗапросовСервера входящее сообщение  \n")

	/*
		П Првоерим ТипЗапроса, если тп запроса обынчый GET или POST - не ajax, то рендерим index.html но в качестве шаблона content создаём новый Шаблон new("content") но из того шаблона который нужно вставить на место content
		Так можно делать с любым блоком
		Например вместо content нужно вставить catalog
		используем
		specificTemplate := tmpl.Lookup("catalog")
		Создаём новый шаблон и передаём в него Дерево
		// newTmpl := template.New("content")
		newTmpl = tmpl.AddParseTree("content", specificTemplate.Tree)
	*/
	ПарсингШаблонов()
	СтруктурироватьДанныеОтветов(&сообщение)
	// ПарсингШаблонов()
	ОтрендеритьОтветКлиенту(&сообщение) // записывает в ответКлиенту отрендеренные html
	ОтрендеритьОтветКлиенту(&сообщение) // записывает в ответКлиенту отрендеренные html
	// html := Рендер()

	// Инфо(" html %+v \n", html)
	// Ответ := сообщение.Ответ[Сервис]
	// Ответ.Данные = "ответ от рендер"
	// сообщение.Ответ[Сервис] = Ответ
	// Инфо(" отправляем ответ %+v \n", сообщение)
	отправить, err := Кодировать(сообщение)
	if err != nil {
		Ошибка("  %+v \n", err)
	}
	поток.Write(отправить)

}

func СтруктурироватьДанныеОтветов(сообщение *Сообщение) {
	// пройдёмся по струтктуре Шаблонизатор, положим данные в соответсвующие структуры для шаблонов,

	//	for _, картаДанныхШаблона := range сообщение.Запрос.Шаблонизатор {
	// 	for ИмяСервиса, _ := range картаДанныхШаблона.Данные {
	// 		ДанныеДляШбалона := сообщение.Ответ[ИмяСервиса].Данные
	// 		картаДанныхШаблона.Данные[ИмяСервиса] = ДанныеДляШбалона
	// 	}
	// }

}

func ОбработчикОтветаРегистрации(сообщение Сообщение) {
	Инфо("  ОбработчикОтветаРегистрации %+v \n", сообщение)
}