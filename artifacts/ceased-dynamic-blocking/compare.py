#!/usr/bin/env python3

import sys
import getopt
import glob

def usage(f=sys.stderr):
    program = sys.argv[0]
    f.write(f"""\
Usage: {program} --actual [FILENAME] --predict [FILENAME]
This script reads and compares two files: one file that contains the results from random probes experiment; the other one contains resutls from zero probes experiment. The program will only print out unmatched results. By default, print results to stdout and log to stderr.

  -h, --help            show this help
  -o, --out             write to file
  -r, --random          path to the result file of random probes experiment
  -z, --zero            path to the result file of zero probes experiment
  -d, --header          print header (default: false)

Example:
  {program} --actual results.csv --predict prediction.csv
""")

def eprint(*args, **kwargs):
    print(*args, file=sys.stderr, **kwargs)

def input_files(args):
    if not args:
        yield sys.stdin.buffer
    else:
        for arg in args:
            if arg == "-":
                yield sys.stdin.buffer
            else:
                for path in glob.glob(arg):
                    with open(path) as f:
                        yield f


def compare(random_file, zero_file):
    # skip headers
    random_file.readline()
    zero_file.readline()

    random = {}
    zero = {}
    for line in random_file:
        line = line.strip()
        if not line:
            continue
        endTime,addr,countSuccess,totalTimeout,consecutiveTimeout,code,affected = line.split(',')
        if affected == "true":
            random[addr] = "true"
        elif affected == "false":
            random[addr] = "false"
        elif affected == "unknown":
            eprint(f"Unknow testing results: {addr}")
            random[addr] = "unknown"
        else:
            raise Exception(f"Unexpected affected value: {affected}")
    for line in zero_file:
        line = line.strip()
        if not line:
            continue
        endTime,addr,countSuccess,totalTimeout,consecutiveTimeout,code,affected = line.split(',')
        if affected == "true":
            zero[addr] = "true"
        elif affected == "false":
            zero[addr] = "false"
        elif affected == "unknown":
            eprint(f"Unknow testing results: {addr}")
            zero[addr] = "unknown"
        else:
            raise Exception(f"Unexpected affected value: {affected}")
    
    for addr in random:
        if addr not in zero:
            eprint(f"Addr {addr} not found in the prediction file.")
            yield addr, random[addr], "missing"
            continue
        if random[addr] != zero[addr]:
            eprint(f"Addr {addr} mismatch: actual {random[addr]}, predict {zero[addr]}")
            yield addr, random[addr], zero[addr]

if __name__ == '__main__':
    try:
        opts, args = getopt.gnu_getopt(sys.argv[1:], "ho:r:z:d", ["help", "out=", "random=", "zero=", "header"])
    except getopt.GetoptError as err:
        eprint(err)
        usage()
        sys.exit(2)

    output_file = sys.stdout
    random_file = None
    zero_file = None
    header = False
    for o, a in opts:
        if o == "-h" or o == "--help":
            usage()
            sys.exit(0)
        if o == "-o" or o == "--out":
            output_file = open(a, 'a+')
        if o == "-r" or o == "--random":
            random_file = open(a, 'r')
        if o == "-z" or o == "--zero":
            zero_file = open(a, 'r')
        if o == "-d" or o == "--header":
            header = True

    if not random_file:
        eprint("Missing --random.")
        usage()
        sys.exit(-1)
    if not zero_file:
        eprint("Missing --zero.")
        usage()
        sys.exit(-1)

    if header:
        print("addr;blocked_in_actual;blocked_in_predict", file=output_file)

    for addr, random_status, zero_status in compare(random_file, zero_file):
        print(f"{addr};{random_status};{zero_status}", file=output_file)
    output_file.close()
