package main

import (
	"bytes"
	"html/template"
	"log"
	"os"
	"path/filepath"

	. "aoanima.ru/ConnQuic"
	. "aoanima.ru/Logger"
)

var СырыеШаблоны *template.Template

func ПостроитьДеревоШаблонов(сообщение *Сообщение) {

	// content := template.New("content")
	// catalog := СырыеШаблоны.Lookup("catalog")
	// Инфо(" СырыеШаблоны %+v \n", СырыеШаблоны)
	// ШаблонДляРендера, err := СырыеШаблоны.Clone()
	// if err != nil {
	// 	Ошибка(" Ошибка клоинровании сырых шаблонов %+v \n", err)
	// }

	// catalog := ШаблонДляРендера.Lookup("catalog")
	// Инфо("  %+v \n catalog  %+v \n", ШаблонДляРендера, catalog)

	// ШаблонДляРендера, err = ШаблонДляРендера.AddParseTree("content", catalog.Tree)
	// if err != nil {
	// 	Ошибка("  %+v \n", err)
	// }
	if сообщение.Запрос.ТипЗапроса == GET || сообщение.Запрос.ТипЗапроса == POST {
		// ШаблонДляРендера, err := СырыеШаблоны.Clone()
		// if err != nil {
		// 	Ошибка("  %+v \n", err)
		// }
		Рендер("index", сообщение.Запрос.Шаблонизатор)
		// найдём шаблон который нужно вставить в блок content
		// ШаблонКонтента := ШаблонДляРендера.Lookup(string(сообщение.Запрос.ИмяШаблона))
		// ШаблонДляРендера, err = ШаблонДляРендера.AddParseTree("content", ШаблонКонтента.Tree)
		// Рендер("index", сообщение.Ответ[ИмяСервиса(сообщение.Запрос.ИмяШаблона)])
	}
	var err error
	var Html []byte
	switch сообщение.Запрос.ТипЗапроса {
	case GET:
		Html, err = Рендер("index", сообщение.Запрос.Шаблонизатор)
	case POST:
		Html, err = Рендер("index", сообщение.Запрос.Шаблонизатор)
	case AJAX:
		Html, err = РендерБлоков(сообщение)
	case AJAXPost:
		Html, err = РендерБлоков(сообщение)
	}

	// Html := new(bytes.Buffer)
	// Кнопка := map[string]map[string]string{
	// 	"Кнопка": {
	// 		"Класс": "success",
	// 		"Тип":   "submit",
	// 		"Текст": "Кнопка волшебная",
	// 	},
	// }

	// Данные := map[string]interface{}{
	// 	"content": Кнопка,
	// }

	if err != nil {
		Ошибка("   %+v \n", err)
	}
	// if errs := ШаблонДляРендера.ExecuteTemplate(Html, "index", Данные); errs != nil {
	// 	Ошибка("%+v\n", errs)

	// }
	Инфо("  %+s \n", Html)

}

func РендерБлоков(сообщение *Сообщение) {

	ответКлиенту := ОтветКлиенту{
		AjaxHTML: make(map[string]ДанныеAjaxHTML),
	}

	for ИмяШаблона, _ := range сообщение.Запрос.Шаблонизатор {
		Html, err := Рендер(string(ИмяШаблона), сообщение.Запрос.Шаблонизатор)
		if err != nil {
			Ошибка("  %+v \n", err)
		}
		ответКлиенту.AjaxHTML[string(ИмяШаблона)] = ДанныеAjaxHTML{
			Цель:          string(ИмяШаблона),
			HTML:          string(Html),
			СпособВставки: Заменить, // способ вставки - нужно придумать где хранить и как определять, либо храним в БД , либо в ajax запросе, например запрос путь в адресной строке catalog/page=2
			// а ajax запрос будет в заивисомсти от События вызвавшее запрос, добавлять каокй нибудь метод в ajax запрос "updateMethod": "replaceWith" ...

			// Хрень а если я буду возвращать несколько бооков... значит способ вставки должен храниться в базе, рядом с данными о шаблонах и сервисах из котрых получаем данные для этих шаблонов
		}
	}
}

func Рендер(имяШаблона string, КартаДанных map[ИмяШаблона]КартаДанныхШаблона) ([]byte, error) {
	Html := new(bytes.Buffer)

	ШаблонДляРендера, err := СырыеШаблоны.Clone()
	if err != nil {
		Ошибка(" Ошибка клоинровании сырых шаблонов %+v \n", err)
		return nil, err
	}
	Инфо(" ШаблонДляРендера  %+v \n", ШаблонДляРендера)
	if errs := ШаблонДляРендера.ExecuteTemplate(Html, имяШаблона, КартаДанных[ИмяШаблона(имяШаблона)].Данные); errs != nil {
		Ошибка("%+v\n", errs)
		return nil, errs
	}

	return Html.Bytes(), nil
}

func ПарсингШаблонов() {
	// "pattern": "../www/tpl/*/*.html",
	var errParseGlob error
	Инфо(" Конфиг.КаталогШаблонов  %+v \n", Конфиг.КаталогШаблонов+"*/*.html")

	// filenames, err := filepath.Glob(Конфиг.КаталогШаблонов + "*/*.html")
	// if err != nil {
	// 	Ошибка(" Ошибка парсинга каталога с шаблонами HTML %+v\n", err)
	// }
	// for _, file := range filenames {
	// 	n, b, err := readFileOS(file)
	// 	if err != nil {
	// 		Ошибка(" Ошибка парсинга каталога с шаблонами HTML %+v\n", err)
	// 	}
	// 	Инфо("  %+v \n", n)
	// 	Инфо("  %+s \n", b)
	// }
	// Инфо(" filenames %+v \n", filenames)
	СырыеШаблоны, errParseGlob = template.New("").ParseGlob(Конфиг.КаталогШаблонов + "*/*.html")
	// СырыеШаблоны = template.Must(template.New("index").Funcs(РендерФункции()).ParseGlob(Конфиг.КаталогШаблонов + "*/*.html"))
	if errParseGlob != nil {
		Ошибка("  %+v \n", errParseGlob)
	}
	log.Printf("СырыеШаблоны %+v \n", СырыеШаблоны.Tree)
	if errParseGlob != nil {
		Ошибка("Ошибка парсинга каталога с шаблонами HTML %+v\n", errParseGlob)
	}

}
func readFileOS(file string) (name string, b []byte, err error) {
	name = filepath.Base(file)
	b, err = os.ReadFile(file)
	return
}
