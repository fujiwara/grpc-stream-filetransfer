package main

import (
	"log"
	"os"

	filetransfer "github.com/fujiwara/grpc-stream-filetransfer"
)

func main() {
	if len(os.Args) != 2 {
		log.Println("Usage: client <filename>")
		os.Exit(1)
	}
	if err := filetransfer.RunClient(os.Args[1]); err != nil {
		log.Fatal(err)
	}
}
