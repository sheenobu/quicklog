
build:
	go build ./cmd/quicklog

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
