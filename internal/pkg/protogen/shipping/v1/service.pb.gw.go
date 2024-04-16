// Code generated by protoc-gen-grpc-gateway. DO NOT EDIT.
// source: stocklet/shipping/v1/service.proto

/*
Package shipping_v1 is a reverse proxy.

It translates gRPC into RESTful JSON APIs.
*/
package shipping_v1

import (
	"context"
	"io"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/grpc-ecosystem/grpc-gateway/v2/utilities"
	"github.com/hexolan/stocklet/internal/pkg/protogen/common/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
)

// Suppress "imported and not used" errors
var _ codes.Code
var _ io.Reader
var _ status.Status
var _ = runtime.String
var _ = utilities.NewDoubleArray
var _ = metadata.Join

func request_ShippingService_ServiceInfo_0(ctx context.Context, marshaler runtime.Marshaler, client ShippingServiceClient, req *http.Request, pathParams map[string]string) (proto.Message, runtime.ServerMetadata, error) {
	var protoReq common_v1.ServiceInfoRequest
	var metadata runtime.ServerMetadata

	msg, err := client.ServiceInfo(ctx, &protoReq, grpc.Header(&metadata.HeaderMD), grpc.Trailer(&metadata.TrailerMD))
	return msg, metadata, err

}

func local_request_ShippingService_ServiceInfo_0(ctx context.Context, marshaler runtime.Marshaler, server ShippingServiceServer, req *http.Request, pathParams map[string]string) (proto.Message, runtime.ServerMetadata, error) {
	var protoReq common_v1.ServiceInfoRequest
	var metadata runtime.ServerMetadata

	msg, err := server.ServiceInfo(ctx, &protoReq)
	return msg, metadata, err

}

func request_ShippingService_ViewShipment_0(ctx context.Context, marshaler runtime.Marshaler, client ShippingServiceClient, req *http.Request, pathParams map[string]string) (proto.Message, runtime.ServerMetadata, error) {
	var protoReq ViewShipmentRequest
	var metadata runtime.ServerMetadata

	var (
		val string
		ok  bool
		err error
		_   = err
	)

	val, ok = pathParams["shipment_id"]
	if !ok {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "missing parameter %s", "shipment_id")
	}

	protoReq.ShipmentId, err = runtime.String(val)
	if err != nil {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "type mismatch, parameter: %s, error: %v", "shipment_id", err)
	}

	msg, err := client.ViewShipment(ctx, &protoReq, grpc.Header(&metadata.HeaderMD), grpc.Trailer(&metadata.TrailerMD))
	return msg, metadata, err

}

func local_request_ShippingService_ViewShipment_0(ctx context.Context, marshaler runtime.Marshaler, server ShippingServiceServer, req *http.Request, pathParams map[string]string) (proto.Message, runtime.ServerMetadata, error) {
	var protoReq ViewShipmentRequest
	var metadata runtime.ServerMetadata

	var (
		val string
		ok  bool
		err error
		_   = err
	)

	val, ok = pathParams["shipment_id"]
	if !ok {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "missing parameter %s", "shipment_id")
	}

	protoReq.ShipmentId, err = runtime.String(val)
	if err != nil {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "type mismatch, parameter: %s, error: %v", "shipment_id", err)
	}

	msg, err := server.ViewShipment(ctx, &protoReq)
	return msg, metadata, err

}

func request_ShippingService_ViewShipmentManifest_0(ctx context.Context, marshaler runtime.Marshaler, client ShippingServiceClient, req *http.Request, pathParams map[string]string) (proto.Message, runtime.ServerMetadata, error) {
	var protoReq ViewShipmentManifestRequest
	var metadata runtime.ServerMetadata

	var (
		val string
		ok  bool
		err error
		_   = err
	)

	val, ok = pathParams["shipment_id"]
	if !ok {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "missing parameter %s", "shipment_id")
	}

	protoReq.ShipmentId, err = runtime.String(val)
	if err != nil {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "type mismatch, parameter: %s, error: %v", "shipment_id", err)
	}

	msg, err := client.ViewShipmentManifest(ctx, &protoReq, grpc.Header(&metadata.HeaderMD), grpc.Trailer(&metadata.TrailerMD))
	return msg, metadata, err

}

