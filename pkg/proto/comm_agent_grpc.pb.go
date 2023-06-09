// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v3.21.12
// source: pkg/proto/comm_agent.proto

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

const (
	Events_GetEvent_FullMethodName = "/Events/getEvent"
)

// EventsClient is the client API for Events service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type EventsClient interface {
	GetEvent(ctx context.Context, in *GetEventRequest, opts ...grpc.CallOption) (Events_GetEventClient, error)
}

type eventsClient struct {
	cc grpc.ClientConnInterface
}

func NewEventsClient(cc grpc.ClientConnInterface) EventsClient {
	return &eventsClient{cc}
}

func (c *eventsClient) GetEvent(ctx context.Context, in *GetEventRequest, opts ...grpc.CallOption) (Events_GetEventClient, error) {
	stream, err := c.cc.NewStream(ctx, &Events_ServiceDesc.Streams[0], Events_GetEvent_FullMethodName, opts...)
	if err != nil {
		return nil, err
	}
	x := &eventsGetEventClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type Events_GetEventClient interface {
	Recv() (*Event, error)
	grpc.ClientStream
}

type eventsGetEventClient struct {
	grpc.ClientStream
}

func (x *eventsGetEventClient) Recv() (*Event, error) {
	m := new(Event)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// EventsServer is the server API for Events service.
// All implementations must embed UnimplementedEventsServer
// for forward compatibility
type EventsServer interface {
	GetEvent(*GetEventRequest, Events_GetEventServer) error
	mustEmbedUnimplementedEventsServer()
}

// UnimplementedEventsServer must be embedded to have forward compatible implementations.
type UnimplementedEventsServer struct {
}

func (UnimplementedEventsServer) GetEvent(*GetEventRequest, Events_GetEventServer) error {
	return status.Errorf(codes.Unimplemented, "method GetEvent not implemented")
}
func (UnimplementedEventsServer) mustEmbedUnimplementedEventsServer() {}

// UnsafeEventsServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to EventsServer will
// result in compilation errors.
type UnsafeEventsServer interface {
	mustEmbedUnimplementedEventsServer()
}

func RegisterEventsServer(s grpc.ServiceRegistrar, srv EventsServer) {
	s.RegisterService(&Events_ServiceDesc, srv)
}

func _Events_GetEvent_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(GetEventRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(EventsServer).GetEvent(m, &eventsGetEventServer{stream})
}

type Events_GetEventServer interface {
	Send(*Event) error
	grpc.ServerStream
}

type eventsGetEventServer struct {
	grpc.ServerStream
}

func (x *eventsGetEventServer) Send(m *Event) error {
	return x.ServerStream.SendMsg(m)
}

// Events_ServiceDesc is the grpc.ServiceDesc for Events service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Events_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "Events",
	HandlerType: (*EventsServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "getEvent",
			Handler:       _Events_GetEvent_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "pkg/proto/comm_agent.proto",
}
