module aoanima.ru/synquestWS

go 1.21.1

replace aoanima.ru/logger => ../logger

require (
	aoanima.ru/logger v0.0.0-00010101000000-000000000000
	github.com/dgrr/fastws v1.0.4
)

require (
	github.com/klauspost/compress v1.10.4 // indirect
	github.com/valyala/bytebufferpool v1.0.0 // indirect
	github.com/valyala/fasthttp v1.12.0 // indirect
)
