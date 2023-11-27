package main

import (
	"crypto/rand"
	"encoding/base64"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"

	. "aoanima.ru/ConnQuic"
	. "aoanima.ru/Logger"
	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/quic-go/quic-go"
)

var клиент = make(Клиент)
var Сервис ИмяСервиса = "Авторизация"

func main() {
	Инфо("  %+v \n", " Запуск сервиса Авторизации")
	сервер := &СхемаСервера{
		Имя:   "SynQuic",
		Адрес: "localhost:4242",
		ДанныеСессии: ДанныеСессии{
			Блок:   &sync.RWMutex{},
			Потоки: []quic.Stream{},
		},
	}
	сообщениеРегистрации := Сообщение{
		Сервис:      Сервис,
		Регистрация: true,
		Маршруты:    []Маршрут{"reg", "auth", "verify"},
	}

	клиент.Соединиться(сервер,
		сообщениеРегистрации,
		ОбработчикОтветаРегистрации,
		ОбработчикЗапросовСервера)
}

func ОбработчикЗапросовСервера(поток quic.Stream, сообщение Сообщение) error {
	Инфо("  ОбработчикЗапросовСервера %+v \n", сообщение)
	var err error
	параметрыЗапроса, err := url.Parse(сообщение.Запрос.МаршрутЗапроса)
	Инфо(" параметрыЗапроса %+v \n", параметрыЗапроса)
	if err != nil {
		Ошибка("Ошибка при парсинге СтрокаЗапроса запроса:", err)
	}

	параметрыЗапроса.Path = strings.Trim(параметрыЗапроса.Path, "/")
	дейсвтия := strings.Split(параметрыЗапроса.Path, "/")
	// TODO: Тут нужно сделать обращение к БД для получения очереди сервисов которые должны обработать запрос еред возвратом клиенту
	var Действие string
	if len(дейсвтия) == 0 {
		Инфо(" Пустой маршрут, добавляем в маршруты обработку по умолчанию.... \n")
		// Читаем заголовки парсим и проверяем JWT
		Действие = "verify"

	} else {
		Действие = дейсвтия[0]
	}

	switch Действие {
	case "reg":
		err = Регистрация(&сообщение)
	case "auth":
		err = Авторизация(&сообщение)
	case "verify":
		err = ВлаидацияТокена(&сообщение)
	}

	if err != nil {
		Ошибка(" Генерируем сообщение ощибки или возвращаем сообщение ошибки %+v \n", err)
	}

	отправить, err := Кодировать(сообщение)
	if err != nil {
		Ошибка("  %+v \n", err)
	}
	поток.Write(отправить)

}

func Регистрация(сообщение *Сообщение) error {

	новыйТокенКлиент := ТокенКлинета{
		ИдКлиента: сообщение.ИдКлиента,
		Роль:      "client",
		Токен:     СоздатьТокенОбновления(16),
		Права:     []string{"client"},
		Истекает:  time.Now().Add(60 * time.Minute),
		Создан:    time.Now(),
	}

	токен, err := СоздатьJWT(новыйТокенКлиент)
	if err != nil {
		Ошибка(" не удалось создать токен  %+v \n", err)
		return err
	}
	сообщение.JWT = токен
	err = СохранитьКлиентаВБД(сообщение)
	if err != nil {
		Ошибка(" не удалось сохранить в БД  %+v \n", err)
		return err
	}
	return nil
}

func СохранитьКлиентаВБД(сообщение *Сообщение) error {
	// TODO: Сохранить клиента в БД
	Инфо(" Сохранить клиента в БД  \n")
}

func Авторизация(сообщение *Сообщение) {

}

func ОбработчикОтветаРегистрации(сообщение Сообщение) {
	Инфо("  ОбработчикОтветаРегистрации %+v \n", сообщение)
}

func СоздатьJWT(данныеТокена ТокенКлинета) (string, error) {

	claims := jwt.MapClaims{
		"UID":     данныеТокена.ИдКлиента,
		"role":    данныеТокена.Роль,
		"token":   данныеТокена.Токен,
		"access":  данныеТокена.Права,
		"expires": данныеТокена.Истекает.Unix(),
		"created": данныеТокена.Создан,
	}
	// Создаем токен с указанными утверждениями (claims)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Устанавливаем время истечения срока действия токена

	// claims["exp"] = времяИстечения.Unix()

	// Подписываем токен с использованием секретного ключа
	tokenString, err := token.SignedString([]byte(ПолучитьСекретныйКлючКлиента(данныеТокена.ИдКлиента)))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func ПолучитьСекретныйКлючКлиента(ИдКлиента uuid.UUID) string {
	// откроем файл в котором храниться скертный ключ
	файл, err := os.Open("secrets/" + ИдКлиента.String())
	defer файл.Close()
	if err != nil {
		Ошибка("  %+v \n", err)
		return СоздатьСеркетКлиента(ИдКлиента)
	}
	секретныйКлюч := make([]byte, 256)
	// прочитаем содержимое файла и вернём
	_, err = файл.Read(секретныйКлюч)
	if err != nil {
		Ошибка("  %+v \n", err)
	}

	return string(секретныйКлюч)
}

func СоздатьТокенОбновления(размер int) string {
	key := make([]byte, размер)
	_, err := rand.Read(key)
	if err != nil {
		return ""
	}

	// Кодируем байты в base64 строку
	keyString := base64.URLEncoding.EncodeToString(key)
	return keyString
}
func СоздатьСеркетКлиента(ИдКлиента uuid.UUID) string {
	// Генерируем байты случайных данных
	key := make([]byte, 256)
	_, err := rand.Read(key)
	if err != nil {
		return ""
	}

	// Кодируем байты в base64 строку
	keyString := base64.URLEncoding.EncodeToString(key)

	// запишем ключ в файл
	err = os.WriteFile("secrets/"+ИдКлиента.String(), []byte(keyString), 0644)
	if err != nil {
		Ошибка("  %+v \n", err)
	}
	return keyString
}
func ВлаидацияТокена(сообщение Сообщение, серкет string) (bool, error) {

	token, err := jwt.Parse(JWT, func(token *jwt.Token) (interface{}, error) {
		return []byte(серкет), nil
	})
	if err != nil {
		return false, err
	}

	// Проверяем валидность токена
	if _, ok := token.Claims.(jwt.Claims); !ok || !token.Valid {
		return false, nil
	}
	Инфо(" Нужно обновить токен и переподписать %+v \n")
	return true, nil
}
