
all: build

bin:
	mkdir -p bin

build: bin
	go build -o bin/quicklog ./cmd/quicklog
	go build -o bin/qlsearch ./cmd/qlsearch
	go build -o bin/ql2etcd ./cmd/ql2etcd

linux: bin
	GOOS=linux go build -o bin/quicklog-linux ./cmd/quicklog
	GOOS=linux go build -o bin/qlsearch-linux ./cmd/qlsearch
	GOOS=linux go build -o bin/ql2etcd-linux ./cmd/ql2etcd

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
