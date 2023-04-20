#!/usr/bin/env python3

import sys
import getopt
import glob

def usage(f=sys.stderr):
    program = sys.argv[0]
    f.write(f"""\
Usage: {program}
This script reads writes the protocol fingerprints that can exempt the blocking by the GFW. By default, print results to stdout and log to stderr.

  -h, --help            show this help
  -o, --out             write to file

Example:
  {program}
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

def protocol_pattern():
    for first in (0x15, 0x16, 0x17):
        for third in range(0x00, 0x10):
            yield f"{first:02x}03{third:02x}"

    for verb in ("GET", "PUT", "POST", "HEAD"):
        for next_byte in range(0x00, 0xff + 1):
            yield f"{verb.encode('utf-8').hex()}{next_byte:02x}"


if __name__ == '__main__':
    try:
        opts, args = getopt.gnu_getopt(sys.argv[1:], "ho:", ["help", "out="])
    except getopt.GetoptError as err:
        eprint(err)
        usage()
        sys.exit(2)

    output_file = sys.stdout
    for o, a in opts:
        if o == "-h" or o == "--help":
            usage()
            sys.exit(0)
        if o == "-o" or o == "--out":
            output_file = open(a, 'a+')

    random_payload = "dadd034913c52da75fd9f05dc76803917134808efed97ef8884f2151b712f60fed634f609f132033a15b77ed3ccaa2d20f5b"
    for line in protocol_pattern():
        print(line + random_payload, file=output_file)
    output_file.close()
