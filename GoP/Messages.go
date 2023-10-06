package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"github.com/gorilla/websocket"
	"golang.org/x/net/context"
	"html/template"
	"log"
	"strconv"
	"strings"
	"time"
)

type Сообщение struct {
	Token struct {
		Hash string
		Истекает string
	} `json:"Token"`
	Id int `json:"Id"`
	Ip string `json:"ip"`
	От string `json:"От"`
	Кому string `json:"Кому"`
	Текст string `json:"Текст"`
	MessageType []string `json:"MessageType"`
	Время string `json:"Время"`
	ОтветНа string `json:"ОтветНа"`
	Файлы []string `json:"Файлы"`
	Offline string `json:"Offline"`
	Online string `json:"Online"`
	AdminMenu []map[string]interface{} `json:"AdminMenu"`// []BotMenuStruct
	ВходящиеАргументы map[string]interface{} `json:"ВходящиеАргументы"`
	Выполнить struct{
		Action string `json:"action"`
		Действие map[string]map[string]interface{} `json:"Действие"`   // "НазваниеДействия" :{"имяАргумента_1":"значение аргумента_1"... }
		Skill int `json:"skill"`
		Навык json.Number `json:"Навык"`
		Cmd string `json:"cmd"`
		Комманду string `json:"Комманду"`
		Arg struct {
			Module string `json:"module"`
			Tables []string `json:"tables"`
			Login string `json:"login"`
			Other map[string]interface{} `json:"other"`
		} `json:"Arg"`
	} `json:"Выполнить"`
	Контэнт *ДанныеОтвета `json:"Контэнт"`
	Content struct {
		Target string `json:"target"`
		Data interface{} `json:"data"`
		Html string `json:"html"`
		Обработчик string `json:"обработчик"` //JS функция или объект/класс/плагин для обработки данных (table..)
	} `json:"Content"`
	UserInfo struct {
		Uid string `json:"uid"`
		Initials string `json:"Initials"`
		FullName string `json:"FullName"`
		FirstName string `json:"FirstName"`
		LastName string `json:"LastName"`
		MiddleName string `json:"MiddleName"`
		OspName string `json:"OspName"`
		OspNum int `json:"OspNum"`
		PostName string `json:"PostName"`
		Инициалы string `json:"Инициалы"`
		ПолноеИмя string `json:"ПолноеИмя"`
		Фамилия string `json:"Фамилия"`
		Имя string `json:"Имя"`
		Отчество string `json:"Отчество"`
		ОСП string `json:"ОСП"`
		КодОСП int `json:"КодОСП"`
		Должность string `json:"Должность"`
	} `json:"UserInfo"`
}

type Message struct {
	From string `json:"from"`
	To string `json:"to"`
	Text string `json:"text"`
	Time string `json:"time"`
	ReaplyTo string `json:"reaply_to"`
	Files []string `json:"files"`
}

type messageRowSql struct {
	Id int `json:"id"`
	Autor string `json:"autor"`
	Recipient string `json:"recipient"`
	ChatRoom sql.NullString `json:"chat_room"`
	ReaplyTo sql.NullString `json:"reaply_to"`
	Date string `json:"date"`
	Files sql.NullString `json:"files"`
	Text sql.NullString `json:"text"`
	TextHtml template.HTML `json:"text_html"`
	AutorName string `json:"autor_name"`
	AutorMiddlename string `json:"autor_middlename"`
	RecipientName string `json:"recipient_name"`
	RecipientMiddlename string `json:"recipient_middlename"`
}

type messageRow struct {
	Id int `json:"id, mes_order"`
	Autor string `json:"autor"`
	Recipient string `json:"recipient"`
	ChatRoom string `json:"chat_room"`
	ReaplyTo string `json:"reaply_to"`
	Date string `json:"mes_date"`
	Files []string `json:"files"`
	Text string `json:"text"`
	TextHtml template.HTML `json:"text_html"`
	AutorName string `json:"autor_name"`
	AutorMiddlename string `json:"autor_middlename"`
	RecipientName string `json:"recipient_name"`
	RecipientMiddlename string `json:"recipient_middlename"`
}

func (mes Сообщение)ИзменитьСообщение() int {

	return 0
}

func (mes *Сообщение)СохранитьЛогСообщения() {
	//log.Printf("СохранитьСообщение mes %+v\n", mes)

	columns :=""
	countColumns:=3
	//if mes.Время == ""{
	ТекущеееВремя := time.Now()
	mes.Время = ТекущеееВремя.Format("2006-01-02T15:04:05.999999")
	//}
	sqlArgStr := []string{
		mes.От,
		mes.Кому,
		mes.Время,
	}
	sqlArgs:=[][]byte{
		[]byte(mes.От),
		[]byte(mes.Кому),
		[]byte(mes.Время),
	}
	//log.Printf("\n mes.Время %+v\n",mes.Время )
	if mes.ОтветНа != ""{
		columns = columns+", reaply_to"
		countColumns++
		sqlArgs = append(sqlArgs, []byte(mes.ОтветНа))
	}

	if mes.Файлы != nil{
		columns = columns+", files"
		countColumns++

		FilesString, err := json.Marshal(mes.Файлы)
		if err != nil {
			log.Printf("err	 %+v\n", err)
		}
		sqlArgs = append(sqlArgs, FilesString)
	}

	columns = columns+", text"
	countColumns++
	if mes.Текст != ""{
		sqlArgs = append(sqlArgs, []byte(mes.Текст))
		sqlArgStr = append(sqlArgStr, mes.Текст)
	} else {
		sqlArgs = append(sqlArgs, []byte(nil))
		sqlArgStr = append(sqlArgStr, mes.Текст)
	}

	if mes.Выполнить.Action!="" ||  mes.Выполнить.Cmd !="" || mes.Выполнить.Skill != 0 || Contains( mes.MessageType, "io_action"){
		columns = columns+", type"
		countColumns++
		sqlArgs = append(sqlArgs, []byte(`["io_action"]`))
		sqlArgStr=append(sqlArgStr, "io_action")

		//
		КомандаБоту := map[string]string{}
		if mes.Выполнить.Action!="" {
			КомандаБоту["Action"] = mes.Выполнить.Action
		}
		if mes.Выполнить.Cmd !="" {
			КомандаБоту["Cmd"] = mes.Выполнить.Cmd
		}
		if mes.Выполнить.Skill != 0  {
			КомандаБоту["Skill"] = strconv.Itoa(mes.Выполнить.Skill)
		}
		byteString, err := json.Marshal(КомандаБоту)
		if err != nil {
			log.Printf("err	 %+v\n", err)
		}
		if len(КомандаБоту) >0{
			columns = columns+", comand_to_io"
			countColumns++
			sqlArgStr=append(sqlArgStr, string(byteString))
			sqlArgs = append(sqlArgs, byteString)
		}

	}

	valuesPlaceholder:=""
	if countColumns>3 {
		for i := 4; i <= countColumns; i++ {
			valuesPlaceholder = valuesPlaceholder + ", $" + strconv.Itoa(i)
		}
	}

	//log.Printf("columns %+v valuesPlaceholder %+v\n",columns,  valuesPlaceholder)

	sqlString := `INSERT INTO messages (autor, recipient, mes_date `+columns+`) VALUES ($1,$2, $3 `+ valuesPlaceholder+`) RETURNING message_id`
	// AND date >= CURRENT_DATE - INTERVAL '1 day'


	sqlQuery := sqlStruct{
		Name:   "messages",
		Sql:    sqlString,
		Values: sqlArgs,
		DBSchema:"iobot",
	}
	//log.Printf("\n\nsqlArgStr >> %+v\n \n", sqlArgStr)

	//messagesArrayMap, _ := ВыполнитьPgSQL(sqlQuery)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	Result,err := sqlQuery.PgSqlResultReader(ctx)
	//log.Printf("\nResultReader %+v\n", ResultReader)

	if err != nil {
		log.Printf("\n\n !!! Ошибка сохранения сообщения %+v\n",err)
		return
	}
	rows := Result.Rows

	var id int
	for _,row := range rows {
		for _ ,cell :=range row {
			id, _ = strconv.Atoi(string(cell))
		}
	}

	mes.Id=id
}

