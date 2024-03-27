package main

import (
	. "aoanima.ru/ConnQuic"
	. "aoanima.ru/DGApi"
	. "aoanima.ru/Logger"
	. "aoanima.ru/QErrors"
	"golang.org/x/crypto/bcrypt"
)

func Регистрация(сообщение *Сообщение) (bool, СтатусСервиса) {
	ответ := сообщение.Ответ[Сервис] // получаем структуру для вставки ответа

	ответ.Сервис = Сервис
	ответ.ЗапросОбработан = true
	var логин, пароль, email string
	форма := сообщение.Запрос.Форма
	if len(форма) > 0 {

		логин, ошибка := ПолучитьЗначениеПоляФормы("логин", форма)
		if ошибка.Код != Ок {
			ответ.Данные = map[string]string{
				"СтатусРегистрации": "Нет логина",
			}
			ответ.СтатусОтвета = ошибка
			return false, ошибка
		}

		пароль, ошибка = ПолучитьЗначениеПоляФормы("пароль", форма)
		if ошибка.Код != Ок {
			ответ.Данные = map[string]string{
				"СтатусРегистрации": "Нет пароля",
			}
			ответ.СтатусОтвета = ошибка
			return false, ошибка
		}

		email, ошибка = ПолучитьЗначениеПоляФормы("email", форма)
		if ошибка.Код != Ок {
			ответ.Данные = map[string]string{
				"СтатусРегистрации": "Нет email",
			}
			ответ.СтатусОтвета = ошибка
			return false, ошибка
		}

		свободен, ошибка := ЛогинСвободен(логин)
		if !свободен && ошибка.Код != Ок {
			ответ.Данные = map[string]string{
				"СтатусРегистрации": "Логин занят",
			}
			ответ.СтатусОтвета = СтатусСервиса{
				Код:   ОшибкаРегистрации,
				Текст: "Логин занят",
			}
			return false, ответ.СтатусОтвета
		}

		emailСвободен, ошибка := EmailСвободен(email)
		if !emailСвободен && ошибка.Код != Ок {
			ответ.Данные = map[string]string{
				"СтатусРегистрации": "email уже заригистрирован",
			}
			ответ.СтатусОтвета = СтатусСервиса{
				Код:   ОшибкаРегистрации,
				Текст: "email уже заригистрирован",
			}
			return false, ответ.СтатусОтвета
		}
	}

	// новыйТокенКлиент := ТокенКлиента{
	// 	ИдКлиента: сообщение.ИдКлиента,
	// 	Роль:      []string{"клиент"},
	// 	Права:    []string{"чтение", "просмотр", "изменение своего"},
	// 	Истекает: time.Now().Add(60 * time.Minute),
	// 	Создан:   time.Now(),
	// }

	хэшПароля, err := bcrypt.GenerateFromPassword([]byte(пароль), bcrypt.DefaultCost)
	if err != nil {
		Ошибка(" %+v \n", err.Error())

	}
	НовыйКлиент := ДанныеКлиента{
		ИдКлиента: сообщение.ИдКлиента,
		Логин:     логин,
		Пароль:    string(хэшПароля),
		Email:     email,
		Роль:      []string{"клиент"},
		Права:     []string{"чтение", "просмотр", "изменение своего"},
	}

	JWT, ошибкаПодписи := СоздатьJWT(НовыйКлиент)
	if ошибкаПодписи.Код != Ок {
		Ошибка(" не удалось создать токен  %+v \n", ошибкаПодписи)
		return false, ошибкаПодписи
	}

	сообщение.JWT = JWT
	// сообщение.ТокенКлиента = новыйТокенКлиент

	ошибкаСохранения := СохранитьКлиентаВБД(НовыйКлиент)
	if ошибкаСохранения.Код != Ок {
		Ошибка(" не удалось сохранить в БД  %+v \n", ошибкаСохранения)
		ответ.Данные = map[string]string{
			"СтатусРегистрации": "не удалось записать данные пользователя в базу",
		}

		ответ.СтатусОтвета = СтатусСервиса{
			Код:   ошибкаСохранения.Код,
			Текст: ошибкаСохранения.Текст,
		}
		return false, СтатусСервиса{
			Код:   ошибкаСохранения.Код,
			Текст: ошибкаСохранения.Текст,
		}
	}

	ответ.Данные = map[string]string{
		"СтатусРегистрации": "успех",
	}
	ответ.СтатусОтвета = СтатусСервиса{
		Код:   Ок,
		Текст: "Успешная регистрация пользователя",
	}
	сообщение.Ответ[Сервис] = ответ
	return true, СтатусСервиса{
		Код:   Ок,
		Текст: "Успешная регистрация пользователя",
	}
}
