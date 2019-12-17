package main

import (
	"fmt"
	flags "github.com/jessevdk/go-flags"
	"io"
	"os"
	"sort"
	"strconv"
)

const AppName = "numstat"

type Options struct {
	Parallels   int  `short:"p" long:"parallels" description:"Parallel degree of execution" default:"1"`
	ShowVersion bool `short:"v" long:"version" description:"Show version"`
}

var opts Options

var numbers sort.Float64Slice

var values = [9]int{50, 66, 75, 80, 90, 95, 98, 99, 100}

func main() {
	parser := flags.NewParser(&opts, flags.Default^flags.PrintErrors)
	parser.Name = AppName
	parser.Usage = "N [OPTIONS] -- COMMAND"

	args, err := parser.Parse()

	if err != nil {
		fmt.Fprint(os.Stderr, err)
		return
	}

	if opts.ShowVersion {
		io.WriteString(os.Stdout, fmt.Sprintf("%s v%s, build %s\n", AppName, Version, GitCommit))
		return
	}

	cnt, err := strconv.Atoi(args[0])
	cmdName := args[1]
	cmdArgs := args[2:]

	stdoutCh := make(chan io.ReadWriter)
	exitCh := make(chan bool)

	if err != nil {
		panic(err)
	}

	go wrapper(os.Stdout, stdoutCh, exitCh)

	Ntimes(cnt, cmdName, cmdArgs, os.Stdin, os.Stderr, stdoutCh, opts.Parallels)

//	r, err := argf.From(args)
//	if err != nil {
//		panic(err)
//	}

//	err = percentile(r, os.Stdout, os.Stderr, opts)
//	if err != nil {
//		panic(err)
//	}

	exitCh <- true
}

func wrapper(stdout io.Writer, stdoutCh chan io.ReadWriter, exitCh chan bool) {
	for {
		select {
		case r := <-stdoutCh:
			err := percentile(r, stdout, os.Stderr, opts)
				if err != nil {
					panic(err)
				}
		case <-exitCh:
			return
		}
	}
}
