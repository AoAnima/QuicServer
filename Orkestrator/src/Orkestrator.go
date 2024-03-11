package main

import (
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	_ "aoanima.ru/ConnQuic"
	_ "aoanima.ru/DGApi"
	. "aoanima.ru/Logger"
	_ "aoanima.ru/QErrors"
)

func main() {
	запускаемыеСервисы := os.Args[1:]
	if len(запускаемыеСервисы) > 0 {
		for i := range запускаемыеСервисы {
			ЗапускСервиса(запускаемыеСервисы[i])
		}
	} else {
		файлСервисов, ошибка := os.Open("./services.txt") // Замените "test.txt" на имя вашего файла
		if ошибка != nil {
			Ошибка("  %+v \n", ошибка)
		}
		defer файлСервисов.Close()

		// Читаем содержимое файла
		данные, ошибка := ioutil.ReadAll(файлСервисов)
		if ошибка != nil {
			Ошибка("  %+v \n", ошибка)
		}

		сервисы := strings.Split(string(данные), ";")
		Инфо("  %+v \n", сервисы)
		if len(сервисы) > 0 {
			for i := range сервисы {
				ЗапускСервиса(сервисы[i])
			}
		}
	}

	// // Замените ./microservice2 на путь к вашему микросервису
	// cmd = exec.Command("cmd", "/C", "start", "cmd.exe", "/K", "D:/QuicMarket/GoP/HTTPServerQuic/bin/HTTPServerQuic.exe")
	// if err := cmd.Run(); err != nil {
	// 	Ошибка(" %+v \n")
	// }
}

func ЗапускСервиса(папка string) {

	открытаяПапка, ошибка := os.Open("D:/QuicMarket/GoP/" + папка + "/bin/")
	if ошибка != nil {
		Ошибка("  %+v \n", ошибка)
	}

	файлы, ошибка := открытаяПапка.Readdir(-1)
	открытаяПапка.Close()
	if ошибка != nil {
		Ошибка("  %+v \n", ошибка)
	}

	for _, файл := range файлы {
		if filepath.Ext(файл.Name()) == ".exe" {
			Инфо(" Запуск приложения%+v \n", "D:/QuicMarket/GoP/"+папка+"/bin/"+файл.Name())

			cmd := exec.Command("cmd", "/C", "start", "cmd.exe", "/K", "D:/QuicMarket/GoP/"+папка+"/bin/"+файл.Name())

			if err := cmd.Run(); err != nil {
				Ошибка(" %+v \n", err)
			}
			break
		}
	}
}

func runPowerShell() {
	cmd := exec.Command("cmd", "/C", "start", "powershell.exe", "-NoExit", "-NoProfile", "-Command", "D:/QuicMarket/GoP/SynQuic/bin/synquic.exe")

	if err := cmd.Run(); err != nil {
		Ошибка(" %+v \n")
	}

	cmd = exec.Command("cmd", "/C", "start", "powershell.exe", "-NoExit", "-NoProfile", "-Command", "D:/QuicMarket/GoP/HTTPServerQuic/bin/HTTPServerQuic.exe")
	if err := cmd.Run(); err != nil {
		Ошибка(" %+v \n")
	}
}
