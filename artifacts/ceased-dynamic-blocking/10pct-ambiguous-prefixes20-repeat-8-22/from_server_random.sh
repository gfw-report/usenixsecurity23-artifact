#!/bin/bash

set -x
set -e

cd "$(dirname "$0")" || exit

# sleep 5

rsync eric6:~/aa.out -auvP . &
rsync eric4:~/ab.out  -auvP . &
rsync eric5:~/ac.out  -auvP . &
rsync random5:~/ad.out  -auvP . &
rsync random6:~/ae.out  -auvP . &
rsync random7:~/af.out  -auvP . &
rsync random8:~/ag.out  -auvP . &
rsync random9:~/ah.out  -auvP . &
rsync random10:~/ai.out -auvP . &

wait

(head -n1 aa.out; cat a*.out | grep -v "endTime") > 10pct-ambiguous-prefixes20-random-9machines-t3-r25-w400-s1s.aggregated.out

rsync eric6:~/aa.pcap -auvP . &
rsync eric4:~/ab.pcap -auvP . &
rsync eric5:~/ac.pcap -auvP . &
rsync random5:~/ad.pcap -auvP . &
rsync random6:~/ae.pcap -auvP . &
rsync random7:~/af.pcap -auvP . &
rsync random8:~/ag.pcap -auvP . &
rsync random9:~/ah.pcap -auvP . &
rsync random10:~/ai.pcap -auvP . &

wait
