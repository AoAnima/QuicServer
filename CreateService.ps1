
$moduleName = $args[0]

New-Item -ItemType Directory -Path ".\$moduleName"
Write-Output "Значение аргумента: $moduleName"

Copy-Item -Path .\Template\* -Destination $moduleName -Recurse
Set-Location ".\$moduleName\src"
Rename-Item -Path "main.go" -NewName "$moduleName.go"
Set-Location "..\"
go mod init aoanima.ru/$moduleName
go work use ./
go mod edit -replace aoanima.ru/Logger=../Logger
go mod edit -replace aoanima.ru/ConnQuic=../ConnQuic
go mod edit -replace aoanima.ru/QErrors=../QErrors
go mod edit -replace aoanima.ru/DataBase=../DataBase
go mod edit -replace aoanima.ru/DGApi=../DGApi
