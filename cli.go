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
	Host string `name:"host" short:"h" default:"localhost" help:"host name"`
	Port int    `name:"port" short:"p" default:"8022" help:"port number"`

	Server bool `name:"server" short:"s" help:"run as server"`
	Quiet  bool `name:"quiet" short:"q" help:"quiet mode for client"`
	Kill   bool `name:"kill" short:"k" help:"kill server"`

	Src  string `arg:"" optional:"" name:"src" short:"s" description:"source file path"`
	Dest string `arg:"" optional:"" name:"dest" short:"d" description:"destination file path"`
}

func RunCLI(ctx context.Context) error {
	var cli CLI
	kong.Parse(&cli)

	if cli.Quiet {
		setLogLevel("warn")
	} else {
		setLogLevel("info")
	}

	if cli.Server {
		opt := &ServerOption{
			Port:   cli.Port,
			Listen: cli.Host,
		}
		return RunServer(ctx, opt)
	} else if cli.Kill {
		opt := &ClientOption{
			Host:  cli.Host,
			Port:  cli.Port,
			Quiet: cli.Quiet,
		}
		clinet := NewClient(opt)
		return clinet.Shutdown(ctx)
	} else if cli.Src != "" && cli.Dest != "" {
		opt := &ClientOption{
			Port:  cli.Port,
			Quiet: cli.Quiet,
		}
		clinet := NewClient(opt)
		return clinet.Copy(ctx, cli.Src, cli.Dest)
	} else {
		return fmt.Errorf("expected: grpcp <src> <dest> or grpcp --server")
	}
}

func setLogLevel(level string) {
	filter := &logutils.LevelFilter{
		Levels: []logutils.LogLevel{"debug", "info", "warn", "error"},
		ModifierFuncs: []logutils.ModifierFunc{
			logutils.Color(color.FgCyan),   // debug
			nil,                            // default
			logutils.Color(color.FgYellow), // warn
			logutils.Color(color.FgRed),    // error
		},
		MinLevel: logutils.LogLevel(level),
		Writer:   os.Stderr,
	}
	log.SetOutput(filter)
}
