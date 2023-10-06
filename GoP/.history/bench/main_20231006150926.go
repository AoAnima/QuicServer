package main

import (
	"fmt"
	"testing"

	"github.com/google/uuid"
)

func main() {
	// Запуск бенчмарка
	result := testing.Benchmark(BenchmarkMapLookup)

	// Вывод результатов
	fmt.Println("Строка", result)
	result1 := testing.Benchmark(BenchmarkArrayLookup)

	// Вывод результатов
	fmt.Println("byte", result1)
}
func BenchmarkMapLookup(b *testing.B) {
	m := make(map[string]int)
	mas := make([]string, 1000)
	for i := 0; i < 1000; i++ {
		id := uuid.New()
		ид := id.String()
		m[ид] = i
		mas[i] = ид
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		fmt.Println("byte", i%1000)
		_ = m[mas[i%1000]]
	}
}

func BenchmarkArrayLookup(b *testing.B) {

	a := make(map[uuid.UUID]int)
	mas := make([]uuid.UUID, 1000)
	for i := 0; i < 1000; i++ {
		id := uuid.New()
		a[id] = i
		mas[i] = id
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = a[mas[i%1000]]
	}
}
