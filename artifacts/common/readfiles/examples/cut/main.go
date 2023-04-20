package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"os"
	"regexp"

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

// cut is a function that uses regex to split the input line into fields,
// and return the fields.
// Example input:
// Email: jo-mitten@hotmail.co.uk - Name: Joe Osullivan - ScreenName: JoeOsullivan2 - Followers: 263 - Created At: Sun Mar 04 14:35:29 +0000 2012
// Example output:
// jo-mitten@hotmail.co.uk, Joe Osullivan, JoeOsullivan2, 263, Sun Mar 04 14:35:29 +0000 2012
var pattern = `Email: (.*) - Name: (.*) - ScreenName: (.*) - Followers: (.*) - Created At: (.*)`
var re = regexp.MustCompile(pattern)

func cut(line string) []string {
	matches := re.FindStringSubmatch(line)
	if len(matches) != 6 {
		log.Println("failed to parse line", line)
		return nil
	}
	return matches[1:]
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
	w := csv.NewWriter(f)

	lines := readfiles.ReadFiles(flag.Args())

	for line := range lines {
		fields := cut(line)
		if fields == nil {
			continue
		}
		// write fields into a csv file with comma as separator
		if err := w.Write(fields); err != nil {
			log.Println("failed to write to csv", err)
		}

		if *flush {
			w.Flush()
		}
	}
	w.Flush()
}
