package main

import (
	"strings"

	. "aoanima.ru/Logger"
)

// import (
// 	"bytes"
// 	"encoding/binary"
// 	"encoding/json"
// 	"fmt"
// 	"testing"

// 	"github.com/google/uuid"
// 	jsoniter "github.com/json-iterator/go"
// )

func main() {
	// // Создайте контекст с таймаутом
	// wg := sync.WaitGroup{}
	// wg.Add(1)
	// go server()

	// wg.Wait()

	объект := map[string]interface{}{
	// 	"Логин": "логин_клиента",
	// 	"Имя":   "Саня",
	// 	"Адрес": map[string]string{
	// 		"Страна":   "Россия",
	// 		"Город":    "Москва",
	// 		"Улица":    "Льва Толстого",
	// 		"Дом":      "16",
	// 		"Квартира": "2",
	// 	},
	// }
	// GetValueFromPath(объект, "Адрес.Город")
}

// func GetValueFromPath[K comparable, V any](m map[K]V, path string) {
// 	keys := strings.Split(path, ".")
// 	current := m
// 	Инфо("key  %+v \n", keys)

// 	for _, key := range keys {
// 		k := key.(K)
// 		next, ok := current[k]

// 		Инфо(" key %+v  next %+v current %+v \n", key, next, current)
// 		if !ok {
// 			Ошибка(" key %+v  next %+v current %+v \n", key, next, current)
// 		}
// 		Инфо("key %+v  next %+v %#T  \n current  %+v \n", key, next, next, current)
// 		switch t := any(next).(type) {
// 		case map[string]T:
// 			current = any(next).(map[string]T)
// 			Инфо(" current %+v \n", current)
// 		default:

// 			Инфо("key %+v  next %+v %#T  \n t  %+s \n", key, next, next, t)
// 		}
// 	}

// }

// func server() {
// 	log.Println("server")
// 	// Создайте TLS конфигурацию

// 	config, err := tlsConfig()
// 	if err != nil {
// 		log.Println(err)
// 	}
// 	// Создайте сервер
// 	listener, err := quic.ListenAddr("localhost:4242", config, nil)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	// Примите соединение
// 	session, err := listener.Accept(context.Background())
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	// Примите поток
// 	stream, err := session.AcceptStream(context.Background())
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	// Читайте данные из потока в цикле
// 	buf := make([]byte, 1024)
// 	for {
// 		n, err := stream.Read(buf)
// 		if err != nil {

// 			log.Println("Stream closed by client")
// 			break

// 		}
// 		log.Printf("Received: %s", buf[:n])
// 	}
// }

// func tlsConfig() (*tls.Config, error) {
// 	caCert, err := os.ReadFile("cert/ca.crt")
// 	if err != nil {
// 		return nil, err
// 	}

// 	caCertPool := x509.NewCertPool()
// 	caCertPool.AppendCertsFromPEM(caCert)
// 	// Инфо("Корневой сертфикат создан?  %v ", ok)

// 	cert, err := tls.LoadX509KeyPair("cert/server.crt", "cert/server.key")
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	return &tls.Config{
// 		// InsecureSkipVerify: true,
// 		RootCAs:      caCertPool,
// 		Certificates: []tls.Certificate{cert},
// 		// NextProtos:   []string{"h3", "quic", "websocket"},
// 		NextProtos: []string{"h3", "quic", "websocket"},
// 	}, nil
// }

// // func main() {
// // 	// Запуск бенчмарка

// // 	// s := []byte("идк:4543543-авы43а-4325а32авуц⁝значение:выавыаавукцп|⁝аавп")
// // 	// data := []byte("какието данные|⁝представленны строкой")

// // 	// // Буфер для чтения данных
// // 	// buffer := bytes.NewBuffer(data)

// // 	// // Читаем байты до символа ⁝
// // 	// bytes, err := buffer.ReadBytes(byte('\u205D'))
// // 	// if err != nil {
// // 	// 	fmt.Println("Error reading bytes:", err)
// // 	// 	return
// // 	// }

// // 	// // Заменяем экранированные символы ⁝ на символ ⁝
// // 	// bytes = bytes.Replace([]byte("|⁝"), []byte("\u205D"), -1)

// // 	// // Вывод результата
// // 	// fmt.Println(string(bytes)) // Вывод: какието данные⁝ представленны строкой
// // 	// буфер := new(bytes.Buffer)
// // 	// binary.Write(буфер, binary.LittleEndian, int32(3))
// // 	// fmt.Printf("0 = %v \n", буфер.Bytes())
// // 	// binary.Write(буфер, binary.LittleEndian, [3]byte{226, 129, 157})

// // 	// fmt.Printf("1 = %s \n", буфер.Bytes())
// // 	// fmt.Printf("2 = %s \n", [3]byte{226, 129, 157})
// // 	// fmt.Printf("3 = %s \n", [4]byte{0, 226, 129, 157})
// // 	// fmt.Printf("3 = %s \n", [4]byte{226, 129, 157, 0})

// // 	// буфер1 := new(bytes.Buffer)

// // 	// binary.Write(буфер1, binary.LittleEndian, []byte("⁝"))
// // 	// fmt.Printf("4 = %v \n", буфер1.Bytes())
// // 	// fmt.Printf("5 = %s \n", "⁝")
// // 	// fmt.Printf("6 = %v \n", []byte("⁝"))

// // 	// fmt.Printf("1 = %08b \n", strconv.FormatUint(255, 2))

// // 	// fmt.Printf("3=  %08b \n", []byte("255"))
// // 	// fmt.Printf("4= %08b \n", int32(255255))

// // 	// буфер := new(bytes.Buffer)
// // 	// binary.Write(буфер, binary.LittleEndian, int32(255))
// // 	// fmt.Printf("число 3  %08b \n", буфер.Bytes())

// // 	// буфер1 := new(bytes.Buffer)
// // 	// binary.Write(буфер1, binary.LittleEndian, []byte("255"))
// // 	// fmt.Printf("Строка 3 %08b \n", буфер1.Bytes())

// // 	// bytes := []byte("255")

// // 	// u := binary.LittleEndian.Uint64(bytes)

// // 	// binaryArray := make([]uint8, binary.Size(bytes))
// // 	// binary.LittleEndian.PutUint64(binaryArray, u)
// // 	// fmt.Printf("Строка 3 %08b \n binaryArray%08b \n", u, binaryArray)

// // 	// fmt.Printf("Строка 2 5 5  %08b  %08b  %08b \n ", "2", "5", "5")

// // 	// result := testing.Benchmark(BenchmarkJsoniter)

// // 	// // // Вывод результатов
// // 	// fmt.Println("Строка", result)

// // 	// result1 := testing.Benchmark(BenchmarkJson)

// // 	// // // Вывод результатов
// // 	// fmt.Println(" \n Строка json", result1)
// // 	// // result1 := testing.Benchmark(BenchmarkArrayLookup)

// // 	// result2 := testing.Benchmark(BenchmarkDecodeStdStructMedium)
// // 	// fmt.Println(" \n Строка BenchmarkDecodeStdStructMedium", result2)

// // 	// result3 := testing.Benchmark(BenchmarkDecodeJsoniterStructMedium)
// // 	// fmt.Println(" \n Строка BenchmarkDecodeJsoniterStructMedium", result3)

// // 	result2 := testing.Benchmark(BenchmarkWrite)
// // 	fmt.Println(" \n Строка BenchmarkWrite", result2)

// // 	result3 := testing.Benchmark(BenchmarkCopy)
// // 	fmt.Println(" \n Строка BenchmarkCopy", result3)

// // 	// // Вывод результатов
// // 	// fmt.Println("byte", result1)
// // }

// func BenchmarkWriteTest(b *testing.B) {
// 	b.ReportAllocs()
// 	data := mediumFixture

// 	буферВОтправку := new(bytes.Buffer)
// 	b.ResetTimer()
// 	for i := 0; i < b.N; i++ {

// 		binary.Write(буферВОтправку, binary.LittleEndian, int32(len(data)))
// 		binary.Write(буферВОтправку, binary.LittleEndian, b)

// 	}

// 	// fmt.Printf("Строка %v", буферВОтправку)

// }
// func BenchmarkCopyTest(b *testing.B) {
// 	b.ReportAllocs()
// 	data := mediumFixture

// 	данные := make([]byte, len(data)+4)
// 	b.ResetTimer()
// 	for i := 0; i < b.N; i++ {

// 		binary.LittleEndian.PutUint32(данные, uint32(len(data)))
// 		copy(данные[4:], data)

// 	}

// }

// func BenchmarkDecodeStdStructMedium(b *testing.B) {
// 	b.ReportAllocs()
// 	var data MediumPayload
// 	for i := 0; i < b.N; i++ {
// 		json.Unmarshal(mediumFixture, &data)
// 	}
// 	fmt.Printf("Строка %v", data)

// }
// func BenchmarkDecodeJsoniterStructMedium(b *testing.B) {
// 	b.ReportAllocs()
// 	var data MediumPayload
// 	for i := 0; i < b.N; i++ {
// 		jsoniter.Unmarshal(mediumFixture, &data)
// 	}
// 	// fmt.Printf("Строка %+s ", data)
// }

// func BenchmarkJsoniter(b *testing.B) {
// 	var data MediumPayload
// 	json.Unmarshal(mediumFixture, &data)
// 	b.ReportAllocs()
// 	var bs []byte
// 	for i := 0; i < 1000; i++ {
// 		bs, _ = jsoniter.Marshal(&data)

// 	}
// 	fmt.Printf("Строка %s", bs)
// }

// func BenchmarkJson(b *testing.B) {
// 	var data MediumPayload
// 	json.Unmarshal(mediumFixture, &data)
// 	b.ReportAllocs()
// 	var bs []byte
// 	for i := 0; i < 1000; i++ {
// 		bs, _ = json.Marshal(&data)

// 	}
// 	fmt.Printf("Строка JSON %s", bs)
// }
// func BenchmarkMapLookup(b *testing.B) {
// 	m := make(map[string]int)
// 	mas := make([]string, 1000)
// 	for i := 0; i < 1000; i++ {
// 		id := uuid.New()
// 		ид := id.String()
// 		m[ид] = i
// 		mas[i] = ид
// 	}

// 	// b.ResetTimer()
// 	for i := 0; i < b.N; i++ {
// 		_ = m[mas[i%1000]]
// 	}
// }

// func BenchmarkArrayLookup(b *testing.B) {

// 	a := make(map[uuid.UUID]int)
// 	mas := make([]uuid.UUID, 1000)
// 	for i := 0; i < 1000; i++ {
// 		id := uuid.New()
// 		a[id] = i
// 		mas[i] = id
// 	}

// 	// b.ResetTimer()
// 	for i := 0; i < b.N; i++ {
// 		_ = a[mas[i%1000]]
// 	}
// }

// var smallFixture []byte = []byte(`{
//     "st": 1,
//     "sid": 486,
//     "tt": "active",
//     "gr": 0,
//     "uuid": "de305d54-75b4-431b-adb2-eb6b9e546014",
//     "ip": "127.0.0.1",
//     "ua": "user_agent",
//     "tz": -6,
//     "v": 1
// }`)

// type SmallPayload struct {
// 	St   int    `json:"st"`
// 	Sid  int    `json:"-"`
// 	Tt   string `json:"-"`
// 	Gr   int    `json:"-"`
// 	Uuid string `json:"uuid"`
// 	Ip   string `json:"-"`
// 	Ua   string `json:"ua"`
// 	Tz   int    `json:"tz"`
// 	V    int    `json:"-"`
// }

// // Reponse from Clearbit API. Size: 2.4kb
// var mediumFixture []byte = []byte(`{
//   "person": {
//     "id": "d50887ca-a6ce-4e59-b89f-14f0b5d03b03",
//     "name": {
//       "fullName": "Leonid Bugaev",
//       "givenName": "Leonid",
//       "familyName": "Bugaev"
//     },
//     "email": "leonsbox@gmail.com",
//     "gender": "male",
//     "location": "Saint Petersburg, Saint Petersburg, RU",
//     "geo": {
//       "city": "Saint Petersburg",
//       "state": "Saint Petersburg",
//       "country": "Russia",
//       "lat": 59.9342802,
//       "lng": 30.3350986
//     },
//     "bio": "Senior engineer at Granify.com",
//     "site": "http://flickfaver.com",
//     "avatar": "https://d1ts43dypk8bqh.cloudfront.net/v1/avatars/d50887ca-a6ce-4e59-b89f-14f0b5d03b03",
//     "employment": {
//       "name": "www.latera.ru",
//       "title": "Software Engineer",
//       "domain": "gmail.com"
//     },
//     "facebook": {
//       "handle": "leonid.bugaev"
//     },
//     "github": {
//       "handle": "buger",
//       "id": 14009,
//       "avatar": "https://avatars.githubusercontent.com/u/14009?v=3",
//       "company": "Granify",
//       "blog": "http://leonsbox.com",
//       "followers": 95,
//       "following": 10
//     },
//     "twitter": {
//       "handle": "flickfaver",
//       "id": 77004410,
//       "bio": null,
//       "followers": 2,
//       "following": 1,
//       "statuses": 5,
//       "favorites": 0,
//       "location": "",
//       "site": "http://flickfaver.com",
//       "avatar": null
//     },
//     "linkedin": {
//       "handle": "in/leonidbugaev"
//     },
//     "googleplus": {
//       "handle": null
//     },
//     "angellist": {
//       "handle": "leonid-bugaev",
//       "id": 61541,
//       "bio": "Senior engineer at Granify.com",
//       "blog": "http://buger.github.com",
//       "site": "http://buger.github.com",
//       "followers": 41,
//       "avatar": "https://d1qb2nb5cznatu.cloudfront.net/users/61541-medium_jpg?1405474390"
//     },
//     "klout": {
//       "handle": null,
//       "score": null
//     },
//     "foursquare": {
//       "handle": null
//     },
//     "aboutme": {
//       "handle": "leonid.bugaev",
//       "bio": null,
//       "avatar": null
//     },
//     "gravatar": {
//       "handle": "buger",
//       "urls": [
//       ],
//       "avatar": "http://1.gravatar.com/avatar/f7c8edd577d13b8930d5522f28123510",
//       "avatars": [
//         {
//           "url": "http://1.gravatar.com/avatar/f7c8edd577d13b8930d5522f28123510",
//           "type": "thumbnail"
//         }
//       ]
//     },
//     "fuzzy": false
//   },
//   "company": null
// }`)

// type CBAvatar struct {
// 	Url string `json:"url"`
// }

// type CBGravatar struct {
// 	Avatars []*CBAvatar `json:"avatars"`
// }

// type CBGithub struct {
// 	Followers int `json:"followers"`
// }

// type CBName struct {
// 	FullName string `json:"fullName"`
// }

// type CBPerson struct {
// 	Name     *CBName     `json:"name"`
// 	Github   *CBGithub   `json:"github"`
// 	Gravatar *CBGravatar `json:"gravatar"`
// }

// type MediumPayload struct {
// 	Person  *CBPerson `json:"person"`
// 	Company string    `json:"compnay"`
// }

// type DSUser struct {
// 	Username string
// }

// type DSTopic struct {
// 	Id   int    `json:"-"`
// 	Slug string `json:"-"`
// }

// type DSTopicsList struct {
// 	Topics        []struct{} `json:"topics"`
// 	MoreTopicsUrl string     `json:"-"`
// }

// type LargePayload struct {
// 	Users  []*DSUser     `json:"-"`
// 	Topics *DSTopicsList `json:"topics"`
// }
