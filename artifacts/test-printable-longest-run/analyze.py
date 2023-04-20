import sys

# 1676376962,REDACTED_US_SERVER_IP:1023,e8e8e8e8e8e8e8e8e8e84b4b4b4b4b4b4b4b4b4b4b4b4b4b4b4b4b4b4b4b4b4b4be8e8e8e8e8e8e8e8e8e8e8e8e8e8e8e8e8e8e8e8e8e8e8e8e8e8e8e8e8e8e8e8e8e8e8e8e8e8e8e8e8e8e8e8e8e8e8e8e8e8e8e8e8e8e8e8e8e8e8e8e8e8e8e8e8e8e8,25,0,0,,false
for line in sys.stdin:
    if 'endTime' in line:
        continue
    ts,host,payload,n_succ,n_fail,n_consec,err,affected = line.split(',')

    try:
        n = (payload.index('e8',20) - 20) / 2
    except:
        n = 90
    print('%d %s,%s,%s,%s' % (n, n_succ, n_fail, n_consec, affected.strip()))
