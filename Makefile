clean:
	rm -fv Mock*.go
	rm -fv **/Mock*.go
	go clean

mocks:
	mockery

code-generation:
	rm -fv **/*_impl.go
	go generate ./adapters/*

test:
	go clean -testcache
	go test ./...

all: clean code-generation mocks
	go clean -testcache
	go test ./...