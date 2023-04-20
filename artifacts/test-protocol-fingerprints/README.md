## Introduction

This directory contains source code and data to check if certain protocol fingerprints can exempt the blocking.

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
