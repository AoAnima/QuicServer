package main

import (
	"github.com/CloudyKit/jet/v6"

	. "aoanima.ru/ConnQuic"
	. "aoanima.ru/Logger"
)

var views = jet.NewSet(
	jet.NewOSFileSystemLoader(ДирректорияЗапуска+"/"+Конфиг.КаталогШаблонов),
	jet.InDevelopmentMode(), // remove or set false in production
	jet.WithTemplateNameExtensions([]string{".html", ".js"}),
)

func JetПарсингШаблонов() {
	// ПатернПарсингаШаблонов := ДирректорияЗапуска + "/" + Конфиг.КаталогШаблонов
	// httpfsLoader, err := httpfs.NewLoader(templates.Assets)
	// if err != nil {
	// 	Ошибка()
	// }

}
