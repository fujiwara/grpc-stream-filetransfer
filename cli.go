package grpcp

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/alecthomas/kong"
)

var LogLevel = new(slog.LevelVar)

type CLI struct {
	Host string `name:"host" short:"h" default:"localhost" help:"host name"`
	Port int    `name:"port" short:"p" default:"8022" help:"port number"`

	Quiet bool `name:"quiet" short:"q" help:"quiet mode for client"`
	Debug bool `name:"debug" short:"d" help:"enable debug log for client and server"`

	Server bool   `name:"server" short:"s" help:"run as server"`
	TLS    bool   `name:"tls" negatable:"" default:"true" help:"enable TLS for server"`
	Cert   string `name:"cert" help:"certificate file for server"`
	Key    string `name:"key" help:"private key file for server"`

	SkipVerify bool `name:"skip-verify" help:"skip TLS verification for client"`

	Kill bool `name:"kill" short:"k" help:"kill server"`

	Src  string `arg:"" optional:"" name:"src" short:"s" description:"source file path"`
	Dest string `arg:"" optional:"" name:"dest" short:"d" description:"destination file path"`
}

func RunCLI(ctx context.Context) error {
	var cli CLI
	kong.Parse(&cli)

	if cli.Quiet {
		slog.SetLogLoggerLevel(slog.LevelWarn)
	} else if cli.Debug {
		slog.SetLogLoggerLevel(slog.LevelDebug)
	} else {
		slog.SetLogLoggerLevel(slog.LevelInfo)
	}

	if cli.Server {
		opt := &ServerOption{
			Port:   cli.Port,
			Listen: cli.Host,
			TLS:    cli.TLS,
		}
		return RunServer(ctx, opt)
	} else if cli.Kill {
		opt := &ClientOption{
			Host:       cli.Host,
			Port:       cli.Port,
			Quiet:      cli.Quiet,
			TLS:        cli.TLS,
			SkipVerify: cli.SkipVerify,
		}
		client := NewClient(opt)
		return client.Shutdown(ctx)
	} else if cli.Src != "" && cli.Dest != "" {
		opt := &ClientOption{
			Port:       cli.Port,
			Quiet:      cli.Quiet,
			TLS:        cli.TLS,
			SkipVerify: cli.SkipVerify,
		}
		client := NewClient(opt)
		return client.Copy(ctx, cli.Src, cli.Dest)
	} else {
		return fmt.Errorf("expected: grpcp <src> <dest> or grpcp --server")
	}
}
