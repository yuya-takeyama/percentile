package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"sort"
	"strconv"

	flags "github.com/jessevdk/go-flags"
	"github.com/yuya-takeyama/argf"
)

const AppName = "numstat"

type Options struct {
	ShowVersion bool   `short:"v" long:"version" description:"Show version"`
	Algorithm   string `short:"a" long:"algorithm" description:"Algorithm to use (simple, linear-interpolation)" default:"linear-interpolation"`
}

var opts Options

var numbers sort.Float64Slice

func main() {
	parser := flags.NewParser(&opts, flags.Default^flags.PrintErrors)
	parser.Name = AppName
	parser.Usage = "[OPTIONS] FILES..."

	args, err := parser.Parse()
	if err != nil {
		fmt.Fprint(os.Stderr, err)
		return
	}

	r, err := argf.From(args)
	if err != nil {
		panic(err)
	}

	err = percentile(r, os.Stdout, os.Stderr, opts)
	if err != nil {
		panic(err)
	}
}

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

	printPercentileN(stdout, &numbers, l, 50, opts.Algorithm)
	printPercentileN(stdout, &numbers, l, 66, opts.Algorithm)
	printPercentileN(stdout, &numbers, l, 75, opts.Algorithm)
	printPercentileN(stdout, &numbers, l, 80, opts.Algorithm)
	printPercentileN(stdout, &numbers, l, 90, opts.Algorithm)
	printPercentileN(stdout, &numbers, l, 95, opts.Algorithm)
	printPercentileN(stdout, &numbers, l, 98, opts.Algorithm)
	printPercentileN(stdout, &numbers, l, 99, opts.Algorithm)
	printPercentileN(stdout, &numbers, l, 100, opts.Algorithm)

	return nil
}

func percentileN(numbers *sort.Float64Slice, l, n int) float64 {
	i := l*n/100 - 1
	ns := *numbers

	return ns[i]
}

// percentileNLinearInterpolation calculates percentile with linear interpolation
// This provides more accurate results, especially for datasets where exact percentiles
// fall between two data points
func percentileNLinearInterpolation(numbers *sort.Float64Slice, l, n int) float64 {
	if l == 0 {
		return 0
	}
	if l == 1 {
		return (*numbers)[0]
	}

	ns := *numbers

	// Calculate the rank using linear interpolation formula
	// rank = p/100 * (n-1) where p is the percentile and n is the number of elements
	rank := float64(n) / 100.0 * float64(l-1)

	lower := int(rank)
	upper := lower + 1

	// Handle edge cases
	if lower < 0 {
		lower = 0
	}
	if upper >= l {
		return ns[l-1]
	}

	// Linear interpolation between the two adjacent values
	fraction := rank - float64(lower)
	return ns[lower] + fraction*(ns[upper]-ns[lower])
}

func printPercentileN(w io.Writer, numbers *sort.Float64Slice, l, n int, algorithm string) {
	var value float64
	switch algorithm {
	case "simple":
		value = percentileN(numbers, l, n)
	case "linear-interpolation":
		value = percentileNLinearInterpolation(numbers, l, n)
	default:
		value = percentileNLinearInterpolation(numbers, l, n)
	}
	fmt.Fprintf(w, "%d%%:\t%s\n", n, strconv.FormatFloat(value, 'g', 16, 64))
}
