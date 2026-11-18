package main

import (
	"fmt"
	"math"
	"math/rand"
	"sync"
	"time"
)

const (
	// SIZE is the number of generated elements.
	SIZE = 100_000_000
	// CHUNKS is the number of concurrent workers used by maxChunks.
	CHUNKS = 8
)

// generateRandomElements returns a slice of size positive random integers.
// For size <= 0 it returns an empty slice.
func generateRandomElements(size int) []int {
	if size <= 0 {
		return []int{}
	}

	data := make([]int, size)
	for i := range data {
		// rand.Intn returns a value in [0, math.MaxInt), so adding 1
		// produces a positive integer in [1, math.MaxInt].
		data[i] = rand.Intn(math.MaxInt) + 1
	}

	return data
}

// maximum returns the maximum value in data.
// For an empty slice it returns 0.
func maximum(data []int) int {
	if len(data) == 0 {
		return 0
	}

	maxValue := data[0]
	for _, value := range data[1:] {
		if value > maxValue {
			maxValue = value
		}
	}

	return maxValue
}

// maxChunks returns the maximum value in data by splitting it into CHUNKS
// parts and searching each part in a separate goroutine.
//
// The last part includes the remainder when len(data) is not divisible by
// CHUNKS. For empty input it returns 0. For input smaller than CHUNKS it
// falls back to a sequential search, which also keeps negative values correct.
func maxChunks(data []int) int {
	if len(data) == 0 {
		return 0
	}
	if len(data) < CHUNKS {
		return maximum(data)
	}

	chunkSize := len(data) / CHUNKS
	maxValues := make([]int, CHUNKS)

	var wg sync.WaitGroup
	for i := 0; i < CHUNKS; i++ {
		start := i * chunkSize
		end := start + chunkSize
		if i == CHUNKS-1 {
			end = len(data)
		}

		wg.Add(1)
		go func(part []int, index int) {
			defer wg.Done()
			// Each goroutine writes only to its own index, so no mutex is needed.
			maxValues[index] = maximum(part)
		}(data[start:end], i)
	}

	wg.Wait()
	return maximum(maxValues)
}

func main() {
	fmt.Printf("Генерируем %d целых чисел\n", SIZE)
	data := generateRandomElements(SIZE)

	fmt.Println("Ищем максимальное значение в один поток")
	start := time.Now()
	maxValue := maximum(data)
	elapsed := time.Since(start).Microseconds()
	fmt.Printf("Максимальное значение элемента: %d\nВремя поиска: %d мкс\n", maxValue, elapsed)

	fmt.Printf("Ищем максимальное значение в %d потоков\n", CHUNKS)
	start = time.Now()
	maxValueChunks := maxChunks(data)
	elapsed = time.Since(start).Microseconds()
	fmt.Printf("Максимальное значение элемента: %d\nВремя поиска: %d мкс\n", maxValueChunks, elapsed)
}
