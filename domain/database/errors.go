package database

import (
	"database/sql"
	stdErrors "errors"

	"github.com/lib/pq"
	"gitlab.com/isteshkov/brute-force-protection/domain/errors"
)

var (
	ErrorProducerGeneral = errors.NewProducer("some error")

	ErrorProducerDoesNotExist = errors.NewProducer("entity does not exist")
	ErrorProducerAlreadyExist = errors.NewProducer("entity already exist")
)

var ErrorsList = []*errors.ErrorProducer{
	ErrorProducerGeneral,
	ErrorProducerDoesNotExist,
	ErrorProducerAlreadyExist,
}

//nolint:errorlint
func processError(errPtr *error) {
	if errPtr == nil || *errPtr == nil {
		return
	}

	err := *errPtr

	if errors.IsProducedBy(err, ErrorsList...) {
		return
	}

	if pqErr, ok := err.(*pq.Error); ok {
		// pq errors code corresponding duplicate primary key
		if pqErr.Code == "23505" {
			*errPtr = ErrorProducerAlreadyExist.Wrap(err, 1)
			return
		}
	}

	if stdErrors.Is(err, sql.ErrNoRows) {
		*errPtr = ErrorProducerDoesNotExist.Wrap(err, 1)
		return
	}

	*errPtr = ErrorProducerGeneral.Wrap(err, 1)
}
