package main

import (
	"encoding/csv"
	"encoding/hex"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"runtime/pprof"
	"strconv"
	"sync"
	"time"

	"common/parseipportargs"
	"common/readfiles"
)

func usage() {
	fmt.Fprintf(os.Stderr, `Usage:
    %[1]s [OPTION]... [FILE]...

Description:
    Test if payloads in FILE(s) are affected by the dynamic blocking. With no FILE, read standard input. By default, print results to stdout and log to stderr.

    * If an ip:port accepts -repeat number of random connections (which doesn't have to be consecutive), then mark it as unaffected.
    * If -try number of consecutive connections to an ip:port all timeout, then mark it as affected.
    * If the total number of successful connections is zero, mark it as unknown (possibly closed/filtered port).
    * If any other error occured, mark it as unknown.

Wait for -timeout second between each connection to ip:port. If timeout already occured, slow down by waiting for -wait seconds between each connection.

Examples:
    Test if payload 00112233 is affected by sending random traffic to the sink server at port 443
	echo "00112233" | %[1]s -host 1.1.1.1 -p 443 
Options:
`, os.Args[0])
	flag.PrintDefaults()
}

type job struct {
	Payload string
	Port    int
}

func worker(id int, jobs chan job, results chan<- []string, dialer *net.Dialer, addr string) {
	for j := range jobs {
		// skip empty addr
		//if len(ad) == 0 {
		//	continue
		//}

		payloadBytes, err := hex.DecodeString(j.Payload)
		if err != nil {
			log.Panic(err)
		}
		log.Println("worker", id, "is testing:", j.Payload)

		code := ""
		// countSuccess the number of connections with random payload to addr
		countSuccess := 0
		consecutiveTimeout := 0
		totalTimeout := 0
		host := fmt.Sprintf("%s:%d", addr, j.Port)
		for countSuccess < repeat && consecutiveTimeout < maxNumTimeout {
			conn, err := dialer.Dial("tcp", host)
			if err != nil {
				if e, ok := err.(net.Error); ok && e.Timeout() {
					code = "timeout"
					consecutiveTimeout++
					totalTimeout++
					time.Sleep(*wait)
					continue
				} else {
					code = err.Error()
					break
				}
			}
			// reset when consecutiveTimeout when connection suceed
			consecutiveTimeout = 0
			countSuccess++
			code = ""

			// send
			conn.Write(payloadBytes)

			go func() {
				time.Sleep(*sleep)
				defer conn.Close()
			}()

			time.Sleep(*interval)
		}

		affected := ""
		if countSuccess == 0 {
			// this is when the port is closed
			affected = "unknown"
		} else {
			if consecutiveTimeout == maxNumTimeout {
				affected = "true"
			} else if countSuccess == repeat {
				affected = "false"
			} else {
				affected = "unknown"
			}
		}

		endTime := strconv.FormatInt(time.Now().Unix(), 10)
		results <- []string{endTime, host, j.Payload, strconv.Itoa(countSuccess), strconv.Itoa(totalTimeout), strconv.Itoa(consecutiveTimeout), code, affected}
		log.Println("worker", id, "finished testing", j.Payload)
	}
}

// global variables
var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file.")
var timeout = flag.Duration("timeout", 6*time.Second, "timeout value of TCP connections.")
var interval = flag.Duration("interval", 1*time.Second, "time interval between each connection to a ip:port.")
var wait = flag.Duration("wait", 3*time.Second, "time interval between each connection, when a ip:port timeout.")
var sleep = flag.Duration("sleep", 0*time.Second, "time interval between sending a probe and closing the connection. This value doesn't affect the -interval between each connection.")
var repeat int
var maxNumTimeout int

func main() {
	flag.Usage = usage
	var maxNumWorkers int
	argPort := flag.String("p", "80", "comma-separated list of ports to which the program sends random payload. eg. 3000,4000-4002")
	flag.IntVar(&maxNumWorkers, "worker", 10, fmt.Sprintf("number of workers in parallel."))
	flag.IntVar(&repeat, "repeat", 25, "repeatedly make up to this number of connections to each ip:port.")
	flag.IntVar(&maxNumTimeout, "try", 5, "mark an ip:port as affected if this number of consecutive connections all timeout.")
	outputFile := flag.String("out", "", "output csv file.  (default stdout)")
	logFile := flag.String("log", "", "log to file.  (default stderr)")
	//payload := flag.String("payload", "dadd034913c52da75fd9f05dc76803917134808efed97ef8884f2151b712f60fed634f609f132033a15b77ed3ccaa2d20f5b", "payload of the probes in hex format")
	host := flag.String("host", "REDACTED_US_SERVER_IP", "host to send to")
	flag.Parse()

	// log, intentionally make it blocking to make sure it got
	// initiliazed before other parts using it
	if *logFile != "" {
		f, err := os.Create(*logFile)
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
		f, err = os.Create(*outputFile)
		if err != nil {
			log.Panicln("failed to open output file", err)
		}
	}
	defer f.Close()
	w := csv.NewWriter(f)

	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Panicln(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	/*payloadBytes, err := hex.DecodeString(*payload)
	if err != nil {
		log.Panic(err)
	}*/

	ports, err := parseipportargs.ParsePortArgs(*argPort)
	if err != nil {
		log.Panic(err)
	}
	// The channel capacity does not have to be equal to the
	// number of workers. It can be much smaller.
	payloads := make(chan job, 100)
	results := make(chan []string, 100)
	results <- []string{"endTime", "addr", "payload", "countSuccess", "totalTimeout", "consecutiveTimeout", "code", "affected"}

	go func() {
		portIdx := 0
		lines := readfiles.ReadFiles(flag.Args())
		for payload := range lines {
			//ips, err := parseipportargs.ParseIPArgs(ipString)
			payloads <- job{payload, ports[portIdx]}
			portIdx += 1
			if portIdx >= len(ports) {
				portIdx = 0
			}

			/*
				if err != nil {
					log.Panic(err)
				}
				for _, ip := range ips {
					for _, port := range ports {
						addrs <- net.JoinHostPort(ip.String(), strconv.Itoa(port))
					}
				}*/
		}
		close(payloads)
	}()

	dialer := &net.Dialer{
		Timeout: *timeout,
	}

	var wg sync.WaitGroup
	wg.Add(maxNumWorkers)
	for i := 0; i < maxNumWorkers; i++ {
		go func(id int) {
			defer wg.Done()
			worker(id, payloads, results, dialer, *host)
		}(i)
	}
	go func() {
		wg.Wait()
		close(results)
	}()

	for r := range results {
		err := w.Write(r)
		if err != nil {
			log.Panicln("error writing result sto file:", err, len(r))
		}
		w.Flush()
	}
}
