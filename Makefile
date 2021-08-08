override GOFLAGS += -mod=vendor

all:

build:
	./build/build_binaries.sh

test:
	go test $(GOFLAGS) ./...

