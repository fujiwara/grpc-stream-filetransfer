build: gen-proto
	go build -o bin/grpcp cmd/grpcp/main.go

gen-proto:
	protoc --go_out=. --go-grpc_out=. filetransfer.proto
