import sys

def popcount(b):
    t = 0
    for c in b:
        t += sum([int(x) for x in bin(c)[2:]])
    return t

# 1676575996,REDACTED_US_SERVER_IP:10054,2800000000a50400020801320000000800044000008480004361010040204093014141000288820050030208200004108800,25,0,0,,false
for line in sys.stdin:
    if 'endTime' in line:
        continue
    ts,host,payload,n_succ,n_fail,n_consec,err,affected = line.split(',')

    pb = bytes.fromhex(payload)
    bits = popcount(pb)
    print(bits, '%.4f' % (bits/len(pb)), n_succ, n_fail, n_consec, affected.strip(), payload)

