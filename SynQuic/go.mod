module aoanima.ru/SynQuic

go 1.21.5

toolchain go1.22.0

replace aoanima.ru/Logger => ../Logger

replace aoanima.ru/ConnQuic => ../ConnQuic

replace aoanima.ru/QErrors => ../QErrors

require (
	aoanima.ru/ConnQuic v0.0.0-00010101000000-000000000000
	aoanima.ru/Logger v0.0.0-00010101000000-000000000000
	aoanima.ru/QErrors v0.0.0-00010101000000-000000000000
	github.com/quic-go/quic-go v0.41.0

)

require (
	github.com/cespare/xxhash/v2 v2.2.0 // indirect
	github.com/dgraph-io/ristretto v0.1.1 // indirect
	github.com/dustin/go-humanize v1.0.0 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/golang/glog v1.0.0 // indirect
	github.com/golang/groupcache v0.0.0-20190702054246-869f871628b6 // indirect
	github.com/golang/protobuf v1.5.3 // indirect
	github.com/golang/snappy v0.0.3 // indirect
	github.com/google/flatbuffers v1.12.1 // indirect
	github.com/gookit/color v1.5.4 // indirect
	github.com/klauspost/compress v1.12.3 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/xo/terminfo v0.0.0-20210125001918-ca9a967f8778 // indirect
	go.opencensus.io v0.22.5 // indirect
	google.golang.org/protobuf v1.28.1 // indirect
)

require (
	github.com/dgraph-io/badger/v4 v4.2.0
	github.com/dgryski/go-metro v0.0.0-20211217172704-adc40b04c140 // indirect
	github.com/go-task/slim-sprig v0.0.0-20230315185526-52ccab3ef572 // indirect
	github.com/google/pprof v0.0.0-20210407192527-94a9f03dee38 // indirect
	github.com/google/uuid v1.3.1 // indirect
	github.com/json-iterator/go v1.1.12
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/onsi/ginkgo/v2 v2.9.5 // indirect
	go.uber.org/mock v0.3.0 // indirect
	golang.org/x/crypto v0.4.0 // indirect
	golang.org/x/exp v0.0.0-20221205204356-47842c84f3db // indirect
	golang.org/x/mod v0.11.0 // indirect
	golang.org/x/net v0.10.0 // indirect
	golang.org/x/sys v0.10.0 // indirect
	golang.org/x/tools v0.9.1 // indirect
)
