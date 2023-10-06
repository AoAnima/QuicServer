package main

import (
	"fmt"
	"strconv"
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
	for i := 0; i < b.N; i++ {
		_ = m[strconv.Itoa(i%1000)]
	}
}

func BenchmarkArrayLookup(b *testing.B) {
	a := [1000][uuid.UUID]uuid.UUID{}
	for i := 0; i < 1000; i++ {
		id := uuid.New()
		a[id]=id
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = a[i%1000]
	}
}
