package main

import (
	"fmt"

	. "aoanima.ru/Logger"

	"github.com/dgraph-io/dgo/v230/protos/api"
	_ "github.com/quic-go/quic-go"
)

func СоздатьФидьтра() {

	запросРолей := `<роли>(func: type(<Роль>)) {
    uid
    <имя.роли>
    <код.роли>
}`

	данные := []int{2, 1}

	var filters []*api.FilterTree
	for _, кодРоли := range данные {
		filters = append(filters, &api.FilterTree{
			Func: &api.Function{
				Name: "eq",
				Args: []*api.Value{
					{
						Val: &api.Value_StrVal{StrVal: "<код.роли>"},
					},
					{
						Val: &api.Value_IntVal{IntVal: int64(кодРоли)},
					},
				},
			},
		})
	}

	filterTree := &api.FilterTree{
		Op:       api.FilterOp_OR,
		Children: filters,
	}
	Инфо(" %+v \n", запросРолей)

	запрос := fmt.Sprintf(запросРолей, dgraph.FilterTree(filterTree))
	Инфо(" %+v \n", запрос)

}
