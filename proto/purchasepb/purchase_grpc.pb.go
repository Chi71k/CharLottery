// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v6.30.2
// source: purchase.proto

package purchasepb

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
	PurchaseService_BuyTicket_FullMethodName         = "/purchasepb.PurchaseService/BuyTicket"
	PurchaseService_ListTicketsByUser_FullMethodName = "/purchasepb.PurchaseService/ListTicketsByUser"
	PurchaseService_UpdatePurchase_FullMethodName    = "/purchasepb.PurchaseService/UpdatePurchase"
	PurchaseService_DeletePurchase_FullMethodName    = "/purchasepb.PurchaseService/DeletePurchase"
)

// PurchaseServiceClient is the client API for PurchaseService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type PurchaseServiceClient interface {
	BuyTicket(ctx context.Context, in *BuyTicketRequest, opts ...grpc.CallOption) (*BuyTicketResponse, error)
	ListTicketsByUser(ctx context.Context, in *ListTicketsByUserRequest, opts ...grpc.CallOption) (*ListTicketsByUserResponse, error)
	UpdatePurchase(ctx context.Context, in *UpdatePurchaseRequest, opts ...grpc.CallOption) (*UpdatePurchaseResponse, error)
	DeletePurchase(ctx context.Context, in *DeletePurchaseRequest, opts ...grpc.CallOption) (*DeletePurchaseResponse, error)
}

type purchaseServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewPurchaseServiceClient(cc grpc.ClientConnInterface) PurchaseServiceClient {
	return &purchaseServiceClient{cc}
}

func (c *purchaseServiceClient) BuyTicket(ctx context.Context, in *BuyTicketRequest, opts ...grpc.CallOption) (*BuyTicketResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(BuyTicketResponse)
	err := c.cc.Invoke(ctx, PurchaseService_BuyTicket_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *purchaseServiceClient) ListTicketsByUser(ctx context.Context, in *ListTicketsByUserRequest, opts ...grpc.CallOption) (*ListTicketsByUserResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ListTicketsByUserResponse)
	err := c.cc.Invoke(ctx, PurchaseService_ListTicketsByUser_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *purchaseServiceClient) UpdatePurchase(ctx context.Context, in *UpdatePurchaseRequest, opts ...grpc.CallOption) (*UpdatePurchaseResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(UpdatePurchaseResponse)
	err := c.cc.Invoke(ctx, PurchaseService_UpdatePurchase_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *purchaseServiceClient) DeletePurchase(ctx context.Context, in *DeletePurchaseRequest, opts ...grpc.CallOption) (*DeletePurchaseResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(DeletePurchaseResponse)
	err := c.cc.Invoke(ctx, PurchaseService_DeletePurchase_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// PurchaseServiceServer is the server API for PurchaseService service.
// All implementations must embed UnimplementedPurchaseServiceServer
// for forward compatibility.
type PurchaseServiceServer interface {
	BuyTicket(context.Context, *BuyTicketRequest) (*BuyTicketResponse, error)
	ListTicketsByUser(context.Context, *ListTicketsByUserRequest) (*ListTicketsByUserResponse, error)
	UpdatePurchase(context.Context, *UpdatePurchaseRequest) (*UpdatePurchaseResponse, error)
	DeletePurchase(context.Context, *DeletePurchaseRequest) (*DeletePurchaseResponse, error)
	mustEmbedUnimplementedPurchaseServiceServer()
}

// UnimplementedPurchaseServiceServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedPurchaseServiceServer struct{}

func (UnimplementedPurchaseServiceServer) BuyTicket(context.Context, *BuyTicketRequest) (*BuyTicketResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method BuyTicket not implemented")
}
func (UnimplementedPurchaseServiceServer) ListTicketsByUser(context.Context, *ListTicketsByUserRequest) (*ListTicketsByUserResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListTicketsByUser not implemented")
}
func (UnimplementedPurchaseServiceServer) UpdatePurchase(context.Context, *UpdatePurchaseRequest) (*UpdatePurchaseResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdatePurchase not implemented")
}
func (UnimplementedPurchaseServiceServer) DeletePurchase(context.Context, *DeletePurchaseRequest) (*DeletePurchaseResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeletePurchase not implemented")
}
func (UnimplementedPurchaseServiceServer) mustEmbedUnimplementedPurchaseServiceServer() {}
func (UnimplementedPurchaseServiceServer) testEmbeddedByValue()                         {}

// UnsafePurchaseServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to PurchaseServiceServer will
// result in compilation errors.
type UnsafePurchaseServiceServer interface {
	mustEmbedUnimplementedPurchaseServiceServer()
}

func RegisterPurchaseServiceServer(s grpc.ServiceRegistrar, srv PurchaseServiceServer) {
	// If the following call pancis, it indicates UnimplementedPurchaseServiceServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&PurchaseService_ServiceDesc, srv)
}

func _PurchaseService_BuyTicket_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(BuyTicketRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PurchaseServiceServer).BuyTicket(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: PurchaseService_BuyTicket_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PurchaseServiceServer).BuyTicket(ctx, req.(*BuyTicketRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _PurchaseService_ListTicketsByUser_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListTicketsByUserRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PurchaseServiceServer).ListTicketsByUser(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: PurchaseService_ListTicketsByUser_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PurchaseServiceServer).ListTicketsByUser(ctx, req.(*ListTicketsByUserRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _PurchaseService_UpdatePurchase_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdatePurchaseRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PurchaseServiceServer).UpdatePurchase(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: PurchaseService_UpdatePurchase_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PurchaseServiceServer).UpdatePurchase(ctx, req.(*UpdatePurchaseRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _PurchaseService_DeletePurchase_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeletePurchaseRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PurchaseServiceServer).DeletePurchase(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: PurchaseService_DeletePurchase_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PurchaseServiceServer).DeletePurchase(ctx, req.(*DeletePurchaseRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// PurchaseService_ServiceDesc is the grpc.ServiceDesc for PurchaseService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var PurchaseService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "purchasepb.PurchaseService",
	HandlerType: (*PurchaseServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "BuyTicket",
			Handler:    _PurchaseService_BuyTicket_Handler,
		},
		{
			MethodName: "ListTicketsByUser",
			Handler:    _PurchaseService_ListTicketsByUser_Handler,
		},
		{
			MethodName: "UpdatePurchase",
			Handler:    _PurchaseService_UpdatePurchase_Handler,
		},
		{
			MethodName: "DeletePurchase",
			Handler:    _PurchaseService_DeletePurchase_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "purchase.proto",
}
