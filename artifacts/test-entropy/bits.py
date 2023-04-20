
import sys

def countBits(n):
    bits = 0
    while n > 0:
        bits += 1
        n &= n-1
    return bits


# 1676342068,REDACTED_US_SERVER_IP:10240,f0f0f0f0f0f0f0f0f0f0f0f0f0f0f0f0f0f0f0f0f0f0f0f0f0f0f0f0f0f0f0f0f0f0f0f0f0f0f0f0f0f0f0f0f0f0f0f0f0f0,1,5,5,timeout,true
for line in sys.stdin:
    if 'payload' in line:
        continue
    ts, host, payload, n_succ, n_fail, n_cons_fail, err, aff = line.split(',')
    b = int(payload[0:2], 16)
    bits = countBits(b)

    print('%02x %d %s %s %s' % (b, bits, n_succ, n_fail, aff.strip()))

