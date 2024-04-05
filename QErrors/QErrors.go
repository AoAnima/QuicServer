package QErrors

type ОшибкаБазы struct {
	Ключ     string // описываем переданые данные запроса из за котрых возникла ошибка
	Индекс   string // описываем переданые данные запроса из за котрых возникла ошибка
	Значение []byte // описываем переданые данные запроса из за котрых возникла ошибка
	Текст    string
	Код      int
}
type СтатусБазы ОшибкаБазы

type ОшибкаСервиса struct {
	Текст string
	Код   int
}
type СтатусСервиса ОшибкаСервиса

const (
	Ок     = iota
	Прочее // Любая другая ошибка
	ОшибкаДубликатКлюча
	ОшибкаДубликатИндекса
	ОшибкаКлючНеНайден
	ОшибкаСоединениеЗакрыто
	ОшибкаФиксацииТранзакции
	ОшибкаЗаписи
	ОшибкаИзмененияСхемы
	ОшибкаПодключения
	ОшибкаПреобразованияДокумента
	ОшибкаРегистрации
	ОшибкаАвторизации
	ОшибкаJSONКодирования
	ОшибкаJSONДеКодирования
	ОшибкаЗаписиВПоток
	ОшибкаПолученияДанныхИзБазы
	ЛогинЗанят
	EmailЗанят

	ПользовательНеОпознан
	НетДанныхКлиента
	ИндексыСуществуют
	БолееОдногоЗначения
	ПустоеПолеФормы
	ОшибкаПодписиJWT
	СекретНеНайден
	ОшибкаВалидацииJWT

	СрокСекретаИстекает
)
