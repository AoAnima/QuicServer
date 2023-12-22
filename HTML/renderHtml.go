package render

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	_ "fmt"
	"html/template"
	"log"
	ctx "main/RWContext"
	. "main/loger"
	"path/filepath"
	_ "path/filepath"
	"strconv"
	"strings"
	"time"
)



// GenHtml парсит шаблоны и генерируют готовый html
/*
	В шаблонах можно использовать другие шаблоны
для этого нужно создать html файл с нужной разметкой, и дать ему имz шаблона равное имени блока через {{define "name"}}{{end}}
после чего вставлять в нушное место через {{template "name" .name}}

Либо использовать имена блоков в качесте ключей с блоками данных
{{width .name}}
	тут будут доступны данные из
BlockDate[name]={
	"key":"value"
}
{{end}}

*/
func GenHtml(RequestCtx ctx.RWContext) ctx.RWContext {

	//var response responseStruct
	//response := ctx.Value("response").
	// Достаём из основного контекста распарсенные html файлы

	log.Printf("RequestCtx %+v\n", Log(RequestCtx))
	tplFiles := ctx.MainCtx.Value("tplFiles").(*template.Template)

	// Клонируем в новую переменную
	TplBlock, err := tplFiles.Clone()
	if err != nil {
		log.Printf("%+v\n TplBlock err", err)
	}

	var goTpl string
	var LayoutName string
	//var BlocksHtml=make(map[string]string)
	RequestCtx.BlocksHtml = make(map[string]string)
	// Если В ответе из БД поля со строками шаблона Page и Layout пусты то проверяем поле с Блоками

	if RequestCtx.PageData.GoTpl == "" && RequestCtx.LayoutData.GoTpl == "" {
		log.Printf("RequestCtx.PageData.GoTpl ==  && RequestCtx.LayoutData.GoTpl ==  %+v", RequestCtx.Blocks != nil)
		if RequestCtx.Blocks != nil {
			//Генерируем HTML для каждго блока и помещаем в карту с ключами = Имени блока
			/*
				получаеться что каждый блок будет отрендерин и помещён в карту BlocksHtml{"ИмяБлока":"html строка") это отлично когда нужно на стороне клинета вставить несколько блоков в разные места с помощью JS

			*/
			for _, BlockName := range RequestCtx.Blocks {

				log.Printf("[BlockName] %+v\n", Log(BlockName))
				log.Printf("RequestCtx.BlocksData[BlockName] %+v\n", Log(RequestCtx.BlocksData[BlockName]))
				resHtml := new(bytes.Buffer)

				//if RequestCtx.BlocksData[BlockName]["block_go_tpl"] != "" {
				//	TplBlock, err = TplBlock.Parse(RequestCtx.BlocksData[BlockName]["block_go_tpl"].(string) )
				//}
				//if RequestCtx.BlocksData[BlockName].BlockGoTpl != "" {
				//	TplBlock, err = TplBlock.Parse(RequestCtx.BlocksData[BlockName].BlockGoTpl)
				//}
				err := TplBlock.ExecuteTemplate(resHtml, BlockName, RequestCtx.BlocksData[BlockName])
				if nil != err {
					log.Printf("%+v\n", err)
				}
				//BlocksHtml[BlockName]=resHtml.String()

				RequestCtx.BlocksHtml[BlockName] = resHtml.String()

			}
		} else {
			log.Printf("%+v\n", "Нет данных для генерации HTML")
		}
	} else {
		// Если в ответе из БД строка шаблона Page и/или Layout не пусты
		//log.Printf("%+v\n %+v","RequestCtx.BlocksData " ,Log(RequestCtx))к

		// если ajax запрос то парсим строку только с контентом (пример: {{define \"page\"}}{{template \"new_reposter\" .new_reposter}}{{template \"reposter_list\" .reposter_list}}{{end}} ) с именем шаблона LayoutName = page

		if RequestCtx.Request.RequestMethod == "ajax/fetch"  {
			if	RequestCtx.Request.LayoutName != "" {
				goTpl = RequestCtx.LayoutData.GoTpl //RequestCtx.PageData.GoTpl
				LayoutName = RequestCtx.LayoutData.Name
			} else {
				goTpl = RequestCtx.PageGotpl //RequestCtx.PageData.GoTpl
				LayoutName = "page"

			}
		} else {
			// иначе если первая загрузка сайта то парсим GoTpl строку всего шаблона LayoutGotpl и выполняем шаблон с именем LayoutName
			goTpl = RequestCtx.LayoutData.GoTpl
			LayoutName = RequestCtx.LayoutData.Name
		}

		log.Printf("goTpl %+v\n", goTpl)
 		log.Printf("LayoutName %+v\n", LayoutName)
 		log.Printf("RequestCtx %+v\n", RequestCtx)
		//log.Printf("RequestCtx.BlocksData %+v\n", Log(RequestCtx.BlocksData))
		TplBlock, err = TplBlock.Parse(goTpl)
		resHtml := new(bytes.Buffer)

		if errs := TplBlock.ExecuteTemplate(resHtml, LayoutName, RequestCtx); errs != nil {
			log.Printf("%+v\n", errs)
		}


		if RequestCtx.Request.RequestMethod == "ajax/fetch" {
			if	RequestCtx.Request.LayoutName != "" {
				RequestCtx.BlocksHtml["main_layout"] = resHtml.String()
			} else {
				RequestCtx.BlocksHtml["page_content"] = resHtml.String()
				}
		} else {
			RequestCtx.HtmlContent = resHtml.String()
		}
	}
	log.Printf("RequestCtx.BlocksHtml 126 %+v\n", Log(RequestCtx.BlocksHtml))

	//ctx = context.WithValue(ctx, "RequestCtx", responseData)

	return RequestCtx
}
type MenuStruct struct {
	Access      []string
	Child       []int
	Description string
	Icon        string
	Id          int
	Layout      []string
	Name        string
	Parrent     int
	Url         string
	Children    map[int]template.HTML `json:"Children"`
}

