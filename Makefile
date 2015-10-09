
build:
	go build ./cmd/quicklog
	go build ./cmd/qlsearch
	go build ./cmd/ql2etcd

linux:
	GOOS=linux go build -o quicklog-linux ./cmd/quicklog
	GOOS=linux go build -o qlsearch-linux ./cmd/qlsearch
	GOOS=linux go build -o ql2etcd-linux ./cmd/ql2etcd

docker: linux
	docker build -t sheenobu/quicklog .

clean:
	rm -f ./quicklog ./quicklog-linux ./qlsearch ./qlsearch-linux ./ql2etcd ./ql2etcd-linux

test:
	go test ./...

fmt:
	go fmt ./...

vet:
	go vet ./...

lint:
	golint ./...
