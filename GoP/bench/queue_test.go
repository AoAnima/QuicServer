// package main

// import (
// 	"strconv"
// 	"sync"
// 	"testing"

// 	. "aoanima.ru/ConnQuic"
// 	quic "github.com/quic-go/quic-go"
// )

// type MyStruct struct {
// 	Field int
// }

// func BenchmarkStruct(b *testing.B) {
// 	m := make(map[string]MyStruct)
// 	for i := 0; i < b.N; i++ {
// 		m["key"+strconv.Itoa(i)] = MyStruct{i}
// 	}
// }

// func BenchmarkPointer(b *testing.B) {
// 	m := make(map[string]*MyStruct)
// 	for i := 0; i < b.N; i++ {
// 		m["key"+strconv.Itoa(i)] = &MyStruct{}
// 		m["key"+strconv.Itoa(i)].Field = i

// 		// m["key"] = &MyStruct{i}
// 	}
// }

// func BenchmarkОчередьПотоковКанал_Взять_Вернуть(b *testing.B) {
// 	очередь := НоваяОчередьПотоковКанал(1)
// 	var поток quic.Stream // Предполагается, что поток инициализирован

// 	b.ResetTimer()
// 	for i := 0; i < b.N; i++ {
// 		очередь.Вернуть(поток)
// 		очередь.Взять()
// 	}
// }

// func BenchmarkОчередьПотоков_Взять_Вернуть(b *testing.B) {
// 	очередь := НоваяОчередьПотоков()
// 	var поток quic.Stream // Предполагается, что поток инициализирован

// 	b.ResetTimer()
// 	for i := 0; i < b.N; i++ {
// 		очередь.Вернуть(поток)
// 		очередь.Взять()
// 	}
// }

// func BenchmarkОчередьПотоковКанал_Взять_Вернуть_MultiGoroutine(b *testing.B) {
// 	очередь := НоваяОчередьПотоковКанал(100)
// 	var поток quic.Stream // Предполагается, что поток инициализирован

// 	var wg sync.WaitGroup
// 	b.ResetTimer()
// 	for i := 0; i < b.N; i++ {
// 		wg.Add(1)
// 		go func() {
// 			defer wg.Done()
// 			for i := 0; i < 1000; i++ {

// 				очередь.Вернуть(поток)
// 				поток, _ = очередь.Взять()
// 				очередь.Вернуть(поток)
// 			}
// 		}()
// 	}
// 	wg.Wait()
// }

// // 2129448               559.4 ns/op            80 B/op          3 allocs/op
// // 1718689               707.7 ns/op            48 B/op          1 allocs/op
// func BenchmarkОчередьПотоков_Взять_Вернуть_MultiGoroutine(b *testing.B) {
// 	очередь := НоваяОчередьПотоков()
// 	var поток quic.Stream // Предполагается, что поток инициализирован

// 	var wg sync.WaitGroup
// 	b.ResetTimer()
// 	for i := 0; i < b.N; i++ {
// 		wg.Add(1)
// 		go func() {
// 			defer wg.Done()
// 			for i := 0; i < 1000; i++ {
// 				очередь.Вернуть(поток)
// 				поток = очередь.Взять()
// 				очередь.Вернуть(поток)
// 			}
// 		}()
// 	}
// 	wg.Wait()
// }
