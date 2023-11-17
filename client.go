package grpcp

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"

	pb "github.com/fujiwara/grpcp/proto"
	"github.com/schollz/progressbar/v3"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func uploadFile(ctx context.Context, client pb.FileTransferServiceClient, remoteFile, localFile string) error {
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
	log.Printf("staring upload: %s -> %s (%d bytes)", localFile, remoteFile, st.Size())
	bar := progressbar.DefaultBytes(
		st.Size(),
		"uploading",
	)
	expectedBytes := st.Size()
	var totalBytes int64
	buf := make([]byte, StreamBufferSize)
	for {
		n, err := file.Read(buf)
		if err == io.EOF {
			log.Printf("[info] upload completed (%d bytes)", totalBytes)
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
	log.Printf("[info] server respone: %s", res.Message)
	return nil
}

func downloadFile(ctx context.Context, client pb.FileTransferServiceClient, remoteFile string, localFile string) error {
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

	log.Printf("staring download: %s -> %s", remoteFile, localFile)

	var once sync.Once
	var bar *progressbar.ProgressBar
	var w io.Writer
	var expectedBytes, totalBytes int64
	for {
		res, err := stream.Recv()
		if err == io.EOF {
			log.Printf("[info] download completed (%d bytes)", totalBytes)
			if totalBytes != expectedBytes {
				return fmt.Errorf("file size mismatch: expected %d bytes, got %d bytes", expectedBytes, totalBytes)
			}
			return nil
		} else if err != nil {
			return fmt.Errorf("failed to receive response: %w", err)
		}
		once.Do(func() {
			expectedBytes = res.Size
			bar = progressbar.DefaultBytes(
				res.Size,
				"downloading",
			)
			w = io.MultiWriter(f, bar)
		})
		if n, err := w.Write(res.Content); err != nil {
			return fmt.Errorf("failed to write file: %w", err)
		} else {
			totalBytes += int64(n)
		}
	}
}

type transferFunc func(ctx context.Context, client pb.FileTransferServiceClient, src, dest string) error

type Client struct {
	Option *Option
}

func NewClient(opt *Option) *Client {
	return &Client{
		Option: opt,
	}
}

func (c *Client) Run(ctx context.Context, src, dest string, quiet bool) error {
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
	conn, err := grpc.Dial(addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return fmt.Errorf("failed to dial server: %w", err)
	}
	defer conn.Close()
	client := pb.NewFileTransferServiceClient(conn)
	log.Printf("[info] connected to %s", addr)

	return transfer(ctx, client, remoteFile, localFile)
}

func parseFilename(filename string) (string, string) {
	p := strings.SplitN(filename, ":", 2)
	if len(p) == 1 {
		return "", p[0] // local
	}
	return p[0], p[1] // remote
}
