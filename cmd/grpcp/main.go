package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/fujiwara/grpcp"
)

func main() {
	if err := grpcp.RunCLI(context.Background()); err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
}
