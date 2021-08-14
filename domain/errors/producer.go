package errors

import (
	"fmt"
	"runtime"
	"strconv"
)

func NewProducer(typ string) *ErrorProducer {
	return &ErrorProducer{
		typ,
	}
}

func IsProducedBy(err error, reasons ...*ErrorProducer) bool {
	base, ok := err.(*baseError)
	if !ok {
		return false
	}

	for _, reason := range reasons {
		if base.producer == reason {
			return true
		}
	}
	return false
}

func GetCode(err error) string {
	switch err.(type) {
	case interface{ Code() string }:
		return err.(interface{ Code() string }).Code()
	default:
		return ErrorCodeGeneral
	}
}

type ErrorProducer struct {
	producingType string
}

func (e *ErrorProducer) New(msg string, args ...interface{}) *baseError {
	return &baseError{
		msg:      fmt.Sprintf(msg, args...),
		typ:      e.producingType,
		stack:    getStackTrace(0),
		producer: e,
	}
}

func (e *ErrorProducer) Wrap(err error, lvl ...int) *baseError {
	level := 0
	if len(lvl) == 1 {
		level = lvl[0]
	}

	return &baseError{
		msg:      err.Error(),
		typ:      e.producingType,
		stack:    getStackTrace(level),
		origin:   err,
		producer: e,
	}
}

func (e *ErrorProducer) WrapF(err error, message string, lvl ...int) *baseError {
	level := 0
	if len(lvl) == 1 {
		level = lvl[0]
	}

	return &baseError{
		msg:      fmt.Sprintf(message),
		typ:      e.producingType,
		stack:    getStackTrace(level),
		origin:   err,
		producer: e,
	}
}

func getStackTrace(lvl int) []string {
	var result []string
	lvl += 2
	for {
		pc, _, line, ok := runtime.Caller(lvl)
		if !ok {
			return result
		}
		file, line := runtime.FuncForPC(pc).FileLine(pc)
		result = append(result, file+":"+strconv.Itoa(line))
		lvl++
	}
}
