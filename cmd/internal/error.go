package internal

import "fmt"

type Error struct {
	orig error
	msg  string
}

func New(msg string) error {
	return Wrap(nil, msg)
}

func Errorf(format string, a ...interface{}) error {
	return Wrapf(nil, format, a...)
}

func Wrap(orig error, msg string) error {
	return &Error{
		orig: orig,
		msg:  msg,
	}
}

func Wrapf(orig error, format string, a ...interface{}) error {
	return &Error{
		orig: orig,
		msg:  fmt.Sprintf(format, a...),
	}
}

func (e *Error) Error() string {
	if e.orig == nil {
		return e.msg
	}

	return e.msg + ": " + e.orig.Error()
}

func (e *Error) Unwrap() error {
	return e.orig
}
