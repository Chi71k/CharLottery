// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v6.30.2
// source: lottery.proto

package lotterypb

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.64.0 or later.
const _ = grpc.SupportPackageIsVersion9

const (
	LotteryService_CreateLottery_FullMethodName = "/lotterypb.LotteryService/CreateLottery"
	LotteryService_GetLottery_FullMethodName    = "/lotterypb.LotteryService/GetLottery"
	LotteryService_ListLotteries_FullMethodName = "/lotterypb.LotteryService/ListLotteries"
)

// LotteryServiceClient is the client API for LotteryService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type LotteryServiceClient interface {
	CreateLottery(ctx context.Context, in *CreateLotteryRequest, opts ...grpc.CallOption) (*CreateLotteryResponse, error)
	GetLottery(ctx context.Context, in *GetLotteryRequest, opts ...grpc.CallOption) (*GetLotteryResponse, error)
	ListLotteries(ctx context.Context, in *ListLotteriesRequest, opts ...grpc.CallOption) (*ListLotteriesResponse, error)
}

type lotteryServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewLotteryServiceClient(cc grpc.ClientConnInterface) LotteryServiceClient {
	return &lotteryServiceClient{cc}
}

func (c *lotteryServiceClient) CreateLottery(ctx context.Context, in *CreateLotteryRequest, opts ...grpc.CallOption) (*CreateLotteryResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(CreateLotteryResponse)
	err := c.cc.Invoke(ctx, LotteryService_CreateLottery_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *lotteryServiceClient) GetLottery(ctx context.Context, in *GetLotteryRequest, opts ...grpc.CallOption) (*GetLotteryResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetLotteryResponse)
	err := c.cc.Invoke(ctx, LotteryService_GetLottery_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *lotteryServiceClient) ListLotteries(ctx context.Context, in *ListLotteriesRequest, opts ...grpc.CallOption) (*ListLotteriesResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ListLotteriesResponse)
	err := c.cc.Invoke(ctx, LotteryService_ListLotteries_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// LotteryServiceServer is the server API for LotteryService service.
// All implementations must embed UnimplementedLotteryServiceServer
// for forward compatibility.
type LotteryServiceServer interface {
	CreateLottery(context.Context, *CreateLotteryRequest) (*CreateLotteryResponse, error)
	GetLottery(context.Context, *GetLotteryRequest) (*GetLotteryResponse, error)
	ListLotteries(context.Context, *ListLotteriesRequest) (*ListLotteriesResponse, error)
	mustEmbedUnimplementedLotteryServiceServer()
}

// UnimplementedLotteryServiceServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedLotteryServiceServer struct{}

func (UnimplementedLotteryServiceServer) CreateLottery(context.Context, *CreateLotteryRequest) (*CreateLotteryResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateLottery not implemented")
}
func (UnimplementedLotteryServiceServer) GetLottery(context.Context, *GetLotteryRequest) (*GetLotteryResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetLottery not implemented")
}
func (UnimplementedLotteryServiceServer) ListLotteries(context.Context, *ListLotteriesRequest) (*ListLotteriesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListLotteries not implemented")
}
func (UnimplementedLotteryServiceServer) mustEmbedUnimplementedLotteryServiceServer() {}
func (UnimplementedLotteryServiceServer) testEmbeddedByValue()                        {}

// UnsafeLotteryServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to LotteryServiceServer will
// result in compilation errors.
type UnsafeLotteryServiceServer interface {
	mustEmbedUnimplementedLotteryServiceServer()
}

func RegisterLotteryServiceServer(s grpc.ServiceRegistrar, srv LotteryServiceServer) {
	// If the following call pancis, it indicates UnimplementedLotteryServiceServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&LotteryService_ServiceDesc, srv)
}

func _LotteryService_CreateLottery_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateLotteryRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LotteryServiceServer).CreateLottery(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: LotteryService_CreateLottery_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LotteryServiceServer).CreateLottery(ctx, req.(*CreateLotteryRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _LotteryService_GetLottery_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetLotteryRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LotteryServiceServer).GetLottery(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: LotteryService_GetLottery_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LotteryServiceServer).GetLottery(ctx, req.(*GetLotteryRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _LotteryService_ListLotteries_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListLotteriesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LotteryServiceServer).ListLotteries(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: LotteryService_ListLotteries_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LotteryServiceServer).ListLotteries(ctx, req.(*ListLotteriesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// LotteryService_ServiceDesc is the grpc.ServiceDesc for LotteryService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var LotteryService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "lotterypb.LotteryService",
	HandlerType: (*LotteryServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateLottery",
			Handler:    _LotteryService_CreateLottery_Handler,
		},
		{
			MethodName: "GetLottery",
			Handler:    _LotteryService_GetLottery_Handler,
		},
		{
			MethodName: "ListLotteries",
			Handler:    _LotteryService_ListLotteries_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "lottery.proto",
}