func local_request_ShippingService_ViewShipmentManifest_0(ctx context.Context, marshaler runtime.Marshaler, server ShippingServiceServer, req *http.Request, pathParams map[string]string) (proto.Message, runtime.ServerMetadata, error) {
	var protoReq ViewShipmentManifestRequest
	var metadata runtime.ServerMetadata

	var (
		val string
		ok  bool
		err error
		_   = err
	)

	val, ok = pathParams["shipment_id"]
	if !ok {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "missing parameter %s", "shipment_id")
	}

	protoReq.ShipmentId, err = runtime.String(val)
	if err != nil {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "type mismatch, parameter: %s, error: %v", "shipment_id", err)
	}

	msg, err := server.ViewShipmentManifest(ctx, &protoReq)
	return msg, metadata, err

}

// RegisterShippingServiceHandlerServer registers the http handlers for service ShippingService to "mux".
// UnaryRPC     :call ShippingServiceServer directly.
// StreamingRPC :currently unsupported pending https://github.com/grpc/grpc-go/issues/906.
// Note that using this registration option will cause many gRPC library features to stop working. Consider using RegisterShippingServiceHandlerFromEndpoint instead.
func RegisterShippingServiceHandlerServer(ctx context.Context, mux *runtime.ServeMux, server ShippingServiceServer) error {

	mux.Handle("GET", pattern_ShippingService_ServiceInfo_0, func(w http.ResponseWriter, req *http.Request, pathParams map[string]string) {
		ctx, cancel := context.WithCancel(req.Context())
		defer cancel()
		var stream runtime.ServerTransportStream
		ctx = grpc.NewContextWithServerTransportStream(ctx, &stream)
		inboundMarshaler, outboundMarshaler := runtime.MarshalerForRequest(mux, req)
		var err error
		var annotatedContext context.Context
		annotatedContext, err = runtime.AnnotateIncomingContext(ctx, mux, req, "/stocklet.shipping.v1.ShippingService/ServiceInfo", runtime.WithHTTPPathPattern("/v1/shipping/service"))
		if err != nil {
			runtime.HTTPError(ctx, mux, outboundMarshaler, w, req, err)
			return
		}
		resp, md, err := local_request_ShippingService_ServiceInfo_0(annotatedContext, inboundMarshaler, server, req, pathParams)
		md.HeaderMD, md.TrailerMD = metadata.Join(md.HeaderMD, stream.Header()), metadata.Join(md.TrailerMD, stream.Trailer())
		annotatedContext = runtime.NewServerMetadataContext(annotatedContext, md)
		if err != nil {
			runtime.HTTPError(annotatedContext, mux, outboundMarshaler, w, req, err)
			return
		}

		forward_ShippingService_ServiceInfo_0(annotatedContext, mux, outboundMarshaler, w, req, resp, mux.GetForwardResponseOptions()...)

	})

	mux.Handle("GET", pattern_ShippingService_ViewShipment_0, func(w http.ResponseWriter, req *http.Request, pathParams map[string]string) {
		ctx, cancel := context.WithCancel(req.Context())
		defer cancel()
		var stream runtime.ServerTransportStream
		ctx = grpc.NewContextWithServerTransportStream(ctx, &stream)
		inboundMarshaler, outboundMarshaler := runtime.MarshalerForRequest(mux, req)
		var err error
		var annotatedContext context.Context
		annotatedContext, err = runtime.AnnotateIncomingContext(ctx, mux, req, "/stocklet.shipping.v1.ShippingService/ViewShipment", runtime.WithHTTPPathPattern("/v1/shipping/shipment/{shipment_id}"))
		if err != nil {
			runtime.HTTPError(ctx, mux, outboundMarshaler, w, req, err)
			return
		}
		resp, md, err := local_request_ShippingService_ViewShipment_0(annotatedContext, inboundMarshaler, server, req, pathParams)
		md.HeaderMD, md.TrailerMD = metadata.Join(md.HeaderMD, stream.Header()), metadata.Join(md.TrailerMD, stream.Trailer())
		annotatedContext = runtime.NewServerMetadataContext(annotatedContext, md)
		if err != nil {
			runtime.HTTPError(annotatedContext, mux, outboundMarshaler, w, req, err)
			return
		}

		forward_ShippingService_ViewShipment_0(annotatedContext, mux, outboundMarshaler, w, req, resp, mux.GetForwardResponseOptions()...)

	})

	mux.Handle("GET", pattern_ShippingService_ViewShipmentManifest_0, func(w http.ResponseWriter, req *http.Request, pathParams map[string]string) {
		ctx, cancel := context.WithCancel(req.Context())
		defer cancel()
		var stream runtime.ServerTransportStream
		ctx = grpc.NewContextWithServerTransportStream(ctx, &stream)
		inboundMarshaler, outboundMarshaler := runtime.MarshalerForRequest(mux, req)
		var err error
		var annotatedContext context.Context
		annotatedContext, err = runtime.AnnotateIncomingContext(ctx, mux, req, "/stocklet.shipping.v1.ShippingService/ViewShipmentManifest", runtime.WithHTTPPathPattern("/v1/shipping/shipment/{shipment_id}/manifest"))
		if err != nil {
			runtime.HTTPError(ctx, mux, outboundMarshaler, w, req, err)
			return
		}
		resp, md, err := local_request_ShippingService_ViewShipmentManifest_0(annotatedContext, inboundMarshaler, server, req, pathParams)
		md.HeaderMD, md.TrailerMD = metadata.Join(md.HeaderMD, stream.Header()), metadata.Join(md.TrailerMD, stream.Trailer())
		annotatedContext = runtime.NewServerMetadataContext(annotatedContext, md)
		if err != nil {
			runtime.HTTPError(annotatedContext, mux, outboundMarshaler, w, req, err)
			return
		}

		forward_ShippingService_ViewShipmentManifest_0(annotatedContext, mux, outboundMarshaler, w, req, resp, mux.GetForwardResponseOptions()...)

	})

	return nil
}

