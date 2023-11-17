package grpcp

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/alecthomas/kong"
	"github.com/fatih/color"
	"github.com/fujiwara/logutils"
)

type CLI struct {
	Port   int  `name:"port" short:"p" default:"8022" help:"port number"`
	Server bool `name:"server" short:"s" help:"run as server"`
	Quiet  bool `name:"quiet" short:"q" help:"quiet mode for client"`

	Src  string `arg:"" optional:"" name:"src" short:"s" description:"source file path"`
	Dest string `arg:"" optional:"" name:"dest" short:"d" description:"destination file path"`
}

func RunCLI(ctx context.Context) error {
	var cli CLI
	kong.Parse(&cli)

	opt := NewDefaultOption()
	opt.Port = cli.Port

	if cli.Server {
		return RunServer(ctx, opt)
	} else if cli.Src != "" && cli.Dest != "" {
		clinet := NewClient(opt)
		return clinet.Run(ctx, cli.Src, cli.Dest, cli.Quiet)
	} else {
		return fmt.Errorf("expected: grpcp <src> <dest> or grpcp --server")
	}
}

func init() {
	filter := &logutils.LevelFilter{
		Levels: []logutils.LogLevel{"debug", "info", "warn", "error"},
		ModifierFuncs: []logutils.ModifierFunc{
			logutils.Color(color.FgCyan),   // debug
			nil,                            // default
			logutils.Color(color.FgYellow), // warn
			logutils.Color(color.FgRed),    // error
		},
		MinLevel: logutils.LogLevel("info"),
		Writer:   os.Stderr,
	}
	log.SetOutput(filter)
}
