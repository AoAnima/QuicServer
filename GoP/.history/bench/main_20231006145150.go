package main

import (
	"fmt"
	"strconv"
	"testing"
)

func main() {
	// Запуск бенчмарка
	result := testing.Benchmark(BenchmarkMapLookup)

	// Вывод результатов
	fmt.Println("Cnhjrf", result)
	result1 := testing.Benchmark(BenchmarkArrayLookup)

	// Вывод результатов
	fmt.Println(result1)
}
func BenchmarkMapLookup(b *testing.B) {
	m := make(map[string]int)
	for i := 0; i < 1000; i++ {
		m[strconv.Itoa(i)] = i
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = m[strconv.Itoa(i%1000)]
	}
}

func BenchmarkArrayLookup(b *testing.B) {
	a := [1000][16]byte{}
	for i := 0; i < 1000; i++ {
		copy(a[i][:], strconv.Itoa(i))
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = a[i%1000]
	}
}