// RegisterShippingServiceHandlerFromEndpoint is same as RegisterShippingServiceHandler but
// automatically dials to "endpoint" and closes the connection when "ctx" gets done.
func RegisterShippingServiceHandlerFromEndpoint(ctx context.Context, mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) (err error) {
	conn, err := grpc.DialContext(ctx, endpoint, opts...)
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

	return RegisterShippingServiceHandler(ctx, mux, conn)
}

// RegisterShippingServiceHandler registers the http handlers for service ShippingService to "mux".
// The handlers forward requests to the grpc endpoint over "conn".
func RegisterShippingServiceHandler(ctx context.Context, mux *runtime.ServeMux, conn *grpc.ClientConn) error {
	return RegisterShippingServiceHandlerClient(ctx, mux, NewShippingServiceClient(conn))
}

// RegisterShippingServiceHandlerClient registers the http handlers for service ShippingService
// to "mux". The handlers forward requests to the grpc endpoint over the given implementation of "ShippingServiceClient".
// Note: the gRPC framework executes interceptors within the gRPC handler. If the passed in "ShippingServiceClient"
// doesn't go through the normal gRPC flow (creating a gRPC client etc.) then it will be up to the passed in
// "ShippingServiceClient" to call the correct interceptors.
func RegisterShippingServiceHandlerClient(ctx context.Context, mux *runtime.ServeMux, client ShippingServiceClient) error {

	mux.Handle("GET", pattern_ShippingService_ServiceInfo_0, func(w http.ResponseWriter, req *http.Request, pathParams map[string]string) {
		ctx, cancel := context.WithCancel(req.Context())
		defer cancel()
		inboundMarshaler, outboundMarshaler := runtime.MarshalerForRequest(mux, req)
		var err error
		var annotatedContext context.Context
		annotatedContext, err = runtime.AnnotateContext(ctx, mux, req, "/stocklet.shipping.v1.ShippingService/ServiceInfo", runtime.WithHTTPPathPattern("/v1/shipping/service"))
		if err != nil {
			runtime.HTTPError(ctx, mux, outboundMarshaler, w, req, err)
			return
		}
		resp, md, err := request_ShippingService_ServiceInfo_0(annotatedContext, inboundMarshaler, client, req, pathParams)
		annotatedContext = runtime.NewServerMetadataContext(annotatedContext, md)
		if err != nil {
			runtime.HTTPError(annotatedContext, mux, outboundMarshaler, w, req, err)
			return
		}

		forward_ShippingService_ServiceInfo_0(annotatedContext, mux, outboundMarshaler, w, req, resp, mux.GetForwardResponseOptions()...)

	})

	mux.Handle("GET", pattern_ShippingService_ViewShipment_0, func(w http.ResponseWriter, req *http.Request, pathParams map[string]string) {
		ctx, cancel := context.WithCancel(req.Context())
		defer cancel()
		inboundMarshaler, outboundMarshaler := runtime.MarshalerForRequest(mux, req)
		var err error
		var annotatedContext context.Context
		annotatedContext, err = runtime.AnnotateContext(ctx, mux, req, "/stocklet.shipping.v1.ShippingService/ViewShipment", runtime.WithHTTPPathPattern("/v1/shipping/shipment/{shipment_id}"))
		if err != nil {
			runtime.HTTPError(ctx, mux, outboundMarshaler, w, req, err)
			return
		}
		resp, md, err := request_ShippingService_ViewShipment_0(annotatedContext, inboundMarshaler, client, req, pathParams)
		annotatedContext = runtime.NewServerMetadataContext(annotatedContext, md)
		if err != nil {
			runtime.HTTPError(annotatedContext, mux, outboundMarshaler, w, req, err)
			return
		}

		forward_ShippingService_ViewShipment_0(annotatedContext, mux, outboundMarshaler, w, req, resp, mux.GetForwardResponseOptions()...)

	})

	mux.Handle("GET", pattern_ShippingService_ViewShipmentManifest_0, func(w http.ResponseWriter, req *http.Request, pathParams map[string]string) {
		ctx, cancel := context.WithCancel(req.Context())
		defer cancel()
		inboundMarshaler, outboundMarshaler := runtime.MarshalerForRequest(mux, req)
		var err error
		var annotatedContext context.Context
		annotatedContext, err = runtime.AnnotateContext(ctx, mux, req, "/stocklet.shipping.v1.ShippingService/ViewShipmentManifest", runtime.WithHTTPPathPattern("/v1/shipping/shipment/{shipment_id}/manifest"))
		if err != nil {
			runtime.HTTPError(ctx, mux, outboundMarshaler, w, req, err)
			return
		}
		resp, md, err := request_ShippingService_ViewShipmentManifest_0(annotatedContext, inboundMarshaler, client, req, pathParams)
		annotatedContext = runtime.NewServerMetadataContext(annotatedContext, md)
		if err != nil {
			runtime.HTTPError(annotatedContext, mux, outboundMarshaler, w, req, err)
			return
		}

		forward_ShippingService_ViewShipmentManifest_0(annotatedContext, mux, outboundMarshaler, w, req, resp, mux.GetForwardResponseOptions()...)

	})

	return nil
}

var (
	pattern_ShippingService_ServiceInfo_0 = runtime.MustPattern(runtime.NewPattern(1, []int{2, 0, 2, 1, 2, 2}, []string{"v1", "shipping", "service"}, ""))

	pattern_ShippingService_ViewShipment_0 = runtime.MustPattern(runtime.NewPattern(1, []int{2, 0, 2, 1, 2, 2, 1, 0, 4, 1, 5, 3}, []string{"v1", "shipping", "shipment", "shipment_id"}, ""))

	pattern_ShippingService_ViewShipmentManifest_0 = runtime.MustPattern(runtime.NewPattern(1, []int{2, 0, 2, 1, 2, 2, 1, 0, 4, 1, 5, 3, 2, 4}, []string{"v1", "shipping", "shipment", "shipment_id", "manifest"}, ""))
)

var (
	forward_ShippingService_ServiceInfo_0 = runtime.ForwardResponseMessage

	forward_ShippingService_ViewShipment_0 = runtime.ForwardResponseMessage

	forward_ShippingService_ViewShipmentManifest_0 = runtime.ForwardResponseMessage
)
