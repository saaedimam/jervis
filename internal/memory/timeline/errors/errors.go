package errors

import "fmt"

var (
	// ErrEventImmutable indicates an attempt to modify an existing event.
	ErrEventImmutable = fmt.Errorf("event is immutable: cannot be modified")
)
