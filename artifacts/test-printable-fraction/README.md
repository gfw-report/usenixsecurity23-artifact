## Introduction

This directory contains source code and data to test if the GFW will exempts a connection if more than 50% of all bytes in its first data packets are ASCII characters.

To generate payloads for testing and ask `detect.py` to predict if each probe will get blocked or not:

```sh
make
```

Run the experiment:

```sh
make test
```

Compare the experiment results with what Algorithm 1 predicts:

```sh
make compare
```

* If there is no output to stdout and stderr at all, it means all experiment results match what `detect.py` (Algorithm 1) predicts.
* If there is any difference, it does not necessarily mean the detection rule does not hold. It may suggest a testing failure due to unreliable network. Simply re-run `make test` and then `make compare` again, to see if the same payload still have inconsistent results.

## Experiment details

This experiment tests sending 100 strings, ranging from 0 printable bytes to 100 printable bytes.

Strings are generated randomly (`test-printable-frac.py`):
for n=0 to 100, choose n printable bytes (0x20-0x7e)
and 100-n non-printable bytes (`0x00-0x1f` or `0x7f-0xff`)
Shuffle those bytes, and send each payload (so the printable bytes are randomly distributed)

```sh
python3 test-printable-fraction/test-printable-frac.py | ./utils/affected-payload -host REDACTED_US_SERVER_IP -p 4000-5000 -worker 30 > ./test-printable-fraction/fraction.out
```

# Analyze results

```sh
python3 ./test-printable-fraction/analyze.py < ./test-printable-fraction/fraction.out | sort -k 1 -n > ./test-printable-fraction/result.out
```
