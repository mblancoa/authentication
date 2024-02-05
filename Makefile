clean:
	rm -fv Mock*.go
	rm -fv **/Mock*.go
	go clean

mocks:
	mockery

code-generation:
	rm -fv **/*_impl.go
	go generate ./adapter/*

test:
	go test ./...