func СохранитьСообщение(mes Сообщение) (int, string) {
	//log.Printf("СохранитьСообщение mes %+v\n", mes)

	columns :=""
	countColumns:=3
	//if mes.Время == ""{
		ТекущеееВремя := time.Now()
		mes.Время = ТекущеееВремя.Format("2006-01-02T15:04:05.999999")
	//}
	sqlArgStr := []string{
		mes.От,
		mes.Кому,
		mes.Время,
	}
	sqlArgs:=[][]byte{
		[]byte(mes.От),
		[]byte(mes.Кому),
		[]byte(mes.Время),
	}
//log.Printf("\n mes.Время %+v\n",mes.Время )
	if mes.ОтветНа != ""{
		columns = columns+", reaply_to"
		countColumns++
		sqlArgs = append(sqlArgs, []byte(mes.ОтветНа))
	}

	if mes.Файлы != nil{
		columns = columns+", files"
		countColumns++

		FilesString, err := json.Marshal(mes.Файлы)
		if err != nil {
			log.Printf("err	 %+v\n", err)
		}
		sqlArgs = append(sqlArgs, FilesString)
	}

			columns = columns+", text"
			countColumns++
		if mes.Текст != ""{
			sqlArgs = append(sqlArgs, []byte(mes.Текст))
			sqlArgStr = append(sqlArgStr, mes.Текст)
		} else {
			sqlArgs = append(sqlArgs, []byte(nil))
			sqlArgStr = append(sqlArgStr, mes.Текст)
		}

		if mes.Выполнить.Action!="" ||  mes.Выполнить.Cmd !="" || mes.Выполнить.Skill != 0 || Contains( mes.MessageType, "io_action"){
			columns = columns+", type"
			countColumns++
			sqlArgs = append(sqlArgs, []byte(`["io_action"]`))
			sqlArgStr=append(sqlArgStr, "io_action")

			//
			КомандаБоту := map[string]string{}
			if mes.Выполнить.Action!="" {
				КомандаБоту["Action"] = mes.Выполнить.Action
			}
			if mes.Выполнить.Cmd !="" {
				КомандаБоту["Cmd"] = mes.Выполнить.Cmd
			}
			if mes.Выполнить.Skill != 0  {
				КомандаБоту["Skill"] = strconv.Itoa(mes.Выполнить.Skill)
			}
			byteString, err := json.Marshal(КомандаБоту)
			if err != nil {
				log.Printf("err	 %+v\n", err)
			}
			if len(КомандаБоту) >0{
				columns = columns+", comand_to_io"
				countColumns++
				sqlArgStr=append(sqlArgStr, string(byteString))
				sqlArgs = append(sqlArgs, byteString)
			}

		}

	valuesPlaceholder:=""
	if countColumns>3 {
		for i := 4; i <= countColumns; i++ {
			valuesPlaceholder = valuesPlaceholder + ", $" + strconv.Itoa(i)
		}
	}

	//log.Printf("columns %+v valuesPlaceholder %+v\n",columns,  valuesPlaceholder)

	sqlString := `INSERT INTO messages (autor, recipient, mes_date `+columns+`) VALUES ($1,$2, $3 `+ valuesPlaceholder+`) RETURNING message_id`
	// AND date >= CURRENT_DATE - INTERVAL '1 day'


	sqlQuery := sqlStruct{
		Name:   "messages",
		Sql:    sqlString,
		Values: sqlArgs,
		DBSchema:"iobot",
	}
	//log.Printf("\n\nsqlArgStr >> %+v\n \n", sqlArgStr)

	//messagesArrayMap, _ := ВыполнитьPgSQL(sqlQuery)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	Result,err := sqlQuery.PgSqlResultReader(ctx)
	//log.Printf("\nResultReader %+v\n", ResultReader)

	if err != nil {
		log.Printf("\n\n !!! Ошибка сохранения сообщения %+v\n",err)
		return 0, ""
	}
	rows := Result.Rows

	var result int
	for _,row := range rows {
		for _ ,cell :=range row {
			result, _ = strconv.Atoi(string(cell))
		}
	}

	return result,mes.Время
}

func (client * Client)ПолучитьЛогТерминала(mes Сообщение){

}

