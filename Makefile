
all: build

bin:
	mkdir -p bin

build: bin
	go build -o bin/quicklog ./cmd/quicklog
	go build -o bin/ql2etcd ./cmd/ql2etcd
	go build -o bin/ql-embedded-example ./examples/embedded

linux: bin
	CGO_ENABLED=0 GOOS=linux go build -ldflags "-s" -a -installsuffix cgo -o bin/quicklog-linux ./cmd/quicklog

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
