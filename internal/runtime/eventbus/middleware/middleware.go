package middleware

import (
	"github.com/saaedimam/jervis/internal/runtime/eventbus/contracts"
)

// Func allows a function with the matching signature to satisfy contracts.Middleware.
type Func func(event contracts.Event, next func(event contracts.Event) error) error

// Execute calls the underlying function implementation.
func (f Func) Execute(event contracts.Event, next func(event contracts.Event) error) error {
	return f(event, next)
}
