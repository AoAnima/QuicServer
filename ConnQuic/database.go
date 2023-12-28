package ConnQuic

import (
	"strings"
	"time"

	. "aoanima.ru/Logger"
	badger "github.com/dgraph-io/badger/v4"
)

type База struct {
	*badger.DB
}

type Транзакция struct {
	*badger.Txn
}

type ОшибкаБазы struct {
	Ключ     string
	Индекс   string
	Значение []byte
	Текст    string
	Код      int
}

const (
	ОшибкаДубликатКлюча = iota + 1
	ОшибкаДубликатИндекса
	ОшибкаСоединениеЗакрыто
	ОшибкаФиксацииТранзакции
	ОшибкаЗаписи
)

func (подключение *База) Подключиться(путь string) error {

	база, err := badger.Open(badger.DefaultOptions(путь))
	if err != nil {
		Ошибка("  %+v \n", err)
	}
	подключение.DB = база
	return err
}

/*
Вставить
таблица - это аналог таблицы (клиент, магазин, товар, заказ, адрес... )
ключ это имя поля в объекте, значение которого будет использоваться в качестве ключа документа в формате:
 коллекция:значениеКлючевогоПоля = объект
 клиент.uid1 = объект

 индексы это поля п которым будут созданы индексы:
  имяПоля.имяВложенногоПоля будет преобразован в
  коллекция.имяПоля.имяВложенногоПоля:значениеПоля = ключОбъекта(имяКлючевогоПоля:значениеКлючевогоПоля)


  тогда мы сможем искать например

  клиент.email = ...
  или
 если создали индекс по городу: клиент.адрес.город:Ставрополь =  ключОбъекта(имяКлючевогоПоля:значениеКлючевогоПоля)
 то сможем найти  клиент.адрес.город = Ставрополь



 Нужно переделать индексацию, и сигнатуру функции, чотбы можно было передавать список обехктов и ключей

*/

// func (подключение *База) Вставить_(таблица string, ключ string, объект interface{}, индекс []string) ОшибкаВставки {

// 	Индексы := make(map[string]string)
// 	if индекс != nil {
// 		Индексы, _ = СоздатьКлючиИндекса(объект, индекс)
// 	}

// 	err := подключение.Update(func(txn *badger.Txn) error {
// 		старт := time.Now()
// 		ok, дубльКлюча := КлючСвободен(таблица, ключ, Индексы, Транзакция{txn})
// 		времяРаботы := time.Since(старт)
// 		Инфо("Время поиска ключа и индекса  %+v \n", времяРаботы)

// 		if !ok {
// 			return дубльКлюча
// 		}
// 		Объект, err := Json(объект)
// 		if err != nil {
// 			Ошибка("  %+v \n", err)
// 		}

// 		КлючОбъекта := таблица + ":" + ключ

// 		if err := txn.Set([]byte(КлючОбъекта), Объект); err != nil {
// 			return err
// 		}

// 		if len(Индексы) > 0 {
// 			for ключИндекса, значениеИндекса := range Индексы {
// 				if err := txn.Set([]byte(таблица+"."+ключИндекса+":"+значениеИндекса), []byte(КлючОбъекта)); err != nil {
// 					return err
// 				}
// 			}
// 		}
// 		return nil
// 	})
// 	if err != nil {
// 		Ошибка("  %+v \n", err)
// 	}
// 	return err
// }

