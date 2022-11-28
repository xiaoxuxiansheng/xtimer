package utils

import "errors"

var ErrRetryable = errors.New("need retry")

func ErrNeedRetry(err error) bool {
	return errors.Is(err, ErrRetryable)
}
