module aoanima.ru/client

go 1.21.2

replace aoanima.ru/ConnQuic => ../ConnQuic

replace aoanima.ru/logger => ../logger

require golang.org/x/net v0.18.0

require (
	golang.org/x/crypto v0.15.0 // indirect
	golang.org/x/sys v0.14.0 // indirect
)
