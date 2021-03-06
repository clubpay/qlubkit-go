package qkit

import (
	"errors"
	"fmt"
)

type qError struct {
	top  error
	down error
}

func (e qError) Error() string {
	if e.down == nil {
		return e.top.Error()
	}

	return fmt.Sprintf("%s: %s", e.top, e.down)
}

func (e qError) Is(err error) bool {
	return errors.Is(e.top, err)
}

func (e qError) Unwrap() error {
	return e.down
}

func WrapError(top, down error) error {
	if top == nil {
		return down
	}

	return qError{
		top:  top,
		down: down,
	}
}
