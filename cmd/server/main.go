package main

import (
	"context"

	filetransfer "github.com/fujiwara/grpc-stream-filetransfer"
)

func main() {
	// TODO flag parser
	opt := filetransfer.NewDefaultOption()
	filetransfer.RunServer(context.Background(), opt)
}
