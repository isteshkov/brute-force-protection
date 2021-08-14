package errors

type HasType interface {
	Type() string
}

type HasCode interface {
	Code() string
}

type HasMsg interface {
	ErrorMessage() string
}
