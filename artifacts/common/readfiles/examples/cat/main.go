package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"

	"common/readfiles"
)

func usage() {
	fmt.Fprintf(os.Stderr, `Usage:
    %[1]s [OPTION]... [FILE]...

Description:
    Print FILE(s). With no FILE, or when FILE is -, read standard input. By default, print results to stdout and log to stderr.

Options:
`, os.Args[0])
	flag.PrintDefaults()
}

func main() {
	flag.Usage = usage
	logFile := flag.String("log", "", "log to file.  (default stderr)")
	outputFile := flag.String("out", "", "output to file.  (default stdout)")
	flush := flag.Bool("flush", true, "flush after every output.")
	flag.Parse()

	// log, intentionally make it blocking to make sure it got
	// initiliazed before other parts using it
	if *logFile != "" {
		f, err := os.OpenFile(*logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		if err != nil {
			log.Panicln("failed to open log file", err)
		}
		defer f.Close()
		log.SetOutput(f)
	}

	// output
	var f *os.File
	var err error
	if *outputFile == "" {
		f = os.Stdout
	} else {
		f, err = os.OpenFile(*outputFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		if err != nil {
			log.Panicln("failed to open output file", err)
		}
	}
	defer f.Close()
	w := bufio.NewWriter(f)

	lines := readfiles.ReadFiles(flag.Args())

	for line := range lines {
		w.WriteString(line + "\n")
		if *flush {
			w.Flush()
		}
	}
	w.Flush()
}
