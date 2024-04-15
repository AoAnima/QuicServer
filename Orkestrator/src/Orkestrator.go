package main

import (
	"bufio"
	"io"
	"os"
	"os/exec"
	"strings"
	"time"

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
		файлСервисов, ошибка := os.Open("services.txt") // Замените "test.txt" на имя вашего файла
		if ошибка != nil {
			Ошибка("  %+v \n", ошибка)
		}
		defer файлСервисов.Close()

		// Читаем содержимое файла
		данные, ошибка := io.ReadAll(файлСервисов)
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

	for {
		time.Sleep(time.Minute * 1)
		Инфо(" прошло 1 минут %+v \n")
	}

	// // Замените ./microservice2 на путь к вашему микросервису
	// cmd = exec.Command("cmd", "/C", "start", "cmd.exe", "/K", "D:/QuicMarket/GoP/HTTPServerQuic/bin/HTTPServerQuic.exe")
	// if err := cmd.Run(); err != nil {
	// 	Ошибка(" %+v \n")
	// }
}

func ЗапускСервиса(папка string) {

	dir, err := os.Getwd()
	Инфо(" %+v  %+v \n", dir, err)
	ДиреткорияПрокета := dir + "/../../"

	открытаяПапка, ошибка := os.Open(ДиреткорияПрокета + папка + "/bin/")
	// открытаяПапка, ошибка := os.Open("D:/QuicMarket/GoP/" + папка + "/bin/")
	if ошибка != nil {
		Ошибка("  %+v \n", ошибка)
	}

	файлы, ошибка := открытаяПапка.Readdir(-1)
	открытаяПапка.Close()
	if ошибка != nil {
		Ошибка("  %+v \n", ошибка)
	}

	for _, файл := range файлы {
		Инфо(" %+v  %+v \n", файл.Name(), папка)

		// if filepath.Ext(файл.Name()) == ".exe" {
		if файл.Name() == папка {
			go func() {
				Инфо(" Запуск приложения%+v \n", ДиреткорияПрокета+папка+"/bin/"+файл.Name())

				// cmd := exec.Command("mate-terminal", "-e", "bash -c '"+ДиреткорияПрокета+папка+"/bin/"+файл.Name()+"; exec bash'")

				cmd := exec.Command("cmd", "/C", "start", "cmd.exe", "/K", "D:/QuicMarket/GoP/"+папка+"/bin/"+файл.Name())

				стандартныйВывод, err := cmd.StdoutPipe()
				if err != nil {
					Ошибка(" %+v \n", err)
				}

				if err := cmd.Start(); err != nil {
					Ошибка(" %+v \n", err)
				}
				scanner := bufio.NewScanner(стандартныйВывод)
				go func() {
					for scanner.Scan() {
						ЧитатьСтандртныйВывод(файл.Name(), scanner.Text())
					}
				}()

				if err := cmd.Wait(); err != nil {
					Ошибка(" %+v \n", err)
				}
			}()
		}
	}

}

func ЧитатьСтандртныйВывод(ИмяСервиса string, лог string) {
	Инфо(" %+v  %+v \n", ИмяСервиса, лог)
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
