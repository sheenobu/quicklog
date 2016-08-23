GOFLAGS ?= $(GOFLAGS:)

TAG := $(VERSION)
ifeq ($(TAG),)
  BRANCH := $(shell git rev-parse --abbrev-ref HEAD)
  DT := $(shell date '+%F_%H%M')
  VSN := $(BRANCH)-$(DT)
else
  VSN := $(TAG)
endif

ENV := $(shell printenv)

GOFLAGS = -ldflags '-X=main.version=$(VSN)'

all: build

bin:
	mkdir -p bin

build: bin
	go build $(GOFLAGS) -o bin/quicklog ./cmd/quicklog
	go build $(GOFLAGS) -o bin/ql2etcd ./cmd/ql2etcd
	go build $(GOFLAGS) -o bin/ql-embedded-example ./examples/embedded

linux: bin
	CGO_ENABLED=0 GOOS=linux go build $(GOFLAGS) -ldflags "-s" -a -installsuffix cgo -o bin/quicklog-linux ./cmd/quicklog

docker: linux
	docker build -t sheenobu/quicklog .

clean:
	rm -f ./bin/

test:
	go test ./...

fmt:
	go fmt ./...

vet:
	go vet ./...

lint:
	golint ./...
