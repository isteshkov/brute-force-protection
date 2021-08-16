package context

import (
	"context"

	"gitlab.com/isteshkov/brute-force-protection/domain/logging"
	"google.golang.org/grpc/metadata"
)

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
