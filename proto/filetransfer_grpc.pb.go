// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.12.4
// source: filetransfer.proto

package proto

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// FileTransferServiceClient is the client API for FileTransferService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type FileTransferServiceClient interface {
	Upload(ctx context.Context, opts ...grpc.CallOption) (FileTransferService_UploadClient, error)
	Download(ctx context.Context, in *FileDownloadRequest, opts ...grpc.CallOption) (FileTransferService_DownloadClient, error)
	Ping(ctx context.Context, in *PingRequest, opts ...grpc.CallOption) (*PingResponse, error)
}

type fileTransferServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewFileTransferServiceClient(cc grpc.ClientConnInterface) FileTransferServiceClient {
	return &fileTransferServiceClient{cc}
}

func (c *fileTransferServiceClient) Upload(ctx context.Context, opts ...grpc.CallOption) (FileTransferService_UploadClient, error) {
	stream, err := c.cc.NewStream(ctx, &FileTransferService_ServiceDesc.Streams[0], "/grpcp.FileTransferService/Upload", opts...)
	if err != nil {
		return nil, err
	}
	x := &fileTransferServiceUploadClient{stream}
	return x, nil
}

type FileTransferService_UploadClient interface {
	Send(*FileUploadRequest) error
	CloseAndRecv() (*FileUploadResponse, error)
	grpc.ClientStream
}

type fileTransferServiceUploadClient struct {
	grpc.ClientStream
}

func (x *fileTransferServiceUploadClient) Send(m *FileUploadRequest) error {
	return x.ClientStream.SendMsg(m)
}

func (x *fileTransferServiceUploadClient) CloseAndRecv() (*FileUploadResponse, error) {
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	m := new(FileUploadResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *fileTransferServiceClient) Download(ctx context.Context, in *FileDownloadRequest, opts ...grpc.CallOption) (FileTransferService_DownloadClient, error) {
	stream, err := c.cc.NewStream(ctx, &FileTransferService_ServiceDesc.Streams[1], "/grpcp.FileTransferService/Download", opts...)
	if err != nil {
		return nil, err
	}
	x := &fileTransferServiceDownloadClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type FileTransferService_DownloadClient interface {
	Recv() (*FileDownloadResponse, error)
	grpc.ClientStream
}

type fileTransferServiceDownloadClient struct {
	grpc.ClientStream
}

func (x *fileTransferServiceDownloadClient) Recv() (*FileDownloadResponse, error) {
	m := new(FileDownloadResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *fileTransferServiceClient) Ping(ctx context.Context, in *PingRequest, opts ...grpc.CallOption) (*PingResponse, error) {
	out := new(PingResponse)
	err := c.cc.Invoke(ctx, "/grpcp.FileTransferService/Ping", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// FileTransferServiceServer is the server API for FileTransferService service.
// All implementations must embed UnimplementedFileTransferServiceServer
// for forward compatibility
type FileTransferServiceServer interface {
	Upload(FileTransferService_UploadServer) error
	Download(*FileDownloadRequest, FileTransferService_DownloadServer) error
	Ping(context.Context, *PingRequest) (*PingResponse, error)
	mustEmbedUnimplementedFileTransferServiceServer()
}

// UnimplementedFileTransferServiceServer must be embedded to have forward compatible implementations.
type UnimplementedFileTransferServiceServer struct {
}

func (UnimplementedFileTransferServiceServer) Upload(FileTransferService_UploadServer) error {
	return status.Errorf(codes.Unimplemented, "method Upload not implemented")
}
func (UnimplementedFileTransferServiceServer) Download(*FileDownloadRequest, FileTransferService_DownloadServer) error {
	return status.Errorf(codes.Unimplemented, "method Download not implemented")
}
func (UnimplementedFileTransferServiceServer) Ping(context.Context, *PingRequest) (*PingResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Ping not implemented")
}
func (UnimplementedFileTransferServiceServer) mustEmbedUnimplementedFileTransferServiceServer() {}

// UnsafeFileTransferServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to FileTransferServiceServer will
// result in compilation errors.
type UnsafeFileTransferServiceServer interface {
	mustEmbedUnimplementedFileTransferServiceServer()
}

func RegisterFileTransferServiceServer(s grpc.ServiceRegistrar, srv FileTransferServiceServer) {
	s.RegisterService(&FileTransferService_ServiceDesc, srv)
}

func _FileTransferService_Upload_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(FileTransferServiceServer).Upload(&fileTransferServiceUploadServer{stream})
}

type FileTransferService_UploadServer interface {
	SendAndClose(*FileUploadResponse) error
	Recv() (*FileUploadRequest, error)
	grpc.ServerStream
}

type fileTransferServiceUploadServer struct {
	grpc.ServerStream
}

func (x *fileTransferServiceUploadServer) SendAndClose(m *FileUploadResponse) error {
	return x.ServerStream.SendMsg(m)
}

func (x *fileTransferServiceUploadServer) Recv() (*FileUploadRequest, error) {
	m := new(FileUploadRequest)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _FileTransferService_Download_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(FileDownloadRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(FileTransferServiceServer).Download(m, &fileTransferServiceDownloadServer{stream})
}

type FileTransferService_DownloadServer interface {
	Send(*FileDownloadResponse) error
	grpc.ServerStream
}

type fileTransferServiceDownloadServer struct {
	grpc.ServerStream
}

func (x *fileTransferServiceDownloadServer) Send(m *FileDownloadResponse) error {
	return x.ServerStream.SendMsg(m)
}

func _FileTransferService_Ping_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PingRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FileTransferServiceServer).Ping(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/grpcp.FileTransferService/Ping",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FileTransferServiceServer).Ping(ctx, req.(*PingRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// FileTransferService_ServiceDesc is the grpc.ServiceDesc for FileTransferService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var FileTransferService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "grpcp.FileTransferService",
	HandlerType: (*FileTransferServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Ping",
			Handler:    _FileTransferService_Ping_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Upload",
			Handler:       _FileTransferService_Upload_Handler,
			ClientStreams: true,
		},
		{
			StreamName:    "Download",
			Handler:       _FileTransferService_Download_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "filetransfer.proto",
}
