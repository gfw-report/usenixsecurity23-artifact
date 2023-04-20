# Introduction

```sh
go build
./sink -h
```
```txt
Usage:
    ./sink [OPTION]...

Description:
    Sink server accepts TCP (or TLS with -tls option) handshakes and data packets. It never sends any data back to server. When timeout, the sink server closes the connection. By default, print CSV results to stdout and log to stderr.

Examples:
    Listening on port 3000, 4000, 4001, and 4002 of 127.0.0.1 as a TCP server. Close each TCP connection after 5 seconds (if the client has not closed it yet):
	./sink -ip 127.0.0.1 -p 3000,4000-4002 -timeout 5s

    Listening on port 4000 of 0.0.0.0 as a TLS server. Write CSV output, including the header, to ouptput.csv. Log to log.txt.
	./sink -tls -tlsCert ./server.crt -tlsKey ./server.key -p 4000 -header=true -out output.csv -log log.txt

Notes:
    The CSV field "len" records the total number of bytes received in a connection.
    The CSV field is called "truncatedPayload" because it only keeps the first -buffSize bytes of payload for each connection, which may not be complete. When len <= -buffSize, the payload is complete without being truncated.

Options:
  -buffSize int
    	Set recv buffer size. This is also the max number of bytes to keep in CSV's truncatedPayload field. (default 2048)
  -flush
    	flush after every output. (default true)
  -header
    	print CSV header.
  -ip string
    	IP address to listen on. (default listen on 0.0.0.0 and ::/0)
  -log string
    	log to file.  (default stderr)
  -out string
    	output csv file.  (default stdout)
  -p string
    	comma-separated list of ports to listen on. eg. 3000,4000-4002 (default "12345")
  -timeout duration
    	timeout value. (default 1m0s)
  -tls
    	listen with TLS.
  -tlsCert string
    	specify TLS certificate file (PEM) for listening. (default "server.crt")
  -tlsKey string
    	specify TLS private key (PEM) for listening. (default "server.key")
```

## Motivation

**go version is significantly more efficient, compared to tcpserver:**

50% of a single core + 100 MB (300 - 200) is enough to handle 1 Million TCP connections per minute.

While the 100% of a single core + 1000 MB is enough to handle 0.3 million TCP connections per minute

## Firewall rules rationale

We need firewall rules to suppress outgoing RST.

This is because the program cannot guarantee not to send RSTs. As introduced in [this post](https://github.com/net4people/bbs/issues/26), "[i]f a user-space process closes a socket without draining the kernel socket buffer, the kernel sends a RST instead of a FIN".

The best a program can do is to try reading from the buffer as fast and as possible; however, as long as the server actively close a connection, it is sometimes unavoidable that some bytes are still left in the buffer.

This is especially likely to happen when the client sends data to the server at a very high speed. For example, the following example will cause the server to send a RST when timeout:

```sh
./sink -ip 0.0.0.0 -p 12345
```

```sh
nc -nv 127.0.0.1 12345 </dev/random
```

This is because, when the server timeout and calls close(), there is a chance that more data has arrived at the server but hasn't been read.

Even if we let the program read forever without timing out, there are still other cases where the server has to actively close the connection. For example, when a "too many open files" error occurs.
