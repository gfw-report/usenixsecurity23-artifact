#!/bin/bash

set -x
set -e

cd "$(dirname "$0")" || exit

# sleep 5

rsync eric6:~/zero-aa.out -auvP . &
rsync eric4:~/zero-ab.out  -auvP . &
rsync eric5:~/zero-ac.out  -auvP . &
rsync random5:~/zero-ad.out  -auvP . &
rsync random6:~/zero-ae.out  -auvP . &
rsync random7:~/zero-af.out  -auvP . &
rsync random8:~/zero-ag.out  -auvP . &
rsync random9:~/zero-ah.out  -auvP . &
rsync random10:~/zero-ai.out -auvP . &

wait

(head -n1 zero-aa.out; cat zero-a*.out | grep -v "endTime") > 10pct-ambiguous-prefixes20-zero-9machines-t5-r25-w400-s1s.aggregated.out
