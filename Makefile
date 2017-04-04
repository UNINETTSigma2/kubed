GO_EXECUTABLE ?= go
PACKAGE_DIRS := $(shell glide nv)
VERSION := $(shell git describe --tags --dirty --always)
DIST_DIRS := find * -type d -exec

all: test build

build:
	mkdir -p dist/
	GOOS=linux GOARCH=amd64 ${GO_EXECUTABLE} build -o dist/kubed-linux-amd64 -ldflags "-X main.version=${VERSION}"
	GOOS=darwin GOARCH=amd64 ${GO_EXECUTABLE} build -o dist/kubed-darwin-amd64 -ldflags "-X main.version=${VERSION}"
	GOOS=windows GOARCH=amd64 ${GO_EXECUTABLE} build -o dist/linux-windows-amd64.exe -ldflags "-X main.version=${VERSION}"

test:
	${GO_EXECUTABLE} test --short $(PACKAGE_DIRS)

clean:
	rm -f ./kubed.test
	rm -f ./kubed
	rm -rf ./dist


