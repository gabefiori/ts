package utils

import (
	"errors"
	"fmt"
)

func NewErrorWithPrefix(prefix string, err error) error {
	errMsg := fmt.Sprintf("Error(%s): %s", prefix, err.Error())

	return errors.New(errMsg)
}
