package main

import (
	"bytes"
	"io"
	"os/exec"
	"sync"
)

func Ntimes(cnt int, cmdName string, cmdArgs []string, stdin io.Reader, stderr io.Writer, stdoutCh chan io.ReadWriter, parallels int) {
	var wg sync.WaitGroup

	sema := make(chan bool, parallels)

	for i := 0; i < cnt; i++ {
		wg.Add(1)
		go func() {
			sema <- true

			defer func() {
				wg.Done()
				<-sema
			}()

			stdoutBuffer := new(bytes.Buffer)

			cmd := exec.Command(cmdName, cmdArgs...)
			cmd.Stdin = stdin
			cmd.Stdout = stdoutBuffer
			cmd.Stderr = stderr
			err := cmd.Run()

			if err != nil {
				panic(err)
			}

			stdoutCh <- stdoutBuffer
		}()
	}

	wg.Wait()
}

func printer(stdout io.Writer, stdoutCh chan io.ReadWriter, exitCh chan bool) {
	for {
		select {
		case r := <-stdoutCh:
			io.Copy(stdout, r)
		case <-exitCh:
			return
		}
	}
}
