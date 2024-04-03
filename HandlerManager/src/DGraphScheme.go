package main

// <маршрут>: string @index(exact) @upsert .
// <обработчик>: string .
// <номер_в_очереди> : int .
// <асинхронный>: bool  .
// <доступ>: [uid] .
// <пользователь>: uid .
// <роль>: string .
// <права>: [string] .
// <дата_создания>: dateTime  .
// <очередь_обоработчиков> [uid]
// <создатель> uid
// type <Доступ> {
// <логин>
// <роль>
// <права>
// }
// type <Обработчик> {
// <обработчик>
// <асинхронный>
// <номер_в_очереди>
// }
// type <Маршрут> {
// <маршут>
// <описание>
// <доступ>
// }
// type <ОчередьОбработчиков> {
// <маршут>
// <очередь_обоработчиков>
// <дата_создания>
// <создатель>
// }
/*
Доступ может складываться из 	<пользователь>+права или <пользователь>+права
*/

var СхемаБазы = `<маршрут>: string @index(exact) @upsert .
				<обработчик>: string .
				<действие>: string .
				<описание>: string .
				<шаблон>: string .
				<номер_в_очереди> : int .
				<ассинхронно>: bool  .
				<доступ>: [uid] .	
				<пользователь>: uid . 
				<роль>: string . 
				<логин>: string . 
				<права>: [string] .
				<дата_создания>: dateTime  .	
				<очередь_обоработчиков>: [uid] .
				<создатель>: uid .
							type <Доступ> {
								<логин>
								<роль>
								<права>
							}
							type <Обработчик> {
								<маршрут>
								<действие>
								<обработчик>
								<Доступ>
								<описание>
								<шаблон>
								<ассинхронно>
							}											
							type <ОчередьОбработчиков> {
									<маршрут>
									<очередь_обоработчиков>
									<дата_создания>
									<создатель>
							}							
						`

/*
						// Определение структуры для типа "ОчередьОбработчиков"
type ОчередьОбработчиков struct {
    UID     string `json:"uid,omitempty"`
    Имя     string `json:"имя,omitempty"`
    Статус  string `json:"статус,omitempty"`
}

// Создание новой записи типа "ОчередьОбработчиков"
newHandler := &ОчередьОбработчиков{
    Имя:    "Обработчик 1",
    Статус: "Активный",
}

// Добавление записи в Dgraph
resp, err := База.Мутация(ДанныеЗапроса{
    Запрос: `
        mutation {
            addОчередьОбработчиков(input: [
                {
                    имя: $имя
                    статус: $статус
                }
            ]) {
                ОчередьОбработчиков {
                    uid
                    имя
                    статус
                }
            }
        }
    `,
    Данные: map[string]interface{}{
        "имя":    newHandler.Имя,
        "статус": newHandler.Статус,
    },
})

if err != nil {
    Ошибка("Ошибка при добавлении данных: %v", err)
    return
}

Инфо("Новая запись добавлена: %+v", resp.Json)

В этом примере:
Определяется структура ОчередьОбработчиков для типа, определенного в схеме.
Создается новый экземпляр ОчередьОбработчиков с заполненными полями.
Используется метод База.Мутация() для добавления новой записи в Dgraph.
В запросе мутации используется ключевое слово addОчередьОбработчиков, которое соответствует типу, определенному в схеме.
Данные для новой записи передаются в виде параметров запроса.
Результат мутации выводится в лог.
Убедитесь, что структура ОчередьОбработчиков соответствует типу, определенному в схеме Dgraph.

*/
