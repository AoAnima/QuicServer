package main

import (
	"encoding/binary"
	"fmt"
	"log"
	"os"
	"runtime"
	"sync"
	"time"

	. "aoanima.ru/Logger"

	dbus "github.com/godbus/dbus/v5"
	"github.com/micmonay/keybd_event"

	evdev "github.com/gvalkov/golang-evdev"
	jsoniter "github.com/json-iterator/go"
)

/*
xset -q | grep -A 0 'LED' | cut -c59-67
It prints 00000002 or 00001002 depending on your current keyboard layout.
*/
var РусскиеБуквы map[string]float64
var РусскиеБиграммы map[string]float64
var РусскиеТриграммы map[string]float64
var АнглийскиеБуквы map[string]float64
var АнглийскиеБиграммы map[string]float64
var АнглийскиеТриграммы map[string]float64

var Клава keybd_event.KeyBonding

func init() {
	ОткрытьКлаву()
	// ОпределитьРаскладку()
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

// keybd_event.KeyBonding
func ОткрытьКлаву() {

	var err error
	Клава, err = keybd_event.NewKeyBonding()
	if err != nil {
		panic(err)
	}
	if runtime.GOOS == "linux" {
		time.Sleep(2 * time.Second)
	}
}

func СменитьРаскладку() {
	// kb, err := keybd_event.NewKeyBonding()
	// if err != nil {
	// 	panic(err)
	// }

	// For linux, it is very important to wait 2 seconds

	// kb.HasCTRL(true)
	// Select keys to be pressed
	Клава.SetKeys(56, 42)

	// Set shift to be pressed

	// Press the selected keys
	err := Клава.Launching()
	if err != nil {
		panic(err)
	}

	// // kb.HasSHIFT(true)
	// time.Sleep(2 * time.Second)
	// Or you can use Press and Release
	// kb.Press()
	// time.Sleep(10 * time.Millisecond)
	// kb.Release()

}

// func aaa(dev *evdev.InputDevice) {
// 	syscall.Setgid(65534)
// 	syscall.Setuid(65534)

// 	Инфо(" %+v \n", dev.File.Name())

// 	syscall.Setgid(0)
// 	syscall.Setuid(0)
// 	// OpenFile(name, O_RDONLY, 0)
// 	input, err := os.OpenFile(dev.Fn, os.O_RDONLY, 0)
// 	if err != nil {
// 		fmt.Println(err)
// 		return
// 	}
// 	defer input.Close()

// 	var buffer = make([]byte, 24)
// 	for {
// 		n, err := input.Read(buffer)
// 		if err != nil {
// 			return
// 		}

// 		if n != 24 {
// 			fmt.Println("Weird Input Event Size: ", n)
// 			continue
// 		}

// 		binary.LittleEndian.Uint64(buffer[0:8])
// 		binary.LittleEndian.Uint64(buffer[8:16])
// 		etype := binary.LittleEndian.Uint16(buffer[16:18])
// 		code := binary.LittleEndian.Uint16(buffer[18:20])
// 		value := int32(binary.LittleEndian.Uint32(buffer[20:24]))

// 		if etype == 1 && value == 1 {
// 			fmt.Println("Key Pressed:", code)
// 		}
// 	}
// }

func dbu() {
	conn, err := dbus.ConnectSessionBus()
	if err != nil {
		log.Fatalf("Failed to connect to D-Bus: %v", err)
	}
	defer conn.Close()

	// Получение объекта IBus
	obj := conn.Object("org.freedesktop.IBus", "/org/freedesktop/IBus")
	Инфо(" obj %+v \n", obj)

	ibus := obj.Call("org.freedesktop.IBus.GetCurrentInputContext", 0)
	Инфо(" ibus %+v \n", ibus)

	if ibus.Err != nil {
		log.Fatalf("Failed to get current input context: %v", ibus.Err)
	}

	// Получение текущей раскладки
	inputContext := conn.Object("org.freedesktop.IBus", ibus.Body[0].(dbus.ObjectPath))
	Инфо("inputContext %+v \n", inputContext)

	layout, err := inputContext.GetProperty("org.freedesktop.IBus.InputContext.InputMode")
	if err != nil {
		log.Fatalf("Failed to get input mode: %v", err)
	}

	fmt.Printf("Текущая раскладка: %s\n", layout.Value().(string))
}

func main() {

	// СменитьРаскладку()
	dbu()
	devices, _ := evdev.ListInputDevices()
	// Dbus()
	var wg sync.WaitGroup
	// lang := language.BCP47.Make("")
	// Инфо(" %+v \n", lang)
	// Клава()
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
			// aaa(dev)
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

	клава, err := evdev.Open(dev.Fn) // Замените X на номер устройства клавиатуры
	if err != nil {
		Ошибка("%+v", err)
	}

	var buffer = make([]byte, 24)
	for {
		n, err := клава.File.Read(buffer)
		if err != nil {
			return
		}

		if n != 24 {
			fmt.Println("Weird Input Event Size: ", n)
			continue
		}

		binary.LittleEndian.Uint64(buffer[0:8])
		binary.LittleEndian.Uint64(buffer[8:16])
		etype := binary.LittleEndian.Uint16(buffer[16:18])
		code := binary.LittleEndian.Uint16(buffer[18:20])
		value := int32(binary.LittleEndian.Uint32(buffer[20:24]))

		// if etype == 1 && value == 1 {
		fmt.Println("Key Pressed:", etype, code, value)
		// }апвук
	}

	Инфо("Читаем клава %+v \n", клава)
	// Читаем события клавиатуры
	for {
		события, err := клава.Read()
		Инфо("события %+v \n", события)

		if err != nil {
			Ошибка("%+v", err)
		}

		for _, событие := range события {
			// Инфо(" %+v Type %+v \n", событие, событие.Type)

			if событие.Type == evdev.EV_KEY || событие.Type == evdev.EV_MSC {
				// Инфо(" %+v \n", format_event(&событие))
				Инфо("Читаем клава %+v \n", клава)

				switch событие.Value {
				case 1:
					Инфо("Клавиша нажата:  Code %+v Value %+v  событие.Type %+v \n", событие.Code, событие.Value, событие.Type)
				case 2:
					Инфо("Клавиша задержанай: Code %+v Value %+v  событие.Type %+v \n", событие.Code, событие.Value, событие.Type)
				case 0:
					Инфо("Клавиша отпущена: %+v - %+v  событие.Type %+v \n", событие.Code, событие.Value, событие.Type)
				default:
					Инфо("default: Code %+v Value %+v  событие.Type %+v \n", событие.Code, событие.Value, событие.Type)

				}

				switch событие.Type {
				case 4:
					Инфо("СОБЫТИЕ4 %+v %+s\n", событие.Value, событие.Value)
					// Rune, size := utf8.DecodeRune(событие.Value)

					// Инфо(" %+v %+s\n", Rune, size)

				}
			}
		}
	}
}

// func format_event(ev *evdev.InputEvent) string {
// 	var res, f, code_name string

// 	code := int(ev.Code)
// 	etype := int(ev.Type)

// 	switch ev.Type {
// 	case evdev.EV_SYN:
// 		if ev.Code == evdev.SYN_MT_REPORT {
// 			f = "time %d.%-8d +++++++++ %s ++++++++"
// 		} else {
// 			f = "time %d.%-8d --------- %s --------"
// 		}
// 		return fmt.Sprintf(f, ev.Time.Sec, ev.Time.Usec, evdev.SYN[code])
// 	case evdev.EV_KEY:
// 		val, haskey := evdev.KEY[code]
// 		if haskey {
// 			code_name = val
// 		} else {
// 			val, haskey := evdev.BTN[code]
// 			if haskey {
// 				code_name = val
// 			} else {
// 				code_name = "?"
// 			}
// 		}
// 	default:
// 		m, haskey := evdev.ByEventType[etype]
// 		if haskey {
// 			code_name = m[code]
// 		} else {
// 			code_name = "?"
// 		}
// 	}

// 	evfmt := "time %d.%-8d type %d (%s), code %-3d (%s), value %d"
// 	res = fmt.Sprintf(evfmt, ev.Time.Sec, ev.Time.Usec, etype,
// 		evdev.EV[int(ev.Type)], ev.Code, code_name, ev.Value)

// 	return res
// }

// func ОпределитьРаскладку() {
// 	cmd := exec.Command("setxkbmap", "-query")
// 	var out bytes.Buffer
// 	cmd.Stdout = &out
// 	err := cmd.Run()
// 	if err != nil {
// 		fmt.Println("Ошибка выполнения команды:", err)
// 		return
// 	}

// 	output := out.String()
// 	lines := strings.Split(output, "\n")
// 	Инфо(" %+v \n", lines)

// 	for _, line := range lines {
// 		if strings.HasPrefix(line, "layout:") {
// 			layout := strings.TrimSpace(strings.TrimPrefix(line, "layout:"))
// 			fmt.Println("Текущая раскладка клавиатуры:", layout)
// 			break
// 		}
// 	}
// }

// func Клава() {
// 	keysEvents, err := keyboard.GetKeys(10)
// 	if err != nil {
// 		panic(err)
// 	}
// 	defer func() {
// 		_ = keyboard.Close()
// 	}()

// 	fmt.Println("Press ESC to quit")
// 	for {
// 		event := <-keysEvents
// 		if event.Err != nil {
// 			panic(event.Err)
// 		}
// 		fmt.Printf("You pressed: rune %q, key %X\r\n", event.Rune, event.Key)
// 		if event.Key == keyboard.KeyEsc {
// 			break
// 		}
// 	}
// }
