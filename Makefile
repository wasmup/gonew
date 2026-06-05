# Makefile
all:
	CGO_ENABLED=0 go install -ldflags=-s -trimpath

build:
	CGO_ENABLED=0 go build -trimpath -ldflags "-s -w -X main.Version=$(shell git describe --tags --always)"

install:
	CGO_ENABLED=0 go install -trimpath -ldflags "-s -w -X main.Version=$(shell git describe --tags --always)"

inspect:
	CGO_ENABLED=0 go install -ldflags "-s -w" -trimpath
	ls -l `which gonew`
	file `which gonew`
	sha256sum `which gonew`
	go version -m   $(which gonew)
	go tool nm $(which gonew)

release:
	go install -trimpath -ldflags "-s -w -X main.Version=v1.0.1"
	git tag v1.0.1
	git push --tags
