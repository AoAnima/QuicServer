package main

import (
	"bytes"
	"html/template"
	"log"
	"os"
	"path/filepath"

	. "aoanima.ru/ConnQuic"
	. "aoanima.ru/Logger"
	fsnotify "github.com/fsnotify/fsnotify"
)

var СырыеШаблоны *template.Template

func ОтрендеритьОтветКлиенту(сообщение *Сообщение) {

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

	if err != nil {
		Ошибка("   %+v \n", err)

	}

}

/*!
Правила построения шаблонов буду описывать json
например /dashboard/profile?name="username"
{
	"content":"dashboard" // он же ИмяБазовогоШаблона
	"subcontent": "path1"
}

*/

func ПолныйРендер(сообщение *Сообщение) error {
	Инфо("  %+v \n", "ПолныйРендер")
	БуферHtml := new(bytes.Buffer)

	ШаблонДляРендера, err := СырыеШаблоны.Clone()
	if err != nil {
		Ошибка(" Ошибка клоинрования сырых шаблонов %+v \n", err)
		return err
	}

	// Так как это полный рендер страницы, а в index.html шаблон для основного контэнта помечен как content которого физически не существует, то необходимо создать новый шаблон с именем content и добавить в него дерефо нужного шаблона, каталог или товар или личный кабинет и т.д.в зависимости от запрошенной страницы
	// имяШаблона := сообщение.Запрос.ИмяБазовогоШаблона
	// Инфо("  %+v Tree %+v \n", имяШаблона, ШаблонДляРендера.Lookup(string(имяШаблона)).Tree)

	/*
	 ДОКУМЕНТАЦИЯ
	 т.к. url.path может быт myt ограниченной вложенности, то в приницпе в теории и глубина вложенных подшаблонов может быть не ограниченной,
	 то буду именовать вложенный конетнт по правилу:
	  имяшаблона_content
	  к примеру /dashboard/adress вложенный контент dashboard будет именован dashboard_content - в котрый будет вставлено html tree из шаблона adress

	  если /dashboard/adress/catalog
	  то вложенный контент будет именован dashboard_content, а вложеный шаблон в adress будет именован adress _content - в котрей будет вставлен html tree из шаблона catalog

	  ПРИ ЭТОМ вложенный шаблн должен быть задан по умолчанию, и существовать, потому что например при входу на dashboard у нас нет второго пути и не известен dashboard_content, но мы создаём файл с dashboard_content в который вставляем шаблон по умолчанию...
	  можно обявляеть его в том же базовом шаблоне, ниже сонвоного

	*/

	// пострим дерево вложенных шаблонов
	картаМаршрута := сообщение.Запрос.КартаМаршрута
	// базовыйШаблон := картаМаршрута[0]
	базовыйШаблон := string(сообщение.Ответ[ИмяСервиса(картаМаршрута[0])].ИмяШаблона)

	switch базовыйШаблон {
	case "/":
		базовыйШаблон = "main"
	case "":
		базовыйШаблон = картаМаршрута[0]
	}
	// if базовыйШаблон == "/" {
	// 	базовыйШаблон = "main"
	// }

	if len(картаМаршрута) > 1 {
		for _, имяСервиса := range картаМаршрута[1:] {

			имяШаблона := string(сообщение.Ответ[ИмяСервиса(имяСервиса)].ИмяШаблона)
			if имяШаблона == "" {
				имяШаблона = string(имяСервиса)
			}
			Инфо("имяСервиса %+v ; имяШаблона %+v ; базовыйШаблон_content = %+v \n", имяСервиса, имяШаблона, базовыйШаблон+"_контент")

			вложенныйШаблон := ШаблонДляРендера.Lookup(string(имяШаблона))
			if вложенныйШаблон == nil {
				Инфо("Не  удаётся найти вложеный шаблон с именем имяШаблона %+v \n", имяШаблона)
				вложенныйШаблон = ШаблонДляРендера.Lookup("неВерноеИмяШаблона")
			}

			ШаблонДляРендера = template.Must(ШаблонДляРендера.AddParseTree(базовыйШаблон+"_контент", вложенныйШаблон.Tree))
			if err != nil {
				Ошибка("  %+v \n", err)
			}
		}

	}
	Инфо("базовыйШаблон  %+v \n", базовыйШаблон)
	деревоБазовогоШаблона := ШаблонДляРендера.Lookup(базовыйШаблон)
	//картаМаршрута[0] - имя базового Шаблона - content
	if деревоБазовогоШаблона == nil {
		Инфо("Не  удаётся найти базовыйШаблон шаблон с именем базовыйШаблон %+v \n", базовыйШаблон)
		деревоБазовогоШаблона = ШаблонДляРендера.Lookup("неВерноеИмяШаблона")
	}
	ШаблонДляРендера, err = ШаблонДляРендера.AddParseTree("content", деревоБазовогоШаблона.Tree)

	if err != nil {
		Ошибка("  %+v \n", err)
	}

	if errs := ШаблонДляРендера.ExecuteTemplate(БуферHtml, "index", сообщение.Ответ); errs != nil {
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
	Инфо(" сообщение.Ответ %+v \n", сообщение.Ответ)
	Инфо(" сообщение.Запрос.КартаМаршрута %+v \n", сообщение.Запрос.КартаМаршрута)

	for индекс, имяСервиса := range сообщение.Запрос.КартаМаршрута[1:] {
		Инфо("  %+v   %+v \n", индекс, имяСервиса)
		имяШаблона := string(сообщение.Ответ[ИмяСервиса(имяСервиса)].ИмяШаблона)

		if имяШаблона == "" {
			имяШаблона = string(имяСервиса)
		}
		Html, err := Рендер(имяШаблона, сообщение.Ответ[ИмяСервиса(имяСервиса)])
		if err != nil {
			Ошибка("  %+v \n", err)
			return err
		}
		ответКлиенту.AjaxHTML[имяШаблона] = ДанныеAjaxHTML{
			Цель: string(имяШаблона),
			HTML: string(Html),
		}
	}

	сообщение.ОтветКлиенту = ответКлиенту
	return nil
}

func Рендер(имяШаблона string, Данные interface{}) ([]byte, error) {
	// func Рендер(имяШаблона string, КартаДанных map[ИмяШаблона]КартаДанныхШаблона) ([]byte, error) {
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
	деревоШаблонаБлока := ШаблонДляРендера.Lookup(имяШаблона)
	//картаМаршрута[0] - имя базового Шаблона - content
	if деревоШаблонаБлока == nil {
		Инфо("Не  удаётся найти имяШаблона шаблон с именем имяШаблона %+v \n", имяШаблона)
		деревоШаблонаБлока = ШаблонДляРендера.Lookup("неВерноеИмяШаблона")
	}
	ШаблонДляРендера, err = ШаблонДляРендера.AddParseTree("content", деревоШаблонаБлока.Tree)
	if err != nil {
		Ошибка("  %+v \n", err)
	}
	Инфо(" ШаблонДляРендера  %+v  \n  имяШаблона  %+v  \n  Данные %+v \n", ШаблонДляРендера, имяШаблона, Данные)

	if errs := ШаблонДляРендера.ExecuteTemplate(Html, имяШаблона, Данные); errs != nil {
		Ошибка("%+v\n", errs)
		return nil, errs
	}

	return Html.Bytes(), nil
}

func ПарсингШаблонов() {
	// "pattern": "../www/tpl/*/*.html",
	// var errParseGlob error
	Инфо(" ПарсингШаблонов ДирректорияЗапуска %+v \n", ДирректорияЗапуска)
	ПатернПарсингаШаблонов := ДирректорияЗапуска + "/" + Конфиг.КаталогШаблонов
	// ПатернПарсингаШаблонов := ДирректорияЗапуска + "/" + Конфиг.КаталогШаблонов + "*/*/*.html"

	СырыеШаблоны = template.New("")
	err := filepath.WalkDir(ПатернПарсингаШаблонов, func(путь string, описание os.DirEntry, err error) error {
		if описание.IsDir() {

			// Инфо("  парсим %+v \n", путь+"/*.html")

			_, err = СырыеШаблоны.ParseGlob(путь + "/*.html")
			if err != nil {
				Ошибка("Ошибка парсинга шаблона : %+s путь %+s ", err.Error(), путь)
			}
		}

		return nil
	})
	if err != nil {
		Ошибка(" WalkDir %+v \n", err)
	}
	// err := filepath.Walk("path/to/directory", func(path string, info os.FileInfo, err error) error {
	// 	// Проверяем, является ли текущий элемент файлом
	// 	if !info.IsDir() {
	// 		// Парсим шаблон из файла
	// 		_, err := СырыеШаблоны.ParseFiles(path)
	// 		if err != nil {
	// 			log.Println("Failed to parse template:", err)
	// 		}
	// 	}
	// 	return nil
	// })

	// Инфо(" Конфиг.КаталогШаблонов  %+v \n", ПатернПарсингаШаблонов)
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
	// СырыеШаблоны, errParseGlob = template.New("").ParseGlob(ПатернПарсингаШаблонов)
	// СырыеШаблоны = template.Must(template.New("index").Funcs(РендерФункции()).ParseGlob(Конфиг.КаталогШаблонов + "*/*.html"))
	// if errParseGlob != nil {
	// 	Ошибка("  %+v \n", errParseGlob)
	// }
	// шаблоны := СырыеШаблоны.Templates()
	// for _, шаблон := range шаблоны {
	// 	Инфо("шаблон  %+v \n", шаблон)
	// }
	// Инфо("СырыеШаблоны  %+v \n")
	// log.Print(СырыеШаблоны.Templates())

	// Инфо("СырыеШаблоны %+v \n", СырыеШаблоны.Lookup("index"))
	// log.Print(СырыеШаблоны.Lookup("index"))
	// log.Print(*СырыеШаблоны)
	// if errParseGlob != nil {
	// 	Ошибка("Ошибка парсинга каталога с шаблонами HTML %+v\n", errParseGlob)
	// }

}

func readFileOS(file string) (name string, b []byte, err error) {
	name = filepath.Base(file)
	b, err = os.ReadFile(file)
	return
}

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
