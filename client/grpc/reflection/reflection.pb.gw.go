// Code generated by protoc-gen-grpc-gateway. DO NOT EDIT.
// source: lbm/base/reflection/v1/reflection.proto

/*
Package reflection is a reverse proxy.

It translates gRPC into RESTful JSON APIs.
*/
package reflection

import (
	"context"
	"io"
	"net/http"

	"github.com/golang/protobuf/descriptor"
	"github.com/golang/protobuf/proto"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/grpc-ecosystem/grpc-gateway/utilities"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/status"
)

// Suppress "imported and not used" errors
var _ codes.Code
var _ io.Reader
var _ status.Status
var _ = runtime.String
var _ = utilities.NewDoubleArray
var _ = descriptor.ForMessage

func request_ReflectionService_ListAllInterfaces_0(ctx context.Context, marshaler runtime.Marshaler, client ReflectionServiceClient, req *http.Request, pathParams map[string]string) (proto.Message, runtime.ServerMetadata, error) {
	var protoReq ListAllInterfacesRequest
	var metadata runtime.ServerMetadata

	msg, err := client.ListAllInterfaces(ctx, &protoReq, grpc.Header(&metadata.HeaderMD), grpc.Trailer(&metadata.TrailerMD))
	return msg, metadata, err

}

func local_request_ReflectionService_ListAllInterfaces_0(ctx context.Context, marshaler runtime.Marshaler, server ReflectionServiceServer, req *http.Request, pathParams map[string]string) (proto.Message, runtime.ServerMetadata, error) {
	var protoReq ListAllInterfacesRequest
	var metadata runtime.ServerMetadata

	msg, err := server.ListAllInterfaces(ctx, &protoReq)
	return msg, metadata, err

}

func request_ReflectionService_ListImplementations_0(ctx context.Context, marshaler runtime.Marshaler, client ReflectionServiceClient, req *http.Request, pathParams map[string]string) (proto.Message, runtime.ServerMetadata, error) {
	var protoReq ListImplementationsRequest
	var metadata runtime.ServerMetadata

	var (
		val string
		ok  bool
		err error
		_   = err
	)

	val, ok = pathParams["interface_name"]
	if !ok {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "missing parameter %s", "interface_name")
	}

	protoReq.InterfaceName, err = runtime.String(val)

	if err != nil {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "type mismatch, parameter: %s, error: %v", "interface_name", err)
	}

	msg, err := client.ListImplementations(ctx, &protoReq, grpc.Header(&metadata.HeaderMD), grpc.Trailer(&metadata.TrailerMD))
	return msg, metadata, err

}

func local_request_ReflectionService_ListImplementations_0(ctx context.Context, marshaler runtime.Marshaler, server ReflectionServiceServer, req *http.Request, pathParams map[string]string) (proto.Message, runtime.ServerMetadata, error) {
	var protoReq ListImplementationsRequest
	var metadata runtime.ServerMetadata

	var (
		val string
		ok  bool
		err error
		_   = err
	)

	val, ok = pathParams["interface_name"]
	if !ok {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "missing parameter %s", "interface_name")
	}

	protoReq.InterfaceName, err = runtime.String(val)

	if err != nil {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "type mismatch, parameter: %s, error: %v", "interface_name", err)
	}

	msg, err := server.ListImplementations(ctx, &protoReq)
	return msg, metadata, err

}

// RegisterReflectionServiceHandlerServer registers the http handlers for service ReflectionService to "mux".
// UnaryRPC     :call ReflectionServiceServer directly.
// StreamingRPC :currently unsupported pending https://github.com/grpc/grpc-go/issues/906.
// Note that using this registration option will cause many gRPC library features (such as grpc.SendHeader, etc) to stop working. Consider using RegisterReflectionServiceHandlerFromEndpoint instead.
func RegisterReflectionServiceHandlerServer(ctx context.Context, mux *runtime.ServeMux, server ReflectionServiceServer) error {

	mux.Handle("GET", pattern_ReflectionService_ListAllInterfaces_0, func(w http.ResponseWriter, req *http.Request, pathParams map[string]string) {
		ctx, cancel := context.WithCancel(req.Context())
		defer cancel()
		inboundMarshaler, outboundMarshaler := runtime.MarshalerForRequest(mux, req)
		rctx, err := runtime.AnnotateIncomingContext(ctx, mux, req)
		if err != nil {
			runtime.HTTPError(ctx, mux, outboundMarshaler, w, req, err)
			return
		}
		resp, md, err := local_request_ReflectionService_ListAllInterfaces_0(rctx, inboundMarshaler, server, req, pathParams)
		ctx = runtime.NewServerMetadataContext(ctx, md)
		if err != nil {
			runtime.HTTPError(ctx, mux, outboundMarshaler, w, req, err)
			return
		}

		forward_ReflectionService_ListAllInterfaces_0(ctx, mux, outboundMarshaler, w, req, resp, mux.GetForwardResponseOptions()...)

	})

	mux.Handle("GET", pattern_ReflectionService_ListImplementations_0, func(w http.ResponseWriter, req *http.Request, pathParams map[string]string) {
		ctx, cancel := context.WithCancel(req.Context())
		defer cancel()
		inboundMarshaler, outboundMarshaler := runtime.MarshalerForRequest(mux, req)
		rctx, err := runtime.AnnotateIncomingContext(ctx, mux, req)
		if err != nil {
			runtime.HTTPError(ctx, mux, outboundMarshaler, w, req, err)
			return
		}
		resp, md, err := local_request_ReflectionService_ListImplementations_0(rctx, inboundMarshaler, server, req, pathParams)
		ctx = runtime.NewServerMetadataContext(ctx, md)
		if err != nil {
			runtime.HTTPError(ctx, mux, outboundMarshaler, w, req, err)
			return
		}

		forward_ReflectionService_ListImplementations_0(ctx, mux, outboundMarshaler, w, req, resp, mux.GetForwardResponseOptions()...)

	})

	return nil
}

