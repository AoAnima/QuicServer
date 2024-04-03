package main

import (
	"net/url"
	"strings"

	. "aoanima.ru/ConnQuic"
	. "aoanima.ru/DGApi"
	. "aoanima.ru/Logger"
	. "aoanima.ru/QErrors"
	"github.com/quic-go/quic-go"
)

/*
Очередь обработчиков
Каждый запрос от клиента может быть обработан одним и более количеством микросервисов, для того чтобы правильно отправлять запрос в сервисы нужно описать последовательность обработки запроса Сервиссами, и указать какой HTNL шаблон рендерить, или куда сделать редирект например после автоизации или регистрации)

КРоме того необходимо учитывать права доступа и роли пользователя, чтобы один и тот же маршрут по разному обрабатывался в заивсимости от роли польтзователя и его прав доступа.

дествие  | сервис | маршрут | роль | права | шаблон |  статусОтвета | редирект | ассинхронно |

	| Рендер | /формаРеситрации | ["гость"] | "формаРегистрации" | Ок | /личныйКабинет | нет

регистрация  | Авторизация | /формаРеситрации (url же не меняется) | ["гость"] |   |  |
*/

/*
Обработчик может указываться если значеине в поле действие называетя иначе чем функция которая должна обработать запрос
например действие: регистрация
Обработчик : сохранитьПользователя
*/
type КонфигурацияОбработчика struct {
	Маршрут      string         `json:"маршрут,omitempty"`
	Действие     string         `json:"действие,omitempty"`
	Обработчик   string         `json:"обработчик,omitempty"`
	ПраваДоступа []ПраваДоступа `json:"права"`
	Описание     string         `json:"описание,omitempty"`
	Шаблон       string         `json:"шаблон,omitempty"`
	Ассинхронно  bool           `json:"ассинхронно"`
}

type ПраваДоступа struct {
	Логин []string `json:"логин"`
	Роль  []string `json:"роль"`
	Права []string `json:"права"`
}

