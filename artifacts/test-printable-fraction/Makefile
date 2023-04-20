PYTHON = python3

ALL = test-probes.txt\
	prediction.csv

.PHONY: all
all: $(ALL)

test-probes.txt: test-printable-frac.py
	$(PYTHON) $^ > "$@"

prediction.csv: ../utils/detect.py test-probes.txt
	$(PYTHON) $^ --header > "$@"

test: test-probes.txt
	../utils/affected-payload < $^ -host REDACTED_US_SERVER_IP -p 10000-11000 | tee result.csv

compare: prediction.csv result.csv
	../utils/compare.py --predict prediction.csv --actual result.csv

.PHONY: clean
clean:
	rm -f $(ALL)

.DELETE_ON_ERROR:
