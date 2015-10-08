
build:
	go build ./cmd/quicklog
	go build ./cmd/qlsearch

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