func (подключение *База) Вставить(таблица string, ключ string, объект interface{}, индекс []string) ОшибкаБазы {

	Индексы := make(map[string]string)
	if индекс != nil {
		Индексы, _ = СоздатьКлючиИндекса(объект, индекс)
	}

	ошибка := подключение.ТранзакцияЗаписи(func(трз *Транзакция) ОшибкаБазы {
		старт := time.Now()
		ok, дубльКлюча := КлючСвободен(таблица, ключ, Индексы, трз)
		времяРаботы := time.Since(старт)
		Инфо("Время поиска ключа и индекса  %+v \n", времяРаботы)

		if !ok {
			return дубльКлюча
		}
		Объект, err := Json(объект)
		if err != nil {
			Ошибка("  %+v \n", err)
		}

		КлючОбъекта := таблица + ":" + ключ

		if err := трз.Set([]byte(КлючОбъекта), Объект); err != nil {
			return ОшибкаБазы{
				Текст: err.Error(),
				Код:   ОшибкаЗаписи,
			}
		}

		//  Получается при каждой вставке  будем проверять таблицу индексов, есть, нету и обновлять... накладно, нужно вынести проверку в отдельную функцию
		if len(Индексы) > 0 {
			for ключИндекса, значениеИндекса := range Индексы {
				if err := трз.Set([]byte(таблица+"."+ключИндекса+":"+значениеИндекса), []byte(КлючОбъекта)); err != nil {
					return ОшибкаБазы{
						Текст: err.Error(),
						Код:   ОшибкаЗаписи,
					}
				}
				// Добавим в справочник индексов, ключ индексируемого поля, чтобы можно было понять есть ли по искомому запросу индекс, дабы не перебирать все записи
				if документИндексов, err := трз.Get([]byte(таблица + ".индексы")); err != nil && err != badger.ErrKeyNotFound {
					var картаИндексов ТаблицаИндексов
					таблицаИндексов, err := документИндексов.ValueCopy(nil)
					if err != nil {
						Ошибка("  %+v \n", err)
					}
					err = ИзJson(таблицаИндексов, &картаИндексов)
					if err != nil {
						Ошибка("  %+v \n", err)
					}
					if _, ok := картаИндексов[таблица+"."+ключИндекса]; ok {
						Инфо("индекс уже существует, ничего не добавляем  %+v \n", картаИндексов)
						continue
					}
					// ключ индекса ещё не был добавлен, добавляем
					картаИндексов[таблица+"."+ключИндекса] = struct{}{}

					таблицаИндексов, err = Json(картаИндексов)
					if err != nil {
						Ошибка("  %+v \n", err)
					}
					if err := трз.Set([]byte(таблица+".индексы"), таблицаИндексов); err != nil {
						return ОшибкаБазы{
							Текст: err.Error(),
							Код:   ОшибкаЗаписи,
						}
					}
				} else {
					картаИндексов := make(ТаблицаИндексов)
					картаИндексов[таблица+"."+ключИндекса] = struct{}{}
					таблицаИНдексов, err := Json(картаИндексов)
					if err != nil {
						Ошибка("  %+v \n", err)
					}
					if err := трз.Set([]byte(таблица+".индексы"), таблицаИНдексов); err != nil {
						return ОшибкаБазы{
							Текст: err.Error(),
							Код:   ОшибкаЗаписи,
						}
					}
				}
			}
		}
		return ОшибкаБазы{}
	})
	if ошибка.Код != 0 {
		Ошибка("  %+v \n", ошибка)
	}
	return ошибка
}

/*
Удалить
Нужно удалять не только ключ, но и индексы связанные с этой записью.

получаем список индексов таблца.индексы
удаляем все индексы связаные с указаным ключём
таблица.индекс:значение

итерируемся по каждому индексу , находим значение в удаляемом докумнте по пути в индексе, адялем запись...
*/
func (подключение *База) Удалить(таблица string, ключ string) ОшибкаБазы {
	if подключение.IsClosed() {
		return ОшибкаБазы{
			Текст: "соединение с базой закрыто",
			Код:   ОшибкаСоединениеЗакрыто,
		}
	}
	подключение.ТранзакцияЗаписи(func(трз *Транзакция) ОшибкаБазы {

		// 1. получаем список индексов таблца.индексы
		// 2. проверяем есть ли запись с ключём таблица:ключ
		// 	если нет то итерируемся по индексамтрз
		// 	проверяем записи в индексах: таблица.индекс:ключ
		// 		если запись найден проверяем есть ли запись с ключём таблица:[таблица.индекс:ключ]
		// 		всё что нашли удаляем

	})

}

func (подключение *База) Изменить(таблиц string, ключ string, изменяемыйПуть string, новоеЗначение any) ОшибкаБазы {
	if подключение.IsClosed() {
		return ОшибкаБазы{
			Текст: "соединение с базой закрыто",
			Код:   ОшибкаСоединениеЗакрыто,
		}
	}
	подключение.ТранзакцияЗаписи(func(трз *Транзакция) ОшибкаБазы {

		// 1. проверяем существует ли обект с ключём таблица:ключ
		// 2. получаем список индексов таблца.индексы
		// 	если обект не существует  то итерируемся по индекса
		// 	проверяем записи в индексах: таблица.индекс:ключ
		// находим обеъкт по любому индексу, проверяем существует ли изменяемыйПуть:
		// 	если есть то меняем значение
		// 	если нет то добавляем новый путь в обхект

	})
}

type ТаблицаИндексов map[string]struct{}

