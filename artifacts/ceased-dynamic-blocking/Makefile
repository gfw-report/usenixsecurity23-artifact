PYTHON = python3

ALL = test

.PHONY: all
all: $(ALL)

2022-08-22-10pct-80-9machines-t5-r25-w500-s1s.aggregated.affected.out: 10pct-ambiguous-prefixes20-repeat-8-22/10pct-ambiguous-prefixes20-random-9machines-t5-r25-w400-s1s.aggregated.out
	grep true $^ > "$@"

retesting-random-all-affected-ip-as-2022-08-22-10pct-80-1machine-t5-r25-w500-s1s.out: 2022-08-22-10pct-80-9machines-t5-r25-w500-s1s.aggregated.affected.out
	cut -d, -f2 $^ | cut -d\: -f1 | ../utils/affected-norand -p 80 -worker 500 | tee "$@"

retesting-zero-all-affected-ip-as-2022-08-22-10pct-80-1machine-t5-r25-w500-s1s.out: 2022-08-22-10pct-80-9machines-t5-r25-w500-s1s.aggregated.affected.out
	cut -d, -f2 $^ | cut -d\: -f1 | ../utils/affected-norand -p 80 -worker 500 -payload 0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000 | tee "$@"

test: retesting-random-all-affected-ip-as-2022-08-22-10pct-80-1machine-t5-r25-w500-s1s.out retesting-zero-all-affected-ip-as-2022-08-22-10pct-80-1machine-t5-r25-w500-s1s.out


compare: compare.py retesting-random-all-affected-ip-as-2022-08-22-10pct-80-1machine-t5-r25-w500-s1s.out retesting-zero-all-affected-ip-as-2022-08-22-10pct-80-1machine-t5-r25-w500-s1s.out
	$(PYTHON) compare.py --random retesting-random-all-affected-ip-as-2022-08-22-10pct-80-1machine-t5-r25-w500-s1s.out --zero retesting-zero-all-affected-ip-as-2022-08-22-10pct-80-1machine-t5-r25-w500-s1s.out | grep "true;false" | cut -d\; -f1 | cut -d\: -f1 > appear_blocked.txt

.PHONY: clean
clean:
	rm -f $(ALL)

.DELETE_ON_ERROR:
