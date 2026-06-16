ROLE := kvpool
.PHONY: build test test-standalone-layout
build:
	go build -o bin/cofiswarm-kvpool ./cmd/cofiswarm-kvpool
test: build test-standalone-layout
test-standalone-layout:
	./test/scripts/assert-layout.sh $(ROLE)
