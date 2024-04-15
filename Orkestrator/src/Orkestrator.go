package main

import (
	"bufio"
	"bytes"
	"io"
	"os"
	"os/exec"
	"runtime"
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
	Инфо("Текущая операционная система: %+v", runtime.GOOS)
	ОС := runtime.GOOS
	for _, файл := range файлы {
		Инфо(" %+v  %+v \n", файл.Name(), папка)

		// if filepath.Ext(файл.Name()) == ".exe" {
		if файл.Name() == папка {
			go func() {
				Инфо(" Запуск приложения%+v \n", ДиреткорияПрокета+папка+"/bin/"+файл.Name())
				var cmd *exec.Cmd
				if ОС == "linux" {

					// Вот так читается лог из сервисов
					// cmd = exec.Command(ДиреткорияПрокета + папка + "/bin/" + файл.Name())
					// Так запускаются отдельные терминалы
					cmd = exec.Command("mate-terminal", "-e", "bash -c '"+ДиреткорияПрокета+папка+"/bin/"+файл.Name()+"; exec bash'")
				} else {
					cmd = exec.Command("cmd", "/C", "start", "cmd.exe", "/K", "D:/QuicMarket/GoP/"+папка+"/bin/"+файл.Name())
				}

				stderr, err := cmd.StderrPipe()
				if err != nil {
					Ошибка("Ошибка при создании StderrPipe:", err)
					return
				}
				Инфо("stderr %+v \n", stderr)

				// reader := bufio.NewReader(stderr)
				scanner := bufio.NewScanner(stderr)

				if err := cmd.Start(); err != nil {
					Инфо("Ошибка при запуске приложения:", err)
					return
				}

				scanner.Split(func(data []byte, atEOF bool) (advance int, token []byte, err error) {
					if atEOF && len(data) == 0 {
						return 0, nil, nil
					}
					if i := bytes.Index(data, []byte("<!>")); i >= 0 {
						// Найдена строка "$EOS", возвращаем все данные до этой строки
						return i + len("<!>"), data[:i+len("<!>")], nil
					}
					// Если "$EOS" не найден, продолжаем чтение
					if atEOF {
						return len(data), data, nil
					}
					return
				})
				go func() {
					// scanner := bufio.NewScanner(stderr)
					for scanner.Scan() {
						Инфо("%+v: %+v", файл.Name(), scanner.Text())
						// Здесь вы можете обработать каждую строку вывода как вам нужно
					}
				}()
				// стандартныйВыводОшибок, err := cmd.StdoutPipe()
				// Инфо("стандартныйВыводОшибок %+v \n", стандартныйВыводОшибок)

				// if err != nil {
				// 	Ошибка(" %+v \n", err.Error())
				// }
				// стандартныйВывод, err := cmd.
				// Инфо("стандартныйВывод %+v \n", стандартныйВывод)

				// if err != nil {
				// 	Ошибка(" %+v \n", err.Error())
				// }
				// bufio.NewReader(stderr).ReadString("\n")

				// go func() {

				// 	output, err := io.ReadAll(stderr)
				// 	if err != nil {
				// 		Ошибка("Ошибка при чтении:", err)
				// 		return
				// 	}

				// 	Инфо("%+v#  %+v", файл.Name(), output)

				// 	// for scanner.Scan() {
				// 	// 	Инфо("%+v#  %+v", файл.Name(), scanner.Text())
				// 	// 	// Здесь вы можете обработать каждую строку вывода как вам нужно
				// 	// }
				// }()
				// go func() {
				// 	scanner := bufio.NewScanner(стандартныйВывод)
				// 	for scanner.Scan() {
				// 		Инфо("STDOUT: %+v \n", scanner.Text())
				// 	}
				// }()

				// go func() {
				// 	scanner := bufio.NewScanner(стандартныйВыводОшибок)
				// 	for scanner.Scan() {
				// 		Инфо("ERR: %+v \n", scanner.Text())
				// 	}
				// }()
				// }(scanner)

				if err := cmd.Wait(); err != nil {
					Ошибка(" %+v \n", err)
				}
			}()
		}
	}

}

func ЧитатьСтандартныйВывод(scanner *bufio.Scanner) {
	for scanner.Scan() {
		Инфо(" %+v \n")

		Инфо(" %+v  %+v \n", scanner.Text(), scanner.Err())
	}
}

// func runPowerShell() {
// 	cmd := exec.Command("cmd", "/C", "start", "powershell.exe", "-NoExit", "-NoProfile", "-Command", "D:/QuicMarket/GoP/SynQuic/bin/synquic.exe")

// 	if err := cmd.Run(); err != nil {
// 		Ошибка(" %+v \n")
// 	}

// 	cmd = exec.Command("cmd", "/C", "start", "powershell.exe", "-NoExit", "-NoProfile", "-Command", "D:/QuicMarket/GoP/HTTPServerQuic/bin/HTTPServerQuic.exe")
// 	if err := cmd.Run(); err != nil {
// 		Ошибка(" %+v \n")
// 	}
// }
