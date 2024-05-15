package main

import (
	"bytes"
	"os"

	jet "gitverse.ru/Ao/jet"

	. "aoanima.ru/ConnQuic"
	. "aoanima.ru/Logger"
)

var НаборШаблонов *jet.Set

// git config --global url."git@gitverse.ru:2222/Ao/jet.git".insteadOf "https://gitverse.ru/sc/Ao/jet.git"

// git config --global url."git@gitverse.ru:2222/Ao/jet".insteadOf "https://gitverse.ru/sc/Ao/jet"

func JetПарсингШаблонов() {

	// Инфо("JetПарсингШаблонов views1 %+s \n", "./jetHTML")
	Инфо("НаборШаблонов %+v \n", НаборШаблонов)

	НаборШаблонов.ПарсингДиректорииШаблонов()
	СобратьJS(НаборШаблонов)

	// НаборШаблонов.ПарсингШаблона("/контент/админ/формы/формаНовогоОбработчика.html")
	Инфо("НаборШаблонов  %+v \n", НаборШаблонов)

}
func СобратьJS(НаборШаблонов *jet.Set) {

	путьJSфайл := ДирректорияЗапуска + "/" + Конфиг.КаталогСтатичныхФайлов + "/js/scripts.js"
	if _, err := os.Stat(путьJSфайл); err == nil {
		if err := os.Remove(путьJSфайл); err != nil {
			Ошибка("Error removing file: %+v", err.Error())
		}
	}
	файл, err := os.OpenFile(путьJSфайл, os.O_CREATE|os.O_WRONLY, 0777)
	defer файл.Close()
	if err != nil {
		Ошибка(" ошибка открытия файла %+v  %+v \n", err.Error(), файл)

	}
	if err := файл.Truncate(0); err != nil {
		Ошибка("Error truncating file: %+v", err.Error())
	}

	JSШаблон := "{{блок JavaScript()}}"

	НаборШаблонов.JsБлоки.Обойти(func(имяБлока, Блок any) bool {

		// Инфо("  %+v %+v \n", имяБлока, Блок.(*jet.BlockNode).String())
		JSШаблон += "{{вставить " + имяБлока.(string) + "()}}"

		return true
	})

	JSШаблон += "{{конец}}"
	Инфо(" JSШаблон %+v \n", JSШаблон)

	шаблон, ошибка := НаборШаблонов.Парсинг("JavaScript", JSШаблон, false)
	if ошибка != nil {
		Ошибка(" ОписаниеОшибки  %+v \n", ошибка.Error())
	}
	БуферHtml := new(bytes.Buffer)
	ошибка = шаблон.Execute(БуферHtml, nil, nil)
	if ошибка != nil {
		Ошибка("  %+s \n", ошибка.Error())
	}
	// Инфо(" %+v \n", БуферHtml.String())
	if _, err := файл.WriteString(string(БуферHtml.String())); err != nil {
		Ошибка(" ошибка записи в файл ]%+v \n", err.Error())
	}

}