func (client * Client)ПолучитьЛогПереписки(mes Сообщение){
	//@language=PostgresSQL
//	sqlString := `select t.id, row_to_json(t.*) from (
//SELECT EXTRACT(EPOCH FROM mes_date) as mes_order , mes_date, iobot.messages.*, user_autor.givenname autor_name, user_autor.initials autor_middlename,
//      user_recipient.givenname recipient_name,user_recipient.initials recipient_middlename
//FROM iobot.messages
//        LEFT JOIN fssp_configs.users user_autor on user_autor.Login = autor
//        LEFT JOIN fssp_configs.users user_recipient on user_recipient.Login =recipient
//WHERE (type != 'io_action' OR type is null) AND ((autor = $1 AND recipient = $2) OR (recipient = $1 AND autor = $2)) ORDER BY id DESC LIMIT 30) t group by t.id, t.*`

sqlString := `select t.message_id, row_to_json(t.*) from (
SELECT EXTRACT(EPOCH FROM mes_date) as mes_order , mes_date, iobot.messages.*, user_autor.givenname autor_name, user_autor.initials autor_middlename,
      user_recipient.givenname recipient_name,user_recipient.initials recipient_middlename
FROM iobot.messages
        LEFT JOIN fssp_configs.users user_autor on user_autor.Login = autor
        LEFT JOIN fssp_configs.users user_recipient on user_recipient.Login =recipient
WHERE ((text IS NOT NULL AND text<>'') OR  files IS NOT NULL) AND ((autor = $1 AND recipient = $2) OR (recipient = $1 AND autor = $2)) ORDER BY message_id DESC LIMIT 30) t group by t.message_id, t.*`

//sqlString := `select 'ChatLog', jsonb_agg(log) as chat_log from (
//select t.id::varchar, row_to_json(t) from (
//SELECT EXTRACT(EPOCH FROM date) as mes_order ,  iobot.messages.*, user_autor.givenname autor_name, user_autor.initials autor_middlename,
//       user_recipient.givenname recipient_name,user_recipient.initials recipient_middlename
//FROM iobot.messages
//         LEFT JOIN fssp_configs.users user_autor on user_autor.Login = autor
//         LEFT JOIN fssp_configs.users user_recipient on user_recipient.Login =recipient
//WHERE (autor = 'maksimchuk@r26' AND recipient = 'kukushkin@r26') OR (recipient =  'kukushkin@r26' AND autor = 'maksimchuk@r26')  ORDER BY mes_order DESC LIMIT 30 ) t group by  t.id, t.*) log
//UNION
//select 'UserInfo', jsonb_agg(t2.*)  from
//    (SELECT givenname, initials, post_name FROM fssp_configs.users JOIN fssp_configs.posts ON fssp_configs.posts.id= fssp_configs.users.post WHERE fssp_configs.users.login = 'maksimchuk@r26') t2`

	// пара автор, получатель
	sqlQuery := sqlStruct{
		Name:   "messages",
		Sql:    sqlString,
		Values: [][]byte{
			[]byte(client.Login), []byte(mes.Выполнить.Arg.Login),
		},
		DBSchema:"iobot",
	}

	//messagesArrayMap, _ := ВыполнитьPgSQL(sqlQuery)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	Result, err := sqlQuery.PgSqlResultReader(ctx)
	if err != nil {
		log.Printf("\n !! ERR %+v\n", err)
	}
	MessagesLog := map[int]messageRow{}

	messagesRows:= Result.Rows

	for _, messageByte := range messagesRows{
		//log.Printf("string(messageByte[0] id%+v\n", id)
		//log.Printf("string(messageByte[0] %+v\n", string(messageByte[0]))
		//log.Printf("string(messageByte[1] %+v\n", string(messageByte[1]))
		messageStruct := messageRow{}

		//message.TextHtml = template.HTML(message.Text.String)
			err := json.Unmarshal(messageByte[1], &messageStruct)

			if err != nil {
				log.Printf("Unmarshal messageByte[1] err  %+v\n", err, string(messageByte[1]))
			}

			mesId,err:=strconv.Atoi(string(messageByte[0]))
			if err != nil {
				log.Printf("err	 %+v\n", err)
			}

			//text := strings.Replace(messageStruct.Text, "\n", "<br>", -1)
			messageStruct.TextHtml=template.HTML(messageStruct.Text)
			MessagesLog[mesId]=messageStruct
	}
	//UserInfo := &UsersStruct{}
	//if mes.Выполнить.Arg.Login != "io"{
		UserInfo := client.ПолучитьДанныеПользователя(mes.Выполнить.Arg.Login)
	//}

	data:=map[string]interface{}{
		"client":client,
		"uid":client.Login,
		"log":MessagesLog,
		"UserInfo": UserInfo,
		"BotMenu":client.ПолучитьМенюБота(),
		"Dialogs":client.ПолучитьБыстрыеДиалоги(),
	}

	var Data map[string]string
	if server.Clients[mes.Выполнить.Arg.Login] != nil{
		Data = map[string]string{"ClientIp":server.Clients[mes.Выполнить.Arg.Login].Ip}
	}

	responseMes := Сообщение{
		Id:     0,
		Ip: 	client.Ip,
		От:     "io",
		Кому:   client.Login,
		MessageType:[]string{"log"},
		Content: struct {
			Target string `json:"target"`
			Data interface{} `json:"data"`
			Html   string `json:"html"`
			Обработчик string `json:"обработчик"`
		}{
			Target: "log_wrapper_"+mes.Выполнить.Arg.Login,
			Html:   string(render("messageLog",data)),
			Data: Data,
		},
	}

	client.Message<-&responseMes

log.Printf("\n\n !!! >>>>. ПолучитьСообщениеИО mes.Выполнить.Arg.Login %+v\n", mes.Выполнить.Arg.Login)
	if mes.Выполнить.Arg.Login == "io"{
		ТекстСообщения := ПолучитьСообщениеИО("приветствие", client)
		СообщениеКлиенту:= &Сообщение{
					Текст:   ТекстСообщения,
					От: "io",
					Кому:client.Login,
				}
				//СообщениеКлиенту.СохранитьЛогСообщения()
				client.Message<-СообщениеКлиенту
	}

	//if client.UserInfo.Info.OspCode == 26911 {

	//}

}

// алгоритм Получить уровень доступа left join fssp_configs.уровень_доступа ON fssp_configs.уровень_доступа.уровень ?| (
//    select array_agg(level.уровень) from (
//           select jsonb_array_elements_text(уровень) уровень  from fssp_configs.уровень_доступа where логин = 'kukushkin@r26' OR отдел = 26911 OR должность = 44 OR (отдел IS NULL AND логин IS NULL)
//                                         )level
//    )
//SELECT distinct iobot.диалоги_ио.* FROM iobot.диалоги_ио
//join fssp_configs.уровень_доступа ON iobot.диалоги_ио.доступ ?| (
//select array_agg(level.уровень) from (
//select jsonb_array_elements_text(уровень) уровень  from fssp_configs.уровень_доступа where логин = 'kukushkin@r26' OR отдел = 26911 OR должность = 44 OR (отдел IS NULL AND логин IS NULL)
//)level
//)
func ПолучитьУровеньДоутспа (client *Client){
	sql :=`SELECT * FROM iobot.диалоги_ио
left join fssp_configs.уровень_доступа ON fssp_configs.уровень_доступа.уровень ?| (
    select array_agg(level.уровень) from (
           select jsonb_array_elements_text(уровень) уровень  from fssp_configs.уровень_доступа where логин = $1 OR отдел = $2 OR должность = $3 OR (отдел IS NULL AND логин IS NULL)
                                         )level
    )`
	_ ,err:= sqlStruct{
			Name:   "уровень_доступа",
			Sql:    sql,
			Values: [][]byte{
				[]byte(client.Login),

			},
		}.RunSQL(nil)
	if err != nil{
	log.Printf(">>>> Ошибка SQL запроса: %+v \n\n",err)
	}
}

