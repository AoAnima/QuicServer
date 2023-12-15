package main

import (
	"bytes"
	"html/template"
	"log"
	"os"
	"path/filepath"

	. "aoanima.ru/ConnQuic"
	. "aoanima.ru/Logger"
	"github.com/fsnotify/fsnotify"
)

var СырыеШаблоны *template.Template

func ПолучитьОтветКлиенту(сообщение *Сообщение) {

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
	// if сообщение.Запрос.ТипЗапроса == GET || сообщение.Запрос.ТипЗапроса == POST {
	// ШаблонДляРендера, err := СырыеШаблоны.Clone()
	// if err != nil {
	// 	Ошибка("  %+v \n", err)
	// }
	// Рендер("index", сообщение.Запрос.Шаблонизатор)
	// найдём шаблон который нужно вставить в блок content
	// ШаблонКонтента := ШаблонДляРендера.Lookup(string(сообщение.Запрос.ИмяШаблона))
	// ШаблонДляРендера, err = ШаблонДляРендера.AddParseTree("content", ШаблонКонтента.Tree)
	// Рендер("index", сообщение.Ответ[ИмяСервиса(сообщение.Запрос.ИмяШаблона)])
	// }
	var err error
	// var Html []byte
	switch сообщение.Запрос.ТипЗапроса {
	case GET:
		err = ПолныйРендер(сообщение)
	case POST:
		err = ПолныйРендер(сообщение)
	case AJAX:
		err = РендерБлоков(сообщение)
	case AJAXPost:
		err = РендерБлоков(сообщение)
	}

	// получается что ничего не возвращаем, сразу пишем изменения в исходное сообщение в раздел сообщение.ОтветКлиенту

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
	// Инфо("  %+s \n", Html)

}

func ПолныйРендер(сообщение *Сообщение) error {
	Инфо("  %+v \n", "ПолныйРендер")
	БуферHtml := new(bytes.Buffer)

	ШаблонДляРендера, err := СырыеШаблоны.Clone()
	if err != nil {
		Ошибка(" Ошибка клоинрования сырых шаблонов %+v \n", err)
		return err
	}

	// Так как это полный рендер страницы, а в index.html шаблон для основного контэнта помечен как content которого физически не существует, то необходимо создать новый шаблон с именем content и добавить в него дерефо нужного шаблона, каталог или товар или личный кабинет и т.д.в зависимости от запрошенной страницы
	имяШаблона := сообщение.Запрос.ИмяБазовогоШаблона
	Инфо("  %+v Tree %+v \n", имяШаблона, ШаблонДляРендера.Lookup(string(имяШаблона)).Tree)
	ШаблонДляРендера.AddParseTree("content", ШаблонДляРендера.Lookup(string(имяШаблона)).Tree)

	// Инфо(" content %+v \n ", ШаблонДляРендера.Lookup("content").Tree)

	КартаДанных := сообщение.Запрос.Шаблонизатор
	// Всегда рендерим шаблон "index" , данные для конетнта будут добавленны из ИмяБазовогоШаблона
	if errs := ШаблонДляРендера.ExecuteTemplate(БуферHtml, "index", КартаДанных[имяШаблона].Данные); errs != nil {
		Ошибка("%+v\n", errs)
		log.Print(errs)
		return errs
	}
	сообщение.ОтветКлиенту = ОтветКлиенту{
		HTML: БуферHtml.Bytes(),
	}
	return nil
}

func РендерБлоков(сообщение *Сообщение) error {

	ответКлиенту := ОтветКлиенту{
		AjaxHTML: make(map[string]ДанныеAjaxHTML),
	}

	for ИмяШаблона, _ := range сообщение.Запрос.Шаблонизатор {
		Html, err := Рендер(string(ИмяШаблона), сообщение.Запрос.Шаблонизатор)
		if err != nil {
			Ошибка("  %+v \n", err)
			return err
		}

		//! TODO: Сбосоп вставки будет определяться на клиенте, в тжгах куда вставлять будем писать data атрибут в котром будем указывать как вставляються данные,
		// data-update-method="replaceWith" или data-update-method="append"
		ответКлиенту.AjaxHTML[string(ИмяШаблона)] = ДанныеAjaxHTML{
			Цель: string(ИмяШаблона),
			HTML: string(Html),
			// СпособВставки: Заменить, // способ вставки - нужно придумать где хранить и как определять, либо храним в БД , либо в ajax запросе, например запрос путь в адресной строке catalog/page=2
			// а ajax запрос будет в заивисомсти от События вызвавшее запрос, добавлять каокй нибудь метод в ajax запрос "updateMethod": "replaceWith" ...

			// Хрень а если я буду возвращать несколько бооков... значит способ вставки должен храниться в базе, рядом с данными о шаблонах и сервисах из котрых получаем данные для этих шаблонов
		}
	}
	сообщение.ОтветКлиенту = ответКлиенту
	return nil
}

