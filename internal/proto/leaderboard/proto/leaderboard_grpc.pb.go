// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.21.12
// source: proto/leaderboard.proto

package leaderboard

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

// LeaderboardServiceClient is the client API for LeaderboardService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type LeaderboardServiceClient interface {
	GetLeaderboard(ctx context.Context, in *GetLeaderboardRequest, opts ...grpc.CallOption) (*GetLeaderboardResponse, error)
	SubmitUserScore(ctx context.Context, in *SubmitUserScoreRequest, opts ...grpc.CallOption) (*SubmitUserScoreResponse, error)
}

type leaderboardServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewLeaderboardServiceClient(cc grpc.ClientConnInterface) LeaderboardServiceClient {
	return &leaderboardServiceClient{cc}
}

func (c *leaderboardServiceClient) GetLeaderboard(ctx context.Context, in *GetLeaderboardRequest, opts ...grpc.CallOption) (*GetLeaderboardResponse, error) {
	out := new(GetLeaderboardResponse)
	err := c.cc.Invoke(ctx, "/leaderboard.LeaderboardService/GetLeaderboard", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *leaderboardServiceClient) SubmitUserScore(ctx context.Context, in *SubmitUserScoreRequest, opts ...grpc.CallOption) (*SubmitUserScoreResponse, error) {
	out := new(SubmitUserScoreResponse)
	err := c.cc.Invoke(ctx, "/leaderboard.LeaderboardService/SubmitUserScore", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// LeaderboardServiceServer is the server API for LeaderboardService service.
// All implementations must embed UnimplementedLeaderboardServiceServer
// for forward compatibility
type LeaderboardServiceServer interface {
	GetLeaderboard(context.Context, *GetLeaderboardRequest) (*GetLeaderboardResponse, error)
	SubmitUserScore(context.Context, *SubmitUserScoreRequest) (*SubmitUserScoreResponse, error)
	mustEmbedUnimplementedLeaderboardServiceServer()
}

// UnimplementedLeaderboardServiceServer must be embedded to have forward compatible implementations.
type UnimplementedLeaderboardServiceServer struct {
}

func (UnimplementedLeaderboardServiceServer) GetLeaderboard(context.Context, *GetLeaderboardRequest) (*GetLeaderboardResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetLeaderboard not implemented")
}
func (UnimplementedLeaderboardServiceServer) SubmitUserScore(context.Context, *SubmitUserScoreRequest) (*SubmitUserScoreResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SubmitUserScore not implemented")
}
func (UnimplementedLeaderboardServiceServer) mustEmbedUnimplementedLeaderboardServiceServer() {}

// UnsafeLeaderboardServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to LeaderboardServiceServer will
// result in compilation errors.
type UnsafeLeaderboardServiceServer interface {
	mustEmbedUnimplementedLeaderboardServiceServer()
}

func RegisterLeaderboardServiceServer(s grpc.ServiceRegistrar, srv LeaderboardServiceServer) {
	s.RegisterService(&LeaderboardService_ServiceDesc, srv)
}

func _LeaderboardService_GetLeaderboard_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetLeaderboardRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LeaderboardServiceServer).GetLeaderboard(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/leaderboard.LeaderboardService/GetLeaderboard",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LeaderboardServiceServer).GetLeaderboard(ctx, req.(*GetLeaderboardRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _LeaderboardService_SubmitUserScore_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SubmitUserScoreRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LeaderboardServiceServer).SubmitUserScore(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/leaderboard.LeaderboardService/SubmitUserScore",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LeaderboardServiceServer).SubmitUserScore(ctx, req.(*SubmitUserScoreRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// LeaderboardService_ServiceDesc is the grpc.ServiceDesc for LeaderboardService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var LeaderboardService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "leaderboard.LeaderboardService",
	HandlerType: (*LeaderboardServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetLeaderboard",
			Handler:    _LeaderboardService_GetLeaderboard_Handler,
		},
		{
			MethodName: "SubmitUserScore",
			Handler:    _LeaderboardService_SubmitUserScore_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/leaderboard.proto",
}
