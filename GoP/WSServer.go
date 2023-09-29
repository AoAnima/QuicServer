package main

import (
	"bytes"
	"context"
	_ "encoding/json"
	"flag"
	"github.com/gorilla/websocket"
	_ "github.com/lib/pq"
	_ "github.com/nakagami/firebirdsql"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"
)



type Аргументы struct{
	Название string `json:"название"`
}

type ДанныеОтвета struct {
	Контейнер string `json:"контейнер"`
	Данные interface{} `json:"данные"`
	HTML string `json:"html"`
	Обработчик string `json:"обработчик"` //JS функция или объект/класс/плагин для обработки данных (table..)
}



var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}


var server = Server{
	Messages:  []*Сообщение{},
	Clients:  map[string]*Client{},
	addCh:     nil,
	delCh:     nil,
	sendAllCh: nil,
	doneCh:    nil,
	errCh:     nil,
}

type Server struct {
	//pattern   string
	Messages  []*Сообщение
	Clients   map[string]*Client
	addCh     chan *Client
	delCh     chan *Client
	sendAllCh chan *Message
	doneCh    chan bool
	errCh     chan error
}


var logFile *os.File
var ProgDir string
var IoLoger *IoLog
func init (){
	var err error
	ProgDir, err = filepath.Abs(filepath.Dir(""+os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	IoLoger = &IoLog{}
	log.SetFlags(log.Ltime|log.Lshortfile)

	/*
		Логируем все log.Printf в базу
	*/
	IoLoger.DB, _ = PGConnect("logs", nil)
	IoLoger.QueryName="main_log"
	IoLoger.ctx, IoLoger.cancel = context.WithCancel(context.Background())
	_, err = IoLoger.DB.Prepare(IoLoger.ctx, IoLoger.QueryName, "INSERT INTO wsserver_log (time,log) VALUES ($1,$2)", nil)
	if err != nil{
		log.Printf("IoLoger.DB.Prepare err %+v IoLoger.QueryName %+v\n", err, IoLoger.QueryName)
	}

	mw := io.MultiWriter(os.Stdout, IoLoger)

	log.SetOutput(mw)
	Actions = map[string]interface{}{
			"getChatLog": (*Client).ПолучитьЛогПереписки,
			"getIoMenu": (*Client).ПолучитьМенюБота,
			"collectData":(*Client).СобратьДанные,
			"GetData":(*Client).ЗагрузитьДанные,
			"SSHConnect":(*Client).SSHConnect,
			"CloseSSH":(*Client).СloseSSH,
			"CreateSkillName":(*Client).CreateSkillName,
			"CreateSkillDescription":(*Client).CreateSkillDescription,
			"CreateSkillCmd":(*Client).CreateSkillCmd,
			"AddSkillTags":(*Client).AddSkillTags,
			"AllowSelfUse":(*Client).AllowSelfUse,
			"AddProblemDescription":(*Client).AddProblemDescription,
			"ShowNewSkill":(*Client).ShowNewSkill,
			"ИзменитьНавык":(*Client).ИзменитьНавык,
			"ПоказатьНавык":(*Client).ПоказатьНавык,
			"ПоказатьНавыки":(*Client).ПоказатьНавыки,
			"РедакторНавыка":(*Client).РедакторНавыка,
			"СохранитьИзмененияНавыка":(*Client).СохранитьИзмененияНавыка,
			"НовыйНавык":(*Client).НовыйНавык,
			"СоздатьНавык":(*Client).СоздатьНавык,
			"УдалитьНавык":(*Client).УдалитьНавык,
			"ПроверитьЛогин":(*Client).ПроверитьЛогин,
			"Авторизация":(*Client).Авторизация,
			"СоздатьЗаявку":(*Client).СоздатьЗаявку,
			"ПоказатьСписокЗаявок":(*Client).ПоказатьСписокЗаявок,
			"СоздатьРабочийСтол":(*Client).СоздатьРабочийСтол,
			"ПоказатьБазуЗнаний":(*Client).ПоказатьБазуЗнаний,
			"ВыполнитьSQLвАИС":(*Client).ВыполнитьSQLвАИС,
			"загрузить патч":(*Client).СохранитьФайл,
			"ОбновитьДанныеПоПК":(*Client).ОбновитьДанныеПоПК,
			"СинхронихироватьКлючи":(*Client).СинхронихироватьКлючи,
			"СинхронихироватьДатыКриптоПро":(*Client).СинхронихироватьДатыКриптоПро,
			//"EndSkillCreate":(*Client).EndSkillCreate,
			//"SkillTags":(*Client).SkillTags,
			//"GetData":DataCollector,
			//"CollectData": DataCollector,
		}



}

func main(){

	flag.Parse()



	//var addr = flag.String("addr", "10.26.6.13:8080", "chat server")
	// рендерим рабочий стол
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request){
		log.Printf("*************Входящий запрос INDEX ***********\n")
		index(w,r)
	})
	http.HandleFunc("/static/", StaticHandler)
	http.HandleFunc("/wsconnect", func(w http.ResponseWriter, r *http.Request){
		log.Printf("*************Входящий запрос wsconnect ***********\n")
		 connected(w,r)
	})

	//GetLdapUsers("")
	//log.Fatal(http.ListenAndServe(*addr, nil))
	serverError := http.ListenAndServe(":80", nil)
	if serverError != nil {
		Ошибка(">>>> serverError ERROR \n %+v \n\n", serverError)
	}
	//log.Printf("serverError %+v\n", serverError)
}

