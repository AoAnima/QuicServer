package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	. "main/RWContext"
	"main/auth"
	"main/connect"
	rpc "main/jsonrpc"
	. "main/loger"
	"main/render"
	"net/http"
	"net/url"
	_ "strconv"
	"strings"
	"time"
)

/*
MainHandler основная точка входа http запроса
Инициирует в памяти контекст запроса RequestCtx RWContext{} и заполняет его данными для отправки в БД, получает ответ из БД, анализирует нужно ли произвести дополнительную обработку с данными на строне сервера, вызывает функцию парсинга и рендера HTML и возвращает ответ клиенту
*/
func MainHandler(w http.ResponseWriter, r *http.Request) {
	r.Host = strings.Replace(r.Host, "www.", "", -1)
	//dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	//if err != nil {
	//	log.Fatal(err)
	//}
	pattern := map[string]string{
		"name":    "tplFiles",
		"pattern": "../www/tpl/*/*.html",
		//"pattern": dir+"/www/tpl/*/*.html",
	}
	render.ParseTplDir(pattern)
	RequestCtx := RWContext{}

	//Doc: Парсим запрос, куки, заголовки, методы и прочее
	RequestCtx.Request.Header = r.Header
	RequestCtx.Request.RequestMethod = r.Header.Get("Request-Method")
	RequestCtx.Request.Cookies = MapCookies(r)
	RequestCtx.Request.Host = r.Host
	RequestCtx.Request.URL = r.URL
	RequestCtx.Request.Method = r.Method
	RequestCtx.Request.RequestURI = r.RequestURI  //сожержит строку после домена, пример /auth?action=actionName
	ParseQuery, _ := url.ParseQuery(r.URL.RawQuery)
	RequestCtx.Request.Get = ParseQuery
		log.Printf("%+v\n", r)
		log.Printf("%+v\n",r.RequestURI)
		log.Printf("%+v\n", ParseQuery.Get("jsonql"))
	if ParseQuery["jsonql"] != nil {

		var JsonQl map[string]interface{}
		json.Unmarshal([]byte(ParseQuery.Get("jsonql")), &JsonQl)
		RequestCtx.Request.JsonQl = JsonQl

	}

	if ParseQuery["jsonrpc"] != nil {
		RequestCtx = rpc.ParseJsonRpc(RequestCtx)
		fmt.Fprint(w, RequestCtx.Request.Get)
	}

	RequestCtx = checkUserAuthCookies(RequestCtx)
	//log.Printf("%+v\n", Log(RequestCtx))
	RequestCtx = checkPostAction(RequestCtx, r)

	RequestCtx = ГенерируемОтвет(ОбменДаннымиСБД(RequestCtx))
	//ResponseResult := ГенерируемОтвет(СтруктурированиеДанных(ОбменДаннымиСБД(checkUserAuthCookies(RequestCtx))))

	// устанавливаем заголовки ответ
	w = setHeaders(w, RequestCtx.Request.RequestMethod)
	if RequestCtx.HttpStatus == 0 {
		RequestCtx.HttpStatus = 200
	}


	// Устанавливаем cookie с токеном авторизации
	if RequestCtx.Request.Post.Action == "auth" && RequestCtx.UserSession.Token != "" {
		//h :=(24 - time.Now().Hour())

		h := time.Duration(23 - time.Now().Hour()) * time.Hour
		m := time.Duration(59 - time.Now().Minute()) * time.Minute
		s := time.Duration(60 - time.Now().Second()) * time.Second
		expiration := time.Now().Add(h+m+s)
		loc, _ := time.LoadLocation("Europe/Moscow")
		cookieToken := http.Cookie{
			Name:    "token",
			Value:   RequestCtx.UserSession.Token,
			Expires: expiration.In(loc),
			RawExpires: expiration.Format(time.UnixDate),
		}
		log.Printf("%+v\n", Log(cookieToken))
		http.SetCookie(w, &cookieToken)
		cookieUid := http.Cookie{
			Name:    "uid",
			Value:   RequestCtx.UserSession.Uid,
			Expires: expiration.In(loc),
			RawExpires: expiration.Format(time.UnixDate),
		}
		http.SetCookie(w, &cookieUid)
	}

	w.WriteHeader(RequestCtx.HttpStatus)
	// отадём данные клиенту
	fmt.Fprint(w, string(RequestCtx.ResponseData))
}

