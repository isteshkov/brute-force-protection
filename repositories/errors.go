package repositories

import (
	"gitlab.com/isteshkov/brute-force-protection/domain/database"
	"gitlab.com/isteshkov/brute-force-protection/domain/errors"
)

const (
	ErrorTypeGeneral       = "SUBNETS_STORES_GENERAL_ERROR"
	ErrorTypeDoesNotExists = "SUBNETS_STORES_ENTITY_DOES_NOT_EXISTS"
	ErrorTypeAlreadyExists = "SUBNETS_STORES_ENTITY_ALREADY_EXISTS"
	ErrorTypeInconsistent  = "SUBNETS_STORES_INCONSISTENT"
)

var (
	ErrorProducerGeneral       = errors.NewProducer(ErrorTypeGeneral)
	ErrorProducerDoesNotExists = errors.NewProducer(ErrorTypeDoesNotExists)
	ErrorProducerAlreadyExists = errors.NewProducer(ErrorTypeAlreadyExists)
	ErrorProducerInconsistent  = errors.NewProducer(ErrorTypeInconsistent)
)

var ErrorsList = []*errors.ErrorProducer{
	ErrorProducerGeneral,
	ErrorProducerDoesNotExists,
	ErrorProducerAlreadyExists,
	ErrorProducerInconsistent,
}

func processError(errPtr *error) {
	if errPtr == nil || *errPtr == nil {
		return
	}

	err := *errPtr

	if errors.IsProducedBy(err, ErrorsList...) {
		return
	}

	if errors.IsProducedBy(err, database.ErrorsList...) {
		switch {
		case errors.IsProducedBy(err, database.ErrorProducerDoesNotExist):
			*errPtr = ErrorProducerDoesNotExists.Wrap(err)
		case errors.IsProducedBy(err, database.ErrorProducerAlreadyExist):
			*errPtr = ErrorProducerAlreadyExists.Wrap(err)
		default:
			*errPtr = ErrorProducerGeneral.Wrap(err)
		}
		return
	}

	*errPtr = ErrorProducerGeneral.Wrap(err)
}
