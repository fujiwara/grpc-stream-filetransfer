package filetransfer

import (
	"io"
	"log"
	"net"
	"os"
	"path/filepath"
	"sync"

	pb "github.com/fujiwara/grpc-stream-filetransfer/proto"
	"github.com/schollz/progressbar/v3"
	"google.golang.org/grpc"
)

var basedir = "/tmp"

type server struct {
	pb.UnimplementedFileTransferServiceServer
}

func (s *server) Upload(stream pb.FileTransferService_UploadServer) error {
	// ファイル受信処理
	var open sync.Once
	var f *os.File
	var bar *progressbar.ProgressBar
	var w io.Writer
	defer func() {
		if f != nil {
			f.Close()
		}
	}()
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			// ファイル受信完了
			log.Printf("EOF")
			return stream.SendAndClose(&pb.FileUploadResponse{Message: "Upload received successfully"})
		} else if err != nil {
			return stream.SendAndClose(&pb.FileUploadResponse{Message: "Failed to receive file"})
		}
		open.Do(func() {
			filename := filepath.Join(basedir, req.Filename)
			log.Println("open file:", filename)
			f, err = os.OpenFile(filename, os.O_WRONLY|os.O_CREATE, 0644)
			if err != nil {
				return
			}
			bar = progressbar.DefaultBytes(
				req.Size,
				"recieving",
			)
			w = io.MultiWriter(f, bar)
		})
		if err != nil {
			log.Printf("Failed to open file: %s", err)
			return stream.SendAndClose(&pb.FileUploadResponse{Message: "Failed to open file"})
		}
		if _, err := w.Write(req.Content); err != nil {
			log.Printf("Failed to write file: %s", err)
			return stream.SendAndClose(&pb.FileUploadResponse{Message: "Failed to write file"})
		}
	}
}

func (s *server) Download(req *pb.FileDownloadRequest, stream pb.FileTransferService_DownloadServer) error {
	// ファイル送信処理
	// ここでファイル内容を読み込んでクライアントに送信する
	return nil
}

func RunServer() {
	lis, err := net.Listen("tcp", ":5000")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	s := grpc.NewServer()
	log.Println("Starting server")
	pb.RegisterFileTransferServiceServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
