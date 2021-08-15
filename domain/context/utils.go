package context

import (
	"context"

	"gitlab.com/isteshkov/brute-force-protection/domain/common"
	"gitlab.com/isteshkov/brute-force-protection/domain/logging"
	"google.golang.org/grpc/metadata"
)

func NewGrpcContext(c *Context) (ctx context.Context) {
	md := metadata.New(map[string]string{KeyRequestID: c.RequestID})
	ctx = metadata.NewOutgoingContext(context.Background(), md)
	return
}

func LogFieldsFromGrpcContext(c context.Context) map[string]interface{} {
	result := make(map[string]interface{})
	meta, ok := metadata.FromIncomingContext(c)
	if ok {
		IncomingRequestID := meta.Get(KeyRequestID)
		if len(IncomingRequestID) > 0 {
			result[logging.FieldKeyRequestID] = IncomingRequestID[0]
		}
	}

	return result
}

func NewGrpcFromGrpc(c context.Context) (ctx context.Context) {
	requestID := common.NewUUIDv4()
	meta, ok := metadata.FromIncomingContext(c)
	if ok {
		IncomingRequestID := meta.Get(KeyRequestID)
		if len(IncomingRequestID) > 0 {
			requestID = IncomingRequestID[0]
		}
	}
	md := metadata.New(map[string]string{KeyRequestID: requestID})
	ctx = metadata.NewOutgoingContext(context.Background(), md)
	return
}
