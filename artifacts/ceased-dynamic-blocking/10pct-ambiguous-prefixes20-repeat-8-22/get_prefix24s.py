import pyasn
import sys
import asname
from collections import defaultdict
import ipaddress
db = pyasn.pyasn('/home/ewust/ipasn.20220130.0600.dat')


SET_PREFIX=20
MASK=0xfffff000

prefixes_aff = defaultdict(int)
prefixes_tot = defaultdict(int)
pasn = defaultdict(int)
asn_aff = defaultdict(int)
asn_tot = defaultdict(int)
for line in sys.stdin.readlines():
    line = line.strip()
    if 'addr,' in line:
        continue
    ts,ip_port,success,tot_to,consec_to,err,affected = line.split(',')
    ip,port = ip_port.split(':')
    asn, prefix = db.lookup(ip)
    prefix24 = str(ipaddress.IPv4Address(int(ipaddress.IPv4Address(ip)) & MASK)) + '/%d' % SET_PREFIX

    if prefix is None:
        continue
    if int(prefix.split('/')[1]) < SET_PREFIX:
        prefix = prefix24

    if affected == 'true':
        prefixes_aff[prefix] += 1
        asn_aff[asn] += 1
    prefixes_tot[prefix] += 1
    asn_tot[asn] += 1
    pasn[prefix] = asn

asnp = defaultdict(list)

for prefix,num in sorted(prefixes_aff.items(), key=lambda x: x[1], reverse=True):
    asn = 0
    if prefix in pasn:
        asn = pasn[prefix]
        asnp[asn].append(prefix)
    if asn == None:
        asn = 0
    name = ''
    if asn != 0:
        name = asname.LookupName(asn)

    tot = prefixes_tot[prefix]
    print('%s  AS%d   %.4f   %d   %d   %s' % (prefix, asn, num/tot, num, tot, name))


# output top 25 ASNs to asns/
plots = []
for asn,num in sorted(asn_aff.items(), key=lambda x: x[1], reverse=True)[:15]:
    fname = 'asns/as%d.out' % (asn)
    with open(fname, 'w') as f:
        for prefix in asnp[asn]:
            paff = prefixes_aff[prefix]
            ptot = prefixes_tot[prefix]
            f.write('%s %.4f %d %d\n' % (prefix, paff/ptot, paff, ptot))

    name = ''
    if asn != 0:
        name = asname.LookupName(asn)

    ' '.join(name.split(' ')[:2])

    plots.append("'%s' u 1:2 w lines title '%s'" % (fname+'.cdf', name))



with open('asns/asn.gnuplot', 'w') as f:
    f.write('plot ' + ',\\\n    '.join(plots))

