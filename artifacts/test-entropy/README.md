## Introduction

This directory contains source code and data to test if the GFW exempts traffic based on mono-bit randomness testing.

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

### Historical experiments

#### Feb 13, 2023

Experiment: (ran on tencent-random-bj3, Feb 13, 2023)

```sh
python3 test-entropy.py | ./affected-payload -p 10000-20000 -repeat 25 -try 5 -worker 30 > test-entropy.out
```

This experiment only tests the strings `0000000...`, `01010101...`, `02020202...`, ..., `fefefefe...`, `ffffffff...`. We perform a more comprehensive "pop-count" experiment, which we include in Section 4.1.

#### Feb 16, 2023

```sh
python3 pops.py | ./affected-payload -p 10000-20000 -repeat 25 -try 5 -worker 30 > test-pops.out
```

```sh
python3 analyze-pops.py < ./test-pops.out | sort -k 1 -n > pop-results.out
```

pop-results.out shows lines that look like:

```txt
0 0.0000 25 0 0 false 0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000
```

This shows the number of bits set (`0`), fraction of bits set per byte (`0.0000`), successfull connections (`25`), failed/failed consecutive connections (`0 0`) and if we would label this connection as affected based on the failed connections (`false`). In this case, we conclude that the payload (`000...000`) is not affected.
We see that this changes to true above `3.40` bits set per byte, and then goes back to false at `4.60` bits set/byte.

When we ran this, we saw one exception (this varies by run):

```txt
213 4.2600 25 0 0 false 756ffb1a4c9f6cef97266e3becaa57f17f625f155f505060754d65749711374d387fa160b49755719d7b5a9e101be2e36301
```

In this case, the payload would be exempted by another GFW rule: more than half of the payload is printable ASCII characters (`[0x20-0x7e]`). Since payloads are generated randomly, each run may produce a few payloads that would be similarly exempt.
