#!/bin/bash

# grep true 10pct-80-9machines-t5-r25-w500-s1s.aggregated.out > 2022-08-22-10pct-80-9machines-t5-r25-w500-s1s.aggregated.affected.out

# scan with 50-byte random probes
cut -d, -f2 2022-08-22-10pct-80-9machines-t5-r25-w500-s1s.aggregated.affected.out | cut -d\: -f1 | ../utils/affected-norand -p 80 -worker 500 | tee retesting-on-2023-03-15-random-all-affected-ip-as-2022-08-22-10pct-80-1machine-t5-r25-w500-s1s.out

# scan with 50-byte zero probes
cut -d, -f2 2022-08-22-10pct-80-9machines-t5-r25-w500-s1s.aggregated.affected.out | cut -d\: -f1 | ../utils/affected-norand -p 80 -worker 500 -payload 0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000 | tee retesting-on-2023-03-15-zero-all-affected-ip-as-2022-08-22-10pct-80-1machine-t5-r25-w500-s1s.out
