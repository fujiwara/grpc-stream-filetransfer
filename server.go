package grpcp

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"sync"

	pb "github.com/fujiwara/grpcp/proto"
	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedFileTransferServiceServer
}

var (
	StreamBufferSize = 1024 * 1024
)

func (s *server) Ping(ctx context.Context, req *pb.PingRequest) (*pb.PingResponse, error) {
	log.Printf("[info] ping message: %s", req.Message)
	return &pb.PingResponse{Message: "pong"}, nil
}

func newUploadResponse(msg string) *pb.FileUploadResponse {
	return &pb.FileUploadResponse{Message: msg}
}

func (s *server) Upload(stream pb.FileTransferService_UploadServer) error {
	if err := s.upload(stream); err != nil {
		log.Printf("[error] %s", err)
		return err
	}
	return nil
}

func (s *server) upload(stream pb.FileTransferService_UploadServer) error {
	var once sync.Once
	var f *os.File
	var totalBytes, expectedSize int64
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			log.Printf("[info] server upload completed (%d bytes)", totalBytes)
			if totalBytes != expectedSize {
				return fmt.Errorf("file size mismatch: expected %d bytes, got %d bytes", expectedSize, totalBytes)
			}
			return stream.SendAndClose(newUploadResponse("Upload received successfully"))
		} else if err != nil {
			return fmt.Errorf("failed to receive file: %w", err)
		}
		once.Do(func() {
			log.Printf("[info] server accepting upload request: %s (%d bytes)", req.Filename, req.Size)
			f, err = os.OpenFile(req.Filename, os.O_WRONLY|os.O_CREATE, 0644)
			expectedSize = req.Size
		})
		if err != nil || f == nil {
			return fmt.Errorf("failed to open file: %w", err)
		}
		if n, err := f.Write(req.Content); err != nil {
			return fmt.Errorf("failed to write file: %w", err)
		} else {
			totalBytes += int64(n)
		}
	}
}

func (s *server) Download(req *pb.FileDownloadRequest, stream pb.FileTransferService_DownloadServer) error {
	if err := s.download(req, stream); err != nil {
		log.Printf("[error] %s", err)
		return err
	}
	return nil
}

func (s *server) download(req *pb.FileDownloadRequest, stream pb.FileTransferService_DownloadServer) error {
	log.Printf("[info] server accepting download request: %s", req.Filename)
	f, err := os.OpenFile(req.Filename, os.O_RDONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer f.Close()
	st, err := f.Stat()
	if err != nil {
		return fmt.Errorf("failed to stat file: %w", err)
	}
	expectedBytes := st.Size()
	totalBytes := int64(0)
	buf := make([]byte, StreamBufferSize)
	for {
		n, err := f.Read(buf)
		if err == io.EOF {
			log.Printf("[info] server download completed (%d bytes)", totalBytes)
			if totalBytes != expectedBytes {
				return fmt.Errorf("file size mismatch: expected %d bytes, got %d bytes", expectedBytes, totalBytes)
			}
			return nil
		} else if err != nil {
			return fmt.Errorf("failed to read file: %w", err)
		}
		if err := stream.Send(&pb.FileDownloadResponse{
			Filename: req.Filename,
			Content:  buf[:n],
			Size:     expectedBytes,
		}); err != nil {
			return fmt.Errorf("failed to send file: %w", err)
		}
		totalBytes += int64(n)
	}
}

func RunServer(ctx context.Context, opt *ServerOption) error {
	addr := fmt.Sprintf("%s:%d", opt.Listen, opt.Port)
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("failed to listen: %w", err)
	}
	s := grpc.NewServer()
	log.Println("[info] Starting server on", addr, "...")
	pb.RegisterFileTransferServiceServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		return fmt.Errorf("failed to serve: %w", err)
	}
	return nil
}
