module aoanima.ru/HTTPServerQuic

go 1.21.5

toolchain go1.22.0

replace aoanima.ru/Logger => ../Logger

replace aoanima.ru/QErrors => ../QErrors

replace aoanima.ru/ConnQuic => ../ConnQuic

require (
	aoanima.ru/ConnQuic v0.0.0-00010101000000-000000000000
	aoanima.ru/Logger v0.0.0-00010101000000-000000000000
	aoanima.ru/QErrors v0.0.0-00010101000000-000000000000
	github.com/golang-jwt/jwt/v5 v5.2.0
	github.com/google/uuid v1.5.0
	github.com/quic-go/quic-go v0.41.0
)

require (
	github.com/dgryski/go-metro v0.0.0-20211217172704-adc40b04c140 // indirect
	github.com/go-task/slim-sprig v0.0.0-20230315185526-52ccab3ef572 // indirect
	github.com/google/pprof v0.0.0-20231229205709-960ae82b1e42 // indirect
	github.com/gookit/color v1.5.4 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/onsi/ginkgo/v2 v2.13.2 // indirect
	github.com/xo/terminfo v0.0.0-20210125001918-ca9a967f8778 // indirect
	go.uber.org/mock v0.4.0 // indirect
	golang.org/x/crypto v0.18.0 // indirect
	golang.org/x/exp v0.0.0-20240110193028-0dcbfd608b1e // indirect
	golang.org/x/mod v0.14.0 // indirect
	golang.org/x/net v0.20.0 // indirect
	golang.org/x/sys v0.16.0 // indirect
	golang.org/x/tools v0.16.1 // indirect
)
