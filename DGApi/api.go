package DGApi

import (
	"context"
	"log"
	"strings"

	. "aoanima.ru/Logger"
	. "aoanima.ru/QErrors"
	dgo "github.com/dgraph-io/dgo/v230"
	"github.com/dgraph-io/dgo/v230/protos/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// https://github.com/dgraph-io/dgo/blob/master/example_set_object_test.go

type ЗакрытьСоединение func()
type КаналДанных struct {
	КаналОтвет chan string
	Ошибка     chan string
	ДанныеЗапроса
}
type ДанныеЗапроса struct {
	ЗакрытьТранзакцию bool
	Запрос            string
	Мутация           []Мутация
	Данные            map[string]string
}

type Мутация struct {
	Условие string
	Мутация []byte
}

type КлиентДГраф *dgo.Dgraph

type СоединениеСДГраф struct {
	Граф         *dgo.Dgraph
	ЗакрытьДГраф func()
	// Транзакция   *dgo.Txn  НАверное не нужно сюда записывать транзакции потому что все обращаються к этому подключению и транзакция будет у всех одна...
}

// func ДГраф(каналДанных chan КаналДанны/х) {
// func ДГраф() (*dgo.Dgraph, ЗакрытьСоединение) {
func ДГраф() СоединениеСДГраф {
	// Dial a gRPC connection. The address to dial to can be configured when
	// setting up the dgraph cluster.
	связь, err := grpc.Dial("localhost:9080", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}

	dc := api.NewDgraphClient(связь)
	граф := dgo.NewDgraphClient(dc)
	// ctx := context.Background()

	// for данные := range каналДанных {
	// 	данные := данные
	// 	go func(граф *dgo.Dgraph, данные КаналДанных) {

	// 		ctx := context.Background()

	// 		мутация := &api.Mutation{
	// 			CommitNow: true,
	// 		}
	// 		// pb, err := json.Marshal(p)
	// 		// if err != nil {
	// 		// 	log.Fatal(err)
	// 		// }

	// 		мутация.SetJson = []byte(данные.Запрос)
	// 		результат, ошибка := граф.NewTxn().Mutate(ctx, мутация)

	// 		if ошибка != nil {
	// 			данные.Ошибка <- ошибка.Error()
	// 		} else {
	// 			данные.КаналОтвет <- результат.String()
	// 		}
	// 		return
	// 	}(граф, данные)
	// }
	// Инфо(" канал закрылся, цикл прервался %+v \n", каналДанных)
	// Авторизация, пока пропустим
	// Perform login call. If the Dgraph cluster does not have ACL and
	// enterprise features enabled, this call should be skipped.
	// for {
	// 	// Keep retrying until we succeed or receive a non-retriable error.
	// 	err = dg.Login(ctx, "groot", "password")
	// 	if err == nil || !strings.Contains(err.Error(), "Please retry") {
	// 		break
	// 	}
	// 	time.Sleep(time.Second)
	// }
	// if err != nil {
	// 	log.Fatalf("While trying to login %v", err.Error())
	// }
	// if err := связь.Close(); err != nil {
	// 	Ошибка(" Ошибка закрытия соединения %+v \n", err)
	// }
	// resp, err := граф.NewTxn().QueryWithVars(ctx, q, variables)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	return СоединениеСДГраф{
		Граф: граф,
		ЗакрытьДГраф: func() {
			if err := связь.Close(); err != nil {
				Ошибка(" Ошибка закрытия соединения %+v \n", err)
			}
		},
	}
	// return граф, func() {
	// // 	if err := связь.Close(); err != nil {
	// // 		Ошибка(" Ошибка закрытия соединения %+v \n", err)
	// // 	}
	// // }
}

/*
Изменить открывает транзакцию на изменение, принимает запрос на измнение, и данные для подстановки
Передавать можно любые запросы на вставку, и чтение.
Берёт соединнение из пула, отправляет запрос и возвращает соелинение в пул
*/
// func Изменить(запрос ДанныеЗапроса, граф *dgo.Dgraph) (string, СтатусСервиса) {
// 	ctx := context.Background()

// 	мутация := &api.Mutation{
// 		CommitNow: true,
// 	}
// 	мутация.SetJson = []byte(запрос.Запрос)
// 	результат, ошибка := граф.NewTxn().Mutate(ctx, мутация)

//		if ошибка != nil {
//			return результат.String(), СтатусСервиса{
//				Код:   ОшибкаЗаписи,
//				Текст: ошибка.Error(),
//			}
//		}
//		return результат.String(), СтатусСервиса{
//			Код: Ок,
//		}
//	}

