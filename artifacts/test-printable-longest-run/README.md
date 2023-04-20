## Introduction

This directory contains source code and data to check if inserting more than 20 bytes of contiguous ASCII character to a random payload can exempt the blocking.

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

## Historical experiment

On tencent-random-bj3, Feb 14, 2023:

`printable.py` will generate longest run printable of length 0 to 90.

```sh
./printable.py | tail -n+2 | ../detection-algorithm/detect.py
```

```sh
python3 printable.py | ./affected-payload -p 1000-2000 -repeat 25 -try 5 -worker 30 > test-printable.out
```

```sh
python3 analyze.py  < ./test-printable.out  | sort -k 1 -n > printable-results.out
```
