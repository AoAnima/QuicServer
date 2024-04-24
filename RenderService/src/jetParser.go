package main

import (
	"bytes"

	"gitverse.ru/Ao/jet"

	. "aoanima.ru/ConnQuic"
	. "aoanima.ru/Logger"
)

var views = jet.NewSet(
	// jet.NewOSFileSystemLoader(ДирректорияЗапуска+"/"+Конфиг.КаталогШаблонов),
	jet.NewOSFileSystemLoader("./jetHTML"),
	jet.InDevelopmentMode(), // remove or set false in production
	jet.WithTemplateNameExtensions([]string{".html", ".js", ".jet"}),
)

// git config --global url."git@gitverse.ru:2222/Ao/jet.git".insteadOf "https://gitverse.ru/sc/Ao/jet.git"

// git config --global url."git@gitverse.ru:2222/Ao/jet".insteadOf "https://gitverse.ru/sc/Ao/jet"

func JetПарсингШаблонов() {
	Инфо("JetПарсингШаблонов views1 %+s \n", ДирректорияЗапуска+"/"+Конфиг.КаталогШаблонов)

	view, err := views.GetTemplate("/index.jet")
	if err != nil {
		Ошибка("Unexpected template err: %+v", err.Error())
	}

	// var w io.Writer
	w := new(bytes.Buffer) // создаем буфер и присваиваем его переменной w

	view.Execute(w, nil, nil)
	// ПатернПарсингаШаблонов := ДирректорияЗапуска + "/" + Конфиг.КаталогШаблонов
	// httpfsLoader, err := httpfs.NewLoader(templates.Assets)
	// if err != nil {
	// 	Ошибка()
	// }

}