// RegisterReflectionServiceHandlerFromEndpoint is same as RegisterReflectionServiceHandler but
// automatically dials to "endpoint" and closes the connection when "ctx" gets done.
func RegisterReflectionServiceHandlerFromEndpoint(ctx context.Context, mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) (err error) {
	conn, err := grpc.Dial(endpoint, opts...)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			if cerr := conn.Close(); cerr != nil {
				grpclog.Infof("Failed to close conn to %s: %v", endpoint, cerr)
			}
			return
		}
		go func() {
			<-ctx.Done()
			if cerr := conn.Close(); cerr != nil {
				grpclog.Infof("Failed to close conn to %s: %v", endpoint, cerr)
			}
		}()
	}()

	return RegisterReflectionServiceHandler(ctx, mux, conn)
}

// RegisterReflectionServiceHandler registers the http handlers for service ReflectionService to "mux".
// The handlers forward requests to the grpc endpoint over "conn".
func RegisterReflectionServiceHandler(ctx context.Context, mux *runtime.ServeMux, conn *grpc.ClientConn) error {
	return RegisterReflectionServiceHandlerClient(ctx, mux, NewReflectionServiceClient(conn))
}

// RegisterReflectionServiceHandlerClient registers the http handlers for service ReflectionService
// to "mux". The handlers forward requests to the grpc endpoint over the given implementation of "ReflectionServiceClient".
// Note: the gRPC framework executes interceptors within the gRPC handler. If the passed in "ReflectionServiceClient"
// doesn't go through the normal gRPC flow (creating a gRPC client etc.) then it will be up to the passed in
// "ReflectionServiceClient" to call the correct interceptors.
func RegisterReflectionServiceHandlerClient(ctx context.Context, mux *runtime.ServeMux, client ReflectionServiceClient) error {

	mux.Handle("GET", pattern_ReflectionService_ListAllInterfaces_0, func(w http.ResponseWriter, req *http.Request, pathParams map[string]string) {
		ctx, cancel := context.WithCancel(req.Context())
		defer cancel()
		inboundMarshaler, outboundMarshaler := runtime.MarshalerForRequest(mux, req)
		rctx, err := runtime.AnnotateContext(ctx, mux, req)
		if err != nil {
			runtime.HTTPError(ctx, mux, outboundMarshaler, w, req, err)
			return
		}
		resp, md, err := request_ReflectionService_ListAllInterfaces_0(rctx, inboundMarshaler, client, req, pathParams)
		ctx = runtime.NewServerMetadataContext(ctx, md)
		if err != nil {
			runtime.HTTPError(ctx, mux, outboundMarshaler, w, req, err)
			return
		}

		forward_ReflectionService_ListAllInterfaces_0(ctx, mux, outboundMarshaler, w, req, resp, mux.GetForwardResponseOptions()...)

	})

	mux.Handle("GET", pattern_ReflectionService_ListImplementations_0, func(w http.ResponseWriter, req *http.Request, pathParams map[string]string) {
		ctx, cancel := context.WithCancel(req.Context())
		defer cancel()
		inboundMarshaler, outboundMarshaler := runtime.MarshalerForRequest(mux, req)
		rctx, err := runtime.AnnotateContext(ctx, mux, req)
		if err != nil {
			runtime.HTTPError(ctx, mux, outboundMarshaler, w, req, err)
			return
		}
		resp, md, err := request_ReflectionService_ListImplementations_0(rctx, inboundMarshaler, client, req, pathParams)
		ctx = runtime.NewServerMetadataContext(ctx, md)
		if err != nil {
			runtime.HTTPError(ctx, mux, outboundMarshaler, w, req, err)
			return
		}

		forward_ReflectionService_ListImplementations_0(ctx, mux, outboundMarshaler, w, req, resp, mux.GetForwardResponseOptions()...)

	})

	return nil
}

var (
	pattern_ReflectionService_ListAllInterfaces_0 = runtime.MustPattern(runtime.NewPattern(1, []int{2, 0, 2, 1, 2, 2, 2, 3, 2, 4}, []string{"lbm", "base", "reflection", "v1", "interfaces"}, "", runtime.AssumeColonVerbOpt(true)))

	pattern_ReflectionService_ListImplementations_0 = runtime.MustPattern(runtime.NewPattern(1, []int{2, 0, 2, 1, 2, 2, 2, 3, 2, 4, 1, 0, 4, 1, 5, 5, 2, 6}, []string{"lbm", "base", "reflection", "v1", "interfaces", "interface_name", "implementations"}, "", runtime.AssumeColonVerbOpt(true)))
)

var (
	forward_ReflectionService_ListAllInterfaces_0 = runtime.ForwardResponseMessage

	forward_ReflectionService_ListImplementations_0 = runtime.ForwardResponseMessage
)
