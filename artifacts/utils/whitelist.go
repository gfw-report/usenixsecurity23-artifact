package main

import (
	"crypto/rand"
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
    Test if prepending the payloads (represented in hex string) in FILE(s) can exempt the GFW's blocking of random traffic. With no FILE, read standard input. By default, print results to stdout and log to stderr.

Examples:
    Test if prepending "QUIT\r\n" (0x51, 0x55, 0x49, 0x54, 0x0d, 0x0a) can exempt the GFW's blocking of random traffic.
	echo -n "515549540d0a" | %[1]s
    Test if prepending "aaaaaa" (0x61, 0x61, 0x61, 0x61, 0x61, 0x61) can exempt the GFW's blocking of random traffic.
	echo -n "616161616161" | %[1]s
Options:
`, os.Args[0])
	flag.PrintDefaults()
}

func worker(id int, jobs chan jobPayload, addrs chan string, results chan<- []string, dialer *net.Dialer) {
	// job is the pre-payload to be tested
	for job := range jobs {

		payload := append(job.bytePrefix, job.byteSuffix...)
		log.Printf("worker %v is testing prefix: %s (% x) with length: %v", id, job.strPrefix, job.bytePrefix, len(payload))

		// pop a addr
		addr := <-addrs

		code := ""
		// countSuccess the number of connections with random payload to addr
		countSuccess := 0
		consecutiveTimeout := 0
		totalTimeout := 0
		for countSuccess < repeat && consecutiveTimeout < maxNumTimeout {
			conn, err := dialer.Dial("tcp", addr)
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
			conn.Write(payload)
			conn.Close()

			time.Sleep(*interval)
		}

		whitelisted := ""
		if countSuccess == 0 {
			// this is when the port is closed
			whitelisted = "unknown"
		} else {
			if consecutiveTimeout == maxNumTimeout {
				whitelisted = "false"
			} else if countSuccess == repeat {
				whitelisted = "true"
			} else {
				whitelisted = "unknown"
			}
		}

		go func() {
			if whitelisted != "true" {
				time.Sleep(*residual)
			}
			addrs <- addr
		}()

		endTime := strconv.FormatInt(time.Now().Unix(), 10)
		results <- []string{endTime, job.strPrefix, addr, strconv.Itoa(countSuccess), strconv.Itoa(totalTimeout), strconv.Itoa(consecutiveTimeout), code, whitelisted, strconv.Itoa(len(payload))}
		log.Printf("worker %v finished testing: %s (% x)", id, job.strPrefix, job.bytePrefix)
	}
}

// global variables
var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file.")
var timeout = flag.Duration("timeout", 6*time.Second, "timeout value of TCP connections.")
var interval = flag.Duration("interval", 1*time.Second, "time interval between each connection to a ip:port.")
var wait = flag.Duration("wait", 3*time.Second, "time interval between each connection, when a ip:port timeout.")
var residual = flag.Duration("redisual", 3*time.Minute, "duration of residual censorship when a ip:port blocked.")
var repeat int
var maxNumTimeout int

// struct for jobs channel
type jobPayload struct {
	strPrefix  string
	bytePrefix []byte
	byteSuffix []byte
}

func main() {
	flag.Usage = usage
	var maxNumWorkers int
	argIP := flag.String("dip", "127.0.0.1", "comma-separated list of destination IP addresses to which the program can send probes. eg. 1.1.1.1,2.2.2.2")
	argPort := flag.String("p", "10000-10200", "comma-separated list of available ports to which the program can send probes. eg. 3000,4000-4002")
	flag.IntVar(&maxNumWorkers, "worker", 10, fmt.Sprintf("number of workers in parallel."))
	flag.IntVar(&maxNumTimeout, "try", 5, "mark an ip:port as affected if this number of consecutive connections all timeout.")
	flag.IntVar(&repeat, "repeat", 25, "repeated make up to this number of connections to each ip:port.")
	outputFile := flag.String("out", "", "output csv file.  (default stdout)")
	logFile := flag.String("log", "", "log to file.  (default stderr)")
	payloadLen := flag.Int("max", 50, "length of payload to send.")
	increment := flag.Int("increment", 0, "increment for generating variable payload sizes. must be used in conjunction with increment. eg. -max -increment 10 will try the same prefix with total payload len 10, 20, ..., 100")
	suffix := flag.String("append", "", "appending a fixed random payload for each probe, in hex format")
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

	ips, err := parseipportargs.ParseIPArgs(*argIP)
	if err != nil {
		log.Panic(err)
	}

	ports, err := parseipportargs.ParsePortArgs(*argPort)
	if err != nil {
		log.Panic(err)
	}
	// The channel capacity does not have to be equal to the
	// number of workers. It can be much smaller.
	jobs := make(chan jobPayload, 100)
	addrs := make(chan string, 100)
	results := make(chan []string, 100)
	results <- []string{"endTime", "prefix", "addr", "countSuccess", "totalTimeout", "consecutiveTimeout", "code", "whitelisted", "payloadLength"}

	lines := readfiles.ReadFiles(flag.Args())
	go func() {
		for line := range lines {
			prefix, err := hex.DecodeString(line)
			if err != nil {
				log.Panic(err)
			}
			if *increment != 0 {
				// Generate incremental payloads
				for i := *increment; i <= *payloadLen; i += *increment {
					randomPayloadLen := i - len(prefix)
					if randomPayloadLen < 0 {
						randomPayloadLen = 0
					}
					randomPayload := make([]byte, randomPayloadLen)
					rand.Read(randomPayload)
					jobs <- jobPayload{line, prefix, randomPayload}
				}
			} else if *suffix != "" {
				suffix, err := hex.DecodeString(*suffix)
				if err != nil {
					log.Panic(err)
				}
				jobs <- jobPayload{line, prefix, suffix}
			} else {
				randomPayloadLen := *payloadLen - len(prefix)
				if randomPayloadLen < 0 {
					randomPayloadLen = 0
				}
				randomPayload := make([]byte, randomPayloadLen)
				rand.Read(randomPayload)
				jobs <- jobPayload{line, prefix, randomPayload}
			}
		}
		close(jobs)
	}()

	go func() {
		for _, port := range ports {
			for _, ip := range ips {
				addrs <- net.JoinHostPort(ip.String(), strconv.Itoa(port))
			}
		}
	}()

	dialer := &net.Dialer{
		Timeout: *timeout,
	}

	var wg sync.WaitGroup
	wg.Add(maxNumWorkers)
	for i := 0; i < maxNumWorkers; i++ {
		go func(id int) {
			defer wg.Done()
			worker(id, jobs, addrs, results, dialer)
		}(i)
	}
	go func() {
		wg.Wait()
		close(results)
	}()

	for r := range results {
		err := w.Write(r)
		if err != nil {
			log.Panicln("error writing results to file:", err, len(r))
		}
		w.Flush()
	}
}
