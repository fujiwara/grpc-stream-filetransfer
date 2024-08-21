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

	Quiet bool `name:"quiet" short:"q" help:"quiet mode"`
	Debug bool `name:"debug" short:"d" help:"enable debug log"`
	TLS   bool `name:"tls" negatable:"" default:"true" help:"enable TLS (default: true)"`

	Server bool   `name:"server" short:"s" help:"run as server"`
	Cert   string `name:"cert" help:"certificate file for server" type:"existingfile"`
	Key    string `name:"key" help:"private key file for server" type:"existingfile"`

	SkipVerify bool `name:"skip-verify" help:"skip TLS verification for client"`
	Kill       bool `name:"kill" help:"send shutdown command to server"`
	Ping       bool `name:"ping" help:"send ping message to server"`

	Src  string `arg:"" optional:"" name:"src" short:"s" description:"source file path"`
	Dest string `arg:"" optional:"" name:"dest" short:"d" description:"destination file path"`
}

func (c *CLI) ClientOption() *ClientOption {
	return &ClientOption{
		Host:       c.Host,
		Port:       c.Port,
		Quiet:      c.Quiet,
		TLS:        c.TLS,
		SkipVerify: c.SkipVerify,
	}
}

func (c *CLI) ServerOption() *ServerOption {
	return &ServerOption{
		Port:     c.Port,
		Listen:   c.Host,
		TLS:      c.TLS,
		CertFile: c.Cert,
		KeyFile:  c.Key,
	}
}

func RunCLI(ctx context.Context) error {
	cli := &CLI{}
	kong.Parse(cli)

	if cli.Quiet {
		slog.SetLogLoggerLevel(slog.LevelWarn)
	} else if cli.Debug {
		slog.SetLogLoggerLevel(slog.LevelDebug)
	} else {
		slog.SetLogLoggerLevel(slog.LevelInfo)
	}

	client := NewClient(cli.ClientOption())
	switch {
	case cli.Server:
		return RunServer(ctx, cli.ServerOption())
	case cli.Ping:
		resp, err := client.Ping(ctx)
		if err != nil {
			return err
		}
		slog.Info("ping", "message", resp.Message)
		return nil
	case cli.Kill:
		return client.Shutdown(ctx)
	case cli.Src != "" && cli.Dest != "":
		return client.Copy(ctx, cli.Src, cli.Dest)
	default:
		return fmt.Errorf("expected: grpcp <src> <dest> or grpcp --server. see --help")
	}
}
