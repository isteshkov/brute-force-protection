package context

import (
	"context"

	"gitlab.com/isteshkov/brute-force-protection/domain/common"
	"gitlab.com/isteshkov/brute-force-protection/domain/logging"

	"google.golang.org/grpc/metadata"
)

func NewGrpcContext(c *Context) (ctx context.Context) {
	md := metadata.New(map[string]string{KeyRequestId: c.RequestId})
	ctx = metadata.NewOutgoingContext(context.Background(), md)
	return
}

func LogFieldsFromGrpcContext(c context.Context) map[string]interface{} {
	result := make(map[string]interface{})
	meta, ok := metadata.FromIncomingContext(c)
	if ok {
		IncomingRequestId := meta.Get(KeyRequestId)
		if len(IncomingRequestId) > 0 {
			result[logging.FieldKeyRequestId] = IncomingRequestId[0]
		}
	}

	return result
}

func NewGrpcFromGrpc(c context.Context) (ctx context.Context) {
	requestId := common.NewUUIDv4()
	meta, ok := metadata.FromIncomingContext(c)
	if ok {
		IncomingRequestId := meta.Get(KeyRequestId)
		if len(IncomingRequestId) > 0 {
			requestId = IncomingRequestId[0]
		}
	}
	md := metadata.New(map[string]string{KeyRequestId: requestId})
	ctx = metadata.NewOutgoingContext(context.Background(), md)
	return
}
