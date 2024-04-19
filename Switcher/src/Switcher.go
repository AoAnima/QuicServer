package main

import (
	"fmt"
	"os"
	"sync"

	. "aoanima.ru/Logger"
	evdev "github.com/gvalkov/golang-evdev"
	jsoniter "github.com/json-iterator/go"
	"golang.org/x/text/language"
)

var РусскиеБуквы map[string]float64
var РусскиеБиграммы map[string]float64
var РусскиеТриграммы map[string]float64
var АнглийскиеБуквы map[string]float64
var АнглийскиеБиграммы map[string]float64
var АнглийскиеТриграммы map[string]float64

func init() {

	дир := "словари"
	словари, err := os.ReadDir(дир)
	if err != nil {
		Ошибка("Ошибка чтения директории %+s", err)
	}

	for _, словарь := range словари {
		Инфо("словарь %+v \n", словарь)

		filePath := дир + "/" + словарь.Name()
		файл, err := os.ReadFile(filePath)
		if err != nil {
			Ошибка("Ошибка открытия файла %+s", err)
			continue
		}
		// defer файл.Close()
		// err = jsoniter.Unmarshal(файл, Конфиг)
		// if err != nil {
		// 	Ошибка("  %+v \n", err)
		// }

		switch имяФайла := словарь.Name(); имяФайла {
		case "РусскиеБуквы.json":
			err = jsoniter.Unmarshal(файл, &РусскиеБуквы)
			if err != nil {
				Ошибка("  %+v \n", err.Error())
			}
		case "АнглийскиеБуквы.json":

			err = jsoniter.Unmarshal(файл, &АнглийскиеБуквы)
			if err != nil {
				Ошибка("  %+v \n", err.Error())
			}
		case "РусскиеТриграммы.json":
			err = jsoniter.Unmarshal(файл, &РусскиеТриграммы)
			if err != nil {
				Ошибка("  %+v \n", err.Error())
			}
		case "АнглийскиеТриграммы.json":
			err = jsoniter.Unmarshal(файл, &АнглийскиеТриграммы)
			if err != nil {
				Ошибка("  %+v \n", err.Error())
			}
		case "РусскиеБиграммы.json":
			err = jsoniter.Unmarshal(файл, &РусскиеБиграммы)
			if err != nil {
				Ошибка("  %+s \n", err.Error())
			}
		case "АнглийскиеБиграммы.json":
			err = jsoniter.Unmarshal(файл, &АнглийскиеБиграммы)
			if err != nil {
				Ошибка("  %+v \n", err.Error())
			}
		default:
			fmt.Println("Неизвестный формат файла", имяФайла)
		}
	}
	// Инфо(" %+v \n %+v \n %+v \n %+v \n %+v \n %+v \n", АнглийскиеБиграммы, РусскиеБиграммы, АнглийскиеТриграммы, РусскиеТриграммы, АнглийскиеБуквы, РусскиеБуквы)

}

func main() {
	devices, _ := evdev.ListInputDevices()

	var wg sync.WaitGroup
	lang := language.BCP47.Make("")
	Инфо(" %+v \n", lang)

	for _, dev := range devices {
		// Инфо("%s %s %s \n", dev.Fn, dev.Name, dev.Phys)
		// Инфо("%s %s %s \n", dev)

		hasSyn := false
		hasKey := false
		hasMsc := false
		hasLed := false
		for ev := range dev.Capabilities {
			switch ev.Type {
			case evdev.EV_SYN:
				hasSyn = true
			case evdev.EV_KEY:
				hasKey = true
			case evdev.EV_MSC:
				hasMsc = true
			case evdev.EV_LED:
				hasLed = true
			}
		}

		if hasSyn && hasKey && hasMsc && hasLed {
			wg.Add(1)
			go ЧитатьУстройство(dev)
		}

		// if strings.Contains(dev.Name, "Keyboard") && !strings.Contains(dev.Name, "Mouse ds") {
		// 	// EV_SYN 0, EV_KEY 1, EV_MSC 4, EV_LED 17
		// 	// evtypes := make([]string, 0)
		// 	Инфо(" %+v \n", dev.Capabilities[evdev.CapabilityType{Type: evdev.EV_KEY}])

		// 	for ev := range dev.Capabilities {
		// 		Инфо(" %+v %+v \n", ev.Name, ev.Type)

		// 	}

		// 	wg.Add(1)
		// 	go ЧитатьУстройство(dev)
		// }
	}
	wg.Wait()
	// Открываем устройство клавиатуры
	// device, err := evdev.Open("/dev/input/event9") // Замените X на номер устройства клавиатуры
	// if err != nil {
	// 	Ошибка("%+v", err)
	// }
	// Инфо(" %+v \n", device)

	// // Читаем события клавиатуры
	// for {
	// 	events, err := device.Read()
	// 	Инфо("events %+v \n", events)

	// 	if err != nil {
	// 		Ошибка("%+v", err)
	// 	}

	// 	for _, event := range events {
	// 		if event.Type == evdev.EV_KEY {
	// 			switch event.Value {
	// 			case 1:
	// 				fmt.Printf("Клавиша нажата: %v\n", event.Code)
	// 			case 0:
	// 				fmt.Printf("Клавиша отпущена: %v\n", event.Code)
	// 			}
	// 		}
	// 	}
	// }
}

func ЧитатьУстройство(dev *evdev.InputDevice) {
	Инфо("Читаем  %+v \n", dev)

	device, err := evdev.Open(dev.Fn) // Замените X на номер устройства клавиатуры
	if err != nil {
		Ошибка("%+v", err)
	}

	// Читаем события клавиатуры
	for {
		events, err := device.Read()
		// Инфо("events %+v \n", events)

		if err != nil {
			Ошибка("%+v", err)
		}

		for _, event := range events {
			// Инфо(" %+v  %+v \n", event)

			if event.Type == evdev.EV_KEY || event.Type == evdev.EV_MSC {
				switch event.Value {
				case 1:
					Инфо("Клавиша нажата:  Code %+v Value %+v \n", event.Code, event.Value)
				case 2:
					Инфо("Клавиша задержанай: Code %+v Value %+v \n", event.Code, event.Value)
				case 0:
					Инфо("Клавиша отпущена: %+v - %+v \n", event.Code, event.Value)
				}
			}
		}
	}
}
