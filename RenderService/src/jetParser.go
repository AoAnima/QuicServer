package main

import (
	"github.com/CloudyKit/jet/examples/asset_packaging/assets/templates"
	"github.com/CloudyKit/jet/v6"
	"github.com/CloudyKit/jet/v6/loaders/httpfs"

	. "aoanima.ru/ConnQuic"
	. "aoanima.ru/Logger"
)

var views = jet.NewSet(
	jet.NewOSFileSystemLoader("./views"),
	jet.DevelopmentMode(true), // remove or set false in production
)
var views *jet.Set



привет 
ыва
ghbdtn 

func JetПарсингШаблонов() {
	ПатернПарсингаШаблонов := ДирректорияЗапуска + "/" + Конфиг.КаталогШаблонов
	httpfsLoader, err := httpfs.NewLoader(templates.Assets)
	if err != nil {
		Ошибка()
	}

	views = jet.NewSet(
		httpfsLoader,
		jet.DevelopmentMode(true), // remove or set false in production
	)
}
