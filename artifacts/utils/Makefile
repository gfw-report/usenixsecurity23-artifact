PYTHON = python3

ALL = affected-norand \
	affected-payload \
	whitelist\
	detect.py\
	compare.py

.PHONY: all
all: $(ALL)

affected-norand: affected-norand.go
	env GOOS=linux GOARCH=amd64 go build $^

affected-payload: affected-payload.go
	env GOOS=linux GOARCH=amd64 go build $^

whitelist: whitelist.go
	env GOOS=linux GOARCH=amd64 go build $^

install: $(ALL)
	rsync -auvP $(ALL) usenix-ae-client-china:~/utils

.PHONY: clean
clean:
	rm -f $(ALL)

.DELETE_ON_ERROR:
