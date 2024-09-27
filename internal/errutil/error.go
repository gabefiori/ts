package errutil

import (
	"fmt"
)

type ErrorType string

const (
	TmuxErr        ErrorType = "Tmux"
	SessionizerErr           = "Sessionizer"
	SelectorErr              = "Selector"
	ConfigErr                = "Config"
)

type Error struct {
	Type ErrorType
	Err  error
}

func NewError(et ErrorType, err error) *Error {
	return &Error{Type: et, Err: err}
}

func (ce *Error) Error() string {
	return fmt.Sprintf("Error(%s): %q", ce.Type, ce.Err)
}