/* алгоритм
// ПолучитьСообщениеИО возвращает сообщение с заданым именем, и в соответствии с уровнем доступа ползователя по отделу и должности
*/
func ПолучитьСообщениеИО(ИмяСообщения string, client *Client) string {


// sql := `SELECT distinct iobot.диалоги_ио.* FROM iobot.диалоги_ио
//                  JOIN fssp_configs.уровень_доступа ON iobot.диалоги_ио.доступ ?| (
//    	SELECT array_agg(level.уровень) from (
//        SELECT jsonb_array_elements_text(уровень) уровень  from fssp_configs.уровень_доступа where логин = $1 OR отдел = $2 OR должность = $3
//        )level
//) WHERE имя_сообщения= $4`

sql := `select distinct iobot.диалоги_ио.* from
    iobot.диалоги_ио
 join (SELECT
       case when
           jsonb_array_length(mainLevel.уровень)>0
           then
                mainLevel.уровень
       else
           childLevel.уровень
       end as уровень
FROM
    (SELECT  jsonb_agg(level.уровень) уровень from (
       SELECT  jsonb_array_elements_text(уровень) уровень  from fssp_configs.уровень_доступа where (логин = $1 OR отдел = $2 OR должность = $3)
    )level) mainLevel,
     (SELECT jsonb_agg(level.уровень) уровень from (
       SELECT  jsonb_array_elements_text(уровень) уровень  from fssp_configs.уровень_доступа where (логин is null AND отдел is null AND должность is null)
    )level) childLevel
) уровеньДоступа ON уровень <@ доступ
WHERE имя_сообщения= $4`

	//"SELECT *FROM iobot.диалоги_ио WHERE имя_сообщения=$1 "
	OSPCode := strconv.Itoa(client.UserInfo.Info.OspCode)
	Post := strconv.Itoa(client.UserInfo.Info.Post)

	log.Printf("OSPCode %+v Post %+v client.Login %+v client.Login %+v\n", OSPCode, Post, client.Login, ИмяСообщения)
	Результат ,err:= sqlStruct{
			Name:   "io_dialogs",
			Sql:    sql,
			Values: [][]byte{
				[]byte(client.Login),
				[]byte(OSPCode),
				[]byte(Post),
				[]byte(ИмяСообщения),
			},
		}.RunSQL(nil)
	if err != nil{
	log.Printf(">>>> Ошибка SQL запроса: %+v \n\n",err)
	}
	//log.Printf(">>>> ПолучитьСообщениеИО  %+v  client %+v \n Результат %+v\n\n ",ИмяСообщения, client.UserInfo.Info, Результат)
	if len(Результат)>0{
		ОтветИО, ЕстьОТветИО := Результат[0]["ответ_ио"]
		if ЕстьОТветИО {
			Утверждение, ЕстьУтверждение:=ОтветИО.(map[string]interface{})["утверждение"]
			if ЕстьУтверждение {
				tpl := template.Must(template.New("Сообщение").Parse(Утверждение.(string)))

				resHtml := new(bytes.Buffer)

				//err = tplFiles.ExecuteTemplate(resHtml, tplName, data)
				err := tpl.Execute(resHtml, client.UserInfo)
				if err != nil {
					log.Println("executing template:", err)
					return ""
				} else {
					return resHtml.String()
				}
			}
			Вопрос, ЕстьВопрос:=ОтветИО.(map[string]interface{})["вопрос"]
			if ЕстьВопрос {
				tpl := template.Must(template.New("Сообщение").Parse(Вопрос.(string)))

				resHtml := new(bytes.Buffer)

				//err = tplFiles.ExecuteTemplate(resHtml, tplName, data)
				err := tpl.Execute(resHtml, client.UserInfo)
				if err != nil {
					log.Println("executing template:", err)
					return ""
				} else {
					return resHtml.String()
				}
			}
		}

	} else {
		return ""
	}
	return ""
}





func (client *Client)ReadMessage(){
	for {
		_, message, err := client.Ws.ReadMessage()
		if err != nil {
			log.Printf("err %+v\n", err)
			//log.Printf("read:err %v, CloseGoingAway %v", err, websocket.CloseGoingAway)

			//log.Printf("websocket.IsCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure): %v",
			//	websocket.IsCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure))

			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				//log.Printf("error: %v", err)
			}
			if websocket.IsCloseError(err, websocket.CloseGoingAway){
				//log.Printf("error: %v", err)
			}
			break
		}

		mes := Сообщение{}

		err = json.Unmarshal(message, &mes)
		if err!=nil{
			log.Printf("err %+v message %+v\n", err, string(message))
		}

		if mes.От == ""{
			mes.От=client.Login
		}
		mes.Id, mes.Время = СохранитьСообщение(mes)
		//СохранитьСообщение(mes)
		if mes.Кому != "io"{
			if server.Clients[mes.Кому] != nil{
				server.Clients[mes.Кому].Message<-&mes
			}
		} else {
			go client.IOHandler(mes)
		}

		//client.message<-message
	}
}

func (client *Client) Write(b []byte)(int, error){
	go func(){
		responseMes := Сообщение{
			Id:      0,
			От:      "io",
			Кому:    client.Login,
			Текст: string(b),
			MessageType:[]string{"server_log"},
		}

// проверим открыт ли канал, если открыт то отправим в него данные иначе удалим получателя из списка отправка ЛОГОВ
//		if канал, ok := (<-client.message); ok {
//			log.Printf("канал %+v\n", канал)
			client.Message<-&responseMes
//		} else {
//			log.Printf("Канал закрыт канал %+v %+v\n",канал, ok)
			//mw := io.MultiWriter(os.Stdout,IoLoger) //, server.Clients["maksimchuk@r26"]
			//log.SetOutput(mw)
		//}
	}()

	return len(b), nil
}





func (client *Client)SendMessage(){
	for {
		select {
			case q := <-client.Message:
				//log.Printf("q %+v,  %+v\n", q)
				response, err := json.Marshal(q)
				if err!=nil{
					log.Printf("SendMessage response err Marshal: %+v\n", err)
				}
				//clinetR:=server.Clients[q["To"]]
				//log.Printf("Проверим активность канала %+v\n",client )
				//err =  client.ws.WriteMessage(1, response)
				//канал, ok  <-client.message
				//log.Printf("канал %+v, ok %+v\n", канал, ok)
				//if  ok {

					err =  client.Ws.WriteMessage(1, response)

					if err!=nil{
						log.Printf("SendMessage client.ws (%+v)\n", client)
						log.Printf("SendMessage client.message response (%+v)\n", response)
						log.Printf("SendMessage WriteMessage err (%+v) response(%+v) client(%+v) \n", err, string(response), client)
						return
					}
				//} else {
				//	log.Printf("Канал пользователя  %+v канал %+v ok %+v закрыт %+v\n", client.Login, ok, канал,client.message)
				//	return
				//}

		}
	}
}