/*
preRenderMenu Преобразует входящий  map[string]interface{} в MenuStruct с помощью JSON сериализации и десириализации - можно в будущем оптимизировать и переписать преобразование по другому.
Для 10 минктов меню время рендера около 5мс. Можно его кэшировать в БД чтобы не делать каждый раз при загрузке/обнволении всей страницы. Но 5 мс это оченьмало и для каждого пользователя делаеься 1 раз.....

*/

func preRenderMenu (menuData map[string]interface{}) map[int]template.HTML {
	start := time.Now()
	menu := map[int]*MenuStruct{}
	jsonMenu, _ := json.Marshal(menuData)
	json.Unmarshal(jsonMenu, &menu)
	menuHtml := map[int]template.HTML{}
	for id := range menu {
		if menu[id].Parrent ==0{
			menuHtml[id] = renderMenu(menu, id)[id]
		}
	}
	fmt.Printf("время рендера меню %v\n",  time.Since(start))
	return menuHtml
}
/*
renderMenu рекурсивная функция, вызывает сама себя если в пункте меню содержаться Child > 0, иначе если Child = 0 то Вызывает renderMenuItem
*/
func renderMenu (menuData map[int]*MenuStruct, itemId int) map[int]template.HTML {
	itemHtml := map[int]template.HTML{}

	if len(menuData[itemId].Child) >0 {
		for _, ChildId := range menuData[itemId].Child{
			itemHtml[ChildId]  =renderMenu(menuData, ChildId)[menuData[ChildId].Id]
			//log.Printf("181 itemHtml%+v\n", itemHtml)
			menuData[menuData[ChildId].Parrent].Children = make(map[int]template.HTML, len(menuData[itemId].Child))
			menuData[menuData[ChildId].Parrent].Children = itemHtml
		}
	}
		itemHtml[itemId] =renderMenuItem(menuData[itemId])
	//log.Printf("187 itemHtml%+v\n", itemHtml)
	return itemHtml
}

/*
renderMenuItem Парсит template файлы, и рендерит HTML
*/
func renderMenuItem(menuElement *MenuStruct) template.HTML {
	tplFiles := ctx.MainCtx.Value("tplFiles").(*template.Template)
	TplBlock, err := tplFiles.Clone()
	CheckErr(err)
	var menuItem template.HTML
			var htmlMenuItem = new(bytes.Buffer)
			TplBlock.ExecuteTemplate(htmlMenuItem, "menu_item", menuElement)
			menuItem = template.HTML(htmlMenuItem.String())
	return menuItem
}

func getUserAgent(httpHeaders map[string][]string) int {
	//log.Printf("%+v\n", string(httpHeaders["User-Agent"][0][0]))
	index := strings.Index(httpHeaders["User-Agent"][0], "rv:")
	//log.Printf("%+v\n", index)
	indexVersion := index+3;
	//log.Printf("%+v\n", Log(httpHeaders["User-Agent"][0][indexVersion:indexVersion+2]))
	ver,_ := strconv.Atoi(httpHeaders["User-Agent"][0][indexVersion:indexVersion+2])
	log.Printf("Версия браузера %+v\n", ver)
	return ver
}

func getStruct(blocks map[string]interface{}, blockStruct interface{}) map[string]interface{}{
	log.Printf("%+v\n", Log(blocks))
	log.Printf("%+v\n", blockStruct)
	nestedBlocks := map[string]interface{}{}

	for _, blockName := range blockStruct.([]interface{}){
		nestedBlocks[blockName.(string)] = blocks[blockName.(string)]
	}

	return nestedBlocks
}

