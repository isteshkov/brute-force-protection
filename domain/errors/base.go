package errors

import (
	"fmt"
)

type baseError struct {
	msg      string
	typ      string
	code     string
	stack    []string
	origin   error
	producer *ErrorProducer
}

func (b baseError) WithCode(code string) *baseError {
	b.code = code
	if coder, ok := b.origin.(interface{ WithCode(string) *baseError }); ok {
		b.origin = coder.WithCode(code)
	}

	return &b
}

func (b *baseError) Code() string {
	if base, ok := b.origin.(interface{ Code() string }); ok {
		return base.Code()
	}

	return b.code
}

func (b *baseError) Type() string {
	return b.typ
}

func (b *baseError) getOrigin() error {
	return b.origin
}

func (b *baseError) getProducer() *ErrorProducer {
	return b.producer
}

func (b *baseError) getPrettyStack() string {
	result := "\n"
	space := ""
	for lvl := len(b.stack) - 1; lvl >= 0; lvl-- {
		result += fmt.Sprintf("%s%s\n", space, b.stack[lvl])
		space += " "
	}
	return result + space
}

func (b *baseError) Error() string {
	switch b.origin.(type) {
	case *baseError:
		return b.origin.Error()
	default:
		if b.origin != nil {
			return fmt.Sprintf("%s:%s", b.msg, b.origin.Error())
		}
		return b.msg
	}
}

func (b *baseError) ErrorMessage() string {
	switch b.origin.(type) {
	case *baseError:
		return b.origin.Error()
	default:
		if b.origin != nil {
			return b.origin.Error()
		} else {
			return b.msg
		}
	}
}

func (b *baseError) Stacktrace() string {
	return b.getPrettyStack()
}
