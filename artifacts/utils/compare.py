#!/usr/bin/env python3

import sys
import getopt
import glob

def usage(f=sys.stderr):
    program = sys.argv[0]
    f.write(f"""\
Usage: {program} --actual [FILENAME] --predict [FILENAME]
This script reads and compares two files: one file that predicts the experiment results; the other one contains actual test results. The program will only print out unmatched results. By default, print results to stdout and log to stderr.

  -h, --help            show this help
  -o, --out             write to file
  -a, --actual          path to the testing reuslt file
  -p, --predict         path to the prediction file
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

# This function compares the prediction and actual results.
# The actual result file contains the following columns:
# endTime,addr,payload,countSuccess,totalTimeout,consecutiveTimeout,code,affected
# The prediction file contains the following columns:
# payload;num_beginnig_ascii;longest_run;popcount_fraction;ascii_fraction;match_fingerprint;exempted
# The function select the (payload, exempted) pair from the prediction file 
# and compare it with the (payloaod, affected) pair from the actual result file.
# If the two pairs match, the function does nothing. Otherwise, it yields the mismatched pair.
def compare(actual_file, predict_file):
    # skip headers
    actual_file.readline()
    predict_file.readline()

    actual = {}
    predict = {}
    for line in actual_file:
        line = line.strip()
        if not line:
            continue
        endTime,addr,payload,countSuccess,totalTimeout,consecutiveTimeout,code,affected = line.split(',')
        if affected == "true":
            actual[payload] = "true"
        elif affected == "false":
            actual[payload] = "false"
        elif affected == "unknown":
            eprint(f"Unknow testing results: {payload}")
            actual[payload] = "unknown"
        else:
            raise Exception(f"Unexpected affected value: {affected}")
    for line in predict_file:
        line = line.strip()
        if not line:
            continue
        payload,num_beginnig_ascii,longest_run,popcount_fraction,ascii_fraction,match_fingerprint,exempted = line.split(';')
        if exempted == "True":
            predict[payload] = "false"
        elif exempted == "False":
            predict[payload] = "true"
        else:
            raise Exception(f"Unexpected exempted value: {exempted}")
    
    for payload in actual:
        if payload not in predict:
            eprint(f"Payload {payload} not found in the prediction file.")
            yield payload, actual[payload], "missing"
            continue
        if actual[payload] != predict[payload]:
            eprint(f"Payload {payload} mismatch: actual {actual[payload]}, predict {predict[payload]}")
            yield payload, actual[payload], predict[payload]

    




if __name__ == '__main__':
    try:
        opts, args = getopt.gnu_getopt(sys.argv[1:], "ho:a:p:d", ["help", "out=", "actual=", "predict=", "header"])
    except getopt.GetoptError as err:
        eprint(err)
        usage()
        sys.exit(2)

    output_file = sys.stdout
    actual_file = None
    predict_file = None
    header = False
    for o, a in opts:
        if o == "-h" or o == "--help":
            usage()
            sys.exit(0)
        if o == "-o" or o == "--out":
            output_file = open(a, 'a+')
        if o == "-a" or o == "--actual":
            actual_file = open(a, 'r')
        if o == "-p" or o == "--predict":
            predict_file = open(a, 'r')
        if o == "-d" or o == "--header":
            header = True

    if not actual_file:
        eprint("Missing --actual.")
        usage()
        sys.exit(-1)
    if not predict_file:
        eprint("Missing --predict.")
        usage()
        sys.exit(-1)

    if header:
        print("payload;blocked_in_actual;blocked_in_predict", file=output_file)

    for payload, actual_status, predict_status in compare(actual_file, predict_file):
        print(f"{payload};{actual_status};{predict_status}", file=output_file)
    output_file.close()
