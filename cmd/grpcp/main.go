package main

import (
	"context"
	"log"
	"os"

	"github.com/fujiwara/grpcp"
)

func main() {
	if err := grpcp.RunCLI(context.Background()); err != nil {
		log.Println("[error]", err)
		os.Exit(1)
	}
}
