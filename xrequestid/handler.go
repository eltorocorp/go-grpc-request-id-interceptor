package xrequestid

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	multiint "github.com/eltorocorp/go-grpc-request-id-interceptor/multiinterceptor"
)

type requestIDKey struct{}

type serverStreamWrapper struct {
	ss  grpc.ServerStream
	ctx context.Context
}

func (w serverStreamWrapper) Context() context.Context        { return w.ctx }
func (w serverStreamWrapper) RecvMsg(msg interface{}) error   { return w.ss.RecvMsg(msg) }
func (w serverStreamWrapper) SendMsg(msg interface{}) error   { return w.ss.SendMsg(msg) }
func (w serverStreamWrapper) SendHeader(md metadata.MD) error { return w.ss.SendHeader(md) }
func (w serverStreamWrapper) SetHeader(md metadata.MD) error  { return w.ss.SetHeader(md) }
func (w serverStreamWrapper) SetTrailer(md metadata.MD)       { w.ss.SetTrailer(md) }

func UnaryServerInterceptor(opt ...Option) grpc.UnaryServerInterceptor {
	var opts options
	opts.validator = defaultReqeustIDValidator
	for _, o := range opt {
		o.apply(&opts)
	}

	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		var requestID string
		if opts.chainRequestID {
			requestID = HandleRequestIDChain(ctx, opts.validator)
		} else {
			requestID = HandleRequestID(ctx, opts.validator)
		}
		if opts.logRequest {
			ctx = addRequestToLogger(ctx, requestID, req)
		}
		ctx = context.WithValue(ctx, requestIDKey{}, requestID)
		for _, header := range opts.persistHeaders {
			headerValue := getStringFromContext(ctx, header)
			if header == DefaultXRequestIDKey {
				headerValue = requestID
			}
			ctx = metadata.AppendToOutgoingContext(ctx, header, headerValue)
		}
		return handler(ctx, req)
	}
}

func StreamServerInterceptor(opt ...Option) grpc.StreamServerInterceptor {
	var opts options
	opts.validator = defaultReqeustIDValidator
	for _, o := range opt {
		o.apply(&opts)
	}

	return func(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) (err error) {
		ctx := stream.Context()
		// tests
		var requestID string
		if opts.chainRequestID {
			requestID = HandleRequestIDChain(ctx, opts.validator)
		} else {
			requestID = HandleRequestID(ctx, opts.validator)
		}
		if opts.logRequest {
			ctx = addRequestToLogger(ctx, requestID, "stream_data")
		}
		ctx = context.WithValue(ctx, requestIDKey{}, requestID)
		stream = multiint.NewServerStreamWithContext(stream, ctx)
		for _, header := range opts.persistHeaders {
			headerValue := getStringFromContext(ctx, header)
			if header == DefaultXRequestIDKey {
				headerValue = requestID
			}
			ctx = metadata.AppendToOutgoingContext(ctx, header, headerValue)
		}
		return handler(srv, serverStreamWrapper{stream, ctx})
	}
}

func FromContext(ctx context.Context) string {
	id, ok := ctx.Value(requestIDKey{}).(string)
	if !ok {
		return ""
	}
	return id
}

// Create a context with the private requestIDKey{} for testing
func ContextWithID(ctx context.Context, requestID string) context.Context {
	return context.WithValue(ctx, requestIDKey{}, requestID)
}

// Gets the key value from the provided context, or returns an empty string
func getStringFromContext(ctx context.Context, key string) string {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return ""
	}

	header, ok := md[key]
	if !ok || len(header) == 0 {
		return ""
	}

	return header[0]
}
