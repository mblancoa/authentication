clean:
#	rm -fv Mock*.go
#	rm -fv **/Mock*.go
	find -type f -name 'Mock*.go' -print -delete
	go clean

mocks:
	mockery

code-generation:
#	rm -fv **/*_impl.go
	find -type f -name '*_impl.go' -print -delete
	go generate ./adapters/*

test:
	go clean -testcache
	go test ./...

all: clean code-generation mocks
	go clean -testcache
	go test ./...