/*
добавить := `[

	    {
	        "логин": "user5",
	        "имя": "Алексей Алексеев",
	        "email": "alexey@example.com",
	        "dgraph.type": "Пользователи"
	    },
	    {
	        "логин": "user6",
	        "имя": "Наталья Натальева", ит
	        "email": "natalya@example.com",
	        "dgraph.type": "Пользователи"
	    }
	]`
	ответиз, статусИзменения := База.Изменить(ДанныеЗапроса{
		Запрос: добавить,
	})
*/
func (база СоединениеСДГраф) Изменить(запрос ДанныеЗапроса) ([]byte, СтатусБазы) {
	граф := база.Граф

	for {
		ctx := context.Background()
		транзакция := граф.NewTxn()
		defer транзакция.Discard(ctx)

		мутация := &api.Mutation{

			CommitNow: true,
		}
		мутация.SetJson = []byte(запрос.Запрос)
		результат, ошибка := транзакция.Mutate(ctx, мутация)

		if ошибка != nil {
			if strings.Contains(ошибка.Error(), "conflict") {
				// Конфликт транзакции, повторяем
				Инфо(" Конфликт транзакции, повторяем %+v \n", ошибка.Error())
				continue
			}
			Ошибка(" ошибка %+s\n", ошибка.Error())

			return nil, СтатусБазы{
				Код:   ОшибкаЗаписи,
				Текст: ошибка.Error(),
			}
		}
		// не делаем комит вручную, так как установлен флаг CommitNow
		// ошибка = транзакция.Commit(ctx)
		// if ошибка != nil {
		// 	if strings.Contains(ошибка.Error(), "conflict") {
		// 		// Конфликт транзакции, повторяем
		// 		continue
		// 	}
		// 	return "", СтатусСервиса{
		// 		Код:   ОшибкаЗаписи,
		// 		Текст: ошибка.Error(),
		// 	}
		// }

		return результат.Json, СтатусБазы{
			Код: Ок,
		}
	}
}

func (база СоединениеСДГраф) Схема(запрос ДанныеЗапроса) СтатусБазы {
	контекст := context.Background()
	операция := &api.Operation{}
	операция.Schema = запрос.Запрос
	ошибка := база.Граф.Alter(контекст, операция)
	if ошибка != nil {
		Ошибка(" Ошибка изменения схемы данных  %+v Запрос изменения схемы: %+v \n", ошибка, ошибка.Error(), запрос.Запрос)
		return СтатусБазы{
			Код:   ОшибкаИзмененияСхемы,
			Текст: ошибка.Error(),
		}
	}
	return СтатусБазы{
		Код:   Ок,
		Текст: "Схема успешно модифицированна",
	}
}

/*
Поллучить открывает транзакцию на выборку данных, отправляет запрос, возвращает результат в  виде json строки
Берёт соединнение из пула, отправляет запрос и возвращает соелинение в пул
*/
func (база СоединениеСДГраф) Получить(запрос ДанныеЗапроса) ([]byte, СтатусБазы) {
	ctx := context.Background()
	транзакция := база.Граф.NewReadOnlyTxn()
	defer транзакция.Discard(ctx)

	ответ, ошибка := транзакция.QueryWithVars(context.Background(), запрос.Запрос, запрос.Данные)
	if ошибка != nil {
		Ошибка("  %+v \n", ошибка.Error())
		return nil, СтатусБазы{
			Код:   ОшибкаПолученияДанныхИзБазы,
			Текст: ошибка.Error(),
		}
	}

	return ответ.Json, СтатусБазы{
		Код:   Ок,
		Текст: "Данные получены",
	}
}

func (база СоединениеСДГраф) Выполнить(запрос string) ([]byte, СтатусБазы) {
	ctx := context.Background()
	транзакция := база.Граф.NewTxn()
	defer транзакция.Discard(ctx)

	ответ, ошибка := транзакция.Query(context.Background(), запрос)
	if ошибка != nil {
		Ошибка("  %+v \n", ошибка.Error())
		return nil, СтатусБазы{
			Код:   ОшибкаПолученияДанныхИзБазы,
			Текст: ошибка.Error(),
		}
	}

	return ответ.Json, СтатусБазы{
		Код:   Ок,
		Текст: "Данные получены",
	}
}

func (база СоединениеСДГраф) ИзмененитьПоУсловию(запрос ДанныеЗапроса) ([]byte, СтатусБазы) {
	ctx := context.Background()
	// транзакция1 := транзакция.Граф.NewTxn()
	// defer транзакция.Discard(ctx)

	мутации := запрос.Мутация

	апиМутации := make([]*api.Mutation, len(мутации))

	for _, мутация := range мутации {
		апиМутации = append(апиМутации, &api.Mutation{
			Cond:    мутация.Условие,
			SetJson: мутация.Мутация,
		})
	}

	запросВБД := &api.Request{
		CommitNow: false,
		Query:     запрос.Запрос,
		Mutations: апиМутации,
		Vars:      запрос.Данные,
	}

	ответ, ошибка := база.Транзакция.Do(ctx, запросВБД)

	if ошибка != nil {
		Ошибка("  %+v \n", ошибка.Error())
		return nil, СтатусБазы{
			Код:   ОшибкаЗаписи,
			Текст: ошибка.Error(),
		}
	}

	return ответ.Json, СтатусБазы{
		Код:   Ок,
		Текст: "Данные получены",
	}
}
