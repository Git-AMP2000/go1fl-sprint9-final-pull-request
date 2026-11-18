package main

import (
	"math"
	"sync"
	"testing"
)

func TestGenerateRandomElements(t *testing.T) {
	tests := []struct {
		name string
		size int
	}{
		{name: "zero size", size: 0},
		{name: "negative size", size: -1},
		{name: "large negative size", size: -1000},
		{name: "one element", size: 1},
		{name: "small slice", size: 10},
		{name: "larger slice", size: 1000},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data := generateRandomElements(tt.size)

			if tt.size <= 0 {
				if len(data) != 0 {
					t.Fatalf("expected empty slice for size %d, got length %d", tt.size, len(data))
				}
				return
			}

			if len(data) != tt.size {
				t.Fatalf("expected length %d, got %d", tt.size, len(data))
			}

			for i, value := range data {
				if value <= 0 {
					t.Fatalf("expected positive value at index %d, got %d", i, value)
				}
			}

			if tt.size > 1 {
				allEqual := true
				for _, value := range data[1:] {
					if value != data[0] {
						allEqual = false
						break
					}
				}
				if allEqual {
					t.Fatalf("expected random values, got %d identical elements", tt.size)
				}
			}
		})
	}
}

func TestMaximum(t *testing.T) {
	tests := []struct {
		name string
		data []int
		want int
	}{
		{name: "nil slice", data: nil, want: 0},
		{name: "empty slice", data: []int{}, want: 0},
		{name: "single positive element", data: []int{42}, want: 42},
		{name: "single zero element", data: []int{0}, want: 0},
		{name: "single negative element", data: []int{-7}, want: -7},
		{name: "multiple elements", data: []int{3, 1, 4, 1, 5, 9, 2, 6}, want: 9},
		{name: "maximum at the beginning", data: []int{100, 2, 3, 4}, want: 100},
		{name: "maximum at the end", data: []int{1, 2, 3, 100}, want: 100},
		{name: "all elements equal", data: []int{5, 5, 5, 5}, want: 5},
		{name: "negative values", data: []int{-10, -3, -42, -1}, want: -1},
		{name: "mixed values", data: []int{-5, 0, 12, -100, 7}, want: 12},
		{name: "large values", data: []int{math.MaxInt - 1, math.MaxInt, 0}, want: math.MaxInt},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := maximum(tt.data); got != tt.want {
				t.Fatalf("maximum(%v) = %d, want %d", tt.data, got, tt.want)
			}
		})
	}
}

func TestMaxChunks(t *testing.T) {
	tests := []struct {
		name string
		data []int
		want int
	}{
		{name: "nil slice", data: nil, want: 0},
		{name: "empty slice", data: []int{}, want: 0},
		{name: "single element", data: []int{42}, want: 42},
		{name: "less than chunks", data: []int{1, 9, 3}, want: 9},
		{name: "less than chunks with negative values", data: []int{-10, -3, -42}, want: -3},
		{name: "exactly chunks", data: []int{1, 2, 3, 4, 5, 6, 7, 8}, want: 8},
		{name: "exactly chunks with maximum in the middle", data: []int{1, 2, 3, 99, 5, 6, 7, 8}, want: 99},
		{name: "with remainder and maximum in the last part", data: []int{1, 2, 3, 4, 5, 6, 7, 8, 100}, want: 100},
		{name: "with remainder and maximum in the first part", data: []int{100, 2, 3, 4, 5, 6, 7, 8, 9}, want: 100},
		{name: "all elements equal", data: []int{7, 7, 7, 7, 7, 7, 7, 7, 7, 7}, want: 7},
		{name: "negative values with full chunks", data: []int{-10, -3, -42, -1, -9, -8, -7, -6}, want: -1},
		{name: "negative values with remainder", data: []int{-10, -3, -42, -1, -9, -8, -7, -6, -100}, want: -1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := maxChunks(tt.data); got != tt.want {
				t.Fatalf("maxChunks(%v) = %d, want %d", tt.data, got, tt.want)
			}
		})
	}
}

func TestMaxChunksMatchesMaximumOnRandomData(t *testing.T) {
	for _, size := range []int{CHUNKS, CHUNKS + 1, 100, 1000, 10_000} {
		data := generateRandomElements(size)
		want := maximum(data)
		if got := maxChunks(data); got != want {
			t.Fatalf("size %d: maxChunks = %d, maximum = %d", size, got, want)
		}
	}
}

func TestMaxChunksDoesNotModifyInput(t *testing.T) {
	data := generateRandomElements(1000)
	original := make([]int, len(data))
	copy(original, data)

	_ = maxChunks(data)

	for i := range data {
		if data[i] != original[i] {
			t.Fatalf("input was modified at index %d: got %d, want %d", i, data[i], original[i])
		}
	}
}

func TestMaxChunksConcurrentCalls(t *testing.T) {
	data := generateRandomElements(1000)
	want := maximum(data)

	const callers = 16
	var wg sync.WaitGroup
	results := make([]int, callers)

	for i := 0; i < callers; i++ {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()
			results[index] = maxChunks(data)
		}(i)
	}

	wg.Wait()

	for i, got := range results {
		if got != want {
			t.Fatalf("caller %d: maxChunks = %d, want %d", i, got, want)
		}
	}
}

func BenchmarkMaximum(b *testing.B) {
	data := generateRandomElements(1_000_000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = maximum(data)
	}
}

func BenchmarkMaxChunks(b *testing.B) {
	data := generateRandomElements(1_000_000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = maxChunks(data)
	}
}
