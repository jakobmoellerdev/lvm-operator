// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.21.2
// source: lvmd/proto/lvmd.proto

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

// LVServiceClient is the client API for LVService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type LVServiceClient interface {
	// Create a logical volume.
	CreateLV(ctx context.Context, in *CreateLVRequest, opts ...grpc.CallOption) (*CreateLVResponse, error)
	// Remove a logical volume.
	RemoveLV(ctx context.Context, in *RemoveLVRequest, opts ...grpc.CallOption) (*Empty, error)
	// Resize a logical volume.
	ResizeLV(ctx context.Context, in *ResizeLVRequest, opts ...grpc.CallOption) (*Empty, error)
	CreateLVSnapshot(ctx context.Context, in *CreateLVSnapshotRequest, opts ...grpc.CallOption) (*CreateLVSnapshotResponse, error)
}

type lVServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewLVServiceClient(cc grpc.ClientConnInterface) LVServiceClient {
	return &lVServiceClient{cc}
}

func (c *lVServiceClient) CreateLV(ctx context.Context, in *CreateLVRequest, opts ...grpc.CallOption) (*CreateLVResponse, error) {
	out := new(CreateLVResponse)
	err := c.cc.Invoke(ctx, "/proto.LVService/CreateLV", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *lVServiceClient) RemoveLV(ctx context.Context, in *RemoveLVRequest, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := c.cc.Invoke(ctx, "/proto.LVService/RemoveLV", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *lVServiceClient) ResizeLV(ctx context.Context, in *ResizeLVRequest, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := c.cc.Invoke(ctx, "/proto.LVService/ResizeLV", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *lVServiceClient) CreateLVSnapshot(ctx context.Context, in *CreateLVSnapshotRequest, opts ...grpc.CallOption) (*CreateLVSnapshotResponse, error) {
	out := new(CreateLVSnapshotResponse)
	err := c.cc.Invoke(ctx, "/proto.LVService/CreateLVSnapshot", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// LVServiceServer is the server API for LVService service.
// All implementations must embed UnimplementedLVServiceServer
// for forward compatibility
type LVServiceServer interface {
	// Create a logical volume.
	CreateLV(context.Context, *CreateLVRequest) (*CreateLVResponse, error)
	// Remove a logical volume.
	RemoveLV(context.Context, *RemoveLVRequest) (*Empty, error)
	// Resize a logical volume.
	ResizeLV(context.Context, *ResizeLVRequest) (*Empty, error)
	CreateLVSnapshot(context.Context, *CreateLVSnapshotRequest) (*CreateLVSnapshotResponse, error)
	mustEmbedUnimplementedLVServiceServer()
}

// UnimplementedLVServiceServer must be embedded to have forward compatible implementations.
type UnimplementedLVServiceServer struct {
}

func (UnimplementedLVServiceServer) CreateLV(context.Context, *CreateLVRequest) (*CreateLVResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateLV not implemented")
}
func (UnimplementedLVServiceServer) RemoveLV(context.Context, *RemoveLVRequest) (*Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RemoveLV not implemented")
}
func (UnimplementedLVServiceServer) ResizeLV(context.Context, *ResizeLVRequest) (*Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ResizeLV not implemented")
}
func (UnimplementedLVServiceServer) CreateLVSnapshot(context.Context, *CreateLVSnapshotRequest) (*CreateLVSnapshotResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateLVSnapshot not implemented")
}
func (UnimplementedLVServiceServer) mustEmbedUnimplementedLVServiceServer() {}

// UnsafeLVServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to LVServiceServer will
// result in compilation errors.
type UnsafeLVServiceServer interface {
	mustEmbedUnimplementedLVServiceServer()
}

func RegisterLVServiceServer(s grpc.ServiceRegistrar, srv LVServiceServer) {
	s.RegisterService(&LVService_ServiceDesc, srv)
}

func _LVService_CreateLV_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateLVRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LVServiceServer).CreateLV(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.LVService/CreateLV",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LVServiceServer).CreateLV(ctx, req.(*CreateLVRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _LVService_RemoveLV_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RemoveLVRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LVServiceServer).RemoveLV(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.LVService/RemoveLV",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LVServiceServer).RemoveLV(ctx, req.(*RemoveLVRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _LVService_ResizeLV_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ResizeLVRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LVServiceServer).ResizeLV(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.LVService/ResizeLV",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LVServiceServer).ResizeLV(ctx, req.(*ResizeLVRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _LVService_CreateLVSnapshot_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateLVSnapshotRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LVServiceServer).CreateLVSnapshot(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.LVService/CreateLVSnapshot",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LVServiceServer).CreateLVSnapshot(ctx, req.(*CreateLVSnapshotRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// LVService_ServiceDesc is the grpc.ServiceDesc for LVService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var LVService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "proto.LVService",
	HandlerType: (*LVServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateLV",
			Handler:    _LVService_CreateLV_Handler,
		},
		{
			MethodName: "RemoveLV",
			Handler:    _LVService_RemoveLV_Handler,
		},
		{
			MethodName: "ResizeLV",
			Handler:    _LVService_ResizeLV_Handler,
		},
		{
			MethodName: "CreateLVSnapshot",
			Handler:    _LVService_CreateLVSnapshot_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "lvmd/proto/lvmd.proto",
}

// VGServiceClient is the client API for VGService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type VGServiceClient interface {
	// Get the list of logical volumes in the volume group.
	GetLVList(ctx context.Context, in *GetLVListRequest, opts ...grpc.CallOption) (*GetLVListResponse, error)
	// Get the free space of the volume group in bytes.
	GetFreeBytes(ctx context.Context, in *GetFreeBytesRequest, opts ...grpc.CallOption) (*GetFreeBytesResponse, error)
	// Stream the volume group metrics.
	Watch(ctx context.Context, in *Empty, opts ...grpc.CallOption) (VGService_WatchClient, error)
}

type vGServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewVGServiceClient(cc grpc.ClientConnInterface) VGServiceClient {
	return &vGServiceClient{cc}
}

func (c *vGServiceClient) GetLVList(ctx context.Context, in *GetLVListRequest, opts ...grpc.CallOption) (*GetLVListResponse, error) {
	out := new(GetLVListResponse)
	err := c.cc.Invoke(ctx, "/proto.VGService/GetLVList", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *vGServiceClient) GetFreeBytes(ctx context.Context, in *GetFreeBytesRequest, opts ...grpc.CallOption) (*GetFreeBytesResponse, error) {
	out := new(GetFreeBytesResponse)
	err := c.cc.Invoke(ctx, "/proto.VGService/GetFreeBytes", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *vGServiceClient) Watch(ctx context.Context, in *Empty, opts ...grpc.CallOption) (VGService_WatchClient, error) {
	stream, err := c.cc.NewStream(ctx, &VGService_ServiceDesc.Streams[0], "/proto.VGService/Watch", opts...)
	if err != nil {
		return nil, err
	}
	x := &vGServiceWatchClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type VGService_WatchClient interface {
	Recv() (*WatchResponse, error)
	grpc.ClientStream
}

type vGServiceWatchClient struct {
	grpc.ClientStream
}

func (x *vGServiceWatchClient) Recv() (*WatchResponse, error) {
	m := new(WatchResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// VGServiceServer is the server API for VGService service.
// All implementations must embed UnimplementedVGServiceServer
// for forward compatibility
type VGServiceServer interface {
	// Get the list of logical volumes in the volume group.
	GetLVList(context.Context, *GetLVListRequest) (*GetLVListResponse, error)
	// Get the free space of the volume group in bytes.
	GetFreeBytes(context.Context, *GetFreeBytesRequest) (*GetFreeBytesResponse, error)
	// Stream the volume group metrics.
	Watch(*Empty, VGService_WatchServer) error
	mustEmbedUnimplementedVGServiceServer()
}

// UnimplementedVGServiceServer must be embedded to have forward compatible implementations.
type UnimplementedVGServiceServer struct {
}

func (UnimplementedVGServiceServer) GetLVList(context.Context, *GetLVListRequest) (*GetLVListResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetLVList not implemented")
}
func (UnimplementedVGServiceServer) GetFreeBytes(context.Context, *GetFreeBytesRequest) (*GetFreeBytesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetFreeBytes not implemented")
}
func (UnimplementedVGServiceServer) Watch(*Empty, VGService_WatchServer) error {
	return status.Errorf(codes.Unimplemented, "method Watch not implemented")
}
func (UnimplementedVGServiceServer) mustEmbedUnimplementedVGServiceServer() {}

// UnsafeVGServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to VGServiceServer will
// result in compilation errors.
type UnsafeVGServiceServer interface {
	mustEmbedUnimplementedVGServiceServer()
}

func RegisterVGServiceServer(s grpc.ServiceRegistrar, srv VGServiceServer) {
	s.RegisterService(&VGService_ServiceDesc, srv)
}

func _VGService_GetLVList_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetLVListRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(VGServiceServer).GetLVList(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.VGService/GetLVList",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(VGServiceServer).GetLVList(ctx, req.(*GetLVListRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _VGService_GetFreeBytes_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetFreeBytesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(VGServiceServer).GetFreeBytes(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.VGService/GetFreeBytes",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(VGServiceServer).GetFreeBytes(ctx, req.(*GetFreeBytesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _VGService_Watch_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(Empty)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(VGServiceServer).Watch(m, &vGServiceWatchServer{stream})
}

type VGService_WatchServer interface {
	Send(*WatchResponse) error
	grpc.ServerStream
}

type vGServiceWatchServer struct {
	grpc.ServerStream
}

func (x *vGServiceWatchServer) Send(m *WatchResponse) error {
	return x.ServerStream.SendMsg(m)
}

// VGService_ServiceDesc is the grpc.ServiceDesc for VGService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var VGService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "proto.VGService",
	HandlerType: (*VGServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetLVList",
			Handler:    _VGService_GetLVList_Handler,
		},
		{
			MethodName: "GetFreeBytes",
			Handler:    _VGService_GetFreeBytes_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Watch",
			Handler:       _VGService_Watch_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "lvmd/proto/lvmd.proto",
}
