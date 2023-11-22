package Logger

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

var StdLog = log.New(os.Stderr, "", log.Lshortfile|log.Ltime)

func Инфо(формат string, данные ...interface{}) {
	StdLog.SetFlags(log.Lshortfile | log.Ltime)
	формат = strings.ReplaceAll(формат, "%+v", "\u001b[38;5;48m %+v  \u001b[0m\u001b[38;5;75m")

	str := fmt.Sprintf("\u001b[0m\u001b[36mИНФО: \u001b[38;5;75m "+формат+" \u001b[0m \n", данные...)
	err := StdLog.Output(2, str)

	if err != nil {
		log.Printf("%+v", err)
	}

}

func Вывод(w io.Writer, формат string, данные ...interface{}) {

	fmt.Fprintf(w, формат, данные)
}

func Ошибка(формат string, данные ...interface{}) {
	//StdLog.SetFlags(log.Lshortfile|log.Ltime)
	формат = strings.ReplaceAll(формат, "%+v", "\u001b[38;5;204m %+v  \u001b[0m\u001b[38;5;1m")

	str := fmt.Sprintf("\u001b[48;5;124m ОШИБКА >> \u001b[0m  \u001b[38;5;1m  "+формат+" \u001b[0m \n", данные...)

	err := StdLog.Output(2, str)
	if err != nil {
		log.Printf("%+v", err)
	}
}

//var Stdout = log.New(os.Stdout, "", log.Llongfile)
/*%		T выводит тип
%t	the word true or false

Инфо(" %+v", runtime.NumGoroutine())
	memStats := new(runtime.MemStats)
	runtime.ReadMemStats(memStats)
	Инфо("Alloc %v", memStats.Alloc)
	Инфо("TotalAlloc %+v", memStats.TotalAlloc)
	Инфо("HeapAlloc %+v", memStats.HeapAlloc)

	a := make([]int, 10000000)
	for k, _ := range a {
		a[k] = rand.Int()
	}
	//_ = [2 << 10]int{}
	runtime.ReadMemStats(memStats)

	Инфо("TotalAlloc %+v", memStats.HeapAlloc)
*/