/*
Добавляет данные об обработчике в БД, особенно важно права доступа
Маршрут может быть пустой если есть действие, и наоборот, если есть маршрут а обработчика нету. не страшно, обработчик будет вычисляться из маршрута.
если есть Действие то оно в приоритете
*/
func ДобавитьОбработчик(поток quic.Stream, сообщение Сообщение) {

	ответ := сообщение.Ответ[Сервис]
	ответ.Сервис = Сервис
	ответ.ЗапросОбработан = true

	маршрут, статусМаршрут := ПолучитьЗначениеПоляФормы("маршрут", сообщение.Запрос.Форма)

	действие, статусДейсвтия := ПолучитьЗначениеПоляФормы("действие", сообщение.Запрос.Форма)

	if статусДейсвтия.Код != Ок && статусМаршрут.Код != Ок {
		ответ.СтатусОтвета = СтатусСервиса{
			Код:   Прочее,
			Текст: "действие и маршрут не заданы, должно быть установлено одно или оба поля",
		}
		сообщение.Ответ[Сервис] = ответ
		ОтправитьСообщение(поток, сообщение)
	}

	//
	обработчик, статусОбработчик := ПолучитьЗначениеПоляФормы("обработчик", сообщение.Запрос.Форма)
	if статусОбработчик.Код != Ок {
		Ошибка(" статус получения обработчика  %+v \n", статусОбработчик)
	}
	роль, статусРоль := ПолучитьВсеЗначенияПоляФормы("роль", сообщение.Запрос.Форма)
	if статусРоль.Код != Ок {
		Ошибка(" статус получения роли  %+v \n", статусРоль)
	}

	права, статусПрав := ПолучитьВсеЗначенияПоляФормы("права", сообщение.Запрос.Форма)
	if статусПрав.Код != Ок {
		Ошибка(" статус получения прав  %+v \n", статусПрав)
	}

	описание, статусОписание := ПолучитьЗначениеПоляФормы("оисание", сообщение.Запрос.Форма)
	if статусОписание.Код != Ок {
		Ошибка(" статус получения описания  %+v \n", статусОписание)
	}

	шаблон, статусШаблон := ПолучитьЗначениеПоляФормы("шаблон", сообщение.Запрос.Форма)
	// ассинхронно, статусАссинхронно := ПолучитьЗначениеПоляФормы("ассинхронно", сообщение.Запрос.Форма)
	if статусШаблон.Код != Ок {
		Ошибка(" статус получения шаблона  %+v \n", статусШаблон)
	}

	новыйОбработчик := &КонфигурацияОбработчика{
		Маршрут:    маршрут,
		Действие:   действие,
		Обработчик: обработчик,
		ПраваДоступа: []ПраваДоступа{
			{
				Роль:  роль,
				Права: права,
			},
		},
		Описание: описание,
		Шаблон:   шаблон,
		// Ассинхронно: ассинхронно,
	}

	обработчикБин, статус := Кодировать(новыйОбработчик)
	if статус != nil {
		Ошибка(" статус %+v \n", статус)
	}
	данные := ДанныеЗапроса{
		Запрос: `query <Обработчики>($hanlder : string, $path : string, $action : string) {
							<Обработчик>(func: eq(<обработчик>, $hanlder)) @filter(eq(<маршрут>, $path) AND eq(<действие>, $action)) {
								<маршрут> @filter(eq(<маршрут>, $path)) {
									<маршрут>
								}
								<действие> @filter(eq(<маршрут>, $action)) {
									<действие>
								}
							}		
			 			}
				`,
		Мутация: []Мутация{
			{
				Условие: "@if(lt(len(loginCount), 1) AND lt(len(email), 1))",
				Мутация: обработчикБин,
			},
		},
		Данные: map[string]string{
			"$login": НовыйКлиент.Логин,
			"$email": НовыйКлиент.Email,
		},
	}
	сохранитьОбработчик := ДанныеЗапроса{
		Мутация: обработчикБин,
	}

	База.Изменить(сохранитьОбработчик)

	// if статус.Код != Ок {
	// 	Ошибка(" статус.Текст %+v \n", статус.Текст)
	// // 	ответ.СтатусОтвета = статус
	// // 	ответ.Данные = map[string]bool{
	// // 		"МаршрутДобавлен": false,
	// // 	}
	// // 	сообщение.Ответ[Сервис] = ответ
	// // 	ОтправитьСообщение(поток, сообщение)
	// // 	return
	// }
	// доступJson := Json(доступ)
	// добавить := `[
	// 		{
	// 			"<маршрут>": "` + маршрут + `",
	// 			"<доступ>": {
	// 				<>
	// 			},
	// 			"<описание>": "",
	// 			"dgraph.type": "Vfhihen"
	// 		}
	// 	]`
	// ответиз, статусИзменения := База.Изменить(ДанныеЗапроса{
	// 	Запрос: добавить,
	// })
	// if статусИзменения.Код != Ок {
	// 	Ошибка(" ОписаниеОшибки %+v \n", статусИзменения.Текст)
	// }
	// Инфо("  %+v \n", ответиз, статусИзменения)

}
func ИзменитьОбработчик(поток quic.Stream, сообщение Сообщение) {

}
func УдалитьОбработчик(поток quic.Stream, сообщение Сообщение) {

}
func ИзменитьОчередьОбработчиков(поток quic.Stream, сообщение Сообщение) {

}
func ДобавитьМаршрут(поток quic.Stream, сообщение Сообщение) {

}
func ИзменитьМаршрут(поток quic.Stream, сообщение Сообщение) {

}
func УдалитьМаршрут(поток quic.Stream, сообщение Сообщение) {

}

func ПолучитьОчередьОбработчиков(поток quic.Stream, сообщение Сообщение) {

	маршрутЗапроса, err := url.Parse(сообщение.Запрос.МаршрутЗапроса)
	Инфо(" маршрутЗапроса %+v \n", маршрутЗапроса)

	if err != nil {
		Ошибка("Parse маршрутЗапроса: ", err)
	}
	маршрутЗапроса.Path = strings.Trim(маршрутЗапроса.Path, "/")
	urlКарта := strings.Split(маршрутЗапроса.Path, "/")

	/*
		Для получения очереди обработчиков нужно проанализировать url и данные формы
		если метод post то аналиируем форму
		если метод get то анализируем Сообщение.ЗАпрос.СтрокаЗапроса содержащую Query часть
		если там не передан параметр "действие" то ищем обработчик из path
	*/
	if сообщение.Запрос.ТипЗапроса == GET || сообщение.Запрос.ТипЗапроса == AJAX {
		// анализируем url параметры
		параметрыЗапроса := маршрутЗапроса.Query()
		дейсвтие, естьДействие := параметрыЗапроса["действие"]
		if естьДействие {
			// получить очередь из БД
			Инфо("получить очередь из БД для: %+v \n", дейсвтие)

		} else {
			if len(urlКарта) > 0 {
				/*
					может пройти по всем частам url и получить очереь обрабочиковдля каждого шага ?
					или брать только первый ?
				*/
				дейсвтие := urlКарта[0]
				Инфо("получить очередь из БД для: %+v \n", дейсвтие)

			}
		}

	}

	if сообщение.Запрос.ТипЗапроса == POST || сообщение.Запрос.ТипЗапроса == AJAXPost {

	}

}

func ПолучитьСписокОчередей(поток quic.Stream, сообщение Сообщение) {

}
