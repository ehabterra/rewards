// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v3.21.12
// source: reward_service.proto

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
	RewardsService_GetPoints_FullMethodName   = "/RewardsService/GetPoints"
	RewardsService_AddActivity_FullMethodName = "/RewardsService/AddActivity"
	RewardsService_SendPoints_FullMethodName  = "/RewardsService/SendPoints"
	RewardsService_SpendPoints_FullMethodName = "/RewardsService/SpendPoints"
)

// RewardsServiceClient is the client API for RewardsService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type RewardsServiceClient interface {
	GetPoints(ctx context.Context, in *GetPointsRequest, opts ...grpc.CallOption) (*GetPointsResponse, error)
	AddActivity(ctx context.Context, in *AddActivityRequest, opts ...grpc.CallOption) (*AddActivityResponse, error)
	SendPoints(ctx context.Context, in *SendPointsRequest, opts ...grpc.CallOption) (*SendPointsResponse, error)
	SpendPoints(ctx context.Context, in *SpendPointsRequest, opts ...grpc.CallOption) (*SpendPointsResponse, error)
}

type rewardsServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewRewardsServiceClient(cc grpc.ClientConnInterface) RewardsServiceClient {
	return &rewardsServiceClient{cc}
}

func (c *rewardsServiceClient) GetPoints(ctx context.Context, in *GetPointsRequest, opts ...grpc.CallOption) (*GetPointsResponse, error) {
	out := new(GetPointsResponse)
	err := c.cc.Invoke(ctx, RewardsService_GetPoints_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *rewardsServiceClient) AddActivity(ctx context.Context, in *AddActivityRequest, opts ...grpc.CallOption) (*AddActivityResponse, error) {
	out := new(AddActivityResponse)
	err := c.cc.Invoke(ctx, RewardsService_AddActivity_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *rewardsServiceClient) SendPoints(ctx context.Context, in *SendPointsRequest, opts ...grpc.CallOption) (*SendPointsResponse, error) {
	out := new(SendPointsResponse)
	err := c.cc.Invoke(ctx, RewardsService_SendPoints_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *rewardsServiceClient) SpendPoints(ctx context.Context, in *SpendPointsRequest, opts ...grpc.CallOption) (*SpendPointsResponse, error) {
	out := new(SpendPointsResponse)
	err := c.cc.Invoke(ctx, RewardsService_SpendPoints_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// RewardsServiceServer is the server API for RewardsService service.
// All implementations must embed UnimplementedRewardsServiceServer
// for forward compatibility
type RewardsServiceServer interface {
	GetPoints(context.Context, *GetPointsRequest) (*GetPointsResponse, error)
	AddActivity(context.Context, *AddActivityRequest) (*AddActivityResponse, error)
	SendPoints(context.Context, *SendPointsRequest) (*SendPointsResponse, error)
	SpendPoints(context.Context, *SpendPointsRequest) (*SpendPointsResponse, error)
	mustEmbedUnimplementedRewardsServiceServer()
}

// UnimplementedRewardsServiceServer must be embedded to have forward compatible implementations.
type UnimplementedRewardsServiceServer struct {
}

func (UnimplementedRewardsServiceServer) GetPoints(context.Context, *GetPointsRequest) (*GetPointsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetPoints not implemented")
}
func (UnimplementedRewardsServiceServer) AddActivity(context.Context, *AddActivityRequest) (*AddActivityResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddActivity not implemented")
}
func (UnimplementedRewardsServiceServer) SendPoints(context.Context, *SendPointsRequest) (*SendPointsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SendPoints not implemented")
}
func (UnimplementedRewardsServiceServer) SpendPoints(context.Context, *SpendPointsRequest) (*SpendPointsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SpendPoints not implemented")
}
func (UnimplementedRewardsServiceServer) mustEmbedUnimplementedRewardsServiceServer() {}

// UnsafeRewardsServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to RewardsServiceServer will
// result in compilation errors.
type UnsafeRewardsServiceServer interface {
	mustEmbedUnimplementedRewardsServiceServer()
}

func RegisterRewardsServiceServer(s grpc.ServiceRegistrar, srv RewardsServiceServer) {
	s.RegisterService(&RewardsService_ServiceDesc, srv)
}

func _RewardsService_GetPoints_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetPointsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RewardsServiceServer).GetPoints(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: RewardsService_GetPoints_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RewardsServiceServer).GetPoints(ctx, req.(*GetPointsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _RewardsService_AddActivity_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AddActivityRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RewardsServiceServer).AddActivity(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: RewardsService_AddActivity_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RewardsServiceServer).AddActivity(ctx, req.(*AddActivityRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _RewardsService_SendPoints_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SendPointsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RewardsServiceServer).SendPoints(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: RewardsService_SendPoints_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RewardsServiceServer).SendPoints(ctx, req.(*SendPointsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _RewardsService_SpendPoints_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SpendPointsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RewardsServiceServer).SpendPoints(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: RewardsService_SpendPoints_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RewardsServiceServer).SpendPoints(ctx, req.(*SpendPointsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// RewardsService_ServiceDesc is the grpc.ServiceDesc for RewardsService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var RewardsService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "RewardsService",
	HandlerType: (*RewardsServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetPoints",
			Handler:    _RewardsService_GetPoints_Handler,
		},
		{
			MethodName: "AddActivity",
			Handler:    _RewardsService_AddActivity_Handler,
		},
		{
			MethodName: "SendPoints",
			Handler:    _RewardsService_SendPoints_Handler,
		},
		{
			MethodName: "SpendPoints",
			Handler:    _RewardsService_SpendPoints_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "reward_service.proto",
}
