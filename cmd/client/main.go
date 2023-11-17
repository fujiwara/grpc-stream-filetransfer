package main

import (
	"log"
	"os"

	filetransfer "github.com/fujiwara/grpc-stream-filetransfer"
)

func main() {
	if len(os.Args) != 3 {
		log.Println("Usage: client <srcfile> <destfile>")
		os.Exit(1)
	}
	if err := filetransfer.RunClient(os.Args[1], os.Args[2]); err != nil {
		log.Fatal(err)
	}
}
