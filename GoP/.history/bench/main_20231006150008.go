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
	for i := 0; i < 1000; i++ {
		id := uuid.New()
		ид := id.String()
		m[ид] = i
	}

	b.ResetTimer()
	for k, _ := range m {
		_ = m[k]
	}
}

func BenchmarkArrayLookup(b *testing.B) {

	a := make(map[uuid.UUID]int)
	mas := make([]uuid.UUID, 1000)
	for i := 0; i < 1000; i++ {
		id := uuid.New()
		a[id] = i
	}

	b.ResetTimer()
	for k, _ := range a {
		_ = a[k]
	}
}
