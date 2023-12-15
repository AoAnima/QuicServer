package main

import (
	"context"
	_ "fmt"
	_ "io"
	_ "io/ioutil"
	"log"
	ctx "main/RWContext"
	"main/connect"
	"main/loger"
	"main/render"
	_ "main/scheduler"
	"net/http"
	"os"
	"path/filepath"
	"time"

	_ "github.com/lib/pq"
	_ "github.com/mysql"
)

var SessionMap = make(map[string]ctx.UserSession)

/*
doc: Config - конфигурация основгого потока.
note: необходимо реалицовать чтение JSON файла и парсинг в эту структуру. Если файл не найден то необходимо сделать чтобы на старте выавлись вопросы с заполнением всех полей и сохранение в конфиг файл.
Пока данные сохраним в коде.
*/

func init() {
	log.SetFlags(log.Lshortfile)
	//loc, err := time.LoadLocation("Europe/Moscow")
	//log.Print(loc, err)
	//log.Print(time.Now().Clock())
	//log.Print(time.Now().Hour())
	//
	//h := time.Duration(23 - time.Now().Hour()) * time.Hour
	//m := time.Duration(59 - time.Now().Minute()) * time.Minute
	//s := time.Duration(60 - time.Now().Second()) * time.Second
	//log.Printf("%+v %+v %+v  %+v\n", h, m, s, h+m+s)

	// Создаём сессию не авторизированного пользователя
	SessionMap["unauth"] = ctx.UserSession{
		UidNumber:    0,
		Uid:          "",
		Token:        "unauth",
		AuthStatus:   "unauth",
		AccessGroups: []string{"unauth"},
		DateAuth:     time.Now(),
		//DateAuth: time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(),0,0,0,0, loc),
		//DB: ctx.UserDataForDb{
		//	"admin",
		//	"admin",
		//},
	}
	ctx.MainCtx = context.WithValue(context.Background(), "db", connect.Dbconnect())
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("filepath %+v\n", dir)
	pattern := map[string]string{
		"name":    "tplFiles",
		"pattern": "../www/tpl/*/*.html",
	}
	render.ParseTplDir(pattern)

}

func main() {

	//log.Printf("Config %+v\n", Log(Config))
	//log.Printf("SessionMap %+v\n", Log(SessionMap))

	http.HandleFunc("/static/", StaticHandler)
	http.HandleFunc("/src/", StaticHandler)

	http.HandleFunc("/", MainHandler)
	http.HandleFunc("/gettpl", JSHandler)

	log.Printf("%+v\n", "Стартуем сервер")
	//authDate := auth.AuthDate{
	//	"maksimchuk@r26",
	//	"Tanos@27$",
	//}
	//auth.Auth(authDate)
	//scheduler.SyncActions()
	//ldapQuery.GetFsspRole(true)
	//ldapQuery.GetLdapUsers()
	//ldapQuery.GetUserRole()

	serverError := http.ListenAndServe(":8080", nil)
	log.Printf("%+v\n", loger.Log(serverError))

}

//
//
//func Dbconnect() *sql.DB {
//
//	var db *sql.DB
////username:password@protocol(address)/dbname
//	dbinfo := fmt.Sprintf("%s:%s@(%s:%d)/%s",
//		"root", "1111", "10.26.4.20", 3306, "fsspsk")
//	log.Printf("%+v\n", dbinfo)
//	var err error
//	//"postgres://user:pass@localhost/bookstore"
//	db, err = sql.Open("mysql", dbinfo)
//	if err != nil {
//		log.Printf("Ошибка подключения к создания подключения к базе данных %+v\n", err)
//	}
//
//	return db
//}
//
//func HostDB() *sql.DB { // *sql.DB //RequestCtx ctx.RWContext
//
//	//host := RequestCtx.Request.R.Host
//	//if	RequestCtx.Request.R.Host == "pro.ru"{
//	//host := "fsspsk"
//	//}
//
//	var db *sql.DB
//	dbinfo := fmt.Sprintf("user=%s password=%s port=%d dbname=%s host=%s sslmode=disable", Config.DBuser, Config.DBpassword, Config.DBport, Config.DBname, Config.DBhost)
//	//dbinfo := fmt.Sprintf("user=%s password=%s port=%d dbname=%s host=%s sslmode=disable", "postgres", "postgres", 5432, host, "10.26.6.15")
//	log.Printf("%+v\n", dbinfo)
//	var err error
//	//"postgres://user:pass@localhost/bookstore"
//	if db, err = sql.Open("postgres", dbinfo); err != nil {
//		log.Printf("%+v\n", err)
//	}
//	errPing := db.Ping()
//	if errPing != nil {
//		db=nil
//		log.Printf("%+v\n", errPing)
//	}
//	//log.Printf("%+v\n", db)
//	//log.Printf("%+v\n",errPing)
//	return db
//}

func Contains(slice []string, item string) bool {
	set := make(map[string]struct{}, len(slice))
	for _, s := range slice {
		set[s] = struct{}{}
	}
	_, ok := set[item]
	return ok
}

func FileType(contentType []string) string {
	contentTyperes := contentType[0]
	var fileType string
	//if contentTyperes == "image/jpeg"{
	//	fileType = "jpeg"
	//}
	//contentType := http.DetectContentType(buffer)
	switch contentTyperes {
	case "image/jpeg":
		fileType = "jpeg"
	case "image/jpg":
		fileType = "jpg"
	case "image/gif":
		fileType = "gif"
	case "image/png":
		fileType = "png"
	}
	return fileType
}
