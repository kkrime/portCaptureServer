// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v3.6.1
// source: portCaptureServer.proto

package pb

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

const (
	PortCaptureService_SavePorts_FullMethodName = "/pb.PortCaptureService/SavePorts"
)

// PortCaptureServiceClient is the client API for PortCaptureService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type PortCaptureServiceClient interface {
	SavePorts(ctx context.Context, opts ...grpc.CallOption) (PortCaptureService_SavePortsClient, error)
}

type portCaptureServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewPortCaptureServiceClient(cc grpc.ClientConnInterface) PortCaptureServiceClient {
	return &portCaptureServiceClient{cc}
}

func (c *portCaptureServiceClient) SavePorts(ctx context.Context, opts ...grpc.CallOption) (PortCaptureService_SavePortsClient, error) {
	stream, err := c.cc.NewStream(ctx, &PortCaptureService_ServiceDesc.Streams[0], PortCaptureService_SavePorts_FullMethodName, opts...)
	if err != nil {
		return nil, err
	}
	x := &portCaptureServiceSavePortsClient{stream}
	return x, nil
}

type PortCaptureService_SavePortsClient interface {
	Send(*Port) error
	CloseAndRecv() (*PortCaptureServiceResponse, error)
	grpc.ClientStream
}

type portCaptureServiceSavePortsClient struct {
	grpc.ClientStream
}

func (x *portCaptureServiceSavePortsClient) Send(m *Port) error {
	return x.ClientStream.SendMsg(m)
}

func (x *portCaptureServiceSavePortsClient) CloseAndRecv() (*PortCaptureServiceResponse, error) {
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	m := new(PortCaptureServiceResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// PortCaptureServiceServer is the server API for PortCaptureService service.
// All implementations must embed UnimplementedPortCaptureServiceServer
// for forward compatibility
type PortCaptureServiceServer interface {
	SavePorts(PortCaptureService_SavePortsServer) error
	mustEmbedUnimplementedPortCaptureServiceServer()
}

// UnimplementedPortCaptureServiceServer must be embedded to have forward compatible implementations.
type UnimplementedPortCaptureServiceServer struct {
}

func (UnimplementedPortCaptureServiceServer) SavePorts(PortCaptureService_SavePortsServer) error {
	return status.Errorf(codes.Unimplemented, "method SavePorts not implemented")
}
func (UnimplementedPortCaptureServiceServer) mustEmbedUnimplementedPortCaptureServiceServer() {}

// UnsafePortCaptureServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to PortCaptureServiceServer will
// result in compilation errors.
type UnsafePortCaptureServiceServer interface {
	mustEmbedUnimplementedPortCaptureServiceServer()
}

func RegisterPortCaptureServiceServer(s grpc.ServiceRegistrar, srv PortCaptureServiceServer) {
	s.RegisterService(&PortCaptureService_ServiceDesc, srv)
}

func _PortCaptureService_SavePorts_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(PortCaptureServiceServer).SavePorts(&portCaptureServiceSavePortsServer{stream})
}

type PortCaptureService_SavePortsServer interface {
	SendAndClose(*PortCaptureServiceResponse) error
	Recv() (*Port, error)
	grpc.ServerStream
}

type portCaptureServiceSavePortsServer struct {
	grpc.ServerStream
}

func (x *portCaptureServiceSavePortsServer) SendAndClose(m *PortCaptureServiceResponse) error {
	return x.ServerStream.SendMsg(m)
}

func (x *portCaptureServiceSavePortsServer) Recv() (*Port, error) {
	m := new(Port)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// PortCaptureService_ServiceDesc is the grpc.ServiceDesc for PortCaptureService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var PortCaptureService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "pb.PortCaptureService",
	HandlerType: (*PortCaptureServiceServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "SavePorts",
			Handler:       _PortCaptureService_SavePorts_Handler,
			ClientStreams: true,
		},
	},
	Metadata: "portCaptureServer.proto",
}