func (server *Server)УведомитьВсехОПодключении (clientOnline *Client, канал chan string){
	for sendLogin, client := range server.Clients{
		//log.Printf("sendLogin %+v clientOnline.Login %+v\n", )
		if sendLogin!=clientOnline.Login{
			responseMes := Сообщение{
				Id:      0,
				От:      "io",
				Кому:    sendLogin,
				Online:  clientOnline.Login,
				Content:  struct {
					Target string `json:"target"`
					Data interface{} `json:"data"`
					Html string `json:"html"`
					Обработчик string `json:"обработчик"` // функция или объект для обработки данных (handsontable..)
				}{
					Data: clientOnline.UserInfo,
				},
			}
			client.Message<-&responseMes
		}
	}
	канал <- "конец УведомитьВсехОПодключении "+clientOnline.Login

}

func (server *Server)УведомитьВсехОбОтключении(code int, text string, login string, канал chan string) error {
	for sendLogin, client := range server.Clients{
		responseMes := Сообщение{
			Id:      0,
			От:      "io",
			Кому:    sendLogin,
			Offline: login,
		}
		client.Message<-&responseMes
	}
	//log.Printf("УведомитьВсехОбОтключении code %+v text  %+v\n", code, text,)
	//log.Printf("УведомитьВсехОбОтключении server %+v\n", server.Clients)

	канал<-"конец УведомитьВсехОбОтключении "+login
	return nil
}

