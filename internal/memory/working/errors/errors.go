package errors

import "fmt"

// ErrInvalidCapacity indicates the provided capacity is not positive.
var ErrInvalidCapacity = fmt.Errorf("invalid capacity: must be greater than zero")
