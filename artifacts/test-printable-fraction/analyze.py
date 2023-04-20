import sys

def numPrintable(d):
    PRINTABLE = set(range(0x20, 0x7e+1))   # Printable ASCII
    n = 0
    for b in d:
        if b in PRINTABLE:
            n += 1
    return n

# 1678257004,REDACTED_US_SERVER_IP:4034,1ffe0adcfc8425688d937827ead8bb5194d7a59665b9617a3267adac44efa5851992ef7ae636a39bf25409dfa5bc9ccf23b6fab9b393792bc07d1d30c4659c9f1ece7136925d5691d94ecec0bf7cbde3034c9cdb4324ae457ecebe3683e7bef24feccc73,10,5,5,timeout,true
for line in sys.stdin:
    if 'endTime' in line:
        continue
    ts,host,payload,n_succ,n_fail,n_consec,err,affected = line.split(',')
    p = bytes.fromhex(payload)
    n = numPrintable(p)

    print('%d %.3f %s,%s,%s %s %s' % (n, float(n)/len(p), n_succ, n_fail, n_consec, affected.strip(), payload))