/*
tplFunc набор функций для выхова из езд шаблонов
*/
func tplFunc() map[string]interface{} { // функции для обработка template
	var funcMap = template.FuncMap{
		"getStruct":getStruct,
		"getUserAgent":getUserAgent,
		//"groupsGetCatalogInfo": func (access_token string)  map[string]interface{}{
		//params:= map[string]string{
		//	"extended":"0",
		//	"subcategories":"1",
		//	}
		//
		//vkCatalogInfo:=Vk("groups","getCatalogInfo","", access_token, params, "")
		//return vkCatalogInfo
		//},
		"renderMenu":preRenderMenu,
		"getGroupsIds": func(data map[string]interface{}) string {
			var result string
			result = ""

			for key, _ := range data {
				//strings.join([],",")
				result += key + ","
			}
			log.Printf("%+v\n", result)
			result = result[:len(result)-1]
			log.Printf("%+v\n", result)
			return result
		},
		"toJSON": func(data map[string]interface{}) string {
			val, _ := json.Marshal(data["attr_params"])
			//log.Printf("%+v\n", data)
			return string(val)
		},
		"getTime": func() int64 {
			val := time.Now().Unix()
			//log.Printf("%+v\n",val)
			return val
		},
		"toString": func(data interface{}) string {

			return data.(string)
		},
		"toInt": func(data float64) int {
			//val, _ := json.Marshal(data["attr_params"])
			//log.Printf("%+v\n", data)
			return int(data)
		},
		"Itoa": func(data float64) string {
			//val, _ := json.Marshal(data["attr_params"])
			//log.Printf("%+v\n", data)

			return strconv.Itoa(int(data))
		},
		"TrimSuffix": func(str string, trimStr string) string {
			result := strings.TrimSuffix(str, trimStr)
			return string(result)
		},
		"firstValue": func(data []interface{}) map[string]interface{} {
			response := make(map[string]interface{})

			for _, value := range data {
				if value.(map[string]interface{})["attr_id"] == 0.0 {
					response = value.(map[string]interface{})
				}
			}

			return response
		},
		"parseJsFiles": func() map[int]string {
			//response := make(map[string]interface{})
			//print("parseJsFiles ")
			filelist, _ := filepath.Glob("./src/www/static/js/*.js")
			response := make(map[int]string)
			for key, value := range filelist {
				response[key] = filepath.ToSlash(value)
				//log.Printf("%+v\n"," filepath.Clean(value)-----------------")
				//log.Printf("%+v\n", filepath.ToSlash(value))

			}
			//log.Printf("%+v\n", Log(response))
			//log.Printf("%+v\n", response)
			return response
		},
		"pastBlock": func(data map[string]interface{}) template.HTML {
			var tplBytesBuffer = new(bytes.Buffer)
			tplFiles := ctx.MainCtx.Value("tplFiles").(*template.Template)
			TplBlock, err := tplFiles.Clone()
			CheckErr(err)
			TplBlock.ExecuteTemplate(tplBytesBuffer, data["tpl"].(string), data)

			return template.HTML(tplBytesBuffer.String())
		},
		"prints": func(data interface{}) string {

			//log.Printf("%+v\n", Log(data))
			return Log(data)
		},
		"jsFunc": func(data interface{}) string {

			//log.Printf("%+v\n", Log(data))
			//log.Printf("%+v\n", Log(data))
			return Log(data)
		},
		"month": func() []int {
			month := make([]int, 31)
			idx := 1
			start := time.Now()
			for idx <= 31 {
				month[idx-1]=idx
				idx++
			}
			fmt.Printf("время рендера дней месяца %v\n",  time.Since(start))
			return month
		},
	}
	return funcMap
}
func ParseTplDir(pattern map[string]string) {

	var errParseGlob error
	tplFiles := ctx.MainCtx.Value(pattern["name"])

		tplFiles = template.Must(template.New("index").Funcs(tplFunc()).ParseGlob(pattern["pattern"]))

		log.Printf("%+v\n", tplFiles)
		if errParseGlob != nil {
			log.Printf("Ошибка парсинга каталога с шаблонами HTML %+v\n", errParseGlob)

		}
	log.Printf("%+v\n", "tplFiles создан и помещён в MainCtx")
	log.Printf("%+v\n", "Нужно периодически обновлять переменную")
	//log.Printf("%+v\n", Log(tplFiles))
	ctx.MainCtx = context.WithValue(ctx.MainCtx, pattern["name"], tplFiles)

}
func check(err error) bool {
	errorResult := false
	if err != nil {
		log.Printf("Ошибка! Необходимо обработать %+v\n", Log(err))
		errorResult = true
	}
	return errorResult
}
