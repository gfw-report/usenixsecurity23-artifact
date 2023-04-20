#!/usr/bin/env python3

import sys
import getopt
import binascii
import glob


def usage(f=sys.stdout):
    program = sys.argv[0]
    f.write(f"""\
Usage: {program} [FILENAME...]
This program simulates the GFW's detection rules specified in Algorithm 1 in the paper. It reads payloads in hex string format from files, and writes its analysis to a CSV file. With no FILE, or when FILE is -, read standard input. By default, print results to stdout and log to stderr.

  -h, --help            show this help
  -o, --out=FILE        write results to FILE
  -d, --header          print a CSV header line (default: False)

Examples:
    {program} --header <<<48656c6c6f2c20776f726c6421
""")

def eprint(*args, **kwargs):
    print(*args, file=sys.stderr, **kwargs)

def input_files(args):
    if not args:
        yield sys.stdin
    else:
        for arg in args:
            if arg == "-":
                yield sys.stdin
            else:
                for path in glob.glob(arg):
                    with open(path) as f:
                        yield f

# This function returns the number of bits set to 1 in n.
def bitcount(n):
    c = 0
    while n:
        c += 1
        n = (n & (n-1))
    return c

# This function returns the fraction of bits set to 1 in data.
def avg_set_bits(data):
    n = 0
    for b in data:
        n += bitcount(b)
    popcount_fraction = n / len(data)
    return popcount_fraction


ascii_set = set(range(0x20,0x7f))

# This function returns the number of ASCII bytes at the beginning of data.
def num_beginning_ascii(data):
    for i, byte in enumerate(data):
        if byte not in ascii_set:
                return i
    else:
        return len(data)

# This function returns the length of the longest run of ASCII bytes in data.
def longest_run(data):
    counting = False
    run_len = 0
    max_run_len = 0
    for b in data:
        if counting:
            if b in ascii_set:
                run_len += 1
            else:
                max_run_len = max(run_len, max_run_len)
                counting = False
                run_len = 0
        else:
            if b in ascii_set:
                counting = True
                run_len += 1
            else:
                counting = False
                run_len = 0
    # this line is important to catch the cases where all bytes are ascii
    max_run_len = max(run_len, max_run_len)
    return max_run_len

# This function returns the fraction of ASCII bytes in data.
def frac_ascii(data):
    c = 0
    for b in data:
        if b in ascii_set:
            c += 1
    return c / len(data)

# This function matches data with any of the following protocol fingerprints:
# TLS: [\x16-\x17]\x03[\x00-\x09]
# HTTP: "GET ", "PUT ", "POST ", or "HEAD ".
# It returns the matched protocol fingerprint.
# If no fingerprint matched, it returns an empty string.
def protocol_exemption(data):
    if data[0] in (0x16, 0x17) and data[1] == 0x03 and data[2] in range(0x00, 0x10):
        return data[0:3]
    if data[0:4] == b"GET ":
        return data[0:4].decode('utf-8')
    if data[0:4] == b"PUT ":
        return data[0:4].decode('utf-8')
    if data[0:5] == b"POST ":
        return data[0:5].decode('utf-8')
    if data[0:5] == b"HEAD ":
        return data[0:5].decode('utf-8')
    return ""


def process(line):
    eprint("data:", line)

    binary = binascii.unhexlify(line)

    n, l, p, f, pf = num_beginning_ascii(binary), longest_run(binary), avg_set_bits(binary), \
        frac_ascii(binary), protocol_exemption(binary)

    eprint('Num of ascii at beginning: %d' % n)
    eprint('Longest run of ascii: %d' % l)
    eprint('Popcount fraction: %.4f' % p)
    eprint('Ascii fraction: %.4f' % f)

    exempted = False
    if n >= 6 or l > 20 or (p <= 3.4 or p >= 4.6) or f > 0.5 or pf:
        exempted = True

    if n >= 6:
        eprint(f"Exempted by num of ascii at beginning: {n} >=6")
    if l > 20:
        eprint(f"Exempted by longest run of ascii: {l} >20")
    if p <= 3.4:
        eprint(f"Exempted by popcount fraction: {p} <=3.4")
    if p >= 4.6:
        eprint(f"Exempted by popcount fraction: {p} >=4.6")
    if f > 0.5:
        eprint(f"Exempted by ascii fraction: {f} >0.5")
    if pf:
        eprint(f"Exempted by the protocol fingerprint: {pf}")

    return str(n), str(l), str(p), str(f), str(pf), str(exempted)

if __name__ == '__main__':
    try:
        opts, args = getopt.gnu_getopt(sys.argv[1:], "ho:d", ["help", "out=", "header"])
    except getopt.GetoptError as err:
        eprint(err)
        usage()
        sys.exit(2)
    output_file = sys.stdout
    header = False
    for o, a in opts:
        if o == "-h" or o == "--help":
            usage()
            sys.exit(0)
        if o == "-o" or o == "--out":
            output_file = open(a, 'w+')
        if o == "-d" or o == "--header":
            header = True

    if header:
        print(";".join(["payload", "num_beginnig_ascii",  "longest_run", "popcount_fraction", "ascii_fraction", "match_fingerprint", "exempted"]), file=output_file)

    for f in input_files(args):
        for line in f:
            line = line.strip()
            if not line:
                continue
            n, l, p, f, pf, exempted = process(line)
            print(";".join([line, n, l , p, f, pf, exempted]), file=output_file)
    output_file.close()
