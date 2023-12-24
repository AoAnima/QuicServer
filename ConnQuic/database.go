package ConnQuic

import (
	. "aoanima.ru/Logger"
	badger "github.com/dgraph-io/badger/v4"
	jsoniter "github.com/json-iterator/go"
)

type База struct {
	*badger.DB
}

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
коллекция - это аналог таблицы
ключ это имя поля в объекте, значение которого будет использоваться в качестве ключа документа в формате:
 коллекция:значениеКлючевогоПоля = объект

 индексы это поля п которым будут созданы индексы:
  имяПоля.имяВложенногоПоля будет преобразован в
  имяПоля.имяВложенногоПоля:значениеПоля = ключОбъекта(имяКлючевогоПоля:значениеКлючевогоПоля)


  тогда мы сможем искать например

  клиент.email = ...
  или
 если создали индекс по городу: клиент.адрес.город:Ставрополь =  ключОбъекта(имяКлючевогоПоля:значениеКлючевогоПоля)
 то сможем найти  клиент.адрес.город = Ставрополь

 соотетсвенно

*/

func (подключение *База) Вставить(коллекция string, ключ string, объект interface{}, индекс ...string) {

	значение, err := Json(объект)
	if err != nil {
		Ошибка("  %+v \n", err)
	}

	err = подключение.Update(func(txn *badger.Txn) error {
		if err := txn.Set([]byte(ключ), значение); err != nil {
			return err
		}
		if err := txn.Set([]byte("email:"+user.Email), user); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		Ошибка("  %+v \n", err)
	}
	return err
}

func НайтиВБазе(путь string) (interface{}, error) {
	db, err := badger.Open(badger.DefaultOptions("/database"))
	if err != nil {
		Ошибка("  %+v \n", err)
	}
	defer db.Close()
	тр := db.NewTransaction(false)
	defer тр.Discard()

	// Use the transaction...
	значение, err := тр.Get([]byte(путь))
	if err != nil {
		return nil, err
	}

	// Commit the transaction and check for error.
	if err := тр.Commit(); err != nil {
		return nil, err

	}
	байтЗначение, err := значение.ValueCopy(nil)
	if err != nil {
		Ошибка("  %+v \n", err)
		return nil, err
	}
	var результат interface{}
	err = jsoniter.Unmarshal(байтЗначение, результат)
	return результат, err
}
