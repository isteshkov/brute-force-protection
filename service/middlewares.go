package service

import (
	"context"
	"math"
	"runtime/debug"
	"time"

	"gitlab.com/isteshkov/brute-force-protection/domain/common"
	myContext "gitlab.com/isteshkov/brute-force-protection/domain/context"
	"gitlab.com/isteshkov/brute-force-protection/domain/contract"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func (s *Service) MiddlewareAccess(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	requestId := common.NewUUIDv4()
	meta, ok := metadata.FromIncomingContext(ctx)
	if ok {
		IncomingRequestId := meta.Get(contract.RequestIdHeader)
		if len(IncomingRequestId) > 0 {
			requestId = IncomingRequestId[0]
		}
	}

	path := info.FullMethod
	generalLoggerFields := map[string]interface{}{
		"path":       path,
		"request_id": requestId,
	}

	meta.Append(myContext.KeyRequestId, requestId)
	ctx = metadata.NewOutgoingContext(ctx, meta)

	logger := s.Logger.WithFields(generalLoggerFields)

	tsBeforeProcess := time.Now().UTC()
	defer func() {
		panicInfo := recover()
		if panicInfo != nil {
			logger.WithFields(
				map[string]interface{}{
					"path":         path,
					"request_body": req,
					"meta":         meta,
					"panic_info":   panicInfo,
				}).Fatal("panic recovered")
		}
	}()
	result, err := handler(ctx, req)

	latency := math.Floor(time.Now().UTC().Sub(tsBeforeProcess).Seconds()*1000) / 1000
	logger.WithFields(
		map[string]interface{}{
			"latency":       latency,
			"path":          path,
			"request_body":  req,
			"response_body": result,
			"meta":          meta,
		}).
		Info("ACCESS")

	return result, err
}

func (s *Service) RecoveryMiddleware(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	logger := s.Logger.WithFields(myContext.LogFieldsFromGrpcContext(ctx))
	defer func() {
		panicInfo := recover()
		if panicInfo != nil {
			logger.WithFields(
				map[string]interface{}{
					"panic_info": panicInfo,
					"trace":      string(debug.Stack()),
				}).Fatal("panic recovered")
			err = ErrorProducerGeneral.New("%+v", panicInfo)
		}
		return
	}()

	result, err := handler(ctx, req)

	return result, err
}