/*
ОбработатьОщибку, получает на вход sql строку, и данные для подстановки в запрос, Функция сама отправляет сообщение клиенту если обработчик умпешно отработал и вернул данные с полем ,message*/
func (client *Client) ОбработатьОшибкуБД (ОбработчикОшибки interface{}, data interface{}) {

	SQLStateСкрипт := ОбработчикОшибки.(map[string]interface{})["sql"]

	SQLStateString := template.Must(template.New("SQLStateСкриптЗапрос").Parse(SQLStateСкрипт.(string)))

	БайтБуферSQLState := new(bytes.Buffer)

	err := SQLStateString.Execute(БайтБуферSQLState, data)

	if err != nil {
		log.Println("executing template:", err)
	} else {
		//SQLСкрипт = БайтБуферSQLState.String()
	}

	Результат ,err:= sqlStruct{
		Name:   "",
		Sql:     БайтБуферSQLState.String(),
		Values: [][]byte{},
	}.RunSQL(nil)

	if err != nil {
		log.Printf(">>>> ERROR \n %+v \n\n", err)
	}

	// алгоритм т.к. это ошибка выполнения запроса, и нам нужно сообщить об этлм клиенту, то выведем сразу сообщение со стандартным шаблном ошибки.
	if len(Результат) ==1 {
		ИнформацияКлиенту := Результат[0]["message"]
		//type ДанныеОтвета struct {
		//	Контейнер string `json:"контейнер"`
		//	Данные interface{} `json:"данные"`
		//	HTML string `json:"html"`
		//	Обработчик string `json:"обработчик"` //JS функция или объект/класс/плагин для обработки данных (table..)
		//}
		ДанныеДляОтвета := &ДанныеОтвета{
			Обработчик:"FloatMessage",
			Данные:ИнформацияКлиенту,
		}
		СообщениеКлиенту:= &Сообщение{
			От: "io",
			Кому:client.Login,
			MessageType: []string{"error"},
			Контэнт:ДанныеДляОтвета,
		}
		log.Printf("СообщениеКлиенту %+v\n", СообщениеКлиенту)
		СообщениеКлиенту.СохранитьИОтправить(client)
	}

}
/*
ОбработатьSQLСкрипты на вход получает объект со скриптами, выполняет каждый из них, сохраняет в карту, и возвращает карту с данными всех запросов
Если в резулттате какогто запроса возникла ошибка, и для этой ошибки есть обработчик, то обработчик ыполниться, и сообщит клиенту об ошибке
*/




func (client *Client)ОбработатьSQLСкрипты(SQLЗапрос interface{}, вопрос *Сообщение) map[string]interface{} {


	log.Printf("Выполним SQLСкрипт %+v\n", SQLЗапрос)

	SQLСкрипт, ЕстьSqlСкрипт := SQLЗапрос.(map[string]interface{})["sql"].(map[string]interface{})

	Данные:= map[string]interface{}{
		"client" : client.UserInfo.Info,
		"data": map[string]interface{}{},
	}

	if вопрос.Выполнить.Действие !=nil{
		for НазваниеДействия, ДанныеДействия := range вопрос.Выполнить.Действие{
			log.Printf("НазваниеДействия %+v ДанныеДействия %+v\n", НазваниеДействия, ДанныеДействия)
			Данные["data"]=ДанныеДействия
		}
	}

	if !ЕстьSqlСкрипт{
		log.Printf("ЕстьSqlСкрипт нет срикптов%+v\n", ЕстьSqlСкрипт)
		return nil
	} else {
		РезультатВсехЗапросов := map[string]interface{}{}
		for ИмяСкрипта, Скрипт := range SQLСкрипт{

			Инфо(" ИмяСкрипта %+v Скрипт %+v",ИмяСкрипта, Скрипт)

			tpl := template.Must(template.New("SqlЗапрос").Funcs(tplFunc()).Parse(Скрипт.(string)))
			БайтБуферSql := new(bytes.Buffer)
			err := tpl.Execute(БайтБуферSql, Данные)

			if err != nil {
				log.Println("executing template:", err)
			} else {
				//SQLСкрипт[ИмяСкрипта] = БайтБуферSql.String()
				РезультатЗапроса ,ОшибкаЗапроса:= sqlStruct{
					Name:   "",
					Sql:    БайтБуферSql.String(),
					Values: [][]byte{},
				}.RunSQL(nil)

				if ОшибкаЗапроса != nil {
					if strings.Contains(ОшибкаЗапроса.Error(), "SQLSTATE 23505"){
						//алгоритм. Если запрос вернул ошибку с ограничением целостности уникального значения , тогда проверим есть ли обработчик SQLSTATE. и Наверное пока sql для обрабтки ошибок будет возвращать сразу текст который будет выводиться клиенту,
						SQLStateОбъект, ЕстьSQLStateОбъект := SQLЗапрос.(map[string]interface{})["SQLState"]

						if !ЕстьSQLStateОбъект{
							log.Printf("ЕстьSQLStateОбъект %+v\n", ЕстьSQLStateОбъект)
						}

						ОбработчикОшибки, ЕстьОбработчикОшибки := SQLStateОбъект.(map[string]interface{})["23505"]

						if ЕстьОбработчикОшибки{
							/*алгоритм Функция ОбработатьОщибку, получает на вход sql строку, и данные для подстановки в запрос, Функция сама отправляет сообщение клиенту если обработчик умпешно отработал и вернул данные с полем ,message*/
							client.ОбработатьОшибкуБД (ОбработчикОшибки , Данные)
						}
					}
					Ошибка("ОшибкаЗАпроса %+v",  ОшибкаЗапроса)
				} else{
					if len(РезультатЗапроса) == 1{
						if сообщениеКлиенту, ЕстьСообщение :=  РезультатЗапроса[0]["message"];ЕстьСообщение {
							ДанныеДляОтвета := &ДанныеОтвета{
								Обработчик:"FloatMessage",
								Данные:сообщениеКлиенту,
							}
							СообщениеКлиенту:= &Сообщение{
								От: "io",
								Кому:client.Login,
								MessageType: []string{"note"},
								Контэнт:ДанныеДляОтвета,
							}
							log.Printf("СообщениеКлиенту %+v\n", СообщениеКлиенту)
							СообщениеКлиенту.СохранитьИОтправить(client)
						}

					}
					РезультатВсехЗапросов[ИмяСкрипта] = РезультатЗапроса
				}
			}
		}
		//log.Printf("РезультатВсехЗапросов %+v\n", РезультатВсехЗапросов)
		return РезультатВсехЗапросов
	}
}

