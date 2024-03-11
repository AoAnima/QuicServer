# Замените 'path_to_microservice1' и 'path_to_microservice2' на пути к исполняемым файлам ваших микросервисов
$services = 'D:/QuicMarket/GoP/SynQuic/bin/synquic.exe', 'D:/QuicMarket/GoP/HTTPServerQuic/bin/HTTPServerQuic.exe'

foreach ($service in $services) {
    Start-Process -FilePath "powershell.exe" -ArgumentList "-NoExit","-Command  &{ $service }"
}