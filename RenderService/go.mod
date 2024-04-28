module RenderService

go 1.21.5

replace aoanima.ru/Logger => ../Logger

replace aoanima.ru/ConnQuic => ../ConnQuic

replace gitverse.ru/Ao/jet => ../../jet

replace aoanima.ru/QErrors => ../QErrors

require (
	aoanima.ru/ConnQuic v0.0.0-00010101000000-000000000000
	aoanima.ru/Logger v0.0.0-00010101000000-000000000000
	github.com/fsnotify/fsnotify v1.7.0
	github.com/gorilla/websocket v1.5.1
	github.com/quic-go/quic-go v0.42.0
	gitverse.ru/Ao/jet v0.0.0-00010101000000-000000000000

)

require github.com/CloudyKit/fastprinter v0.0.0-20200109182630-33d98a066a53 // indirect

require (
	aoanima.ru/QErrors v0.0.0-00010101000000-000000000000 // indirect
	github.com/dgryski/go-metro v0.0.0-20211217172704-adc40b04c140 // indirect
	github.com/go-task/slim-sprig v0.0.0-20230315185526-52ccab3ef572 // indirect
	github.com/google/pprof v0.0.0-20240424215950-a892ee059fd6 // indirect
	github.com/google/uuid v1.3.1 // indirect
	github.com/gookit/color v1.5.4 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/onsi/ginkgo/v2 v2.17.1 // indirect
	github.com/xo/terminfo v0.0.0-20220910002029-abceb7e1c41e // indirect
	gitverse.ru/Ao/Logger v0.0.0-20240425082056-42ac73a57c7e // indirect
	go.uber.org/mock v0.4.0 // indirect
	golang.org/x/crypto v0.22.0 // indirect
	golang.org/x/exp v0.0.0-20240416160154-fe59bbe5cc7f // indirect
	golang.org/x/mod v0.17.0 // indirect
	golang.org/x/net v0.24.0 // indirect
	golang.org/x/sys v0.19.0 // indirect
	golang.org/x/tools v0.20.0 // indirect
)
