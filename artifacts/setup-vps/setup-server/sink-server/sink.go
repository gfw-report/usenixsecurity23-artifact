package main

import (
	"crypto/tls"
	"encoding/csv"
	"encoding/hex"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strconv"
	"sync"
	"time"

	"common/parseipportargs"
	"common/readfiles"
)

func usage() {
	fmt.Fprintf(os.Stderr, `Usage:
    %[1]s [OPTION]...

Description:
    Sink server accepts TCP (or TLS with -tls option) handshakes and data packets. It never sends any data back to server. When timeout, the sink server closes the connection. By default, print CSV results to stdout and log to stderr.

Examples:
    Listening on port 3000, 4000, 4001, and 4002 of 127.0.0.1 as a TCP server. Close each TCP connection after 5 seconds (if the client has not closed it yet):
	%[1]s -ip 127.0.0.1 -p 3000,4000-4002 -timeout 5s

    Listening on port 4000 of 0.0.0.0 as a TLS server. Write CSV output, including the header, to ouptput.csv. Log to log.txt.
	%[1]s -tls -tlsCert ./server.crt -tlsKey ./server.key -p 4000 -header=true -out output.csv -log log.txt

Notes:
    The CSV field "len" records the total number of bytes received in a connection.
    The CSV field is called "truncatedPayload" because it only keeps the first -buffSize bytes of payload for each connection, which may not be complete. When len <= -buffSize, the payload is complete without being truncated.

	To generate TLS certificate and key files:
		openssl genrsa 2048 > server.key && chmod 400 server.key
		openssl req -new -x509 -nodes -sha256 -days 365 -key server.key -out server.crt

    # redirect all tcp packets from testIP to server's port 0-65535 to a single port 12345
    sudo iptables -t nat -A PREROUTING -i eth0 -p tcp -s "$testIP" --dport 0:65535 -j REDIRECT --to-port 12345
    # drop packets with RST flag set, regardless of other TCP flags
    iptables -I OUTPUT -p tcp -d "$testIP" --tcp-flags RST RST -j DROP


Options:
`, os.Args[0])
	flag.PrintDefaults()
}

func handleConnection(c net.Conn, results chan<- []string) {
	defer c.Close()

	startTime := time.Now()

	var err error
	localIP, localPort, err := net.SplitHostPort(c.LocalAddr().String())
	if err != nil {
		fmt.Println(err)
	}
	remoteIP, remotePort, err := net.SplitHostPort(c.RemoteAddr().String())
	if err != nil {
		fmt.Println(err)
	}

	log.Printf("Open:  %v <-- %v\n", c.LocalAddr(), c.RemoteAddr())

	totalLength := 0
	buff := make([]byte, buffSize)
	var length int
	var truncatedPayload string
	for {
		length, err = c.Read(buff)
		// first round
		if totalLength == 0 {
			truncatedPayload = hex.EncodeToString(buff[:length])
		}

		totalLength += length
		log.Printf("Recv:  %v <-- %v: %q (%x), %v bytes\n", c.LocalAddr(), c.RemoteAddr(), string(buff[:length]), buff[:length], length)
		if err != nil {
			break
		}
	}
	endTime := time.Now()
	duration := endTime.Sub(startTime).Seconds()
	if err != nil {
		if err == io.ErrUnexpectedEOF || err == io.EOF {
			log.Printf("Close: %v <-- %v, by client after %v seconds. Received %v bytes in total\n", c.LocalAddr(), c.RemoteAddr(), duration, totalLength)
		} else {
			log.Printf("Close: %v --> %v, by server after %v seconds, due to %v. Received %v bytes in total\n", c.LocalAddr(), c.RemoteAddr(), duration, err.Error(), totalLength)
		}
	} else {
		log.Printf("Close: %v --> %v, by server after %v seconds. Received %v bytes in total\n", c.LocalAddr(), c.RemoteAddr(), duration, totalLength)
	}

	results <- []string{strconv.FormatInt(startTime.UnixMilli(), 10), localIP, localPort, remoteIP, remotePort, truncatedPayload, strconv.Itoa(totalLength), fmt.Sprintf("%.3f", duration)}

	return
}

// global variables
var buffSize int

func main() {
	flag.Usage = usage
	argIP := flag.String("ip", "", "IP address to listen on. (default listen on 0.0.0.0 and ::/0)")
	argPort := flag.String("p", "12345", "comma-separated list of ports to listen on. eg. 3000,4000-4002")
	timeout := flag.Duration("timeout", 60*time.Second, "timeout value.")
	flag.IntVar(&buffSize, "buffSize", 2048, "Set recv buffer size. This is also the max number of bytes to keep in CSV's truncatedPayload field.")
	outputFile := flag.String("out", "", "output csv file.  (default stdout)")
	logFile := flag.String("log", "", "log to file.  (default stderr)")
	flush := flag.Bool("flush", true, "flush after every output.")
	header := flag.Bool("header", false, "print CSV header.")
	useTLS := flag.Bool("tls", false, "listen with TLS.")
	ssLCert := flag.String("tlsCert", "server.crt", "specify TLS certificate file (PEM) for listening.")
	tlsKey := flag.String("tlsKey", "server.key", "specify TLS private key (PEM) for listening.")
	flag.Parse()

	if *logFile != "" {
		f, err := os.Create(*logFile)
		if err != nil {
			log.Panicln("Failed to open log file", err)
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

	ports, err := parseipportargs.ParsePortArgs(*argPort)
	if err != nil {
		log.Panic(err)
	}

	// The channel capacity does not have to be equal to the
	// number of workers. It can be much smaller.
	jobs := make(chan string, 100)
	results := make(chan []string, 100)
	// CSV header
	if *header {
		results <- []string{"ts", "localIP", "localPort", "remoteIP", "remotePort", "truncatedPayload", "len", "duration"}
	}

	lines := readfiles.ReadFiles(flag.Args())

	go func() {
		for line := range lines {
			jobs <- line
		}
		close(jobs)
	}()

	addrs := make([]string, 0)
	for _, port := range ports {
		addr := net.JoinHostPort(*argIP, strconv.Itoa(port))
		addrs = append(addrs, addr)
	}

	var wg sync.WaitGroup
	for _, addr := range addrs {
		var l net.Listener
		var err error

		if !*useTLS {
			log.Printf("TCP server is listening on %s, timeout value: %v\n", addr, *timeout)
			l, err = net.Listen("tcp", addr)
		} else {
			log.Printf("TLS server is listening on %s, timeout value: %v, cert: %v, key: %v\n", addr, *timeout, *ssLCert, *tlsKey)
			cer, err := tls.LoadX509KeyPair(*ssLCert, *tlsKey)
			if err != nil {
				log.Panicln(err)
			}
			config := &tls.Config{Certificates: []tls.Certificate{cer}}

			l, err = tls.Listen("tcp", addr, config)
		}

		if err != nil {
			fmt.Println(err)
			return
		}
		defer l.Close()

		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				c, err := l.Accept()
				if err != nil {
					fmt.Println(err)
					return
				}
				c.SetDeadline(time.Now().Add(*timeout))
				go handleConnection(c, results)
			}
		}()

	}
	go func() {
		wg.Wait()
		close(results)
	}()

	for r := range results {
		if err := w.Write(r); err != nil {
			log.Panicln("error writing results to file", err)
		}
		if *flush {
			w.Flush()
		}
	}
	w.Flush()
}
