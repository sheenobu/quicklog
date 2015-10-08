
build:
	go build ./cmd/quicklog
	go build ./cmd/qlsearch
	go build ./cmd/ql2etcd

clean:
	rm -f ./quicklog

test:
	go test ./...

fmt:
	go fmt ./...

vet:
	go vet ./...

lint:
	golint ./...