func index(writer http.ResponseWriter, request *http.Request) {

	Token, err := request.Cookie("Token")
	if err!= nil{
		log.Printf("err Token не найден, рендерим форму авторизации %+v\n", err)
		//tplName = "auth"
	}
	tplName := "index"
	var data interface{}

	if Token != nil{
		//log.Printf(" Token %+v\n", Token.Value)
		ПолучитьДанныеСессии, err := sqlStruct{
			Name:   "user_session",
			Sql:    "select * from iobot.user_session where token = $1",
			Values: [][]byte{
				[]byte(Token.Value),
			},
		}.RunSQL(nil)
		if err != nil {
			log.Printf(">>>> ERROR \n %+v \n\n", err)
		}
		if len(ПолучитьДанныеСессии)>0{
			ДанныеСессии := ПолучитьДанныеСессии[0]
			log.Printf("ДанныеСессии %+v\n", ДанныеСессии)

							data = map[string]interface{}{
								"ContentData":map[string]interface {}{
									"tplName":"dashboard",
									"tplData":map[string]string{"login":ДанныеСессии["uid"].(string)},
								},
							}

			//if ДанныеАвторизации["date_auth"] != ""{
			//	log.Printf("ДанныеАвторизации[date_auth] %+v\n", ДанныеАвторизации["date_auth"])
			//	времяАвторизации, errPARSETIME := time.Parse("2006-01-02T15:04:05.999999", ДанныеАвторизации["date_auth"].(string))
			//
			//	if errPARSETIME != nil {
			//		log.Printf(">>>> ERROR \n %+v \n\n", err)
			//	}
			//	сейчас := time.Now()
			//
			//	log.Printf("сейчас %+v\n", сейчас)
			//
			//	сейчас.Format("2006-01-02T15:04:05.999999")
			//
			//	//log.Printf("времяАвторизации - сейчас %+v\n", времяАвторизации - сейчас)
			//	разница := сейчас.Unix()-времяАвторизации.Unix()
			//	log.Printf("разница %+v разница / 60 %+v\n", разница, разница / 60 )
			//	// Если разница больше 5 дней, то просим ввести пароль но подставляем последнее фио
			//	if разница > 5 {
			//		Результат ,err:= sqlStruct{
			//				Name:   "fssp_configs",
			//				Sql:    "SELECT second_name,givenname,initials,login  FROM fssp_configs.users WHERE login = $1",
			//				Values: [][]byte{
			//					[]byte(ДанныеАвторизации["uid"].(string)),
			//				},
			//				DBSchema:"fssp_configs",
			//			}.RunSQL(nil)
			//		if err != nil{
			//
			//		log.Printf(">>>> Ошибка SQL запроса: %+v \n\n",err)
			//			//ContentHtml = render("auth", nil)
			//		} else {
			//			 if len(Результат)>0{
			//				data = map[string]interface{}{
			//					"ContentData":map[string]interface {}{
			//						"tplName":"auth",
			//						"tplData":Результат[0],
			//					},
			//				}
			//				 //ContentHtml = render("auth", Результат[0])
			//			 }
			//		}
			//
			//
			//	}
			//}
		} else {
			log.Printf("Данные авторизации не найдены в базе данных, токен не найден %+v\n", Token)
			IP := strings.Split(request.RemoteAddr,":")[0]

			data = map[string]interface{}{
				"ContentData":map[string]interface {}{
					"tplName":"auth",
					"tplData":ОпределитьПользователяПоIp(IP),
				},
			}
			//ContentHtml = render("auth", nil)
		}
	} else {
		IP := strings.Split(request.RemoteAddr,":")[0]
		data = map[string]interface{}{
			"ContentData":map[string]interface {}{
				"tplName":"auth",
				"tplData":ОпределитьПользователяПоIp(IP),
			},
		}
		//ContentHtml = render("auth", nil)
	}

	_, errWrite := writer.Write(render(tplName, data))
	if errWrite != nil {
		log.Print("render:", errWrite)
	}
}

func connected (w http.ResponseWriter, r *http.Request){
	канал := make(chan string, 3)
	go wsConnector(w,r,канал)
	for {
		_ = <-канал
		//log.Printf("connected result %+v\n", result)
	}
}

