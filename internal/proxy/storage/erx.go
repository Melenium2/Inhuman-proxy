package storage

import (
	"fmt"
)

var (
	ErrOnPush = func(expected, n int) error {
		return fmt.Errorf("expected %d to push, but only %d pushed to redis", expected, n)
	}
)
