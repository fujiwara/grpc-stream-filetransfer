package filetransfer

import (
	"io"
	"log"
	"net"
	"os"
	"sync"

	pb "github.com/fujiwara/grpc-stream-filetransfer/proto"
	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedFileTransferServiceServer
}

func (s *server) Upload(stream pb.FileTransferService_UploadServer) error {
	// ファイル受信処理
	var open sync.Once
	var f *os.File
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
			log.Println("open file:", req.Filename)
			f, err = os.OpenFile(req.Filename, os.O_WRONLY|os.O_CREATE, 0644)
			if err != nil {
				return
			}
		})
		if err != nil {
			log.Printf("Failed to open file: %s", err)
			return stream.SendAndClose(&pb.FileUploadResponse{Message: "Failed to open file"})
		}
		if _, err := f.Write(req.Content); err != nil {
			log.Printf("Failed to write file: %s", err)
			return stream.SendAndClose(&pb.FileUploadResponse{Message: "Failed to write file"})
		}
	}
}

func (s *server) Download(req *pb.FileDownloadRequest, stream pb.FileTransferService_DownloadServer) error {
	f, err := os.OpenFile(req.Filename, os.O_RDONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	st, err := f.Stat()
	if err != nil {
		return err
	}
	for {
		buf := make([]byte, 4096)
		n, err := f.Read(buf)
		if err == io.EOF {
			// ファイル読み込み完了
			break
		}
		if err != nil {
			log.Println("Failed to read file:", err)
			return err
		}
		// 読み込んだデータをクライアントに送信
		if err := stream.Send(&pb.FileDownloadResponse{
			Filename: req.Filename,
			Content:  buf[:n],
			Size:     st.Size(),
		}); err != nil {
			log.Println("Failed to send file:", err)
			return err
		}
	}

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
