package errors

import "fmt"

// ErrEventImmutable indicates an attempt to modify an existing event.
var ErrEventImmutable = fmt.Errorf("event is immutable: cannot be modified")
