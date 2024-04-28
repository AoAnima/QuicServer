package main

import (
	jet "gitverse.ru/Ao/jet"

	// . "aoanima.ru/ConnQuic"
	. "aoanima.ru/Logger"
)

var НаборШаблонов *jet.Set

// (
// 	// jet.NewOSFileSystemLoader(ДирректорияЗапуска+"/"+Конфиг.КаталогШаблонов),
// 	jet.NewOSFileSystemLoader(ДирректорияЗапуска+"/"+Конфиг.КаталогШаблонов),
// 	jet.InDevelopmentMode(), // remove or set false in production
// 	jet.WithTemplateNameExtensions([]string{"", ".html", ".js"}),
// )

// git config --global url."git@gitverse.ru:2222/Ao/jet.git".insteadOf "https://gitverse.ru/sc/Ao/jet.git"

// git config --global url."git@gitverse.ru:2222/Ao/jet".insteadOf "https://gitverse.ru/sc/Ao/jet"

func JetПарсингШаблонов() {

	// Инфо("JetПарсингШаблонов views1 %+s \n", "./jetHTML")
	Инфо(" %+v \n", НаборШаблонов)
	// НаборШаблонов.ПарсингДиректорииШаблонов()
	НаборШаблонов.ПарсингШаблона("/контент/админ/формы/формаНовогоОбработчика.html")
	НаборШаблонов.ПоказатьШаблоныйВКэше()
	// шаблон, err := Шаблоны.GetTemplate("/index.jet")
	// if err != nil {
	// 	Ошибка("Unexpected template err: %+v", err.Error())
	// }

	// // var w io.Writer
	// w := new(bytes.Buffer) // создаем буфер и присваиваем его переменной w

	// шаблон.Execute(w, nil, nil)

	// Инфо(" %+v \n  шаблон %+v \n", w.String(), шаблон)

	// ПатернПарсингаШаблонов := ДирректорияЗапуска + "/" + Конфиг.КаталогШаблонов
	// httpfsLoader, err := httpfs.NewLoader(templates.Assets)
	// if err != nil {
	// 	Ошибка()
	// }

}
