package main

import (
	"fmt"
	"io"
	"sort"
	"strconv"
)

func percentile(r float64, stdout io.Writer, stderr io.Writer, opts Options) error {
	if opts.ShowVersion {
		io.WriteString(stdout, fmt.Sprintf("%s v%s, build %s\n", AppName, Version, GitCommit))
		return nil
	}

	numbers = append(numbers, r)

	sort.Sort(numbers)
	l := len(numbers)

	if l > 2 {
		deletePreviouslines(stdout, 9)
	}
	if l > 1 {
		for _, value := range values {
			printPercentileN(stdout, &numbers, l, value)
		}
	}

	return nil
}

func percentileN(numbers *sort.Float64Slice, l, n int) float64 {
	i := l*n/100 - 1
	ns := *numbers

	return ns[i]
}

func printPercentileN(w io.Writer, numbers *sort.Float64Slice, l, n int) {
	fmt.Fprintf(w, "%d%%:\t%s\n", n, strconv.FormatFloat(percentileN(numbers, l, n), 'g', 5, 64))
}

func deletePreviouslines(w io.Writer, count int) {
	for i := 0; i < count; i++ {
		fmt.Fprint(w, "\033[F\033[K")
	}
}
