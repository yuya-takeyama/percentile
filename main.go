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
	ShowVersion bool `short:"v" long:"version" description:"Show version"`
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

	printPercentileN(stdout, &numbers, l, 50)
	printPercentileN(stdout, &numbers, l, 66)
	printPercentileN(stdout, &numbers, l, 75)
	printPercentileN(stdout, &numbers, l, 80)
	printPercentileN(stdout, &numbers, l, 90)
	printPercentileN(stdout, &numbers, l, 95)
	printPercentileN(stdout, &numbers, l, 98)
	printPercentileN(stdout, &numbers, l, 99)
	printPercentileN(stdout, &numbers, l, 100)

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
