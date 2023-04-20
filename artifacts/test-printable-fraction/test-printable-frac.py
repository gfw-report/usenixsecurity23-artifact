
import random
import os

N = 100 # bytes to send
ALL       = set(range(256))          # All characters
PRINTABLE = set(range(0x20,0x7e+1))  # Printable ASCII
NONPRINT  = ALL - PRINTABLE          # Everything but printable

# Get n random bytes in the allowed set
# (e.g. set(range(256)) for all bytes,
# or set(range(0x20,0x7e+1)) for printable ASCII, etc)
def getRand(n, allowed=set()):
    buf = b''
    while len(buf) < n:
        b = os.urandom(1)
        if b[0] in allowed:
            buf += b
    return buf


# Create N strings, with 0 to N bytes of printable characters in each, and print them
for i in range(N):
    s = getRand(i, PRINTABLE) + getRand(N-i, NONPRINT)
    l = list(s)
    random.shuffle(l)
    print(bytes(l).hex())
