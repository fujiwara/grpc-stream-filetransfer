package filetransfer

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	pb "github.com/fujiwara/grpc-stream-filetransfer/proto"
	"github.com/schollz/progressbar/v3"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func uploadFile(client pb.FileTransferServiceClient, filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	st, err := file.Stat()
	if err != nil {
		return fmt.Errorf("failed to stat file: %w", err)
	}

	stream, err := client.Upload(context.Background())
	if err != nil {
		return fmt.Errorf("failed to new upload stream: %w", err)
	}
	log.Println("staring upload:", filename)
	basename := filepath.Base(filename)

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
			Filename: basename,
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

func RunClient(filename string) error {
	conn, err := grpc.Dial("localhost:5000",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return fmt.Errorf("failed to dial server: %w", err)
	}
	defer conn.Close()

	client := pb.NewFileTransferServiceClient(conn)
	if err := uploadFile(client, filename); err != nil {
		return fmt.Errorf("failed to upload file: %w", err)
	}
	return nil
}
