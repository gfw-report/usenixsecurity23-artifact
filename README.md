# USENIX Security 2023 AE

This repo includes the source code, data, and documentation to reproduce the major claims in the USENIX Security 2023 paper [*How China Detects and Blocks Fully Encrypted Traffic*](https://www.usenix.org/system/files/sec23fall-prepub-234-wu-mingshi.pdf). 

It is designed for anyone who is curious about the methodologies in our study and wants to reproduce the claims in our paper independently.

Note that it is possible that you cannot reproduce any of the experiment results because the GFW has stopped dynamic blocking since March, 2023. See [this documentation](./artifacts/ceased-dynamic-blocking) for more information.


## Overview of the Repo Structure

```txt
.
├── ae-appendix
├── artifacts
│   ├── ceased-dynamic-blocking
│   ├── common
│   ├── setup-vps
│   ├── sink-server
│   ├── test-entropy
│   ├── test-printable-fraction
│   ├── test-printable-longest-run
│   ├── test-printable-prefixes
│   ├── test-protocol-fingerprints
│   └── utils
├── CHECKLIST
├── LICENSE
└── README.md
```

* `ceased-dynamic-blocking` contains the source code, data, and documentation on the observation that the GFW of China has stopped blocking random traffic dynamically at least since March 15, 2023.
* `ae-appendix` contains the source code and Makefile to generate the [artifact appendix](./ae-appendix/usenix23.pdf).
* `artifacts/setup-vps` contains the source code to set up remote VPSes.
* `artifacts/sink-server` contains the source code for [a sink server](./artifacts/sink-server/), which runs on the server side.
* `artifacts/utils` contains [client-side testing tools](./artifacts/utils/).
* `artifacts/test-*` contain five different tests. Each of them corresponds to a claim of the GFW's traffic exemption rules.
* `artifacts/common-*` is a module that contains code on which measurement tools are built.

## VPS Information and Configuration

To conduct the measurement experiments described in this repo, it requires at least one host in China and one host outside of China.

To assist the [USENIX SECURITY'23 Artifact Evaluation](https://secartifacts.github.io/usenixsec2023/), we provided the reviewers with two VPSes below. 

| SSH Nickname                              | Location                              | ASN     | CPU Model                | # Core(s) | RAM  | OS                                                      |
|---------------------------------------|---------------------------------------|---------|--------------------------|-----------|------|---------------------------------------------------------|
| usenix-ae-client-china      | AlibabaCloud Beijing Datacenter       | AS37963 | Intel Xeon Platinum 8163 | 1         | 1GB  | Ubuntu 22.04.2 LTS (GNU/Linux 5.15.0-56-generic x86_64)   |
| usenix-ae-client-us | DigitalOcean San Francisco Datacenter | AS14061 | Intel DO-Regular         | 1         | 1GB  | Ubuntu 20.04.3 LTS (GNU/Linux 5.4.0-88-generic x86_64) |

If you are not an AE reviewer, but simply want to repeat some of the experiments yourself, you need to purchase and set up the two servers yourself. 

1. We refer you to this [README](./artifacts/setup-vps/README.md) for detailed instructions.

2. To set up the client (VPS in China), execute:

```sh
./artifacts/setup-vps/setup-client/to_alibaba_server.sh
```

3. To set up the server (VPS in the US), execute:

```sh
./artifacts/setup-vps/setup-server/to_digitalocean_server.sh
```

4. Note that we have replaced the IP addresses of the two machines with strings of `REDACTED_CN_SERVER_IP` and `REDACTED_US_SERVER_IP` in our code and documentation. You may want to replace them with your servers' IP addresses (which are `1.1.1.1` and `2.2.2.2` in the below example), using some commands like these:

```sh
find . -type f ! -name "*.pcap" ! -path '*/\.*' -exec sed -i "s#REDACTED_US_SERVER_IP#1.1.1.1#g" {} \;
find . -type f ! -name "*.pcap" ! -path '*/\.*' -exec sed -i "s#REDACTED_CN_SERVER_IP#2.2.2.2#g" {} \;
```

## Minimal Working Example

1. First login to the VPS in China:

```sh
ssh usenix-ae-client-china
```

2. Send some random probes from `usenix-ae-client-china` to the port `2` of `usenix-ae-server-us` by repetitively executing the following command:

```sh
head -c200 /dev/urandom | nc -vn REDACTED_US_SERVER_IP 2
```

3. After executing the command a few times (1 time to 15 times), if you notice that the `nc` cannot connect to `REDACTED_US_SERVER_IP:2` anymore. Congratulations! The blocking is triggered (and will residually last for up to three minutes).
You should still be able to connect to other ports of the same server, for example, `REDACTED_US_SERVER_IP:3`. It is also likely that you cannot trigger the blocking, because the GFW has stopped dynamic blocking since March, 2023. See this documentation for more information: [./artifacts/ceased-dynamic-blocking].

4. (Optional) Alternatively, one can use the triggering tools:

```sh
echo REDACTED_US_SERVER_IP | ./utils/affected-norand -p 2 -log /dev/null
```

This tool will take a list of IPs on stdin, and perform (default 25) repeated connections to
the specified port, sending the
same (configurable) random payload in each connection. If the tool is unable to connect for
(default 5) consecutive connections in a row, the tool labels the IP as `affected` by
blocking (`true` in the `affected` column):

```txt
endTime,addr,countSuccess,totalTimeout,consecutiveTimeout,code,affected
1678258922,REDACTED_US_SERVER_IP:2,2,5,5,timeout,true
```

This output means that connecting to the endpoint (REDACTED_US_SERVER_IP:2)
succeeded in 2 connections, but then
had 5 consecutive connections timeout in a row (and a total of 5 failed).
Because there was at least 5 consecutive timeouts, our tool labels this
endpoint/payload combination as affected (true).

## Estimated Required Time to Reproduce Experiments

We provide a list of estimated required time to reproduce different experiments.

| Experiments                     | Human Time (minutes) | Compute Time (minutes) |
|---------------------------------|------------|--------------|
| Ex0: test-random                |      5      |        5      |
| [Ex1: confirm-ceased-blocking](./artifacts/ceased-dynamic-blocking/)                |      15      |        2 days      |
| [Ex2: test-entropy](./artifacts/test-entropy/)                |      30      |       30       |
| [Ex3: test-printable-prefixes](./artifacts/test-printable-prefixes/)     |    15        |     30         |
| [Ex4: test-printable-fraction](./artifacts/test-printable-fraction/)     |     15       |     30         |
| [Ex5: test-printable-longest-run](./artifacts/test-printable-longest-run/) |     15       |   15           |
| [Ex6: test-protocol-fingerprints](./artifacts/test-protocol-fingerprints/) |      15      |       240       |
