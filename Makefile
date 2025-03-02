VERSION = $(shell godzil show-version)
CURRENT_REVISION = $(shell git rev-parse --short HEAD)
BUILD_LDFLAGS = "-s -w -X github.com/getperf/getperf2.revision=$(CURRENT_REVISION)"
u := $(if $(update),-u)

export GO111MODULE=on

.PHONY: deps
deps:
	go get ${u} -d
	go mod tidy

.PHONY: devel-deps
devel-deps:
	sh -c '\
      tmpdir=$$(mktemp -d); \
      cd $$tmpdir; \
      go get ${u} \
        golang.org/x/lint/golint            \
        github.com/Songmu/godzil/cmd/godzil \
        github.com/tcnksm/ghr; \
      rm -rf $$tmpdir'

.PHONY: test
test:
	go test

.PHONY: lint
lint: devel-deps
	golint -set_exit_status

.PHONY: buildgetperf
buildgetperf:
	go build -ldflags=$(BUILD_LDFLAGS)  -o ./bin/_getperf ./cmd/getperf2
	GOOS=windows GOARCH=386 go build -ldflags=$(BUILD_LDFLAGS) -o ./bin/getperf.exe ./cmd/getperf2

.PHONY: build
build:
	go build -ldflags=$(BUILD_LDFLAGS) -o ./bin/_getperf ./cmd/getperf2
	go build -ldflags=$(BUILD_LDFLAGS) ./cmd/gconf
	GOOS=windows GOARCH=386 go build -ldflags=$(BUILD_LDFLAGS) -o ./bin/getperf.exe ./cmd/getperf2
	GOOS=windows GOARCH=386 go build -ldflags=$(BUILD_LDFLAGS) ./cmd/gconf

.PHONY: install
install:
	go install -ldflags=$(BUILD_LDFLAGS) ./cmd/getperf2
	go install -ldflags=$(BUILD_LDFLAGS) ./cmd/gconf

.PHONY: release
release: devel-deps
	godzil release

CREDITS: go.sum deps devel-deps
	godzil credits -w

.PHONY: crossbuild
crossbuild: CREDITS
	godzil crossbuild -pv=v$(VERSION) -build-ldflags=$(BUILD_LDFLAGS) \
      -os=linux,darwin -d=./dist/v$(VERSION) ./cmd/*

.PHONY: upload
upload:
	ghr -body="$$(godzil changelog --latest -F markdown)" v$(VERSION) dist/v$(VERSION)
