module aoanima.ru/synqTCP

go 1.21.3

replace aoanima.ru/logger => ../logger

replace aoanima.ru/connector => ../connector

require (
	aoanima.ru/connector v0.0.0-00010101000000-000000000000
	aoanima.ru/logger v0.0.0-00010101000000-000000000000
	github.com/json-iterator/go v1.1.12
)

require github.com/stretchr/testify v1.6.1 // indirect

require (
	github.com/dgryski/go-metro v0.0.0-20211217172704-adc40b04c140 // indirect
	github.com/google/uuid v1.3.1 // indirect
	github.com/modern-go/concurrent v0.0.0-20180228061459-e0a39a4cb421 // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
)
