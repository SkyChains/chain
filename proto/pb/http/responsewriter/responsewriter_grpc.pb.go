// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             (unknown)
// source: http/responsewriter/responsewriter.proto

package responsewriter

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

const (
	Writer_Write_FullMethodName       = "/http.responsewriter.Writer/Write"
	Writer_WriteHeader_FullMethodName = "/http.responsewriter.Writer/WriteHeader"
	Writer_Flush_FullMethodName       = "/http.responsewriter.Writer/Flush"
	Writer_Hijack_FullMethodName      = "/http.responsewriter.Writer/Hijack"
)

// WriterClient is the client API for Writer service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type WriterClient interface {
	// Write writes the data to the connection as part of an HTTP reply
	Write(ctx context.Context, in *WriteRequest, opts ...grpc.CallOption) (*WriteResponse, error)
	// WriteHeader sends an HTTP response header with the provided
	// status code
	WriteHeader(ctx context.Context, in *WriteHeaderRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
	// Flush is a no-op
	Flush(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*emptypb.Empty, error)
	// Hijack lets the caller take over the connection
	Hijack(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*HijackResponse, error)
}

type writerClient struct {
	cc grpc.ClientConnInterface
}

func NewWriterClient(cc grpc.ClientConnInterface) WriterClient {
	return &writerClient{cc}
}

func (c *writerClient) Write(ctx context.Context, in *WriteRequest, opts ...grpc.CallOption) (*WriteResponse, error) {
	out := new(WriteResponse)
	err := c.cc.Invoke(ctx, Writer_Write_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *writerClient) WriteHeader(ctx context.Context, in *WriteHeaderRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, Writer_WriteHeader_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *writerClient) Flush(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, Writer_Flush_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *writerClient) Hijack(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*HijackResponse, error) {
	out := new(HijackResponse)
	err := c.cc.Invoke(ctx, Writer_Hijack_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// WriterServer is the server API for Writer service.
// All implementations must embed UnimplementedWriterServer
// for forward compatibility
type WriterServer interface {
	// Write writes the data to the connection as part of an HTTP reply
	Write(context.Context, *WriteRequest) (*WriteResponse, error)
	// WriteHeader sends an HTTP response header with the provided
	// status code
	WriteHeader(context.Context, *WriteHeaderRequest) (*emptypb.Empty, error)
	// Flush is a no-op
	Flush(context.Context, *emptypb.Empty) (*emptypb.Empty, error)
	// Hijack lets the caller take over the connection
	Hijack(context.Context, *emptypb.Empty) (*HijackResponse, error)
	mustEmbedUnimplementedWriterServer()
}

// UnimplementedWriterServer must be embedded to have forward compatible implementations.
type UnimplementedWriterServer struct {
}

func (UnimplementedWriterServer) Write(context.Context, *WriteRequest) (*WriteResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Write not implemented")
}
func (UnimplementedWriterServer) WriteHeader(context.Context, *WriteHeaderRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method WriteHeader not implemented")
}
func (UnimplementedWriterServer) Flush(context.Context, *emptypb.Empty) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Flush not implemented")
}
func (UnimplementedWriterServer) Hijack(context.Context, *emptypb.Empty) (*HijackResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Hijack not implemented")
}
func (UnimplementedWriterServer) mustEmbedUnimplementedWriterServer() {}

// UnsafeWriterServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to WriterServer will
// result in compilation errors.
type UnsafeWriterServer interface {
	mustEmbedUnimplementedWriterServer()
}

func RegisterWriterServer(s grpc.ServiceRegistrar, srv WriterServer) {
	s.RegisterService(&Writer_ServiceDesc, srv)
}

func _Writer_Write_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(WriteRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(WriterServer).Write(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Writer_Write_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(WriterServer).Write(ctx, req.(*WriteRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Writer_WriteHeader_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(WriteHeaderRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(WriterServer).WriteHeader(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Writer_WriteHeader_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(WriterServer).WriteHeader(ctx, req.(*WriteHeaderRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Writer_Flush_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(WriterServer).Flush(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Writer_Flush_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(WriterServer).Flush(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _Writer_Hijack_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(WriterServer).Hijack(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Writer_Hijack_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(WriterServer).Hijack(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

// Writer_ServiceDesc is the grpc.ServiceDesc for Writer service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Writer_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "http.responsewriter.Writer",
	HandlerType: (*WriterServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Write",
			Handler:    _Writer_Write_Handler,
		},
		{
			MethodName: "WriteHeader",
			Handler:    _Writer_WriteHeader_Handler,
		},
		{
			MethodName: "Flush",
			Handler:    _Writer_Flush_Handler,
		},
		{
			MethodName: "Hijack",
			Handler:    _Writer_Hijack_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "http/responsewriter/responsewriter.proto",
}
