package Logger

import (
	// "encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	json "github.com/json-iterator/go"
	"github.com/quic-go/quic-go"
)

var StdLog = log.New(os.Stderr, "", log.Lshortfile|log.Ltime)

func Инфо(формат string, данные ...interface{}) {
	StdLog.SetFlags(log.Lshortfile | log.Ltime)
	формат = strings.ReplaceAll(формат, "%+v", "\u001b[38;5;48m %+v  \u001b[0m\u001b[38;5;75m")

	// log.Printf("1 %+v ", len(данные))
	данные = КрасивыйВывод(данные...)

	str := fmt.Sprintf("\u001b[0m\u001b[36mИНФО: \u001b[38;5;75m "+формат+" \u001b[0m \n", данные...)
	err := StdLog.Output(2, str)

	if err != nil {
		log.Printf("%+v", err)
	}

}

func КрасивыйВывод(данные ...interface{}) []interface{} {
	данныеДляВывода := []interface{}{}
	// log.Print(данные)
	for _, д := range данные {
		// log.Printf("1 %+v %#T", д, д)
		switch д.(type) {
		case byte:
			данныеДляВывода = append(данныеДляВывода, д)
			continue
		case bool:
			данныеДляВывода = append(данныеДляВывода, д)
			continue
		case []uint8:
			данныеДляВывода = append(данныеДляВывода, д)
			continue
		case []uint64:
			данныеДляВывода = append(данныеДляВывода, д)
			continue
		case uint64:
			данныеДляВывода = append(данныеДляВывода, д)
			continue
		case []uint32:
			данныеДляВывода = append(данныеДляВывода, д)
			continue
		case uint32:
			данныеДляВывода = append(данныеДляВывода, д)
			continue
		case float64:
			данныеДляВывода = append(данныеДляВывода, д)
			continue
		case []float64:
			данныеДляВывода = append(данныеДляВывода, д)
			continue
		case float32:
			данныеДляВывода = append(данныеДляВывода, д)
			continue
		case []float32:
			данныеДляВывода = append(данныеДляВывода, д)
			continue
		case int:
			данныеДляВывода = append(данныеДляВывода, д)
			continue
		case int32:
			данныеДляВывода = append(данныеДляВывода, д)
			continue
		case int64:
			данныеДляВывода = append(данныеДляВывода, д)
			continue

		case http.Request:
			// log.Printf("2 %+v ", д)
			empJSON, err := json.MarshalIndent(д, "", "  ")
			if err != nil {
				данныеДляВывода = append(данныеДляВывода, д)
			}
			данныеДляВывода = append(данныеДляВывода, string(empJSON))
			continue
		case quic.Stream:
			// log.Printf("2 %+v ", д)
			empJSON, err := json.MarshalIndent(д, "", "  ")
			if err != nil {
				log.Printf("101 *Stream %#T %+v %+v \n", д, д, err)

				данныеДляВывода = append(данныеДляВывода, д)
			}
			данныеДляВывода = append(данныеДляВывода, string(empJSON))
			continue
		case *quic.Stream:
			log.Printf("101 *Stream %#T %+v \n", д, д)

			// log.Printf("2 %+v %+v", тип, д)
			// log.Printf("2 %+v ", д)
			empJSON, err := json.MarshalIndent(д, "", "  ")
			if err != nil {
				log.Printf("101 *Stream %#T %+v %+v \n", д, д, err)
				данныеДляВывода = append(данныеДляВывода, д)
			}
			данныеДляВывода = append(данныеДляВывода, string(empJSON))
			continue
		case string:
			// log.Printf("2 %+v %+v", тип, д)
			// log.Printf("2 %+v ", д)
			данныеДляВывода = append(данныеДляВывода, д)
			continue
		default:
			// log.Printf("3 %+v ", д)
			// log.Printf("default %+v %+v", тип, д)
			empJSON, err := json.MarshalIndent(д, "", "  ")
			if err != nil {
				log.Printf("default %+v %#T", д, д)
				данныеДляВывода = append(данныеДляВывода, д)
			}
			данныеДляВывода = append(данныеДляВывода, string(empJSON))
		}

	}
	// log.Print(данныеДляВывода)
	return данныеДляВывода
}

func Вывод(w io.Writer, формат string, данные ...interface{}) {

	fmt.Fprintf(w, формат, данные)
}

func Ошибка(формат string, данные ...interface{}) {
	//StdLog.SetFlags(log.Lshortfile|log.Ltime)
	формат = strings.ReplaceAll(формат, "%+v", "\u001b[38;5;204m %+v  \u001b[0m\u001b[38;5;1m")
	данные = КрасивыйВывод(данные...)
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
