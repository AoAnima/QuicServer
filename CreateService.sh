#!/bin/bash

moduleName=$1
mkdir "./$moduleName"
echo "Значение аргумента: $moduleName"
cp -r ./Template/* "$moduleName"
cd "./$moduleName/src"
mv "main.go" "$moduleName.go"
cd ".."
go mod init "aoanima.ru/$moduleName"
go work use ./
go mod edit -replace aoanima.ru/Logger=../Logger
go mod edit -replace aoanima.ru/ConnQuic=../ConnQuic
go mod edit -replace aoanima.ru/QErrors=../QErrors
go mod edit -replace aoanima.ru/DataBase=../DataBase