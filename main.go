package main

import (
	"strings"
	"rqueue/config"
	"rqueue/worker"
	"rqueue/queue"
	"fmt"
	"flag"
	"os"
	"bufio"
	"runtime"
)

func main() {
	var (
		conf string
		test string
		nCPU int
	)
	flag.StringVar(&conf, "c", "./config.toml", "Path to config file")
	flag.BoolVar(&config.Verbose, "v", false, "Verbose")
	flag.IntVar(&nCPU, "n", 1, "Maxprocs (N-system threads)")
	flag.Parse()

	max := runtime.NumCPU()
	runtime.GOMAXPROCS(nCPU)
	fmt.Printf("Number of CPUs: %d/%d\n", nCPU, max)

	test = flag.Arg(0)
	if test == "" {
		fmt.Fprintf(os.Stderr, "No testfile given\n")
		os.Exit(1)
	}
	tests, e := os.Open(test)
	if e != nil {
		panic(e)		
	}
	defer tests.Close()

	// arg0 = msgids
	// -c = servers
	if e := config.Open(conf); e != nil {
		panic(e)
	}
	fmt.Printf("%+v\n", config.C)
	if e := worker.Init(config.C.Servers); e != nil {
		panic(e)
	}

	// fake traffic
	scanner := bufio.NewScanner(tests)
    for scanner.Scan() {
    	msgid := scanner.Text()
    	tok := strings.Split(msgid, " ")
    	if len(tok) != 2 {
    		// Skip line
    		continue
    	}
        go queue.Add(tok[0], tok[1])
    }

    if e := scanner.Err(); e != nil {
        panic(e)
    }

    queue.Wait()
    fmt.Printf("Perf %+v", queue.Stats())
}