package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	. "aoanima.ru/ConnQuic"
	. "aoanima.ru/Logger"
)

func ПарсингШаблонов() {
	// "pattern": "../www/tpl/*/*.html",
	// var errParseGlob error
	Инфо(" ПарсингШаблонов ДирректорияЗапуска %+v \n", ДирректорияЗапуска)
	ПатернПарсингаШаблонов := ДирректорияЗапуска + "/" + Конфиг.КаталогШаблонов
	// ПатернПарсингаШаблонов := ДирректорияЗапуска + "/" + Конфиг.КаталогШаблонов + "*/*/*.html"

	СырыеШаблоны = template.New("")
	// JavaScript = js.New("JS")
	СырыеШаблоны.Funcs(РендерФункции())
	// JavaScript.Funcs(РендерФункции())

	err := filepath.WalkDir(ПатернПарсингаШаблонов, func(путь string, описание os.DirEntry, err error) error {
		if описание.IsDir() {
			// Инфо("  парсим %+v \n", путь+"/*.html")

			_, err := СырыеШаблоны.ParseGlob(путь + "/*.html")
			// Инфо("следующийШаблон %+v \n", СырыеШаблоны.Templates())
			if err != nil {
				Ошибка("Ошибка парсинга шаблона : %+s путь %+s ", err.Error(), путь)
			}
		}

		return nil
	})

	if err != nil {
		Ошибка(" WalkDir %+v \n", err.Error())
	}
	РендерJS()
}

// првоеряет все ноды в СырыеШаблоны, если имя шаблона оканчивается на .js, то выполняет этот шаблон, и результат записыват в файл scripts.js
func РендерJS() {
	// пройди по дереву СырыеШАблоны, проверь define имя каждого шаблона, если оно содержит .js, выполни этот шабон и верни результат в качестве template.JS , запиши в файл
	путьJSфайл := ДирректорияЗапуска + "/" + Конфиг.КаталогСтатичныхФайлов + "/js/scripts.js"
	if _, err := os.Stat(путьJSфайл); err == nil {
		if err := os.Remove(путьJSфайл); err != nil {
			Ошибка("Error removing file: %+v", err.Error())
		}
	}
	файл, err := os.OpenFile(путьJSфайл, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0777)

	if err != nil {
		Ошибка(" ошибка открытия файла %+v \n", err.Error())

	}
	клонШаблонов, err := СырыеШаблоны.Clone()
	if err != nil {
		Ошибка(" ОписаниеОшибки %+v \n", err.Error())
	}

	defer файл.Close()
	for _, шаблон := range клонШаблонов.Templates() {
		if strings.HasSuffix(шаблон.Name(), ".js") {
			var буфер bytes.Buffer
			// Инфо(" %+v \n", шаблон.Name())
			// Инфо(" %+v \n", шаблон.Tree.Root.Nodes)
			// for _, n := range шаблон.Tree.Root.Nodes {
			// 	if defineNode, ok := n.(*parse.TemplateNode); ok {
			// 		// Инфо("Имя определения: %s\n", defineNode.Name)
			// 	}
			// }

			if err := шаблон.Execute(&буфер, nil); err != nil {
				Ошибка(" Ошибка выполнения шаблона %+v \n", err.Error())
			}

			// jsШаблон := template.JS(буфер.String())
			// htmlШаблон := template.HTML(буфер.String())
			Инфо("jsШаблhtmlШаблонон %+v \n", буфер.String())
			// Инфо("уфбуферер.String() %+v \n", буфер.String())

			if _, err := файл.WriteString(string(буфер.String())); err != nil {
				Ошибка(" ошибка записи в файл ]%+v \n", err.Error())
			}
		}
	}

	файл.Close()

}

func writeFileOS(name string, data []byte, perm os.FileMode) error {
	return os.WriteFile(name, data, perm)
}

func readFileOS(file string) (name string, b []byte, err error) {
	name = filepath.Base(file)
	b, err = os.ReadFile(file)
	return
}
