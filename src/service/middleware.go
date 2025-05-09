package service

import (
	"context"
	"fmt"
	"io"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	pb "google.golang.org/grpc/examples/features/proto/echo"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

var errMissingMetadata = status.Errorf(codes.InvalidArgument, "no incoming metadata in rpc context")

func (s *usersService) BidirectionalStreamingEcho(stream pb.Echo_BidirectionalStreamingEchoServer) error {
	fmt.Printf("--- BidirectionalStreamingEcho ---\n")

	md, ok := metadata.FromIncomingContext(stream.Context())
	if !ok {
		return status.Errorf(codes.Internal, "BidirectionalStreamingEcho: missing incoming metadata in rpc context")
	}

	// Read and print metadata added by the interceptor.
	if v, ok := md["key1"]; ok {
		fmt.Printf("key1 from metadata: \n")
		for i, e := range v {
			fmt.Printf(" %d. %s\n", i, e)
		}
	}

	// Read requests and send responses.
	for {
		in, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		if err = stream.Send(&pb.EchoResponse{Message: in.Message}); err != nil {
			return err
		}
	}
}

func UnaryInterceptor(ctx context.Context, req any, _ *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, errMissingMetadata
	}

	// Create and set metadata from interceptor to server.
	md.Append("key1", "value1")
	ctx = metadata.NewIncomingContext(ctx, md)

	// Call the handler to complete the normal execution of the RPC.
	resp, err := handler(ctx, req)

	// Create and set header metadata from interceptor to client.
	header := metadata.Pairs("header-key", "val")
	grpc.SetHeader(ctx, header)

	// Create and set trailer metadata from interceptor to client.
	trailer := metadata.Pairs("trailer-key", "val")
	grpc.SetTrailer(ctx, trailer)

	return resp, err
}
