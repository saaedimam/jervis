package middleware

import (
	"fmt"

	"github.com/saaedimam/jervis/internal/runtime/eventbus/contracts"
	errs "github.com/saaedimam/jervis/internal/runtime/eventbus/errors"
)

// Chain manages and executes an ordered sequence of middleware interceptors.
type Chain struct {
	middlewares []contracts.Middleware
}

// NewChain constructs an initialized Chain containing the provided middleware.
func NewChain(mw ...contracts.Middleware) *Chain {
	c := &Chain{
		middlewares: make([]contracts.Middleware, 0, len(mw)),
	}
	c.Use(mw...)
	return c
}

// Use appends one or more middleware interceptors to the chain in FIFO order.
func (c *Chain) Use(mw ...contracts.Middleware) {
	for _, m := range mw {
		if m != nil {
			c.middlewares = append(c.middlewares, m)
		}
	}
}

// Middlewares returns a defensive copy slice of all registered middleware in FIFO order.
func (c *Chain) Middlewares() []contracts.Middleware {
	if len(c.middlewares) == 0 {
		return nil
	}
	cp := make([]contracts.Middleware, len(c.middlewares))
	copy(cp, c.middlewares)
	return cp
}

// Count returns the total number of registered middleware interceptors.
func (c *Chain) Count() int {
	return len(c.middlewares)
}

// Execute constructs the recursive onion closure chain and executes it, ending at terminal.
func (c *Chain) Execute(event contracts.Event, terminal func(event contracts.Event) error) error {
	if terminal == nil {
		return fmt.Errorf("%w: terminal execution function cannot be nil", errs.ErrValidationFailed)
	}

	next := terminal
	for i := len(c.middlewares) - 1; i >= 0; i-- {
		mw := c.middlewares[i]
		currNext := next

		next = func(evt contracts.Event) (err error) {
			defer func() {
				if r := recover(); r != nil {
					err = fmt.Errorf("%w: middleware panicked: %v", errs.ErrHandlerFailure, r)
				}
			}()
			return mw.Execute(evt, currNext)
		}
	}

	return next(event)
}
