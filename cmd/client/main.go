package main

import (
	"context"
	"log"
	"os"

	filetransfer "github.com/fujiwara/grpc-stream-filetransfer"
)

func main() {
	// TODO flag parser
	if len(os.Args) != 3 {
		log.Println("Usage: client <srcfile> <destfile>")
		os.Exit(1)
	}
	opt := filetransfer.NewDefaultOption()
	client := filetransfer.NewClient(opt)
	if err := client.Run(context.Background(), os.Args[1], os.Args[2]); err != nil {
		log.Fatal(err)
	}
}
