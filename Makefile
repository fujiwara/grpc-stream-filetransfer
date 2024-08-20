build: gen-proto
	go build -o bin/grpcp \
		-ldflags '-s -w' \
		cmd/grpcp/main.go

clean:
	rm -rf bin/*

gen-proto:
	protoc --go_out=. --go-grpc_out=. filetransfer.proto

test: gen-proto
	go test -v ./...