func  wsConnector(w http.ResponseWriter, r *http.Request, канал chan string) {
	log.Printf("http.Request %+v\n", r.Header.Get("Origin"))
	queryArgs, _ := url.ParseQuery(r.URL.RawQuery)

	var ЛогинСПортала string
	if _, ЕстьЛогинССайта := queryArgs["login"];ЕстьЛогинССайта{
		ЛогинСПортала = queryArgs["login"][0]
	}

	//log.Printf("ЛогинСПортала %+v\n", ЛогинСПортала)

	conn, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		log.Print("upgrade:", err)
		return
	}

	IP := strings.Split(r.RemoteAddr,":")[0]


	client := &Client{
		Ip: 	 IP,
		Ws:      conn,
		Message: make(chan *Сообщение),
	}

	if ЛогинСПортала == "" {
		client.ОпределитьПользователяПоIp(IP)
		ЛогинСПортала = client.Login
	} else {
		if ЛогинСПортала == "auth"{
			client.Login =  IP
		} else {
			client.Login= ЛогинСПортала
			//log.Printf("ПолучитьДанныеПользователя %+v\n", )
			client.ПолучитьДанныеПользователя("")
			//log.Printf("ПолучитьНaстройки %+v\n", )
			client.ПолучитьНaстройки()
			//log.Printf("закончили ПолучитьНaстройки %+v\n", )
		}

	}

	client.Ws.SetCloseHandler(func(code int, text string)error{
		//server.Clients[userName]=nil
		//log.Printf("Удаляем запись о пользователи из памяти %+v server.Clients %+v\n", userName , server.Clients)
		if ПоказыватьЛогиВБраузере, ok := client.Setting["show_logs"];ok{
			if ПоказыватьЛогиВБраузере==true {
				mw := io.MultiWriter(os.Stdout,IoLoger) //, server.Clients["maksimchuk@r26"]
				log.SetOutput(mw)
			}
		}
		if server.Clients[client.Login] != nil && server.Clients[client.Login].Ssh == nil{
			// тут возможно ислкючение когда струткура Ssh ещё пустая но ведёться попытка  оединиться с клиентом, тогда может выпасть ошибка

			// note: возможно нужно сделать задержку перед удалением из памяти, и добавить ключ онлайн оффлан, чтобы не писать/удалять постоянно клитента если он просто перешёл на новую страницу или обновил страницу...
			delete(server.Clients, client.Login)
		}
		//log.Printf("server.Clients %+v . %+v удалён %+v\n", server.Clients, userName, server.Clients[userName]==nil)

		 go server.УведомитьВсехОбОтключении(code , text , client.Login, канал)

		//client = nil
		//close(client.message)
		return nil
	})

	server.Clients[client.Login] = client

	//log.Printf("server.Clients %+v\n", server.Clients[client.Login])
	//fooType := reflect.TypeOf(client)
	//log.Printf(" fooType %+v\n", fooType)
	//log.Printf(" fooType %+v\n", fooType.NumMethod())
	//
	//for i := 0; i < fooType.NumMethod(); i++ {
	//	method := fooType.Method(i)
	//	log.Printf("method.Name %+v\n", method.Name)
	//	log.Printf("method.Name %+v\n", method.Type)
	//	//log.Printf("method.Name %+v\n", method.Func)
	//
	//}

	go client.SendMessage()
	go client.ReadMessage()

	Token, errToken := r.Cookie("Token")

	if errToken!= nil{
		log.Printf("err Token не найден %+v\n", err)
		//tplName = "auth"
	}


	if queryArgs["reconect"] != nil && queryArgs["reconect"][0]=="true"{
	} else {
		if ЛогинСПортала != "auth" && Token != nil {
			client.Token= struct {
				Hash     string
				Истекает string
			}{
				Hash:Token.Value,
				Истекает: Token.Expires.String(),
			}
			go client.СоздатьМессенджер() //server.Clients[userName]
		} else if Token == nil && (ЛогинСПортала != "auth" && ЛогинСПортала !="") && r.Header.Get("Origin") == "http://10.26.4.20"{
				client.СозадтьИСохранитьТокен()
			go client.СоздатьМессенджер()
		}
	}
	if ЛогинСПортала != "auth" && Token !=nil{
		go server.УведомитьВсехОПодключении(client, канал)
	}
	if client.Setting != nil{
		if client.Setting["show_logs"] == true{
			mw := io.MultiWriter(os.Stdout,IoLoger, client) //, server.Clients["maksimchuk@r26"]
			log.SetOutput(mw)
		}
		//log.Printf("client.Login %+v\n", client.Login)
		//mw := io.MultiWriter(os.Stdout, client) //, server.Clients["maksimchuk@r26"]
		//log.SetPrefix(">> ")
		//log.SetOutput(mw)
	}
	//log.Printf("server %+v\n", server)

	канал<-"конец"
}


func StaticHandler(w http.ResponseWriter, req *http.Request) {
	var static_file string
	//log.Printf("ProgDir %+v\n", ProgDir)

	//log.Printf("http.Dir(static) %+v\n", http.Dir(ProgDir+"/static/"))
	static_file = req.URL.Path[len("/static/"):]

	if len(static_file) != 0 {

		if static_file == "js/tpl.js"{
			//log.Printf("static_file %+v\n", static_file)
			jsFile := renderJS("tplsJs", nil)
			fileBytes := bytes.NewReader(jsFile)
			content := io.ReadSeeker(fileBytes)

			http.ServeContent(w, req, static_file, time.Now(), content)
			return
		}


		f, err := http.Dir(ProgDir+"/static/").Open(static_file)

		if err == nil {
			content := io.ReadSeeker(f)
			http.ServeContent(w, req, static_file, time.Now(), content)
			return
		} else {
			log.Printf("%+v\n", err)
		}
	}
	http.NotFound(w, req)
}



//func renderContent (tplName string, ContentData interface{}) *template.Template {
//	//ContentHtml := `{{define "content"}}`
//	//content := render(tplName, ContentData)
//	//ContentHtml += string(content)
//	//ContentHtml += `{{end}}`
//	pattern := map[string]string{
//		"name":    "tplFiles",
//		"pattern": ProgDir+"/html/*.*",
//	}
//
//
//	return ContentHtml
//}

