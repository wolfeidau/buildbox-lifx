GOPATH=$(shell pwd)/.gopath

debug:
	scripts/build.sh

clean:
	rm -f bin/buildbox-lifx || true
	rm -rf .gopath || true

test:
	cd .gopath/src/github.com/wolfeidau/buildbox-lifx && go get -t ./...
	cd .gopath/src/github.com/wolfeidau/buildbox-lifx && go test ./...

vet:
	go vet ./...

.PHONY: debug clean test vet