func (подключение *База) ТранзакцияЗаписи(функ func(трз *Транзакция) ОшибкаБазы) ОшибкаБазы {

	if подключение.IsClosed() {
		return ОшибкаБазы{
			Текст: "соединение с базой закрыто",
			Код:   ОшибкаСоединениеЗакрыто,
		}
	}

	трз := подключение.NewTransaction(true)
	defer трз.Discard()

	if ошибкаВставки := функ(&Транзакция{трз}); ошибкаВставки.Код != 0 {
		return ошибкаВставки
	}
	err := трз.Commit()
	if err != nil {
		Ошибка("  %+v \n", err)
		return ОшибкаБазы{
			Текст: "не удалось зафиксировать транзакцию, причина: " + err.Error(),
			Код:   ОшибкаФиксацииТранзакции,
		}
	}

	return ОшибкаБазы{}
}

/*
КлючСвободен - проверяет существует ли в базе ключ и индексы

принимает теже аргументы что и функция Вставить,
-таблица,
-первичный ключ для записи данных,
-индексы созданные функцией СоздатьКлючиИндекса
- транзакция
если ключ и индексы не найдены то возращает true и пустой объект ОшибкаДублированиеКлюча
Иначе false и ОшибкаДублированиеКлюча с найденным индексом и значением записанным в этот индекс
*/
func КлючСвободен(таблица string, ключ string, индексы map[string]string, тр *Транзакция) (bool, ОшибкаБазы) {
	КлючОбъекта := таблица + ":" + ключ
	найденноеЗначение, err := тр.Get([]byte(КлючОбъекта))

	if err == badger.ErrKeyNotFound {
		// Ключа нет, запись производится
		if len(индексы) > 0 {
			for ключИндекса, значениеИндекса := range индексы {
				if индексированноеЗначение, err := тр.Get([]byte(таблица + "." + ключИндекса + ":" + значениеИндекса)); err != nil {
					if err == badger.ErrKeyNotFound {
						// индекса нет всё ок, проверяем дальше
					} else {
						// индекс  существует не нужно записывать
						данные, err := индексированноеЗначение.ValueCopy(nil)
						if err != nil {
							Ошибка("  %+v \n", err)
						}
						return false, ОшибкаБазы{
							Индекс:   таблица + "." + ключИндекса + ":" + значениеИндекса,
							Значение: данные,
							Текст:    "индекс существует",
							Код:      ОшибкаДубликатИндекса,
						}
					}
				}
			}
		}
		return true, ОшибкаБазы{}
	} else if err != nil {
		if err != nil {
			Ошибка("  %+v \n", err.Error())
		}
		return false, ОшибкаБазы{
			Текст: err.Error(),
		}
	} else {
		данные, err := найденноеЗначение.ValueCopy(nil)
		if err != nil {
			Ошибка("  %+v \n", err)
		}
		// Ключ существует, запись не производится
		return false, ОшибкаБазы{
			Ключ:     КлючОбъекта,
			Значение: данные,
			Текст:    "ключ существует",
			Код:      ОшибкаДубликатКлюча,
		}
	}
}

func СоздатьКлючиИндекса(объект interface{}, индекс []string) (map[string]string, error) {

	копияОбъекта := объект
	ЗначенияИндекса := make(map[string]string)

	for _, ключИндекса := range индекс {

		ПутьКЗначению := strings.Split(ключИндекса, ".")

		for _, ключ := range ПутьКЗначению {
			if значение, ok := копияОбъекта.(map[string]interface{})[ключ]; ok {
				копияОбъекта = значение
			} else {
				// return nil, errors.New(fmt.Sprintf("нет такого ключа %+v индекс %+v", ключ, индекс))
				Ошибка("нет такого ключа %+v индекс %+v", ключ, индекс)
				копияОбъекта = nil
			}
		}
		if копияОбъекта != nil {
			switch копияОбъекта.(type) {
			case []byte:
				ЗначенияИндекса[ключИндекса] = копияОбъекта.(string)
			case string:
				ЗначенияИндекса[ключИндекса] = копияОбъекта.(string)
			default:
				строка, err := Json(копияОбъекта)
				if err != nil {
					Ошибка("  %+v \n", err)
				}
				ЗначенияИндекса[ключИндекса] = string(строка)

			}
		}
	}
	return ЗначенияИндекса, nil
}

/*
Найти
-таблица - это аналог таблицы (клиент, магазин, товар, заказ, адрес... )
-значение -
- индекс это индексируемое поле, или поле внутри объекта, если не было проиндексированно то нужно будет итерироваться по всей таблице, в искомом объекте, если не задан (="") значит не требуется посик по индексу

если задан то ищем по индекса, и затем ищем по значению индекса.
-
*/

