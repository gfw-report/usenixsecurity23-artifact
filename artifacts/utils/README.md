## Introduction

This directory contains client-side tools to measure and analyze the dynamic blocking by the GFW.

## Compilation and installation

To build utils:

```
make
```

To install them to the remote server in China:

```
make install
```

## detect.py

```sh
./detect.py -h
```
```txt
Usage: ./detect.py [FILENAME...]
This program simulates the GFW's detection rules specified in Algorithm 1 in the paper. It reads payloads in hex string format from files, and writes its analysis to a CSV file. With no FILE, or when FILE is -, read standard input. By default, print results to stdout and log to stderr.

  -h, --help            show this help
  -o, --out=FILE        write results to FILE
  -d, --header          print a CSV header line (default: False)

Examples:
    ./detect.py --header <<<48656c6c6f2c20776f726c6421
```

## affected-norand

```sh
./affected-norand -h
```
```txt
Usage:
    ./affected-norand [OPTION]... [FILE]...

Description:
    Test if IP addresses in FILE(s) are affected by the dynamic blocking. With no FILE, read standard input. By default, print results to stdout and log to stderr.

    * If an ip:port accepts -repeat number of random connections (which doesn't have to be consecutive), then mark it as unaffected.
    * If -try number of consecutive connections to an ip:port all timeout, then mark it as affected.
    * If the total number of successful connections is zero, mark it as unknown (possibly closed/filtered port).
    * If any other error occured, mark it as unknown.

Wait for -timeout second between each connection to ip:port. If timeout already occured, slow down by waiting for -wait seconds between each connection.

Examples:
    Test if 1.1.1.1 is affected by sending random traffic to its port 80
	echo "1.1.1.1" | ./affected-norand -p 443
Options:
  -cpuprofile string
    	write cpu profile to file.
  -interval duration
    	time interval between each connection to a ip:port. (default 1s)
  -log string
    	log to file.  (default stderr)
  -out string
    	output csv file.  (default stdout)
  -p string
    	comma-separated list of ports to which the program sends random payload. eg. 3000,4000-4002 (default "80")
  -payload string
    	payload of the probes in hex format (default "dadd034913c52da75fd9f05dc76803917134808efed97ef8884f2151b712f60fed634f609f132033a15b77ed3ccaa2d20f5b")
  -repeat int
    	repeatedly make up to this number of connections to each ip:port. (default 25)
  -sleep duration
    	time interval between sending a probe and closing the connection. This value doesn't affect the -interval between each connection.
  -timeout duration
    	timeout value of TCP connections. (default 6s)
  -try int
    	mark an ip:port as affected if this number of consecutive connections all timeout. (default 5)
  -wait duration
    	time interval between each connection, when a ip:port timeout. (default 3s)
  -worker int
    	number of workers in parallel. (default 10)
```

## affected-payload

```sh
./affected-payload -h
```
```txt
Usage:
    ./affected-payload [OPTION]... [FILE]...

Description:
    Test if payloads in FILE(s) are affected by the dynamic blocking. With no FILE, read standard input. By default, print results to stdout and log to stderr.

    * If an ip:port accepts -repeat number of random connections (which doesn't have to be consecutive), then mark it as unaffected.
    * If -try number of consecutive connections to an ip:port all timeout, then mark it as affected.
    * If the total number of successful connections is zero, mark it as unknown (possibly closed/filtered port).
    * If any other error occured, mark it as unknown.

Wait for -timeout second between each connection to ip:port. If timeout already occured, slow down by waiting for -wait seconds between each connection.

Examples:
    Test if payload 00112233 is affected by sending random traffic to the sink server at port 443
	echo "00112233" | ./affected-payload -host 1.1.1.1 -p 443
Options:
  -cpuprofile string
    	write cpu profile to file.
  -host string
    	host to send to (default "REDACTED_US_SERVER_IP")
  -interval duration
    	time interval between each connection to a ip:port. (default 1s)
  -log string
    	log to file.  (default stderr)
  -out string
    	output csv file.  (default stdout)
  -p string
    	comma-separated list of ports to which the program sends random payload. eg. 3000,4000-4002 (default "80")
  -repeat int
    	repeatedly make up to this number of connections to each ip:port. (default 25)
  -sleep duration
    	time interval between sending a probe and closing the connection. This value doesn't affect the -interval between each connection.
  -timeout duration
    	timeout value of TCP connections. (default 6s)
  -try int
    	mark an ip:port as affected if this number of consecutive connections all timeout. (default 5)
  -wait duration
    	time interval between each connection, when a ip:port timeout. (default 3s)
  -worker int
    	number of workers in parallel. (default 10)
```


## whitelist

```sh
./whitelist -h
```
```txt
Usage:
    ./whitelist [OPTION]... [FILE]...

Description:
    Test if prepending the payloads (represented in hex string) in FILE(s) can exempt the GFW's blocking of random traffic. With no FILE, read standard input. By default, print results to stdout and log to stderr.

Examples:
    Test if prepending "QUIT\r\n" (0x51, 0x55, 0x49, 0x54, 0x0d, 0x0a) can exempt the GFW's blocking of random traffic.
	echo -n "515549540d0a" | ./whitelist
    Test if prepending "aaaaaa" (0x61, 0x61, 0x61, 0x61, 0x61, 0x61) can exempt the GFW's blocking of random traffic.
	echo -n "616161616161" | ./whitelist
Options:
  -append string
    	appending a fixed random payload for each probe, in hex format
  -cpuprofile string
    	write cpu profile to file.
  -dip string
    	comma-separated list of destination IP addresses to which the program can send probes. eg. 1.1.1.1,2.2.2.2 (default "127.0.0.1")
  -increment int
    	increment for generating variable payload sizes. must be used in conjunction with increment. eg. -max -increment 10 will try the same prefix with total payload len 10, 20, ..., 100
  -interval duration
    	time interval between each connection to a ip:port. (default 1s)
  -log string
    	log to file.  (default stderr)
  -max int
    	length of payload to send. (default 50)
  -out string
    	output csv file.  (default stdout)
  -p string
    	comma-separated list of available ports to which the program can send probes. eg. 3000,4000-4002 (default "10000-10200")
  -redisual duration
    	duration of residual censorship when a ip:port blocked. (default 3m0s)
  -repeat int
    	repeated make up to this number of connections to each ip:port. (default 25)
  -timeout duration
    	timeout value of TCP connections. (default 6s)
  -try int
    	mark an ip:port as affected if this number of consecutive connections all timeout. (default 5)
  -wait duration
    	time interval between each connection, when a ip:port timeout. (default 3s)
  -worker int
    	number of workers in parallel. (default 10)
```

## compare.py

```sh
./compare.py -h
```
```txt
Usage: ./compare.py --actual [FILENAME] --predict [FILENAME]
This script reads and compares two files: one file that predicts the experiment results; the other one contains actual test results. The program will only print out unmatched results. By default, print results to stdout and log to stderr.

  -h, --help            show this help
  -o, --out             write to file
  -a, --actual          path to the testing reuslt file
  -p, --predict         path to the prediction file
  -d, --header          print header (default: false)

Example:
  ./compare.py --actual results.csv --predict prediction.csv
```
