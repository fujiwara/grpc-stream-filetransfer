build: gen-proto
	go build -o bin/server cmd/server/main.go
	go build -o bin/client cmd/client/main.go

gen-proto:
	protoc --go_out=. --go-grpc_out=. filetransfer.proto