func (подключение *База) Найти(таблица string, значение string, индекс string) (interface{}, error) {

	трз := Транзакция{подключение.NewTransaction(false)}
	defer трз.Discard()

	var ключПоиска string
	if индекс != "" {
		ключПоиска = таблица + "." + индекс + ":" + значение
	} else {
		ключПоиска = таблица + ":" + значение
	}
	найденоеЗначение, err := трз.Get([]byte(ключПоиска))
	if err != nil {
		Ошибка(" не нашли ничего по ключу %+v err  %+v \n", ключПоиска, err.Error())
	}

	байтЗначение, err := найденоеЗначение.ValueCopy(nil)
	if err != nil {
		Ошибка("  %+v \n", err)
	}

	var результат interface{}
	err = ИзJson(байтЗначение, &результат)
	// нашли значение, если индекс не пустой, значит нашли ключ в индексе, получим объект по значению индекса.
	//

	if индекс != "" {
		найденоеЗначение, err := трз.Get([]byte(ключПоиска))
		if err == badger.ErrKeyNotFound {
			// Индекса не существует
			/*
				Проверим какие индексы есть в таблице, если в указаный индекс есть , и значение не найдено то всё ок, записи прост онет. Иначе если индекса нет, то нужно итерироваться по всей таблице и искать внутри документа
			*/
			индексыТаблицы, err := трз.Get([]byte(таблица + "." + "индексы"))
			if err == badger.ErrKeyNotFound {
				Ошибка(" индексы не существуют для   %+v %+v ?  будем перебирать все документы\n", таблица+"."+"индексы", err.Error())
			}
			копияДанных, err := индексыТаблицы.ValueCopy(nil)
			if err != nil {
				Ошибка("  %+v \n", err)
			}
			префиксИндекса := НайтиВJson(копияДанных, таблица+"."+индекс)
			if префиксИндекса != nil {
				// префиксИндекса если индекс есть в списке индексов таблицы, то т.к. ранее мы не нашли значение по индексу. значит такого документа нету, возращаем ответ что документа нету, потому что инчае документ был бы проиндексирован
				return nil, nil
			} else {
				// индекса нету в списке индексов таблицы, значит нужно итерироваться по всем документам и искать внутри документа
				найденныеДанные, err := трз.НайтиЗначение(таблица, значение, индекс)
				if err != nil {
					Ошибка("  %+v \n", err)
				}
				return найденныеДанные, err
			}

		} else if err != nil {
			Ошибка("  %+v \n", err.Error())
			return nil, err
		}

		байтЗначение, err := найденоеЗначение.ValueCopy(nil)
		if err != nil {
			Ошибка("  %+v \n", err)
			return nil, err
		}

		err = ИзJson(байтЗначение, &результат)
	}
	return результат, err
}

/*
НайтиЗначение - ищет значение по указаному пути внутри документа, дессериализация не происходит, по идее jsoniter умеет искать внутри бинарного представления
таблица - префикс по которому будем итерировать документы
путь - путь до значения внутри объекта
*/
func (трз *Транзакция) НайтиЗначение(таблица string, значение string, путь string) (map[string][]byte, error) {
	итератор := трз.NewIterator(badger.DefaultIteratorOptions)

	defer итератор.Close()

	префикс := []byte(таблица)

	найденныеДанные := make(map[string][]byte)

	for итератор.Seek(префикс); итератор.ValidForPrefix(префикс); итератор.Next() {
		объект := итератор.Item()
		ключ := объект.Key()
		err := объект.Value(func(знач []byte) error {
			значениеВОбъекте := НайтиВJson(знач, путь)
			if значениеВОбъекте == nil {
				// продолжаем итерации по документам
			} else {
				// значение найдено, кладём ключ документа и документ в карту
				Инфо("key=%s, value=%s\n", ключ, знач)
				найденныеДанные[string(ключ)] = знач

			}
			return nil
		})
		if err != nil {

			return nil, err
		}
	}
	return найденныеДанные, nil
}

/*
ДобавитьЗначениеВМультиИндекс
МультИНдекс - индекс не уникальных значений, если часто ищу значение какого то поля, то можно создать индекс в котором будут храниться ссылки на документы содержащие нужное значение поля

например адрес.город:Ставрополь= ["клиент.ид1", "клиент.ид2", "клиент.ид3"]
и не нужно будет перебирать все документы nxj,s найти у кого в поле город есть Ставрополь
*/
func ДобавитьЗначениеВМультиИндекс() {

}