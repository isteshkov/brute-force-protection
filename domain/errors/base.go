package errors

import (
	"fmt"
)

type withCoder interface {
	WithCode(string) *baseError
}

type coder interface {
	Code() string
}

type baseError struct {
	msg      string
	typ      string
	code     string
	stack    []string
	origin   error
	producer *ErrorProducer
}

//nolint:errorlint
func (b baseError) WithCode(code string) *baseError {
	b.code = code
	if coder, ok := b.origin.(withCoder); ok {
		b.origin = coder.WithCode(code)
	}

	return &b
}

//nolint:errorlint
func (b *baseError) Code() string {
	if base, ok := b.origin.(coder); ok {
		return base.Code()
	}

	return b.code
}

func (b *baseError) Type() string {
	return b.typ
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

//nolint:errorlint
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

//nolint:errorlint
func (b *baseError) ErrorMessage() string {
	switch b.origin.(type) {
	case *baseError:
		return b.origin.Error()
	default:
		if b.origin != nil {
			return b.origin.Error()
		}
		return b.msg
	}
}

func (b *baseError) Stacktrace() string {
	return b.getPrettyStack()
}
