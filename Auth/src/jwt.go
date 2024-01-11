package main

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"time"

	. "aoanima.ru/ConnQuic"
	. "aoanima.ru/DataBase"
	. "aoanima.ru/Logger"
	. "aoanima.ru/QErrors"
	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func СоздатьJWT(данныеТокена ТокенКлиента) (string, ОшибкаСервиса) {

	claims := jwt.MapClaims{
		"UID":     данныеТокена.ИдКлиента,
		"role":    данныеТокена.Роль,
		"token":   данныеТокена.Токен,
		"access":  данныеТокена.Права,
		"expires": данныеТокена.Истекает,
		"created": данныеТокена.Создан,
	}
	секрет, статус := СоздатьСекретКлиента(данныеТокена.ИдКлиента)
	if статус.Код != Ок {
		return "", статус
	}
	return ПодписатьJWT(claims, секрет)
}

func ПодписатьJWT(данныеJWT jwt.MapClaims, секрет string) (string, ОшибкаСервиса) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, данныеJWT)

	// Подписываем токен с использованием секретного ключа
	подписаннаяСтрока, err := token.SignedString([]byte(секрет))
	if err != nil {
		return "", ОшибкаСервиса{
			Код:   ОшибкаПодписиJWT,
			Текст: "Не удалось подписать токен: " + err.Error(),
		}
	}

	return подписаннаяСтрока, ОшибкаСервиса{
		Код:   Ок,
		Текст: "JWT успешно подписан",
	}
}

func ПолучитьСекретныйКлючКлиента(ИдКлиента uuid.UUID) (string, ОшибкаСервиса) {
	// откроем файл в котором храниться скертный ключ
	данные, статус := БазаСекретов.Найти("ИдКлиента", ИдКлиента.String())
	if статус.Код != Ок {
		Инфо("  %+v \n", статус)
		return "", ОшибкаСервиса{
			Код:   СекретНеНайден,
			Текст: "Секретный ключ не найден",
		}
	}
	секрет := данные["ИдКлиента"].Данные["секрет"]
	return секрет.(string), ОшибкаСервиса{
		Код:   Ок,
		Текст: "Секретный ключ получен",
	}
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
func СоздатьСекретКлиента(ИдКлиента uuid.UUID) (string, ОшибкаСервиса) {
	// Генерируем байты случайных данных
	key := make([]byte, 256)
	_, err := rand.Read(key)
	if err != nil {
		return "", ОшибкаСервиса{
			Код:   Прочее,
			Текст: "Не удалось сгенерировать случайное значение: " + err.Error(),
		}
	}

	// Кодируем байты в base64 строку
	секрет := base64.URLEncoding.EncodeToString(key)

	документСекрет := Документ{
		ПервичныйКлюч: ПервичныйКлюч(ИдКлиента.String()),
		Данные: map[string]interface{}{
			"секрет": секрет,
		},
	}

	статус := БазаСекретов.ВставитьДокумент(&документСекрет, true)
	if статус.Код != Ок {
		Инфо("  %+v \n", статус)
		return "", ОшибкаСервиса{
			Код:   ОшибкаЗаписи,
			Текст: "Секретный ключ не удалось записать в базу",
		}
	}
	// запишем ключ в файл
	// err = os.WriteFile("secrets/"+ИдКлиента.String(), []byte(keyString), 0644)
	// if err != nil {
	// 	Ошибка("  %+v \n", err)
	// }

	return секрет, ОшибкаСервиса{
		Код:   Ок,
		Текст: "Секретный ключ создан",
	}
}
func ВлаидацияТокена(сообщение *Сообщение) (bool, ОшибкаСервиса) {
	секрет, статус := ПолучитьСекретныйКлючКлиента(сообщение.ИдКлиента)
	if статус.Код == Ок {
		return false, ОшибкаСервиса{
			Код:   СекретНеНайден,
			Текст: "не удалось получить секретный ключ клиента",
		}
	}
	token, err := jwt.Parse(сообщение.JWT, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("Не известный метод шифрования")
		}
		return []byte(секрет), nil
	})
	if err != nil {
		return false, ОшибкаСервиса{
			Код:   ОшибкаВалидацииJWT,
			Текст: "не удалось валидировать токен: " + err.Error(),
		}
	}

	// Проверяем валидность токена
	if токен, ok := token.Claims.(jwt.MapClaims); !ok || !token.Valid {
		Ошибка(" токен не валидный %+v \n", токен)
		сообщение.JWT = "invaild"
		return false, ОшибкаСервиса{
			Код:   ОшибкаВалидацииJWT,
			Текст: "Токен не валидный",
		}
	} else {
		истекает := time.Unix(токен["expires"].(int64), 0)

		Инфо(" токен Валидный, подпишем заново %+v \n", токен)
		// если осталось менее 5 минут переподпишем токен
		if осталосьВремениДоИстечения := time.Now().Sub(истекает); time.Duration(осталосьВремениДоИстечения.Minutes()) < 5*time.Minute {

			токен["token"] = СоздатьТокенОбновления(16)
			токен["expires"] = time.Now().Add(60 * time.Minute).Unix()
			токен["created"] = time.Now().Unix()

			новыйСекрет, статус := СоздатьСекретКлиента(токен["UID"].(uuid.UUID))
			if статус.Код != Ок {
				Ошибка(" %+v \n", статус)
				return false, статус
			}
			сообщение.ТокенКлиента.Токен = токен["token"].(string)
			сообщение.ТокенКлиента.Истекает = токен["expires"].(int64)
			сообщение.ТокенКлиента.Создан = токен["created"].(int64)

			новыйJWT, ошибкаСервиса := ПодписатьJWT(токен, новыйСекрет)
			if ошибкаСервиса.Код != Ок {
				Ошибка("  %+v \n", err)
			}
			сообщение.JWT = новыйJWT

			ответ := сообщение.Ответ[Сервис]

			ответ.Сервис = Сервис
			ответ.ЗапросОбработан = true
			ответ.Данные = map[string]bool{
				"ТокенВерный": true,
			}
			ответ.ОшибкаСервиса = ОшибкаСервиса{
				Код:   Ок,
				Текст: "Токен валидный",
			}
			сообщение.Ответ[Сервис] = ответ
		}
		return true, ОшибкаСервиса{
			Код:   Ок,
			Текст: "Токен валидный",
		}
	}
}