/*
checkPostAction Проверяет пришёл ли POST запрос, если да то парсит тело запроса, проверяет Action поле, если его значение удовлетворяет одному из условий то выполняет соответсвующие действия с данными

Action =
	auth - авторизация пользователя через LDAP
	logout - выход пользователя с портала
*/


func checkPostAction(RequestCtx RWContext, r *http.Request) RWContext{
	if r.Method == http.MethodPost {
		//Если пришёл POST запрос парсим форму
		r.ParseMultipartForm(1024)
		RequestCtx.Request.Post.PostedData = r.Form
		//RequestCtx.Request.Post.Action = strings.Trim(r.RequestURI, "/")
		RequestCtx.Request.Post.Action = r.Form["action"][0]


		/*
			doc: auth единственное действие которое выполняется до отправки запроса в локальную БД т.к. обращаеться к внешнему LDAP серверу
		*/
		if RequestCtx.Request.Post.Action == "auth" {
			authData := auth.AuthDate{
				RequestCtx.Request.Post.PostedData.Get("uid")+"@r26",
				RequestCtx.Request.Post.PostedData.Get("password"),
			}
			RequestCtx.UserSession.LdapUser = auth.Auth(authData)
			log.Printf("%+v\n", Log(RequestCtx.UserSession))

			if RequestCtx.UserSession.LdapUser.Error != "" {
				log.Printf("Ошибка авторизации ldap  %+v\n", RequestCtx.UserSession.LdapUser.Error)
			}
		} else if RequestCtx.Request.Post.Action == "logout"{

		}
	}
	return RequestCtx
}

//func СтруктурированиеДанных (RequestCtx RWContext) (RWContext) {
//	for blockKey, block := range RequestCtx.BlocksData {
//		//log.Printf("%+v\n", Log(RequestCtx.BlocksData[block]))
//		log.Printf("%+v\n", Log(blockKey))
//		log.Printf("%+v\n", Log(block))
//		//if block["inner_blocks"]!=nil {
//		//	for _, blockName := range block["inner_blocks"].([]interface{}){
//		//		block[blockName.(string)]=RequestCtx.BlocksData[blockName.(string)]
//		//	}
//		//}
//		if block.Struct!=nil {
//			for _, blockName := range block.Struct{
//				block.Content[blockName]=RequestCtx.BlocksData[blockName]
//			}
//			log.Printf("block %+v\n", Log(block))
//		}
//
//	}
//	log.Printf("%+v\n", Log(RequestCtx))
//return RequestCtx
//}
func ГенерируемОтвет(RequestCtx RWContext) RWContext {
	/**
	Нужно сделать так что если в запросе флоормат json то не генерировать html иначе генерировать
	*/
	RequestCtx = render.GenHtml(RequestCtx)
	//log.Printf("%+v\n", Log(RequestCtx))
	var ResponseResult []byte
	if RequestCtx.Request.RequestMethod != "ajax/fetch" {
		RequestCtx.ResponseData = []byte(RequestCtx.HtmlContent)
	} else {
		type AjaxResponseStruct struct {
			MetaData struct {
				Keywords    []string `json:"keywords"`
				Description string   `json:"description"`
				Title       string   `json:"title"`
				HashTags    []string `json:"HashTags"`
			}
			JsFunc []struct {
				Func   string `json:"Func"`
				Required []string `json:"Required"`
				Params struct {
					Target string                   `json:"target"`
					Data   []map[string]interface{} `json:"data"`
					Blocks []string                 `json:"blocks"`
				} `json:"Params"`
			} `json:"JsFunc"`
			BlocksOrder []string          `json:"BlocksOrder"` // Порядок следования блоков на странице соответсвует struct полю в таблицах
			BlocksHtml  map[string]string `json:"BlocksHtml"`
			Messages    []MessagesType `json:"Messages"`
			Error       bool
		}
		// если AJAX запрос с параметром только HTML
		var AjaxResponse AjaxResponseStruct

		AjaxResponse.MetaData = RequestCtx.PageData.MetaData

		if RequestCtx.JsFunc != nil {
			//log.Printf("%+v\n", Log(RequestCtx.JsFunc))
			AjaxResponse.JsFunc = RequestCtx.JsFunc
		}
		if RequestCtx.BlocksHtml != nil {
			AjaxResponse.BlocksHtml = RequestCtx.BlocksHtml
		}
		if RequestCtx.Blocks != nil {
			AjaxResponse.BlocksOrder = RequestCtx.Blocks
		}
		if RequestCtx.Messages != nil {
			AjaxResponse.Messages = RequestCtx.Messages
		}
		if RequestCtx.Error != false {
			AjaxResponse.Error = RequestCtx.Error
		}

		log.Printf("%+v\n", Log(AjaxResponse))
		var err error
		ResponseResult, err = jsonStringify(AjaxResponse)
		//ResponseResult = string(responseByte)
		if nil != err {
			log.Printf("%+v\n", Log(err))
		} else {
			RequestCtx.ResponseData = ResponseResult
		}

	}
	log.Printf("Ответ сгенерирован \n")
	//log.Printf("ResponseResult %+v\n", string(ResponseResult))
	return RequestCtx
}



