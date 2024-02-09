package DGApi

import (
	"log"

	. "aoanima.ru/Logger"

	dgo "github.com/dgraph-io/dgo/v230"
	"github.com/dgraph-io/dgo/v230/protos/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// https://github.com/dgraph-io/dgo/blob/master/example_set_object_test.go

type ФункцияОтмены func()

func Граф() (*dgo.Dgraph, ФункцияОтмены) {
	// Dial a gRPC connection. The address to dial to can be configured when
	// setting up the dgraph cluster.
	связь, err := grpc.Dial("localhost:9080", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}

	dc := api.NewDgraphClient(связь)
	dg := dgo.NewDgraphClient(dc)
	// ctx := context.Background()

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

	return dg, func() {
		if err := связь.Close(); err != nil {
			Ошибка(" Ошибка закрытия соединения %+v \n", err)
		}
	}
}

func Вставить() {

}
