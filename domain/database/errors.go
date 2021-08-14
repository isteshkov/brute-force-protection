package database

import (
	"database/sql"

	"gitlab.com/isteshkov/brute-force-protection/domain/errors"

	"github.com/lib/pq"
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

func processError(errPtr *error) {
	if errPtr == nil || *errPtr == nil {
		return
	}

	err := *errPtr

	if errors.IsProducedBy(err, ErrorsList...) {
		return
	}

	if pqErr, ok := err.(*pq.Error); ok {
		switch pqErr.Code {
		//pq errors code corresponding duplicate primary key
		case "23505":
			*errPtr = ErrorProducerAlreadyExist.Wrap(err, 1)
			return
		}
	}

	if err == sql.ErrNoRows {
		*errPtr = ErrorProducerDoesNotExist.Wrap(err, 1)
		return
	}

	*errPtr = ErrorProducerGeneral.Wrap(err, 1)
	return
}