func ОбменДаннымиСБД(RequestCtx RWContext) RWContext {

	//log.Printf("%+v\n", RequestCtx)
	//var db *sql.DB

	//if RequestCtx.Request.R.Host != Config.Host {
	//	log.Printf("HostDB \n", )
	db := connect.MainDB("")
	//} else {
	//	log.Printf("Dbconnect \n", )
	//	db = MainCtx.Value("db").(*sql.DB)
	//}

	if db == nil {

		RequestCtx = RequestCtx.SetDeffault(RequestCtx)
		log.Printf("RequestCtx %+v\n", Log(RequestCtx))

		RequestCtx = RequestCtx.SetError(RequestCtx, map[string]interface{}{"Message":"Отсутствует подклюбчение к базе данных", "page":"error"})
		return RequestCtx
	}
	RequestCtxJson, errStringify := json.Marshal(RequestCtx)
	if errStringify != nil {
		log.Printf("%+v\n", errStringify)
	}

	log.Printf("Запрос в БД %+v\n", string(RequestCtxJson))

	rows, errDB := db.Query("select rout($1)", RequestCtxJson)
	defer db.Close()

	defer rows.Close()

	var result []byte

	//doc: обрабатываем ответ от БД
	if errDB == nil {
		for rows.Next() {
			//log.Printf("%+v\n", rows)
			//log.Printf("%+v\n", Log(rows))
			rows.Scan(&result)
			//log.Printf("Ответ из БД: %+v\n", string(result))
		}
		errorUnmarshal := json.Unmarshal(result, &RequestCtx)
		if errorUnmarshal != nil {
			log.Printf("errorUnmarshal ДБ result %+v\n", errorUnmarshal)

		}
	} else if errDB != nil {
		db.Close()
		//log.Panic(errDB)
		log.Printf("%+v %+v %+v\n", "Ошибка запроса ", errDB.Error())
		log.Printf("%+v %+v %+v\n", "Ошибка запроса ", errDB, rows)
	}



	if RequestCtx.Error == true {
		log.Printf("%+v\n", Log(RequestCtx.ErrorMessages))
		for _, dbError := range RequestCtx.ErrorMessages{
			if dbError.Code == 401 {

			}
		}
	}

	//log.Printf("%+v\n", string(result))
	//var test interface{}
	//errror := json.Unmarshal(result, &test)
	//log.Printf("%+v\n", errror)
	//log.Printf("%+v\n", test)



	return RequestCtx
}

//func ResponseWrite(w http.ResponseWriter, r *http.Request, ctx context.Context) {
//
//	w = setHeaders(w, ctx)
//
//
//	response := ctx.Value("responseData").(response)
//
//	w.WriteHeader(response.HttpStatus)
//
//	var res string
//	//log.Printf("%+v\n response206 ", Log(response))
//	if r.Header.Get("request-method") != "ajax/fetch" {
//		res = response.HtmlContent
//	} else {
//		type AjaxResponse struct {
//			PageHtml string `json:"pageHtml"'`
//			JsFunc map[string]map[string]interface{} `json:"JsFunc"'`
//			Blocks map[string]string `json:"blocks"'`
//		}
//		// если AJAX запрос с параметром только HTML
//		var ajaxResponse AjaxResponse
//
//		//if response.PageData !=nil {
//		//	ajaxResponse.PageHtml=response.PageData["html"].(string)
//		//}
//		if response.JsFunc != nil {
//			log.Printf("%+v\n", Log(response.JsFunc))
//			ajaxResponse.JsFunc=response.JsFunc
//		}
//		if response.BlocksHtml!= nil {
//			ajaxResponse.Blocks=response.BlocksHtml
//		}
//
//		log.Printf("%+v\n", Log(ajaxResponse))
//		responseByte, err := jsonStringify(ajaxResponse);
//		res = string(responseByte)
//		if nil!=err{log.Printf("%+v\n", Log(err))}
//	}
//	log.Printf("%+v\n", Log(res))
//	//startTime  := ctx.Value("startTime").(time.Time)
//	log.Printf("Время ответа: %v\n", time.Since(ctx.Value("startTime").(time.Time)))
//	fmt.Fprint(w, res)
//}

