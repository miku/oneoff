SHELL := /bin/bash
TARGETS := colorchars dijkstra65 slowdown tarcheck webshare localbin uppsala safename

.PHONY: all
all: $(TARGETS)

%: cmd/%/main.go
	go build -ldflags="-s -w" -o $@ $<

.PHONY: clean
clean:
	rm -f $(TARGETS)


