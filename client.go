package grpcp

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"sync"

	pb "github.com/fujiwara/grpcp/proto"
	"github.com/schollz/progressbar/v3"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func uploadFile(ctx context.Context, client pb.FileTransferServiceClient, remoteFile, localFile string, opt *ClientOption) error {
	file, err := os.Open(localFile)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()
	st, err := file.Stat()
	if err != nil {
		return fmt.Errorf("failed to stat file: %w", err)
	}

	// if remoteFile is directory, use localFile's basename
	if strings.HasSuffix(remoteFile, "/") {
		remoteFile = filepath.Join(remoteFile, filepath.Base(localFile))
	}

	stream, err := client.Upload(ctx)
	if err != nil {
		return fmt.Errorf("failed to new upload stream: %w", err)
	}
	slog.Info("staring upload", "local", localFile, "remote", remoteFile, "bytes", st.Size())
	var bar io.Writer
	if opt.Quiet {
		bar = io.Discard
	} else {
		bar = progressbar.DefaultBytes(st.Size(), "uploading")
	}
	expectedBytes := st.Size()
	var totalBytes int64
	buf := make([]byte, StreamBufferSize)
	for {
		n, err := file.Read(buf)
		if err == io.EOF {
			slog.Info("client upload completed", "bytes", totalBytes)
			if totalBytes != expectedBytes {
				return fmt.Errorf("file size mismatch: expected %d bytes, got %d bytes", st.Size(), totalBytes)
			}
			break
		} else if err != nil {
			return fmt.Errorf("failed to read file: %w", err)
		}

		req := &pb.FileUploadRequest{
			Filename: remoteFile,
			Content:  buf[:n],
			Size:     expectedBytes,
		}
		if err := stream.Send(req); err != nil {
			return fmt.Errorf("failed to send file: %w", err)
		}
		bar.Write(req.Content)
		totalBytes += int64(n)
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		return fmt.Errorf("failed to receive response: %w", err)
	}
	slog.Info("server response", "message", res.Message)
	return nil
}

func downloadFile(ctx context.Context, client pb.FileTransferServiceClient, remoteFile, localFile string, opt *ClientOption) error {
	stream, err := client.Download(ctx, &pb.FileDownloadRequest{
		Filename: remoteFile,
	})
	if err != nil {
		return fmt.Errorf("failed to new download stream: %w", err)
	}

	// if localFile is directory, use remoteFile's basename
	if st, err := os.Stat(localFile); err == nil && st.IsDir() {
		localFile = filepath.Join(localFile, filepath.Base(remoteFile))
	}

	f, err := os.OpenFile(localFile, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer f.Close()

	slog.Info("staring download", "remote", remoteFile, "local", localFile)

	var once sync.Once
	var w io.Writer
	var expectedBytes, totalBytes int64
	for {
		res, err := stream.Recv()
		if err == io.EOF {
			slog.Info("client download completed", "bytes", totalBytes)
			if totalBytes != expectedBytes {
				return fmt.Errorf("file size mismatch: expected %d bytes, got %d bytes", expectedBytes, totalBytes)
			}
			return nil
		} else if err != nil {
			return fmt.Errorf("failed to receive response: %w", err)
		}
		once.Do(func() {
			expectedBytes = res.Size
			if opt.Quiet {
				w = f
			} else {
				bar := progressbar.DefaultBytes(res.Size, "downloading")
				w = io.MultiWriter(f, bar)
			}
		})
		if n, err := w.Write(res.Content); err != nil {
			return fmt.Errorf("failed to write file: %w", err)
		} else {
			totalBytes += int64(n)
		}
	}
}

type transferFunc func(ctx context.Context, client pb.FileTransferServiceClient, src, dest string, opt *ClientOption) error

type Client struct {
	Option *ClientOption
}

func NewClient(opt *ClientOption) *Client {
	return &Client{
		Option: opt,
	}
}

func (c *Client) Ping(ctx context.Context) (*pb.PingResponse, error) {
	addr := fmt.Sprintf("%s:%d", c.Option.Host, c.Option.Port)
	conn, err := grpc.DialContext(ctx, addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to dial server: %w", err)
	}
	defer conn.Close()
	client := pb.NewFileTransferServiceClient(conn)
	return client.Ping(ctx, &pb.PingRequest{Message: "ping"})
}

func (c *Client) Copy(ctx context.Context, src, dest string) error {
	var transfer transferFunc
	var remoteHost, remoteFile, localFile string

	srcHost, srcFile := parseFilename(src)
	destHost, destFile := parseFilename(dest)
	if srcHost != "" && destHost != "" {
		return fmt.Errorf("both src and dest are remote")
	}
	if srcHost != "" && destHost == "" {
		// remote to local (download)
		transfer = downloadFile
		remoteHost = srcHost
		remoteFile = srcFile
		localFile = destFile
	} else if srcHost == "" && destHost != "" {
		// local to remote (upload)
		transfer = uploadFile
		remoteHost = destHost
		remoteFile = destFile
		localFile = srcFile
	} else {
		return fmt.Errorf("both src and dest are local")
	}

	addr := fmt.Sprintf("%s:%d", remoteHost, c.Option.Port)
	conn, err := grpc.DialContext(ctx, addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return fmt.Errorf("failed to dial server: %w", err)
	}
	slog.Info("connected", "remote", addr)
	defer conn.Close()
	client := pb.NewFileTransferServiceClient(conn)

	return transfer(ctx, client, remoteFile, localFile, c.Option)
}

func (c *Client) Shutdown(ctx context.Context) error {
	addr := fmt.Sprintf("%s:%d", c.Option.Host, c.Option.Port)
	conn, err := grpc.DialContext(ctx, addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return fmt.Errorf("failed to dial server: %w", err)
	}
	defer conn.Close()
	client := pb.NewFileTransferServiceClient(conn)
	_, err = client.Shutdown(ctx, &pb.ShutdownRequest{})
	return err
}

func parseFilename(filename string) (string, string) {
	p := strings.SplitN(filename, ":", 2)
	if len(p) == 1 {
		return "", p[0] // local
	}
	return p[0], p[1] // remote
}
