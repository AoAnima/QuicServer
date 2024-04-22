module aoanima.ru/Switcher

go 1.22.0

replace aoanima.ru/Logger => ../Logger

replace aoanima.ru/ConnQuic => ../ConnQuic

replace aoanima.ru/QErrors => ../QErrors

replace aoanima.ru/DGApi => ../DGApi

require (
	aoanima.ru/Logger v0.0.0-00010101000000-000000000000
	github.com/godbus/dbus/v5 v5.1.0
	github.com/gvalkov/golang-evdev v0.0.0-20220815104727-7e27d6ce89b6
	github.com/json-iterator/go v1.1.12
	github.com/micmonay/keybd_event v1.1.2
)

require (
	github.com/go-task/slim-sprig v0.0.0-20230315185526-52ccab3ef572 // indirect
	github.com/google/pprof v0.0.0-20210407192527-94a9f03dee38 // indirect
	github.com/gookit/color v1.5.4 // indirect
	github.com/modern-go/concurrent v0.0.0-20180228061459-e0a39a4cb421 // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/onsi/ginkgo/v2 v2.9.5 // indirect
	github.com/quic-go/quic-go v0.41.0 // indirect
	github.com/xo/terminfo v0.0.0-20210125001918-ca9a967f8778 // indirect
	go.uber.org/mock v0.3.0 // indirect
	golang.org/x/crypto v0.4.0 // indirect
	golang.org/x/exp v0.0.0-20221205204356-47842c84f3db // indirect
	golang.org/x/mod v0.11.0 // indirect
	golang.org/x/net v0.10.0 // indirect
	golang.org/x/sys v0.12.0 // indirect
	golang.org/x/text v0.13.0 // indirect
	golang.org/x/tools v0.9.1 // indirect
)
