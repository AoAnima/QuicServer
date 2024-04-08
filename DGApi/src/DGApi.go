package main

import (
	"context"
	"fmt"
	"log"

	. "aoanima.ru/ConnQuic"
	. "aoanima.ru/DGApi"
	. "aoanima.ru/Logger"
	. "aoanima.ru/QErrors"
	dgo "github.com/dgraph-io/dgo/v230"
	"github.com/dgraph-io/dgo/v230/protos/api"
	"github.com/google/uuid"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// type Адрес struct {
// 	Страна        string `json:"страна,omitempty"`
// 	Город         string `json:"город,omitempty"`
// 	Район         string `json:"район,omitempty"`
// 	ТипУлицы      string `json:"тип_улицы,omitempty"`
// 	НазваниеУлицы string `json:"название_улицы,omitempty"`
// 	НомерДома     string `json:"номер_дома,omitempty"`
// 	Корпус        string `json:"корпус,omitempty"`
// 	НомерКвартиры string `json:"номер_квартиры,omitempty"`
// }
// type Секрет struct {
// 	ИдКлиента string    `json:"ид_клиента"`
// 	Секрет    string    `json:"секрет"`
// 	Обновлен  time.Time `json:"обновлен"`
// }
// type ДанныеКлиента struct {
// 	Имя       string    `json:"имя,omitempty"`
// 	Фамилия   string    `json:"фамилия,omitempty"`
// 	Отчество  string    `json:"отчество,omitempty"`
// 	ИдКлиента uuid.UUID `json:"ид_клиента,omitempty"`
// 	Роль      []string  `json:"роль,omitempty"`
// 	Права     []string  `json:"права_доступа,omitempty"`
// 	Статус    string    `json:"статус,omitempty"`
// 	Аватар    string    `json:"аватар,omitempty"`
// 	Email     string    `json:"email,omitempty"`
// 	Логин     string    `json:"логин,omitempty"`
// 	Пароль    string    `json:"пароль,omitempty"`
// 	JWT       string    `json:"jwt,omitempty"`
// 	Телефон   string    `json:"телефон,omitempty"`
// 	Адрес     Адрес     `json:"адрес,omitempty"`
// 	Создан    time.Time `json:"создан,omitempty"`
// 	Обновлен  time.Time `json:"обновлен,omitempty"`
// 	ОСебе     string    `json:"о_себе,omitempty"`
// 	СоцСети   []string  `json:"социальные_ссылки,omitempty"`
// 	Профиль   map[string]interface{}
// }

var База = ДГраф()

func main() {

	defer База.ЗакрытьДГраф()

	статус := База.Схема(СхемаБазы)
	Инфо("  %+v \n", статус)
	// маршрут := "рабочийСтол1"
	// добавить := `[
	// 	{
	// 		"маршрут": "` + маршрут + `",
	// 		"uid": "_:` + маршрут + `",
	// 		"доступ":[{
	// 			"роль" : "администратор1",
	// 			"права": ["чтение", "изменение","удаление"],
	// 			"dgraph.type": "Доступ"
	// 		},
	// 		{
	// 			"роль" : "пользователь",
	// 			"права": ["чтение"],
	// 			"dgraph.type": "Доступ"
	// 		}],
	// 		"описание": "описание маргрута",
	// 		"dgraph.type": "Маршрут1"
	// 	}
	// ]`

	// query:{
	//     checkRoute(func: has(<маршрут>)) {
	//  		   expand(_all_)
	//  		}
	// }

	/*
		РАбочий запрос на выбор с вложенной структурой
	*/
	// ответ, статусхемы := граф.Получить(ДанныеЗапроса{
	// 	Запрос: `{
	// 		checkRoute(func: has(<доступ>)) {
	// 			<маршрут>,
	// 			<доступ> {
	// 				<роль>,
	// 				<права>
	// 				}
	// 			<описание>
	// 			}
	// 	 	 }`,
	// 	Данные: nil,
	// }
	// Так же рабочий пример фильтрации запроса по вложенной структуре доступ
	// query:{
	// 	checkRoute(func: eq(<маршрут>, "рабочийСтол")) {
	// 	<маршрут>,
	// 	<доступ> @filter(eq(<роль>, "пользователь"))  {
	// 		  <роль>,
	// 		  <права>
	// 	},
	// 			  <описание>
	// 				 }
	// 	}

	// ответ, статусхемы := граф.Получить(ДанныеЗапроса{
	// 	Запрос: `{
	// 		checkRoute(func: has(<доступ>)) {
	// 			<маршрут>,
	// 			<доступ>  @filter(eq(<роль>, "пользователь")){
	// 				<роль>,
	// 				<права_доступа>
	// 				}
	// 			<описание>
	// 			}
	// 	 	 }`,
	// 	Данные: nil,
	// })
	// Инфо("  %+s %+v \n", ответ, статусхемы)

	// for i := 0; i < 1000; i++ {
	// 	go func() {
	// 		Даннные := ДанныеЗапроса{
	// 			Запрос: "set ",
	// 			Данные: make(map[string]string),
	// 		}
	// 		результат, статус := граф.Изменить(Даннные)
	// 		if статус.Код != Ок {
	// 			Инфо(" %+v \n", статус.Текст)
	// 		}
	// 		Инфо(" результат %+v \n", результат)
	// 	}()
	// }
	// Аутентификация()
	// ИзменитьКоличествоПопытокАутентификации()
	// Аутентификация()
	// ПримерРегистрацииКлиента()
	ТестОбработчиков()
	// свободен, статуслогин := ЛогинСвободен("anima_ao")
	// Инфо("ЛогинСвободен %+v  статуслогин %+v \n", свободен, статуслогин)

	// свободен, статуслогин = EmailСвободен("nefri@ya.ru")
	// Инфо("EmailСвободен %+v  статус  %+v \n", свободен, статуслогин)

	// InitSchema()
	// // SetClient()
	// Auth()
}

// type КонфигурацияОбработчика struct {
// 	UID          string         `json:"uid,omitempty"`
// 	Маршрут      string         `json:"маршрут,omitempty"`
// 	Действие     string         `json:"действие,omitempty"`
// 	Обработчик   string         `json:"обработчик,omitempty"`
// 	ПраваДоступа []ПраваДоступа `json:"доступ"`
// 	Описание     string         `json:"описание,omitempty"`
// 	Шаблонизатор []Шаблон       `json:"шаблонизатор,omitempty"`
// 	Ассинхронно  bool           `json:"ассинхронно,omitempty"`
// 	Тип          string         `json:"dgraph.type,omitempty"`
// }

// type Шаблон struct {
// 	UID          string         `json:"uid,omitempty"`
// 	Тип          string         `json:"dgraph.type,omitempty"`
// 	Код          int            `json:"код,omitempty"` // статус ответа сервиса QErrors
// 	Шаблон       string         `json:"имя_шаблона,omitempty"`
// 	ПраваДоступа []ПраваДоступа `json:"доступ,omitempty"`
// }
// type ПраваДоступа struct {
// 	UID   string   `json:"uid,omitempty"`
// 	Тип   string   `json:"dgraph.type,omitempty"`
// 	Логин []string `json:"пользователи"`
// 	Роль  []string `json:"роль"`
// 	Права []string `json:"права"`
// }

func ТестОбработчиков() {

	// Для сохранение данных одна функция, для изменения данных нужно делать другую функцию чтобы изменять связанные узлы корректно

	новыйОбработчик := &КонфигурацияОбработчика{
		UID:        "_:новыйОбработчик",
		Тип:        "Обработчик",
		Маршрут:    "/редакторОбработчиков",
		Действие:   "сохранитьОбработчик",
		Обработчик: "сохранитьОбработчик",
		ПраваДоступа: []ПраваДоступа{
			{
				Тип:   "ПраваДоступа",
				Роль:  []string{"админ"},
				Права: []string{"чтение", "создание", "изменение", "удаление"},
			},
			{
				Тип:   "ПраваДоступа",
				Роль:  []string{"модератор"},
				Права: []string{"чтение", "создание", "изменение своего", "удаление своего"},
			},
		},
		Описание: "обработчик для создания очереди обработчиков новое",
		Шаблонизатор: []Шаблон{{
			Тип:    "Шаблон",
			Код:    Ок,
			Шаблон: "всплывающееУведомление",
		},
			{
				Тип:    "Шаблон",
				Код:    Прочее,
				Шаблон: "всплывающаяОшибка",
			},
		},
		// Ассинхронно: ассинхронно,
	}

	обработчикБин, статус := Json(новыйОбработчик)
	if статус != nil {
		Ошибка(" статус %+v \n", статус, обработчикБин)
	}
	// @filter(eq(<маршрут>, $path) OR eq(<действие>, $action))
	данные := ДанныеЗапроса{
		Запрос: `query <УдалитьОбработчик>($uid : string) {
							<УдаляемыеУзлы>(func: uid($uid)) {				
								<обработчик_ид> as uid
								<доступ> {
									<доступ_ид> as uid
								}
								<шаблонизатор> {
									<шаблонизатор_ид> as uid
								}
							}
							<об>(func: uid(<обработчик_ид>)) {
								val(<доступ_ид>)
								val(<обработчик_ид>)
								val(<шаблонизатор_ид>)
								uid
								<маршрут>
	 							<действие>
	 							<обработчик>
	 							<доступ>{
 									<пользователи>
	 								<права>
	 								<роль>
	 								uid
									dgraph.type
								}
								<описание>
								<шаблонизатор> {
 									<имя_шаблона>
								 	uid
	 								<код>
									dgraph.type
	 							}
							}
			 			} `,
		Мутация: []Мутация{
			{
				// Условие: "@if(eq(len(handlers), 0))",
				// Мутация: обработчикБин,
				Удалить: []byte(`[
					{"uid": "uid(доступ_ид)"},
					{"uid": "uid(обработчик_ид)"},
					{"uid": "uid(шаблонизатор_ид)"}
					]`),
				// Удалить: []byte(`{
				// 				"uid": "uid(handler)",
				// 				"обработчик": null,
				// 				"маршрут": null,
				// 				"действие": null,
				// 				"описание": null,
				// 				"шаблонизатор": null,
				// 				"ассинхронно": null,
				// 				"права": null
				// 			}`),
			},
		},
		Данные: map[string]string{
			"$uid": "0x4e31",
		},
	}
	// данные := ДанныеЗапроса{
	// 	Запрос: `query <Обработчики>($hanlder : string, $path : string, $action : string) {
	// 						handlers as var(func: type(<Обработчики>)) @filter(eq(<обработчик>, $hanlder) AND eq(<маршрут>, $path) AND eq(<действие>, $action)){
	// 							<обработчик>
	// 							<маршрут>
	// 							<действие>
	// 						}

	// 						 <Обработчики>(func: uid(<handlers>)){
	// 							count(handlers)
	// 							uid
	// 							<маршрут>
	// 							<действие>
	// 							<обработчик>
	// 							<права>{
	// 								<пользователи>
	// 								<права>
	// 								<роль>
	// 								uid
	// 								dgraph.type
	// 							}
	// 							<описание>
	// 							<шаблонизатор> {
	// 								<имя_шаблона>
	// 								<код>
	// 								dgraph.type
	// 							}
	// 							<ассинхронно>
	// 							dgraph.type
	// 						}

	// 		 			} `,
	// 	Мутация: []Мутация{
	// 		{
	// 			// Условие: "@if(gt(len(handlers), 0))",
	// 			// Мутация: обработчикБин,
	// 			// Удалить: []byte(`{
	// 			// 	"uid": "uid(handler)",
	// 			// 	}`),
	// 			// Удалить: []byte(`{
	// 			// 				"uid": "uid(handler)",
	// 			// 				"обработчик": null,
	// 			// 				"маршрут": null,
	// 			// 				"действие": null,
	// 			// 				"описание": null,
	// 			// 				"шаблонизатор": null,
	// 			// 				"ассинхронно": null,
	// 			// 				"права": null
	// 			// 			}`),
	// 		},
	// 	},
	// 	Данные: map[string]string{
	// 		"$hanlder": новыйОбработчик.Обработчик,
	// 		"$path":    новыйОбработчик.Маршрут,
	// 		"$action":  новыйОбработчик.Действие,
	// 	},
	// }

	ответ, статусБазы := База.Изменить(данные)
	if статусБазы.Код != Ок {
		Ошибка(" статус %+v \n данные %+v \n", статусБазы, данные)
	}
	var данныеОтвета interface{}
	ИзJson(ответ, &данныеОтвета)
	Инфо("Исходные данные %+v \n ответ %+s \n", данные, данныеОтвета)

	// данные = ДанныеЗапроса{
	// 	Запрос: `query <Обработчики>($hanlder : string, $path : string, $action : string) {
	// 						var(func: eq(<обработчик>, $hanlder)) {
	// 							handler as uid
	// 						}
	// 						<Обработчики>(func: eq(<обработчик>, $hanlder))  {
	// 							<обработчик>
	// 							<маршрут>
	// 							<действие>
	// 							<описание>
	// 							<шаблон>
	// 							<ассинхронно>
	// 							<права> {
	// 								<логин>
	// 								<роль>
	// 								<права>
	// 							}
	// 						}
	// 		 			} `,
	// 	Мутация: []Мутация{
	// 		{
	// 			// Условие: "@if(gt(len(a), 0))",
	// 			Удалить: []byte(`{
	// 							"uid": "uid(handler)",
	// 							"обработчик": null,
	// 							"маршрут": null,
	// 							"действие": null,
	// 							"описание": null,
	// 							"шаблон": null,
	// 							"ассинхронно": null,
	// 							"права": null
	// 						}`),
	// 		},
	// 	},
	// 	Данные: map[string]string{
	// 		"$hanlder": новыйОбработчик.Обработчик,
	// 		"$path":    новыйОбработчик.Маршрут,
	// 		"$action":  новыйОбработчик.Действие,
	// 	},
	// }

	// // ответ, статусБазы := База.Изменить(данные)
	// // if статусБазы.Код != Ок {
	// // 	Ошибка(" статус %+v \n", статусБазы)
	// // }
	// // Инфо(" ответ %+s \n", ответ)

	данные = ДанныеЗапроса{
		Запрос: `query <Обработчики>($hanlder : string, $path : string, $action : string) {					
							<Обработчик>(func: has(<обработчик>)){	
								uid			
								<маршрут>
								<действие>
								<обработчик>
								<доступ>{
									<пользователи>
									<права>
									<роль>
									uid
									dgraph.type
								}
								<описание>
								<шаблонизатор> {
									uid
									<имя_шаблона>
									<код>
									dgraph.type
								}
								<ассинхронно>
								dgraph.type
							
							}	
			 			} `,

		Данные: map[string]string{
			"$hanlder": "создатьОбработчик",
			"$path":    "/редакторОбработчиков",
			"$action":  "создатьОбработчик",
		},
	}

	ответ, статусПолучения := База.Получить(данные)
	if статусПолучения.Код != Ок {
		Ошибка(" статус %+v \n", статусПолучения)
	}

	ИзJson(ответ, &данныеОтвета)
	Инфо(" ответ %+s \n", данныеОтвета)

}

func ИзменитьКоличествоПопытокАутентификации() {
	// НовыйКлиент := ДанныеКлиента{
	// 	ИдКлиента: uuid.New(),
	// 	Логин:     "aaaaa1",
	// 	Пароль:    string("password"),
	// 	Email:     "aaaa1@ya.ru",
	// 	Роль:      []string{"клиент"},
	// 	Права:     []string{"чтение", "просмотр", "изменение своего"},
	// }
	// новыйКлиентаСтрока, ошибка := Json(НовыйКлиент)
	// if ошибка != nil {
	// 	Ошибка(" сериализаии нового клиента  %+v  новыйКлиент %+v \n", ошибка, НовыйКлиент)
	// }

	ctx := context.Background()
	var транзакция Транзакция
	транзакция.Txn = База.Граф.NewTxn()
	defer транзакция.Discard(ctx)
	данные := ДанныеЗапроса{
		Запрос: `query User($login: string, $pass: string) {
					 var(func: eq(<логин>, $login) ) {
						<статус_пароля> as checkpwd(<пароль>, $pass)
					}
					<ПарольВерный>(func: eq(val(<статус_пароля>), 1)) {
						<aутентифицирован>: val(<статус_пароля>)
						uid
						<ид_клиента>
						<имя>
						<отчество>
						<логин>						email
						<права_доступа>
						<статус>
						jwt					
						<количество_неудачных_попыток_входа>
						<удача> as <количество_удачных_попыток_входа>
						<количество_удач> as math(<удача>+1)	
					}

					<ПарольНеВерный>(func: eq(val(<статус_пароля>), 0)) {
						<aутентифицирован>: val(<статус_пароля>)
						<неудача> as <количество_неудачных_попыток_входа>
						<количество_неудач> as math(<неудача>+1)	
					}
				}
				`,
		Мутация: []Мутация{
			{
				Условие: "@if(ge(len(<удача>), 1))",
				Мутация: []byte(`
									{
									"uid": "uid(статус_пароля)",
									"количество_удачных_попыток_входа": "val(количество_удач)",
									"количество_неудачных_попыток_входа": 0
									}
								`),
			},
			{
				Условие: "@if(ge(len(<неудача>), 1))",
				Мутация: []byte(`
					{
					"uid": "uid(статус_пароля)",
					"количество_неудачных_попыток_входа": "val(количество_неудач)"
					}
				`),
			},
			// {
			// 	Условие: "@if(eq(val(<статус_пароля>), 1))",
			// 	Мутация: []byte(`{
			// 				"uid" : uid(v),
			// 				"количество_неудачных_попыток_входа": 0
			// 			}
			// 		`),
			// },
		},
		Данные: map[string]string{
			"$login": "aaaaa",
			"$pass":  "password",
		},
	}
	ответ, статус := транзакция.Измененить(данные, ctx)
	Инфо(" %+v  %+v \n", ответ, статус)
	транзакция.Commit(ctx)
}
func Аутентификация() {

	запрос := ДанныеЗапроса{
		Запрос: `query User($login: string, $pass: string) {

			<ПроверитьПароль>(func: eq(<логин>, $login) ) {						
				<статус_пароля> as   checkpwd(<пароль>, $pass)	
			}
			<ПарольВерный>(func: eq(val(<статус_пароля>), 1)) {
				<статус_пароля>: val(<статус_пароля>)
				uid
				<ид_клиента>
				<имя>
				<фамилия>
				<отчество>
				<логин>		
				email					
				<права_доступа>				
				<статус>							
				jwt	
				<количество_неудачных_попыток_входа>
				<количество_удачных_попыток_входа>
			  }
			
			  <ПарольНеВерный>(func: eq(val(<статус_пароля>), 0)) {
				<статус_пароля>: val(<статус_пароля>)
				expand(_all_)
			  }
		}
		`,
		Данные: map[string]string{
			"$login": "aaaaa",
			"$pass":  "password",
		},
	}

	ответ, статус := База.Получить(запрос)

	Инфо("Аутентификация ответ %+s; статус %+v \n", ответ, статус)

}

func ЛогинСвободен(логин string) (bool, СтатусСервиса) {
	Инфо("ЛогинСвободен , нужно проверить логин на свободу \n")

	ответ, статус := База.Получить(ДанныеЗапроса{
		Запрос: `
			query checkLogin($login: string){
				<Логин>(func: eq(<логин>, $login)) {
				  <занято>:count(uid)
				}
			}`,
		Данные: map[string]string{
			"$login": логин,
		},
	})

	if статус.Код != Ок {
		Инфо(" %+v \n", статус.Текст)
		return false, СтатусСервиса{
			Код:   статус.Код,
			Текст: статус.Текст,
		}
	}

	картаОтвета := ОтветИзБазы{}
	ошибкаРазбора := ИзJson(ответ, &картаОтвета)
	if ошибкаРазбора != nil {
		Ошибка(" Не удалось разобрать ответ %+v \n", ошибкаРазбора.Error())
	}

	for имя, массивДанных := range картаОтвета {
		if len(массивДанных) == 1 {
			Инфо(" %+v  %+v \n", имя, массивДанных)
			Занято := массивДанных[0]["занято"]
			Инфо(" %+v  %+v \n", Занято, uint8(Занято.(float64)) > 0)

			if uint8(Занято.(float64)) > 0 {

				return false, СтатусСервиса{
					Код:   ЛогинЗанят,
					Текст: "Логин занят",
				}
			} else {
				return true, СтатусСервиса{
					Код:   Ок,
					Текст: "Логин свободен",
				}
			}
		} else if len(массивДанных) > 1 {
			Ошибка(" Возвращено более 1 записи %+v \n", массивДанных)
			return false, СтатусСервиса{
				Код:   ЛогинЗанят,
				Текст: "Логин занят",
			}
		}
	}

	Инфо(" %+s %+v \n", ответ, статус)

	return true, СтатусСервиса{
		Код:   Ок,
		Текст: "Логин свободен",
	}
}

func EmailСвободен(email string) (bool, СтатусСервиса) {

	ответ, статус := База.Получить(ДанныеЗапроса{
		Запрос: `
			query checkEmail ($email: string){
				Emails(func: eq(email, $email)) {
					<занято>: count(uid)
				}
			}`,
		Данные: map[string]string{
			"$email": email,
		},
	})
	Инфо(" %+s %+v \n", ответ, статус)

	if статус.Код != Ок {
		Инфо(" %+v \n", статус.Текст)
		return false, СтатусСервиса{
			Код:   статус.Код,
			Текст: статус.Текст,
		}
	}

	картаОтвета := ОтветИзБазы{}
	ошибкаРазбора := ИзJson(ответ, &картаОтвета)
	if ошибкаРазбора != nil {
		Ошибка(" Не удалось разобрать ответ %+v \n", ошибкаРазбора.Error())
	}

	for имя, массивДанных := range картаОтвета {
		if len(массивДанных) == 1 {
			Инфо(" %+v  %+v \n", имя, массивДанных)
			Занято := массивДанных[0]["занято"]
			Инфо(" %+v  %+v \n", Занято, uint8(Занято.(float64)) > 0)

			if uint8(Занято.(float64)) > 0 {

				return false, СтатусСервиса{
					Код:   EmailЗанят,
					Текст: "Email занят",
				}
			} else {
				return true, СтатусСервиса{
					Код:   Ок,
					Текст: "Email свободен",
				}
			}
		} else if len(массивДанных) > 1 {
			Ошибка(" Возвращено более 1 записи %+v \n", массивДанных)
			return false, СтатусСервиса{
				Код:   EmailЗанят,
				Текст: "Email занят",
			}
		}
	}

	Инфо(" %+s %+v \n", ответ, статус)

	return true, СтатусСервиса{
		Код:   Ок,
		Текст: "Email свободен",
	}
}
func ПримерРегистрацииКлиента() {
	НовыйКлиент := ДанныеКлиента{
		ИдКлиента: uuid.New(),
		Логин:     "aaaaa",
		Пароль:    string("password"),
		Email:     "aaaa@ya.ru",
		Роль:      []string{"клиент"},
		Права:     []string{"чтение", "просмотр", "изменение своего"},
	}
	новыйКлиентаСтрока, ошибка := Json(НовыйКлиент)
	if ошибка != nil {
		Ошибка(" сериализаии нового клиента  %+v  новыйКлиент %+v \n", ошибка, НовыйКлиент)
	}
	// Инфо(" новыйКлиентаСтрока %+s \n", новыйКлиентаСтрока)

	// le меньше или равно
	// lt меньше, чем
	// ge больше или равно
	// gt больше, чем
	// регистраци := `upsert{
	// 	query {
	// 	  v as var(func: eq(<логин>, "anima"))
	// 	}
	// 	mutation @if(lt(len(v), 1)) {
	// 	  "set": {
	// 		"логин":"login"
	// 	  }
	// 	}
	//   }`
	// Инфо("регистраци %+v \n", регистраци)
	// м := Мутация{
	// 	Условие: "@if(lt(len(v), 1))",
	// 	Мутация: новыйКлиентаСтрока,
	// }

	ctx := context.Background()
	var транзакция Транзакция
	транзакция.Txn = База.Граф.NewTxn()
	defer транзакция.Discard(ctx)
	данные := ДанныеЗапроса{
		Запрос: `query User($login : string, $email : string) {
			 				logins(func: eq(<логин>, $login)){								
								loginCount as <логин>
							
			 				}	
							emails(func: eq(<email>, $email)){										
								email as <email>		
			 				}				
			 			}
				`,
		Мутация: []Мутация{
			{
				Условие: "@if(lt(len(loginCount), 1) AND lt(len(email), 1))",
				Мутация: новыйКлиентаСтрока,
			},
		},
		Данные: map[string]string{
			"$login": НовыйКлиент.Логин,
			"$email": НовыйКлиент.Email,
		},
	}

	ответ, статус := транзакция.Измененить(данные, ctx)
	// Инфо(" ответи %+s  ( статусИзменения %+v ) данные %+s \n", ответ, статусИзменения, данные)

	if статус.Код != Ок {
		Ошибка(" ответ %+v \n", ответ)
		// Инфо("СуществующиеЗаписи %+v \n", СуществующиеЗаписи)
		if ответ["logins"] != nil {
			Ошибка(" логин занят %+v \n", ответ["logins"])

		}
		if ответ["emails"] != nil {
			Ошибка(" email занят %+v \n", ответ["emails"])
		}
	}

	данные = ДанныеЗапроса{
		Запрос: `query User($login : string) {
					us(func: eq(<логин>, $login)){				
						<ид_клиента> : <uid>
						<логин>
						<email>
						<роль>
						<права>						
					}				
				}`,
		Данные: map[string]string{
			"$login": НовыйКлиент.Логин,
		},
	}

	данныеПользователя, статусПолучения := транзакция.Получить(данные)
	Инфо(" данныеПользователя %+s  статусПолученияия %+v  данные %+v \n", данныеПользователя, статусПолучения, данные)
	транзакция.Commit(ctx)
}

// ИНФО: Получается dgraph  не умеет поддерживать уникальные узлы, нужно самому контролировате перед записью, тоесть каждая запись требующая уникальоного контроля, должна быть предварительно проверена что такого узла нет, и только потом добавлена или обновлена. тоесть оснвоной механизм UPSERT
// проверяем существует ли узел с полем которое для нас должно быть уникальным, если его нет то добавляем узел, иначе обновляем.
var СхемаБазы = `<маршрут>: string @index(exact) @upsert .
<обработчик>:  string @index(exact) @upsert .
<действие>:  string @index(exact) @upsert .
<описание>: string .
<шаблонизатор>: [uid] .
<имя_шаблона>: string .
<код>: int .
<номер_в_очереди> : int .
<ассинхронно>: bool  .
<доступ>: [uid] .	
<пользователи>: [uid] . 
<роль>: [string] . 
<права>: [string] .
<дата_создания>: dateTime  .	
<очередь_обоработчиков>: [uid] .
<создатель>: uid .
			type <Шаблон>{
				<код>
				<имя_шаблона>
				<доступ>
			}
			type <ПраваДоступа> {
				<пользователи>
				<роль>
				<права>
			}
			type <Обработчик> {
				<маршрут>
				<действие>
				<обработчик>
				<права>
				<описание>
				<шаблонизатор>
				<ассинхронно>
			}											
			type <ОчередьОбработчиков> {
					<маршрут>
					<очередь_обоработчиков>
					<дата_создания>
					<создатель>
			}						
							<количество_неудачных_попыток_входа> : int .
							<количество_удачных_попыток_входа> : int .
							<aутентифицирован>: bool .
							<ид_клиента>: string @index(exact) @upsert .
							<права_доступа>: [string] .
							<секрет> : string .
							<имя>: string  .
							<фамилия>: string  .
							<отчество>: string  .
							<логин>: string @index(exact) @upsert .
							<пароль>: password .
							email: string @index(exact) @upsert .
							<телефон>: string @index(exact)  @upsert .
							<создан>: datetime .
							<обновлен>: datetime  .
							<статус>: string  .
							<аватар>: string .
							<о_себе>: string .
							<социальные_ссылки>: [string] .	
							<адрес>: uid .
							<страна>: string  .
							<город>: string  .
							<район>: string  .
							<тип_улицы>: string  .
							<название_улицы>: string  .
							<номер_дома>: string  .
							<корпус>: string  .
							<номер_квартиры>: string  .
							jwt: string .
									type <Секрет> {	
											<ид_клиента> 
											<секрет> 
											<обновлен> 							
										}
										type <Адрес> {
												<страна>
												<город>
												<район>
												<тип_улицы>
												<название_улицы>
												<номер_дома>
												<корпус>
												<номер_квартиры>
										}
										type <Пользователь> {
												<aутентифицирован>
												<количество_неудачных_попыток_входа>
												<количество_удачных_попыток_входа>
												<ид_клиента>
												<имя>
												<фамилия>
												<отчество>
												<логин>
												<пароль>
												email
												<телефон>
												<адрес>
												<права_доступа>
												<создан>
												<обновлен>
												<статус>
												<аватар>
												<о_себе>
												<социальные_ссылки>							
												jwt
										}`

// var демоНоваяЗапись = `{"set":[{"маршрут": "маршрут3",
// 							"доступ":[{
// 								"роль" : "администратор",
// 								"права": ["чтение", "изменение","удаление"],
// 						"описание":"описание достпа",
// 								"dgraph.type": "Доступ"
// 				},{
// 								"роль" : "клиент",
// 								"права": ["чтение"],
// 								"dgraph.type": "Доступ"
// 							}],
// 							"описание": "описание маргрута",
// 							"dgraph.type": "Маршрут"
// 				}]}`

var Schema = `
	login: string @index(exact) @upsert .
	password: password .
	authorized: bool .
	borken_authorization: int .
	type Uuser {
		login
		password
		authorized
		borken_authorization
	}
`

func DgraphConnect() (*dgo.Dgraph, func()) {
	conn, err := grpc.Dial("localhost:9080", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}

	dc := api.NewDgraphClient(conn)
	return dgo.NewDgraphClient(dc), func() {
		if err := conn.Close(); err != nil {
			log.Printf("Error while closing connection:%v", err)
		}
	}

}

func InitSchema() {
	dg, cancel := DgraphConnect()
	defer cancel()
	err := dg.Alter(context.Background(), &api.Operation{
		Schema: Schema,
	})
	if err != nil {
		panic(err)
	}
}

func Auth() {
	dg, cancel := DgraphConnect()
	defer cancel()
	ctx := context.Background()
	req := &api.Request{
		CommitNow: true,
		Query: `query User($login: string, $pass: string) {
						var(func: eq(login, $login)){								
							status as checkpwd(password, $pass)
						}	
						User(func: eq(val(status), 1)) {
							login
							email
							authorized : val(status)
						}	
						BrokenAuth(func: eq(val(status), 0)) {
							login							
							authorized : val(status)
							countBroken as borken_authorization
							countBrokenAuth as math(countBroken+1)

						}				
					}
			`,
		Mutations: []*api.Mutation{
			{
				Cond: "@if(eq(val(status), 1))",
				SetJson: []byte(`{
					"uid": "uid(status)",
					"authorized": "true",
					"borken_authorization": 0
					}`),
			},
			{
				Cond: "@if(eq(val(status), 0))",
				SetJson: []byte(`{
					"uid": "uid(status)",
					"authorized": true,
					"borken_authorization": "val(countBrokenAuth)"
					}`),
			},
		},
		Vars: map[string]string{
			"$login": "UserName",
			"$pass":  "Password",
		},
	}
	fmt.Printf(" req %+v \n", req)
	res, err := dg.NewTxn().Do(ctx, req)

	fmt.Printf("%+v \n", res)
	if err != nil {
		fmt.Printf("err %+v; res %+v \n", err, res)
	}
}

func SetClient() {
	dg, cancel := DgraphConnect()
	defer cancel()
	ctx := context.Background()
	req := &api.Request{
		CommitNow: true,
		Query: `query User($login : string, $email : string) {
						logins(func: eq(login, $login)){								
							loginCount as login
						}	
						emails(func: eq(email, $email)) {
							emailCount as email		
						}				
					
					}
			`,
		Mutations: []*api.Mutation{
			{
				Cond: "@if(lt(len(loginCount), 1) OR lt(len(emailCount), 1))",
				SetJson: []byte(`{
					"borken_authorization": 0,
					"login": "UserName",
					"email": "userName@mail.com",
					"password": "Password"
					}`),
			},
		},
		Vars: map[string]string{
			"$login": "UserName",
			"$email": "userName@mail.com",
		},
	}
	fmt.Printf(" req %+v \n", req)
	res, err := dg.NewTxn().Do(ctx, req)

	fmt.Printf("%+v \n", res)
	if err != nil {
		fmt.Printf("err %+v; res %+v \n", err, res)
	}

}

// query User($login : string, $email : string) {
// 	logins(func: eq(login, $login)){\
// 		loginCount as login
// 		}
// 		emails: eq(email, $email)) {
// 			emailCount as email
// 			}
// 		}
// 		" vars:<key:"$email" value:"userName@mail.com" > vars:<key:"$login" value:"UserName" >
// 		mutations:<set_json:"{
// 			"login": "UserName",
// 			"email": "userName@mail.com",
// 			"password": "Password"
// 			}">
// 		cond:"@if(lt(len(loginCount), 1) AND lt(len(email), 1))"
