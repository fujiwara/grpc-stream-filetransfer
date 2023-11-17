package filetransfer

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"sync"

	pb "github.com/fujiwara/grpc-stream-filetransfer/proto"
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

	stream, err := client.Upload(ctx)
	if err != nil {
		return fmt.Errorf("failed to new upload stream: %w", err)
	}
	log.Printf("staring upload: %s -> %s", localFile, remoteFile)

	bar := progressbar.DefaultBytes(
		st.Size(),
		"uploading",
	)
	buf := make([]byte, 4096)
	for {
		n, err := file.Read(buf)
		if err == io.EOF {
			// ファイル読み込み完了
			break
		}
		if err != nil {
			return fmt.Errorf("failed to read file: %w", err)
		}

		// 読み込んだデータをサーバーに送信
		req := &pb.FileUploadRequest{
			Filename: remoteFile,
			Content:  buf[:n],
			Size:     st.Size(),
		}
		bar.Write(req.Content)
		if err := stream.Send(req); err != nil {
			if err == io.EOF {
				log.Println("EOF")
				break
			}
			return fmt.Errorf("failed to send file: %w", err)
		}
	}

	// アップロード終了をサーバーに通知
	response, err := stream.CloseAndRecv()
	if err != nil {
		return fmt.Errorf("failed to receive response: %w", err)
	}
	log.Printf("Upload response: %s", response.Message)
	return nil
}

func downloadFile(ctx context.Context, client pb.FileTransferServiceClient, remoteFile string, localFile string) error {
	stream, err := client.Download(ctx, &pb.FileDownloadRequest{
		Filename: remoteFile,
	})
	if err != nil {
		return fmt.Errorf("failed to new download stream: %w", err)
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
	for {
		res, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("failed to receive response: %w", err)
		}
		once.Do(func() {
			bar = progressbar.DefaultBytes(
				res.Size,
				"downloading",
			)
			w = io.MultiWriter(f, bar)
		})
		if _, err := w.Write(res.Content); err != nil {
			return fmt.Errorf("failed to write file: %w", err)
		}
	}
	return nil
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

func (c *Client) Run(ctx context.Context, src, dest string) error {
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
		// local to local and remote to remote are not supported
		return fmt.Errorf("both src and dest are local or remote")
	}

	conn, err := grpc.Dial(fmt.Sprintf("%s:%d", remoteHost, c.Option.Port),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return fmt.Errorf("failed to dial server: %w", err)
	}
	defer conn.Close()
	client := pb.NewFileTransferServiceClient(conn)

	return transfer(ctx, client, remoteFile, localFile)
}

func parseFilename(filename string) (string, string) {
	p := strings.SplitN(filename, ":", 2)
	if len(p) == 1 {
		return "", p[0] // local
	}
	return p[0], p[1] // remote
}
