package main

import (
	"bufio"
	"fmt"
	"io"
	"sort"
	"strconv"
)

func percentile(r io.Reader, stdout io.Writer, stderr io.Writer, opts Options) error {
	if opts.ShowVersion {
		io.WriteString(stdout, fmt.Sprintf("%s v%s, build %s\n", AppName, Version, GitCommit))
		return nil
	}

	reader := bufio.NewReader(r)
	var line []byte
	var err error
	for {
		if line, _, err = reader.ReadLine(); err == io.EOF {
			break
		} else if err != nil {
			panic(err)
		}

		f, convErr := strconv.ParseFloat(string(line), 64)
		if convErr != nil {
			fmt.Fprintf(stderr, "number conversion error: %s\n", convErr)
			continue
		}

		numbers = append(numbers, f)
	}

	sort.Sort(numbers)
	l := len(numbers)

	if l > 2 {
		deletePreviouslines(stdout, 10)
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
	fmt.Fprintf(w, "%d%%:\t%s\n", n, strconv.FormatFloat(percentileN(numbers, l, n), 'g', 16, 64))
}

func deletePreviouslines(w io.Writer, count int) {
	for i := 0; i < count; i++ {
		fmt.Fprintf(w, "\r\b")
	}
}