//  ОтветитьКлиенту - обрабатывает входящий запрос от клиента в соответсвии со сценаием из диалоги_ио
func (client *Client) ОтветитьКлиенту (Диалог map[string]interface{}, вопрос *Сообщение){
	_, ОтветЕсть := Диалог["ответ_ио"]
	log.Printf("ответ_ио %+v\n", Диалог["ответ_ио"])
	log.Printf("ОтветЕсть %+v\n", ОтветЕсть)
	if ОтветЕсть {
		if Диалог["ответ_ио"] != nil {
			ОтветИО := Диалог["ответ_ио"].(map[string]interface{})

			log.Printf("ОтветИО %+v\n", ОтветИО)
			/*алгоритм  наличие утверждения в ответе исключает вопрос и  ожидание ответа от клиента
						Если есть утвверждение то сразу обращаемся к полю далее, и проверяем наличие "ДействиеПередОтветом"
					   Если есть "ДействиеПередОтветом" то выполняем его, и если ДействиеПередОтветом вернуло данные с полем "ОтветКлиенту" то добавляем ответ в текст сообщения, инчале если ДействиеПередОтветом вернуло данные с полями HTML и таргет то помщяем эти данные в соответсвтующие поля в сообщение клиенту
			  алгоритм
						цель: куда вставлять хтмл с результатами,
						всегда должен быть полный путь IDs блоков через точку (НЕ КЛАССОВ А ID)
					   например main_content.tickets_list   (если на странице есть ID то данные в этом блоке обновяться)
					   обработка блоков для вставки идёт с права на лево, если есть tickets_list обновляем данные внутри,
						если нет tickets_list  , то ищем main_content и вставляем данные в него , если нет main_content то ничего не делаем
					   имя хтмл шаблона который будет возвращён клиенту в виде ответа с результатом sql запроса,
					   если нет sql запроса то возращаеться  просто хтмл
			*/

			утверждение, ЕстьУтверждение := ОтветИО["утверждение"].(string)
			log.Printf("ЕстьУтверждение %+v\n", ЕстьУтверждение)
			if ЕстьУтверждение {

				номер_диалога:= strconv.Itoa(Диалог["номер_диалога"].(int))
				номер_сообщения:= strconv.Itoa(Диалог["номер_сообщения"].(int))

				_, err := sqlStruct{
					Name: "история_диалогов_ио",
					Sql:  "INSERT INTO iobot.история_диалогов_ио (клиент, номер_диалога, номер_сообщения, время_сообщения, завершено, сообщение_клиента) VALUES ($1,$2,$3,NOW(),true, $4)",
					Values: [][]byte{
						[]byte(client.Login),
						[]byte(номер_диалога),
						[]byte(номер_сообщения),
						[]byte(вопрос.Текст),
					},
				}.RunSQL(nil)

				if err != nil {
					log.Printf(">>>> ERROR \n %+v \n\n", err)
				}
				tpl := template.Must(template.New("Сообщение").Parse(утверждение))
				resHtml := new(bytes.Buffer)
				err = tpl.Execute(resHtml, client.UserInfo.Info)
				if err != nil {
					log.Printf("executing template:", err)
				} else {
					утверждение = resHtml.String()
				}

				СообщениеКлиенту := &Сообщение{
					Текст: утверждение,
					От:    "io",
					Кому:  client.Login,
				}

				log.Printf("ОтветИО %+v\n", ОтветИО)

				if ДействиеДляОТвета, ЕстьДействиеДляОТвета := ОтветИО["ДействиеДляОтвета"]; ЕстьДействиеДляОТвета {
					log.Printf("ОтветИО[ДействиеДляОтвета] %+v\n", ОтветИО["ДействиеДляОтвета"])
					client.ВыполнитьДействиеДляОтвета(ДействиеДляОТвета.(string), вопрос, СообщениеКлиенту)
				}

				СообщениеКлиенту.СохранитьИОтправить(client)
				//client.Message<-СообщениеКлиенту

			}
			вопросКлиенту, ЕстьВопрос := ОтветИО["вопрос"]
			if ЕстьВопрос {
				log.Printf("ЕстьВопрос %+v\n", ЕстьВопрос)
				ожидает, ОжидаетОтвета := ОтветИО["ожидает"]
				if ОжидаетОтвета {

					номер_диалога := strconv.Itoa(Диалог["номер_диалога"].(int))
					номер_сообщения := strconv.Itoa(Диалог["номер_сообщения"].(int))

					_, _ = sqlStruct{
						Name: "история_диалогов_ио",
						Sql:  "INSERT INTO iobot.история_диалогов_ио (клиент, номер_диалога, номер_сообщения, время_сообщения, ожидает, завершено, сообщение_клиента) VALUES ($1,$2,$3,NOW(),$4,false, $5)",
						Values: [][]byte{
							[]byte(client.Login),
							[]byte(номер_диалога),
							[]byte(номер_сообщения),
							[]byte(ожидает.(string)),
							[]byte(вопрос.Текст),
						},
					}.RunSQL(nil)
				}
				ВариантыОжидаемыхОтветов := []string{}
				if ожидает.(string) == "вариант_ответа" {

					Далее, ЕстьСледующийШаг := Диалог["далее"]
					if ЕстьСледующийШаг && Далее != nil {
						for _, Вариант := range Далее.([]interface{}){
							if ОжидаемыйОтвет, ЕстьОжидаемыйВариант := Вариант.(map[string]interface{})["ответ"];ЕстьОжидаемыйВариант{
								ВариантыОжидаемыхОтветов= append(ВариантыОжидаемыхОтветов, ОжидаемыйОтвет.(string))
							}
						}
					}
				}

				tpl := template.Must(template.New("Сообщение").Parse(вопросКлиенту.(string)))
				resHtml := new(bytes.Buffer)

				ДанныеДляГенерацииОтвета := map[string]interface{}{
					"client":client.UserInfo.Info,
					"ВариантыОжидаемыхОтветов":ВариантыОжидаемыхОтветов,
				}

				err := tpl.Execute(resHtml, ДанныеДляГенерацииОтвета)
				if err != nil {
					log.Println("executing template:", err)
				} else {
					вопросКлиенту = resHtml.String()
				}

				if len(ВариантыОжидаемыхОтветов)>0{
					вопросКлиенту = вопросКлиенту.(string)+string(render("variantsArray", ДанныеДляГенерацииОтвета))
				}

				СообщениеКлиенту := &Сообщение{
					Текст: вопросКлиенту.(string),
					От:    "io",
					Кому:  client.Login,
				}

				if ДействиеДляОТвета, ЕстьДействиеДляОТвета := ОтветИО["ДействиеДляОтвета"]; ЕстьДействиеДляОТвета {
					client.ВыполнитьДействиеДляОтвета(ДействиеДляОТвета.(string), вопрос, СообщениеКлиенту)
				}
				СообщениеКлиенту.СохранитьИОтправить(client)
				//client.Message<-СообщениеКлиенту
			}

			if !ЕстьУтверждение && !ЕстьВопрос {
				if ДействиеДляОТвета, ЕстьДействиеДляОТвета := ОтветИО["ДействиеДляОтвета"]; ЕстьДействиеДляОТвета {
					log.Printf("ОтветИО[ДействиеДляОтвета] %+v\n", ОтветИО["ДействиеДляОтвета"])
					СообщениеКлиенту := &Сообщение{
						Текст: "",
						От:    "io",
						Кому:  client.Login,
					}
					client.ВыполнитьДействиеДляОтвета(ДействиеДляОТвета.(string), вопрос, СообщениеКлиенту)
					СообщениеКлиенту.СохранитьИОтправить(client)
				}
			}
		}


		SQLЗапрос, ЕстьСкрипт := Диалог["sql_запрос"]
		ДанныеHTMLШаблона, ЕстьHTMLШаблон := Диалог["html_шаблон"]
		JSONДанные, ЕстьJSONДанные := Диалог["json_данные"]

		log.Printf("ЕстьСкрипт %+v SQLЗапрос %+v\n", ЕстьСкрипт, SQLЗапрос)
		log.Printf("ЕстьHTMLШаблон %+v Диалог %+v\n", ЕстьHTMLШаблон,  Диалог)



		if ЕстьСкрипт && SQLЗапрос != nil {

			ДанныеЗапросов := client.ОбработатьSQLСкрипты(SQLЗапрос, вопрос)

			//log.Printf("Выполним SQLСкрипт %+v\n", SQLЗапрос)
			//
			//SQLСкрипт, ЕстьSqlСкрипт := SQLЗапрос.(map[string]interface{})["sql"].(map[string]interface{})
			//data:= map[string]interface{}{
			//	"client" : client.UserInfo.Info,
			//	"data": map[string]interface{}{},
			//}
			//if вопрос.Выполнить.Действие !=nil{
			//	for НазваниеДействия, ДанныеДействия := range вопрос.Выполнить.Действие{
			//		log.Printf("НазваниеДействия %+v\n", НазваниеДействия)
			//		data["data"]=ДанныеДействия
			//	}
			//}
			//
			//if !ЕстьSqlСкрипт{
			//	log.Printf("ЕстьSqlСкрипт %+v\n", ЕстьSqlСкрипт)
			//} else {
			//	for ИмяСкрипта, Скрипт := range SQLСкрипт{
			//		tpl := template.Must(template.New("SqlЗапрос").Parse(Скрипт.(string)))
			//		БайтБуферSql := new(bytes.Buffer)
			//		err := tpl.Execute(БайтБуферSql, data)
			//
			//		if err != nil {
			//			log.Println("executing template:", err)
			//		} else {
			//			//SQLСкрипт[ИмяСкрипта] = БайтБуферSql.String()
			//			РезультатЗапроса ,ОшибкаЗАпроса:= sqlStruct{
			//				Name:   "",
			//				Sql:    БайтБуферSql.String(),
			//				Values: [][]byte{},
			//			}.RunSQL(nil)
			//
			//			if ОшибкаЗАпроса != nil {
			//
			//				if strings.Contains(ОшибкаЗАпроса.Error(), "SQLSTATE 23505"){
			//					//алгоритм. Если запрос вернул ошибку с ограничением целостности уникального значения , тогда проверим есть ли обработчик SQLSTATE. и Наверное пока sql для обрабтки ошибок будет возвращать сразу текст который будет выводиться клиенту
			//					SQLStateОбъект, ЕстьSQLStateОбъект := SQLЗапрос.(map[string]interface{})["SQLState"]
			//
			//					if !ЕстьSQLStateОбъект{
			//						log.Printf("ЕстьSQLStateОбъект %+v\n", ЕстьSQLStateОбъект)
			//					}
			//
			//					ОбработчикОшибки, ЕстьОбработчикОшибки := SQLStateОбъект.(map[string]interface{})["23505"]
			//
			//					if ЕстьОбработчикОшибки{
			//						client.ОбработатьОшибкуБД (ОбработчикОшибки , data)
			//					}
			//
			//
			//				}
			//				log.Printf(">>>> Ошибка SQL запроса: %+v \n\n",err)
			//			}
			//
			//
			//
			//		}
			//	}
			//}
			//tpl := template.Must(template.New("SqlЗапрос").Parse(SQLСкрипт.(string)))
			//БайтБуферSql := new(bytes.Buffer)
			/*
			Для постановки данных в sql передадим в data данные клиента, и данные из входящего вопроса: ВходящиеАргументы и вопрос.Выполнить.Действие
			*/
			//data:= map[string]interface{}{
			//	"client" : client.UserInfo.Info,
			//	"data": map[string]interface{}{},
			//}
			//if вопрос.ВходящиеАргументы!=nil{
			//	data["data"].(map[string]interface{})["ВходящиеАргументы"]=вопрос.ВходящиеАргументы
			//}
			//for НазваниеДействия, _ := range вопрос.Выполнить.Действие{
			//if вопрос.Выполнить.Действие !=nil{
			//	for НазваниеДействия, ДанныеДействия := range вопрос.Выполнить.Действие{
			//		log.Printf("НазваниеДействия %+v\n", НазваниеДействия)
			//		data["data"]=ДанныеДействия
			//	}
			//}
			//err := tpl.Execute(БайтБуферSql, data)
			//
			//if err != nil {
			//	log.Println("executing template:", err)
			//} else {
			//	SQLСкрипт = БайтБуферSql.String()
			//}
			//log.Printf("SQLСкрипт %+v\n", SQLСкрипт)
			//Данные ,err:= sqlStruct{
			//		Name:   "",
			//		Sql:    SQLСкрипт.(string),
			//		Values: [][]byte{},
			//	}.RunSQL(nil)
			//
			//if err != nil{
			//	if strings.Contains(err.Error(), "SQLSTATE 23505"){
			//		//алгоритм. Если запрос вернул ошибку с ограничением целостности уникального значения , тогда проверим есть ли обработчик SQLSTATE. и Наверное пока sql для обрабтки ошибок будет возвращать сразу текст который будет выводиться клиенту
			//		SQLStateОбъект, ЕстьSQLStateОбъект := SQLЗапрос.(map[string]interface{})["SQLState"]
			//
			//		if !ЕстьSQLStateОбъект{
			//			log.Printf("ЕстьSQLStateОбъект %+v\n", ЕстьSQLStateОбъект)
			//		}
			//		ОбработчикОшибки, ЕстьОбработчикОшибки := SQLStateОбъект.(map[string]interface{})["23505"]
			//		if ЕстьОбработчикОшибки{
			//			SQLStateСкрипт := ОбработчикОшибки.(map[string]interface{})["sql"]
			//			SQLStateString := template.Must(template.New("SQLStateСкриптЗапрос").Parse(SQLStateСкрипт.(string)))
			//			БайтБуферSQLState := new(bytes.Buffer)
			//			err := SQLStateString.Execute(БайтБуферSQLState, data)
			//			if err != nil {
			//				log.Println("executing template:", err)
			//			} else {
			//				SQLСкрипт = БайтБуферSQLState.String()
			//			}
			//			Результат ,err:= sqlStruct{
			//				Name:   "",
			//				Sql:    SQLСкрипт.(string),
			//				Values: [][]byte{},
			//			}.RunSQL(nil)
			//			if err != nil {
			//				log.Printf(">>>> ERROR \n %+v \n\n", err)
			//			}
			//
			//			// алгоритм т.к. это ошибка выполнения запроса, и нам нужно сообщить об этлм клиенту, то выведем сразу сообщение со стандартным шаблном ошибки.
			//			if len(Результат) ==1 {
			//				ИнформацияКлиенту := Результат[0]["message"]
			//				//type ДанныеОтвета struct {
			//				//	Контейнер string `json:"контейнер"`
			//				//	Данные interface{} `json:"данные"`
			//				//	HTML string `json:"html"`
			//				//	Обработчик string `json:"обработчик"` //JS функция или объект/класс/плагин для обработки данных (table..)
			//				//}
			//				ДанныеДляОтвета := &ДанныеОтвета{
			//					Обработчик:"FloatMessage",
			//					Данные:ИнформацияКлиенту,
			//				}
			//				СообщениеКлиенту:= &Сообщение{
			//					От: "io",
			//					Кому:client.Login,
			//					MessageType: []string{"error"},
			//					Контэнт:ДанныеДляОтвета,
			//				}
			//				log.Printf("СообщениеКлиенту %+v\n", СообщениеКлиенту)
			//				СообщениеКлиенту.СохранитьИОтправить(client)
			//			}
			//
			//		}
			//
			//
			//	}
			//	log.Printf(">>>> Ошибка SQL запроса: %+v \n\n",err)
			//}
				//log.Printf("Данные %+v\n", Данные)

			ДанныеДляРендераСИнфоКлиента := map[string]interface{}{
				"client":client.UserInfo.Info,
				"data":ДанныеЗапросов,
				"OspList":client.ПолучитьСписокОтделов(),
				"Osp": ПолучитьСписокОСП(),
			}
			if ЕстьHTMLШаблон && ДанныеHTMLШаблона != nil {

				НазваниеШаблона := ДанныеHTMLШаблона.(map[string]interface{})["HTML"].(string)
				ЦельДляВставки := ДанныеHTMLШаблона.(map[string]interface{})["цель"].(string)
				//, {"DATA": "rutokens_count", "цель": "rutokens.caption.rutokens_count"}]



				//log.Printf("ДанныеДляРендераСИнфоКлиента %+v\n", ДанныеДляРендераСИнфоКлиента)

				ДанныеДляОтвета :=  &ДанныеОтвета{
					Контейнер:ЦельДляВставки,
					HTML:string(render(НазваниеШаблона, ДанныеДляРендераСИнфоКлиента)),
				}

				СообщениеКлиенту:= &Сообщение{
					От: "io",
					Кому:client.Login,
					//MessageType: []string{"irritation","io_action"},
					Контэнт:ДанныеДляОтвета,
				}
				log.Printf("СообщениеКлиенту %+v\n", СообщениеКлиенту)
				СообщениеКлиенту.СохранитьИОтправить(client)
			}

			if ЕстьJSONДанные && JSONДанные != nil{
				log.Printf("ЕстьJSONДанные %+v JSONДанные %+v\n", ЕстьJSONДанные, JSONДанные)
				//пройдём по массиву , отрендерим и отправим клиенту
				for _, ЭлементДанных :=  range JSONДанные.([]interface{}) {

					Цель := ЭлементДанных.(map[string]interface{})["цель"].(string)

					tpl := template.Must(template.New("цель").Funcs(tplFunc()).Parse(Цель))
					БайтБуферSql := new(bytes.Buffer)
					err := tpl.Execute(БайтБуферSql, ДанныеДляРендераСИнфоКлиента)
					if err != nil {
						log.Printf(">>>> ERROR \n %+v \n\n", err)
					} else {
						Цель=БайтБуферSql.String()
					}
					ОбъектДанных := ЭлементДанных.(map[string]interface{})["DATA"].(string)

					if ДанныеЗапросов[ОбъектДанных] != nil{
						Данные := ДанныеЗапросов[ОбъектДанных]
						log.Printf("ОбъектДанных %+v\n", ОбъектДанных)
						log.Printf("получатель== 'сlient' %+v\n", ЭлементДанных.(map[string]interface{})["получатель"]== "сlient")

						if ЭлементДанных.(map[string]interface{})["получатель"] == "сlient"{
							ДанныеДляОтвета :=  &ДанныеОтвета{
								Контейнер:Цель,
								Данные:Данные,
							}
							СообщениеКлиенту:= &Сообщение{
								От: "io",
								Кому:client.Login,
								Контэнт:ДанныеДляОтвета,
							}
							log.Printf("СообщениеКлиенту %+v\n", СообщениеКлиенту)
							СообщениеКлиенту.СохранитьИОтправить(client)
						} else if  ЭлементДанных.(map[string]interface{})["получатель"] == "всем" { // всем кто онлайн
							for Логин, Получатель := range server.Clients{
								//log.Printf("sendLogin %+v clientOnline.Login %+v\n", )
								if Логин!=client.Login{
									СообщениеКлиенту:= &Сообщение{
										От: "io",
										Кому:Логин,
										Контэнт:&ДанныеОтвета{
											Контейнер:Цель,
											Данные:Данные,
										},
									}
									СообщениеКлиенту.СохранитьИОтправить(Получатель)
								}
							}
						} else if  ЭлементДанных.(map[string]interface{})["получатель"] == "админу" { // всем админам онлайн
							// получить список админов, проверить кто из них онлайн и отправить ему данные

						} else if  ЭлементДанных.(map[string]interface{})["получатель"] == "всем.на_странице" { // всем кто онлайн на той же странице

						} else if  ЭлементДанных.(map[string]interface{})["получатель"] == "админу.на_странице" { // админам онлайн на той же странице

						} else {

							Получатель, ЕстьПолучатель := server.Clients[ЭлементДанных.(map[string]interface{})["получатель"].(string)]
							if ЕстьПолучатель{
								СообщениеКлиенту:= &Сообщение{
									От: "io",
									Кому:Получатель.Login,
									Контэнт:&ДанныеОтвета{
										Контейнер:Цель,
										Данные:Данные,
									},
								}
								СообщениеКлиенту.СохранитьИОтправить(Получатель)

							}

						}

					}

				}
			}

		} else {
			if ЕстьHTMLШаблон && ДанныеHTMLШаблона != nil{

				log.Printf("Рендерим щаблон %+v\n", ДанныеHTMLШаблона)

				НазваниеШаблона := ДанныеHTMLШаблона.(map[string]interface{})["HTML"].(string)
				ЦельДляВставки := ДанныеHTMLШаблона.(map[string]interface{})["цель"].(string)

				ДанныеДляРендераСИнфоКлиента := map[string]interface{}{
					"client":client.UserInfo.Info,
					"OspList":client.ПолучитьСписокОтделов(),
					"Osp": ПолучитьСписокОСП(),
				}
log.Printf("ДанныеДляРендераСИнфоКлиента %+v\n", ДанныеДляРендераСИнфоКлиента)

				ДанныеДляОтвета :=  &ДанныеОтвета{
					Контейнер:ЦельДляВставки,
					HTML:string(render(НазваниеШаблона, ДанныеДляРендераСИнфоКлиента)),
				}
				СообщениеКлиенту:= &Сообщение{
							От: "io",
							Кому:client.Login,
							//MessageType: []string{"irritation","io_action"},
							Контэнт:ДанныеДляОтвета,
						}
				СообщениеКлиенту.СохранитьИОтправить(client)
			}
		}


	} else {
		// отвечать нечего
		log.Printf("ОтветитьКлиенту Диалог %+v ОтветЕсть %+v\n", Диалог, ОтветЕсть)
	}
}