func Рендер(имяШаблона string, КартаДанных map[ИмяШаблона]КартаДанныхШаблона) ([]byte, error) {
	Html := new(bytes.Buffer)

	ШаблонДляРендера, err := СырыеШаблоны.Clone()
	if err != nil {
		Ошибка(" Ошибка клоинровании сырых шаблонов %+v \n", err)
		return nil, err
	}
	// content - виртуальный шаблон, который вставляеться в body.html , его содержимое меняется в зависимости от загружаемой страницы, при условии что требуется полная загрузка HTML , тоесть был обычный запрос GET или POST не ajax
	// если имяШаблона передано как
	// if имяШаблона == "index" {
	// 	ШаблонДляРендера.AddParseTree("content", ШаблонДляРендера.Lookup("main").Tree)
	// }

	Инфо(" ШаблонДляРендера  %+v имяШаблона  %+v   КартаДанных %+v \n", ШаблонДляРендера, имяШаблона, КартаДанных)

	if errs := ШаблонДляРендера.ExecuteTemplate(Html, имяШаблона, КартаДанных[ИмяШаблона(имяШаблона)].Данные); errs != nil {
		Ошибка("%+v\n", errs)
		return nil, errs
	}

	return Html.Bytes(), nil
}

func ПарсингШаблонов() {
	// "pattern": "../www/tpl/*/*.html",
	var errParseGlob error
	Инфо(" ДирректорияЗапуска %+v \n", ДирректорияЗапуска)
	ПатернПарсингаШаблонов := ДирректорияЗапуска + "/" + Конфиг.КаталогШаблонов + "*/*.html"

	Инфо(" Конфиг.КаталогШаблонов  %+v \n", ПатернПарсингаШаблонов)
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
	СырыеШаблоны, errParseGlob = template.New("").ParseGlob(ПатернПарсингаШаблонов)
	// СырыеШаблоны = template.Must(template.New("index").Funcs(РендерФункции()).ParseGlob(Конфиг.КаталогШаблонов + "*/*.html"))
	if errParseGlob != nil {
		Ошибка("  %+v \n", errParseGlob)
	}
	// Инфо("СырыеШаблоны %+v \n", СырыеШаблоны.Lookup("index"))
	// log.Print(СырыеШаблоны.Lookup("index"))
	// log.Print(*СырыеШаблоны)
	if errParseGlob != nil {
		Ошибка("Ошибка парсинга каталога с шаблонами HTML %+v\n", errParseGlob)
	}

}

// func readFileOS(file string) (name string, b []byte, err error) {
// 	name = filepath.Base(file)
// 	b, err = os.ReadFile(file)
// 	return
// }

// func получитьВложенныеДиректории(directory string) ([]string, error) {
// 	var подкаталоги []string

// 	err := filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
// 		if err != nil {
// 			return err
// 		}
// 		if info.IsDir() && path != directory {
// 			подкаталоги = append(подкаталоги, path)
// 		}
// 		return nil
// 	})

// 	return подкаталоги, err
// }

// Наблюдает за изменениями HTMl шаблонов и перечитывает их
func наблюдатьЗаИзменениямиШаблонов() {
	каталогНаблюдения := ДирректорияЗапуска + "/" + Конфиг.КаталогШаблонов
	Инфо("наблюдатьЗаИмзенениями   %+v \n", каталогНаблюдения)
	// Создаем новый Watcher
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	// Устанавливаем новый лимит наблюдаемых директорий

	// Добавляем директорию для наблюдения
	err = watcher.Add(каталогНаблюдения)
	if err != nil {
		Ошибка("  %+v \n", err)
	}

	// Инфо("watcher  %+v \n", watcher.WatchList())

	err = filepath.Walk(каталогНаблюдения, func(директория string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() && директория != каталогНаблюдения {
			err = watcher.Add(директория)
			if err != nil {
				Ошибка("  %+v \n", err)
			}
		}
		return nil
	})
	if err != nil {
		Ошибка("  %+v \n", err)
	}
	// Инфо("  %+v \n", подкаталоги)

	// Бесконечный цикл для обработки событий
	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return
			}
			if event.Op&fsnotify.Write == fsnotify.Write {
				Инфо("Изменен файл:", event.Name)
				hub.broadcast <- []byte("reload")
				ПарсингШаблонов()
				// Делайте необходимые действия при изменении файла
			}
		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}
			Ошибка("  %+v \n", err)
		}
	}
}

func наблюдатьЗаИзменениямиСтатичныхФайлов() {
	каталогНаблюдения := ДирректорияЗапуска + "/" + Конфиг.КаталогСтатичныхФайлов
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()
	// Добавляем директорию для наблюдения
	err = watcher.Add(каталогНаблюдения)
	if err != nil {
		Ошибка("  %+v \n", err)
	}

	err = filepath.Walk(каталогНаблюдения, func(директория string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() && директория != каталогНаблюдения {
			err := watcher.Add(директория)
			if err != nil {
				Ошибка("  %+v \n", err)
			}
		}
		return nil
	})
	if err != nil {
		Ошибка("  %+v \n", err)
	}

	// Инфо("watcher  %+v \n", watcher.WatchList())
	// Бесконечный цикл для обработки событий
	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return
			}
			if event.Op&fsnotify.Write == fsnotify.Write {
				Инфо("Изменен файл отправить сообщение в браузер:", event.Name)
				hub.broadcast <- []byte("reload")
				// Делайте необходимые действия при изменении файла
			}
		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}
			Ошибка("  %+v \n", err)
		}
	}

}
