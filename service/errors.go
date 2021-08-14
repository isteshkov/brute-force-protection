package service

import (
	"gitlab.com/isteshkov/brute-force-protection/domain/errors"
	"gitlab.com/isteshkov/brute-force-protection/error_codes"
)

var (
	ErrorProducerGeneral = errors.NewProducer(error_codes.CodeGeneral)

	ErrorsList = []*errors.ErrorProducer{
		ErrorProducerGeneral,
	}
)

func processError(errPtr *error) {
	if errPtr == nil || *errPtr == nil {
		return
	}

	err := *errPtr

	if errors.IsProducedBy(err, ErrorsList...) {
		return
	}

	*errPtr = ErrorProducerGeneral.Wrap(err)
	return
}

func (s *Service) processRpcError(errPtr *error, errField *string) {
	if errPtr == nil || *errPtr == nil {
		return
	}
	s.Logger.Error(*errPtr)
	defer func() {
		*errPtr = nil
	}()

	*errField = error_codes.CodeGeneral
	if err, ok := (*errPtr).(errors.HasType); ok {
		if !ok {
			return
		}

		if withCode, ok := err.(errors.HasCode); ok {
			if withCode.Code() != "" {
				*errField = withCode.Code()
			} else {
				*errField = err.Type()
			}
		}
		return
	}
}
