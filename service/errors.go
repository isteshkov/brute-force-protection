package service

import (
	"gitlab.com/isteshkov/brute-force-protection/domain/errors"
	"gitlab.com/isteshkov/brute-force-protection/errorcodes"
)

var (
	ErrorProducerGeneral = errors.NewProducer(errorcodes.CodeGeneral)

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
}

//nolint:errorlint
func (s *Service) processRPCError(errPtr *error, errField *string) {
	if errPtr == nil || *errPtr == nil {
		return
	}
	s.Logger.Error(*errPtr)
	defer func() {
		*errPtr = nil
	}()

	*errField = errorcodes.CodeGeneral
	if err, ok := (*errPtr).(errors.HasType); ok {
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
