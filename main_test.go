package main

import (
	"bytes"
	"io"
	"math"
	"sort"
	"strings"
	"testing"
)

const float64EqualityThreshold = 1e-9

func almostEqual(a, b float64) bool {
	return math.Abs(a-b) <= float64EqualityThreshold
}

func TestPercentileN(t *testing.T) {
	tests := []struct {
		name       string
		numbers    []float64
		percentile int
		expected   float64
	}{
		{
			name:       "simple dataset 50th percentile",
			numbers:    []float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
			percentile: 50,
			expected:   5,
		},
		{
			name:       "simple dataset 90th percentile",
			numbers:    []float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
			percentile: 90,
			expected:   9,
		},
		{
			name:       "simple dataset 100th percentile",
			numbers:    []float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
			percentile: 100,
			expected:   10,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			numbers := sort.Float64Slice(tt.numbers)
			sort.Sort(numbers)
			result := percentileN(&numbers, len(numbers), tt.percentile)
			if result != tt.expected {
				t.Errorf("percentileN() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestPercentileNLinearInterpolation(t *testing.T) {
	tests := []struct {
		name       string
		numbers    []float64
		percentile int
		expected   float64
	}{
		{
			name:       "simple dataset 50th percentile",
			numbers:    []float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
			percentile: 50,
			expected:   5.5,
		},
		{
			name:       "simple dataset 90th percentile",
			numbers:    []float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
			percentile: 90,
			expected:   9.1,
		},
		{
			name:       "simple dataset 95th percentile",
			numbers:    []float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
			percentile: 95,
			expected:   9.55,
		},
		{
			name:       "simple dataset 100th percentile",
			numbers:    []float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
			percentile: 100,
			expected:   10,
		},
		{
			name:       "empty dataset",
			numbers:    []float64{},
			percentile: 50,
			expected:   0,
		},
		{
			name:       "single element",
			numbers:    []float64{5},
			percentile: 50,
			expected:   5,
		},
		{
			name:       "two elements 50th percentile",
			numbers:    []float64{1, 2},
			percentile: 50,
			expected:   1.5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			numbers := sort.Float64Slice(tt.numbers)
			sort.Sort(numbers)
			result := percentileNLinearInterpolation(&numbers, len(numbers), tt.percentile)
			if !almostEqual(result, tt.expected) {
				t.Errorf("percentileNLinearInterpolation() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestPercentileFunction(t *testing.T) {
	tests := []struct {
		name             string
		input            string
		opts             Options
		expectedContains []string
	}{
		{
			name:  "basic input with default algorithm",
			input: "1\n2\n3\n4\n5\n6\n7\n8\n9\n10\n",
			opts: Options{
				Algorithm: "linear-interpolation",
			},
			expectedContains: []string{"50%:", "5.5", "90%:", "9.1", "100%:", "10"},
		},
		{
			name:  "basic input with simple algorithm",
			input: "1\n2\n3\n4\n5\n6\n7\n8\n9\n10\n",
			opts: Options{
				Algorithm: "simple",
			},
			expectedContains: []string{"50%:", "5", "90%:", "9", "100%:", "10"},
		},
		{
			name:  "version flag",
			input: "",
			opts: Options{
				ShowVersion: true,
			},
			expectedContains: []string{"numstat v0.0.1"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset global numbers slice for each test
			numbers = sort.Float64Slice{}

			reader := strings.NewReader(tt.input)
			stdout := &bytes.Buffer{}
			stderr := &bytes.Buffer{}

			err := percentile(reader, stdout, stderr, tt.opts)
			if err != nil {
				t.Fatalf("percentile() error = %v", err)
			}

			output := stdout.String()
			for _, expected := range tt.expectedContains {
				if !strings.Contains(output, expected) {
					t.Errorf("percentile() output missing %q, got: %q", expected, output)
				}
			}
		})
	}
}

func TestPercentileWithInvalidInput(t *testing.T) {
	// Reset global numbers slice
	numbers = sort.Float64Slice{}

	input := "1\n2\ninvalid\n3\n4\n5\n"
	reader := strings.NewReader(input)
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}

	opts := Options{
		Algorithm: "linear-interpolation",
	}

	err := percentile(reader, stdout, stderr, opts)
	if err != nil {
		t.Fatalf("percentile() error = %v", err)
	}

	// Check that error message was written to stderr
	if !strings.Contains(stderr.String(), "number conversion error") {
		t.Errorf("Expected error message in stderr, got: %q", stderr.String())
	}

	// Should still process valid numbers
	if stdout.Len() == 0 {
		t.Error("Expected output for valid numbers")
	}
}

func TestPrintPercentileN(t *testing.T) {
	numbers := sort.Float64Slice([]float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10})
	sort.Sort(numbers)

	tests := []struct {
		name       string
		percentile int
		algorithm  string
		expected   string
	}{
		{
			name:       "50th percentile with linear interpolation",
			percentile: 50,
			algorithm:  "linear-interpolation",
			expected:   "50%:\t5.5\n",
		},
		{
			name:       "90th percentile with simple algorithm",
			percentile: 90,
			algorithm:  "simple",
			expected:   "90%:\t9\n",
		},
		{
			name:       "default to linear interpolation",
			percentile: 75,
			algorithm:  "unknown",
			expected:   "75%:\t7.75\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			printPercentileN(&buf, &numbers, len(numbers), tt.percentile, tt.algorithm)

			if buf.String() != tt.expected {
				t.Errorf("printPercentileN() = %q, want %q", buf.String(), tt.expected)
			}
		})
	}
}

// Benchmark tests
func BenchmarkPercentileNLinearInterpolation(b *testing.B) {
	numbers := make(sort.Float64Slice, 1000)
	for i := 0; i < 1000; i++ {
		numbers[i] = float64(i)
	}
	sort.Sort(numbers)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		percentileNLinearInterpolation(&numbers, len(numbers), 95)
	}
}

func BenchmarkPercentileN(b *testing.B) {
	numbers := make(sort.Float64Slice, 1000)
	for i := 0; i < 1000; i++ {
		numbers[i] = float64(i)
	}
	sort.Sort(numbers)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		percentileN(&numbers, len(numbers), 95)
	}
}

func BenchmarkPercentile(b *testing.B) {
	input := strings.Repeat("1\n2\n3\n4\n5\n6\n7\n8\n9\n10\n", 100)
	opts := Options{
		Algorithm: "linear-interpolation",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		numbers = sort.Float64Slice{}
		reader := strings.NewReader(input)
		percentile(reader, io.Discard, io.Discard, opts)
	}
}
