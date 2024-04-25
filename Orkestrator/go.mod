module aoanima.ru/Orkestrator

go 1.22.0

replace aoanima.ru/Logger => ../Logger

replace aoanima.ru/ConnQuic => ../ConnQuic

replace aoanima.ru/QErrors => ../QErrors

// replace aoanima.ru/DataBase => ../DataBase

replace aoanima.ru/DGApi => ../DGApi

require (
	aoanima.ru/ConnQuic v0.0.0-00010101000000-000000000000
	aoanima.ru/DGApi v0.0.0-00010101000000-000000000000
	aoanima.ru/Logger v0.0.0-00010101000000-000000000000
	aoanima.ru/QErrors v0.0.0-00010101000000-000000000000
)

require (
	github.com/dgraph-io/dgo/v230 v230.0.1 // indirect
	github.com/dgryski/go-metro v0.0.0-20211217172704-adc40b04c140 // indirect
	github.com/go-task/slim-sprig v0.0.0-20230315185526-52ccab3ef572 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/golang/protobuf v1.5.3 // indirect
	github.com/google/pprof v0.0.0-20240424215950-a892ee059fd6 // indirect
	github.com/google/uuid v1.4.0 // indirect
	github.com/gookit/color v1.5.4 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/onsi/ginkgo/v2 v2.17.1 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/quic-go/quic-go v0.42.0 // indirect
	github.com/xo/terminfo v0.0.0-20220910002029-abceb7e1c41e // indirect
	go.uber.org/mock v0.4.0 // indirect
	golang.org/x/crypto v0.22.0 // indirect
	golang.org/x/exp v0.0.0-20240416160154-fe59bbe5cc7f // indirect
	golang.org/x/mod v0.17.0 // indirect
	golang.org/x/net v0.24.0 // indirect
	golang.org/x/sys v0.19.0 // indirect
	golang.org/x/text v0.14.0 // indirect
	golang.org/x/tools v0.20.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20231106174013-bbf56f31fb17 // indirect
	google.golang.org/grpc v1.61.0 // indirect
	google.golang.org/protobuf v1.31.0 // indirect
)