func checkUserAuthCookies(RequestCtx RWContext) RWContext {

	if uid := RequestCtx.Request.Cookies["uid"]; uid != nil {
		RequestCtx.Uid = uid.Value
		RequestCtx.UserSession.Uid = uid.Value
	} else {
		RequestCtx.UserSession.Uid = "unauth"
		RequestCtx.Uid = "unauth"
	}

	if token := RequestCtx.Request.Cookies["token"]; token != nil {
		RequestCtx.UserSession.Token = token.Value
		//if SessionMap[token.Value].Uid == RequestCtx.Uid {
		//	RequestCtx.UserSession = SessionMap[token.Value]
		//} else {
		//	log.Printf("Имя пользователя '%+v' не совпадает с токеном  '%+v', необходима повторная авторизация \n", RequestCtx.Uid, token)
		//	RequestCtx.UserSession = SessionMap["unauth"]
		//	RequestCtx.UserSession.HostAuth = RequestCtx.Request.Host
		//}
	} else {
		RequestCtx.UserSession.Token = "unauth"

		//RequestCtx.UserSession = SessionMap["unauth"]
		//log.Printf("%+v\n", RequestCtx.Request.Host)
		RequestCtx.UserSession.HostAuth = RequestCtx.Request.Host
	}



	return RequestCtx
}

/*Вспомогательные функции */

func StaticHandler(w http.ResponseWriter, req *http.Request) {
	var static_file string
	reqFile := req.URL.Path

	if strings.Contains(reqFile, "src") {
		static_file = req.URL.Path[len("/srv/src/www/static/"):]
	} else {
		static_file = req.URL.Path[len("/static/"):]
	}

	if len(static_file) != 0 {

		f, err := http.Dir("../www/static/").Open(static_file)
		//fmt.Println(f)
		if err == nil {
			content := io.ReadSeeker(f)
			http.ServeContent(w, req, static_file, time.Now(), content)
			return
		} else {
			log.Printf("%+v\n", Log(err))
		}
	}
	http.NotFound(w, req)
}

func setHeaders(w http.ResponseWriter, RequestMethod string) http.ResponseWriter {
	//userQuery := ctx.Value("userQuery").(QueryData)
	//if r.Header.Get("request-method") != "ajax/fetch" {
	var contentType string
	if RequestMethod != "ajax/fetch" {
		contentType = "text/html"
	} else {
		contentType = "application/json"
	}

	w.Header().Set("Accept-Charset", "utf-8")
	w.Header().Set("Access-Control-Allow-Origin", "true")
	w.Header().Set("Access-Control-Allow-Credentials", "include")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
	//w.Header().Set("Access-Control-Allow-Headers","Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	w.Header().Set("Accept-Encoding", "gzip, deflate")
	w.Header().Set("Accept-Language", " ru-RU,ru;q=0.8")
	w.Header().Set("Content-Type", strings.Join([]string{contentType, "; charset=utf-8"}, ""))
	//lenght := unsafe.Sizeof(ctx.Value("response").(responseStruct))

	//w.Header().Set("Content-Length",  string(lenght))

	return w
}
func getCookie(r *http.Request, name string) string {
	var value string
	// cookie НЕ установлены
	if len(r.Cookies()) <= 0 {
		value = "nil"
	} else {
		// cookie установлены
		//log.Printf("%+v\n", name)
		if cookieName, err := r.Cookie(name); cookieName != nil {
			//log.Printf("%+v\n", cookieName.Value)
			value = cookieName.Value
			//user = username.Value
		} else if err != nil {
			//log.Printf("%+v\n", "cookie", cookieName.Value, "не установлены, возвращаем nil")
			log.Printf("%+v\n", err)
			value = "nil"
		}
	}
	//log.Printf("%+v\n", value)
	return value
}
